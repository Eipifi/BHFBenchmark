package main
import (
    "fmt"
    "github.com/Eipifi/aldeon/bhf"
)

const ALLOW_SUGGEST = true
const ITERATIONS = 1000

func main() {
    bhf.LoggingEnabled = false
    bhf.NewSeed()

    var total_req, total_num int

    for i := 0; i < ITERATIONS; i += 1 {
        req, num := run(4, 4, 0.5)
        total_req += req
        total_num += num
    }

    fmt.Println("Requests sent: ", float64(total_req) / float64(ITERATIONS))
    fmt.Println("Posts copied: ", float64(total_num) / float64(ITERATIONS))
}

func run(depth, width int, probability float64) (int, int) {
    // Two databases
    db_local := bhf.NewDB()
    db_remote := bhf.NewDB()

    // The root of the conversation
    root := bhf.Post{Id: 42, Parent: 0}

    // Generate the test instance
    db_remote.Put(root)
    bhf.GenerateRandomBalanced(db_remote, root.Id, depth, width)
    bhf.CopyPartially(db_local, db_remote, root.Id, probability)

    local_size := db_local.Size()
    remote_size := db_remote.Size()

    requests_sent := bhf.Synchronize(db_local, db_remote, root.Id, ALLOW_SUGGEST)
    return requests_sent, remote_size - local_size
}