package typing

import (
	"math/rand"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(1))
}
