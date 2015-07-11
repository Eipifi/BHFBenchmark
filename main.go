package main
import (
    "aldeon/aldeon"
    "fmt"
)


func main() {
    db := aldeon.NewDB()
    db.Put(aldeon.Post{5, 0})
    db.Put(aldeon.Post{6, 5})

    fmt.Println(db)
}
