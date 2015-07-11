package bhf
import (
    "math/rand"
    "time"
)

var LoggingEnabled = true

func random_uint64() uint64 {
    return uint64(rand.Uint32()) << 32 + uint64(rand.Uint32())
}

func random_int(min, max int) int {
    return min + rand.Intn(max - min)
}

func random_float64() float64 {
    return rand.Float64()
}

func NewSeed() {
    rand.Seed(time.Now().UTC().UnixNano())
}