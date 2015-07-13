package main
import (
    "fmt"
    "github.com/Eipifi/bhf_benchmark/bhf"
    "orwell/lib/utils"
)

const ROOT_ID uint64 = 42

func run_test_from_csv(posts []bhf.CSVPost, threshold uint64, suggests_allowed bool) (reqs int, size_start int, size_end int) {
    db_local := bhf.NewDB()
    db_remote := bhf.NewDB()
    root := bhf.Post{ROOT_ID, 0}
    db_local.Put(root)
    db_remote.Put(root)
    bhf.PutCSVPosts(db_remote, posts, 999999999999)
    bhf.PutCSVPosts(db_local, posts, threshold)
    local_size := db_local.Size()
    return bhf.Synchronize(db_local, db_remote, root.Id, suggests_allowed), local_size, db_local.Size()
}

func run_both_tests(posts []bhf.CSVPost, threshold uint64) {
    t0, s0, f0 := run_test_from_csv(posts, threshold, true)
    t1, _, _   := run_test_from_csv(posts, threshold, false)
    fmt.Printf("%v,%v,%v,%v,%v \n", threshold, t0, t1, s0, f0)
}

func main() {
    bhf.LoggingEnabled = false
    posts, err := bhf.ParseCSV(ROOT_ID, "data/comments.csv")
    utils.Ensure(err)
    ts_min := bhf.MinTS(posts)
    ts_max := bhf.MaxTS(posts)

    for i := ts_min; i < ts_max; i += 300 { // 5 mins
        run_both_tests(posts, i)
    }
}
