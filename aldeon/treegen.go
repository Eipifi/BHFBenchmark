package aldeon

func CopyPartially(a, b DB, id uint64, probability float64) {
    a.Put(*(b.Get(id)))
    for _, child := range b.Children(id) {
        if random_float64() < probability {
            CopyPartially(a, b, child, probability)
        }
    }
}

func GenerateRandomBalanced(db DB, id uint64, height int, branching int) {
    if height < 1 { return }
    for i := 0; i < branching; i += 1 {
        new_id := random_uint64()
        db.Put(Post{Id: new_id, Parent: id})
        GenerateRandomBalanced(db, new_id, height-1, branching)
    }
}