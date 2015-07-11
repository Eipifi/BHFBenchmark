package aldeon
import "math/rand"


func assert(condition bool) {
    if ! condition {
        panic("Condition failed")
    }
}

func random_uint64() uint64 {
    return uint64(rand.Uint32()) << 32 + uint64(rand.Uint32())
}

func random_int(min, max int) int {
    return min + rand.Intn(max - min)
}