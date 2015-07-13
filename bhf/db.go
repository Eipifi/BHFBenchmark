package bhf
import (
    "errors"
    "sync"
    "bytes"
    "fmt"
)

var ErrPostPresent = errors.New("Post present")
var ErrPostUnknown = errors.New("Post unknown")

type ActualDB struct {
    posts map[uint64] Post
    hashes map[uint64] uint64
    mtx sync.RWMutex
}

func NewDB() *ActualDB {
    db := &ActualDB{}
    db.posts = make(map[uint64] Post)
    db.hashes = make(map[uint64] uint64)
    return db
}

func (d *ActualDB) Get(id uint64) *Post {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    if msg, ok := d.posts[id]; ok { return &msg }
    return nil
}

func (d *ActualDB) contains(id uint64) bool {
    if id == 0 { return true }
    _, ok := d.posts[id]
    return ok
}

func (d *ActualDB) Put(p Post) {
    d.mtx.Lock()
    defer d.mtx.Unlock()
    if d.contains(p.Id) { panic(ErrPostPresent) }
    if ! d.contains(p.Parent) { panic(ErrPostUnknown) }

    d.posts[p.Id] = p
    d.hashes[p.Id] = p.Id
    current := p.Parent

    for current != 0 {
        d.hashes[current] = d.hashes[current] ^ p.Id
        current = d.posts[current].Parent
    }
}

func (d *ActualDB) Find(hash uint64) []uint64 {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    found := make([]uint64, 0)
    for k, v := range d.hashes {
        if hash == v {
            found = append(found, k)
        }
    }
    return found
}

func (d *ActualDB) Hash(id uint64) uint64 {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    hash, ok := d.hashes[id]
    if ok { return hash }
    return 0
}

func (d *ActualDB) Children(id uint64) []uint64 {
    d.mtx.RLock()
    defer d.mtx.RUnlock()
    found := make([]uint64, 0)
    for _, v := range d.posts {
        if v.Parent == id {
            found = append(found, v.Id)
        }
    }
    return found
}

func (d *ActualDB) String() string {
    var buffer bytes.Buffer
    buffer.WriteString("DB: ")
    for _, post := range d.posts {
        buffer.WriteString(fmt.Sprintf("%v/%v/%v, ", post.Id, post.Parent, d.hashes[post.Id]))
    }
    return buffer.String()
}

func (d *ActualDB) Size() int {
    return len(d.posts)
}