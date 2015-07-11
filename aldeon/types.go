package aldeon

type IdAndHash struct {
    Id uint64
    Hash uint64
}

type Post struct {
    Id uint64
    Parent uint64
}

type DB interface {
    // Fetch the post
    Get(uint64) *Post

    // Insert the post
    Put(Post)

    // Find all posts with a given hash value
    Find(uint64) []uint64

    // Fetch the hash of the branch
    Hash(uint64) uint64

    // Fetch all the children of a given post
    Children(uint64) []uint64
}