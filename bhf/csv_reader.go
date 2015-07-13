package bhf
import (
    "encoding/csv"
    "os"
    "io"
    "errors"
    "strconv"
    "hash/fnv"
    "strings"
    "math"
)

type CSVPost struct {
    Post
    Timestamp uint64
}

func ParseCSV(root uint64, path string) ([]CSVPost, error) {
    file, err := os.Open(path)
    if err != nil { return nil, err }
    reader := csv.NewReader(file)
    reader.Comma = ';'
    var result []CSVPost
    for {
        record, err := reader.Read()
        if err == io.EOF { break }
        if err != nil { return nil, err }
        if len(record) != 3 { return nil, errors.New("Invalid line format: 3 strings expected") }

        post := CSVPost{}

        timestamp_float, err := strconv.ParseFloat(record[2], 64)
        if err != nil { return nil, err }

        post.Timestamp = uint64(timestamp_float)
        post.Id = hash64(record[0])
        if strings.HasPrefix(record[1], "t3") {
            post.Parent = root
        } else {
            post.Parent = hash64(record[1])
        }
        result = append(result, post)
    }
    return result, nil
}

func PutCSVPosts(db DB, posts []CSVPost, threshold uint64) {
    for _, v := range posts {
        if v.Timestamp <= threshold {
            db.Put(v.Post)
        }
    }
}

func MinTS(posts []CSVPost) uint64 {
    var min uint64 = math.MaxUint64
    for i := 1; i < len(posts); i += 1 {
        if posts[i].Timestamp < min {
            min = posts[i].Timestamp
        }
    }
    return min
}

func MaxTS(posts []CSVPost) uint64 {
    var max uint64 = 0
    for i := 0; i < len(posts); i += 1 {
        if posts[i].Timestamp > max {
            max = posts[i].Timestamp
        }
    }
    return max
}

func hash64(s string) uint64 {
    h := fnv.New64a()
    h.Write([]byte(s))
    return h.Sum64()
}