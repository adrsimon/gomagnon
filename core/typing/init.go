package typing

import (
	"math/rand"
)

var Seed int64 = 010101
var Randomizer *rand.Rand

func init() {
	Randomizer = rand.New(rand.NewSource(Seed))
}
