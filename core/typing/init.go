package typing

import (
	"github.com/adrsimon/gomagnon/settings"
	"math/rand"
)

var Seed int64 = settings.Setting.World.Seed
var Randomizer *rand.Rand

func init() {
	Randomizer = rand.New(rand.NewSource(Seed))
}
