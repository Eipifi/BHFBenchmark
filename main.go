package main
import (
    "aldeon/aldeon"
    "fmt"
    "math/rand"
    "time"
)


func main() {

    rand.Seed(time.Now().UTC().UnixNano())

    // Two databases
    db_local := aldeon.NewDB()
    db_remote := aldeon.NewDB()

    // The root of the conversation
    root := aldeon.Post{Id: 42, Parent: 0}

    // Generate the test instance
    db_remote.Put(root)
    aldeon.GenerateRandomBalanced(db_remote, root.Id, 5, 4) // tree of depth 5 and width 4
    aldeon.CopyPartially(db_local, db_remote, root.Id, 0.5) // copy each branch with 50% chance

    local_size := db_local.Size()
    remote_size := db_remote.Size()

    requests_sent := aldeon.Synchronize(db_local, db_remote, root.Id, true)
    fmt.Println("Requests sent: ", requests_sent)
    fmt.Println("Posts copied: ", remote_size - local_size)
}
