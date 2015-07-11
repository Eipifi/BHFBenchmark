package bhf

func CopyPartially(a, b DB, id uint64, probability float64) {
    a.Put(*(b.Get(id)))
    for _, child := range b.Children(id) {
        if random_float64() < probability {
            CopyPartially(a, b, child, probability)
        }
    }
}

func GenerateBalanced(db DB, id uint64, height int, branching int) {
    if height < 1 { return }
    for i := 0; i < branching; i += 1 {
        new_id := random_uint64()
        db.Put(Post{Id: new_id, Parent: id})
        GenerateBalanced(db, new_id, height-1, branching)
    }
}

func GenerateRandomBalanced(db DB, id uint64, height int, b_min, b_max int) {
    if height < 1 { return }
    branching := random_int(b_min, b_max)
    for i := 0; i < branching; i += 1 {
        new_id := random_uint64()
        db.Put(Post{Id: new_id, Parent: id})
        GenerateRandomBalanced(db, new_id, height-1, b_min, b_max)
    }
}

func GenerateFurryList(db DB, id uint64, height int, fur_length, fur_density int) {
    for i := 0; i < height; i += 1 {
        GenerateRandomBalanced(db, id, fur_length, 0, fur_density)
        new_id := random_uint64()
        db.Put(Post{Id: new_id, Parent: id})
        id = new_id
    }
}