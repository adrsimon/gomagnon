package typing

import (
	"math/rand"
)

var Seed int64 = 10
var Randomizer *rand.Rand

func init() {
	Randomizer = rand.New(rand.NewSource(Seed))
}
