package bhf
import (
    "log"
    "os"
    "reflect"
)

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

// This method simulates the branch hash function approach
// to efficiently copy the branch ID from node B to A.
func Synchronize(a, b DB, id uint64, allow_suggest bool) (requests_sent int) {
    //logger.Printf("SYNC id=%v suggest=%v", id, allow_suggest)
    if a.Get(id) == nil {
        // The fetch the branch
        requests_sent += 1
        for _, post := range HandleGetBranch(b, id) {
            a.Put(post)
        }
        return 1
    } else {
        requests_sent += 1
        response := HandleCompareBranches(b, id, a.Hash(id), allow_suggest)
        switch rsp := response.(type) {
            case *RspSuggest:
                if a.Get(rsp.parent) != nil {
                    // we have the parent, download the suggested message and retry
                    requests_sent += 1
                    for _, post := range HandleGetBranch(b, rsp.id) {
                        a.Put(post)
                    }
                    requests_sent += Synchronize(a, b, id, true)
                } else {
                    // Invalid parent, retry without suggestions
                    requests_sent += Synchronize(a, b, id, false)
                }

            case *RspChildren:
                ack_chan := make(chan int)
                acks := 0
                for _, child := range rsp.children {
                    if child.Hash != a.Hash(child.Id) {
                        // This child needs synchronization. Run the procedure asynchronously
                        acks += 1
                        tmp := child.Id
                        go func(){
                            ack_chan <- Synchronize(a, b, tmp, true)
                        }()
                    }
                }
                // Await for the asynchronous calls to finish
                for i := 0; i < acks; i += 1 {
                    requests_sent += (<- ack_chan)
                }

            case *RspBranchInSync:
            // do nothing

            case *RspBranchNotFound:
            // also do nothing

        }
    }
    return
}

func HandleCompareBranches(db DB, id, hash uint64, allow_suggest bool) Response {
    result := do_compare_branches(db, id, hash, allow_suggest)
    logger.Printf("CMP id=%v hash=%v allow=%v | Result=(%v) %+v", id, hash, allow_suggest, reflect.TypeOf(result), result)
    return result
}

func do_compare_branches(db DB, id, hash uint64, allow_suggest bool) Response {
    if db.Get(id) == nil { return &RspBranchNotFound{} }

    diff := db.Hash(id) ^ hash
    if diff == 0 { return &RspBranchInSync{} }

    if allow_suggest {
        hits := db.Find(diff)
        if len(hits) > 0 {
            return &RspSuggest{ hits[0], db.Get(hits[0]).Parent }
        }
    }

    children := db.Children(id)
    foo := make([]IdAndHash, len(children))
    for i := 0; i < len(foo); i += 1 {
        foo[i].Id = children[i]
        foo[i].Hash = db.Hash(foo[i].Id)
    }

    return &RspChildren{foo}
}

func HandleGetBranch(db DB, id uint64) []Post {
    result := do_get_branch(db, id)
    logger.Printf("GET id=%v | Result=%+v", id, result)
    return result
}

func do_get_branch(db DB, id uint64) []Post {
    result := make([]Post, 0)
    if post := db.Get(id); post != nil {
        result = append(result, *post)
        for _, child := range db.Children(id) {
            result = append(result, do_get_branch(db, child)...)
        }
    }
    return result
}


///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

type Request interface{}

// Check if the branches are the same
type ReqCompareBranches struct {
    id uint64
    hash uint64
    allow_suggest bool
}

// Fetch the post from remote node
type ReqGetBranch struct {
    id uint64
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

type Response interface{}

// Branch appears the same
type RspBranchInSync struct { }

// Suggested difference
type RspSuggest struct {
    id uint64
    parent uint64
}

// List of all children of the post
type RspChildren struct {
    children []IdAndHash
}

type RspBranchFound struct {
    posts []Post
}

type RspBranchNotFound struct {}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
