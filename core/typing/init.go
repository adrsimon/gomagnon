package typing

import (
	"math/rand"
)

var Randomizer *rand.Rand

func init() {
	Randomizer = rand.New(rand.NewSource(1))
}
