package typing

type Human struct {
	id          int
	Position    string
	Type        rune // can be cromagnon or neandertal
	Hungriness  int  // 0 to 100
	Thirstiness int  // 0 to 100
	Age         int
	Gender      rune
	Strength    int // 0 to 100
	Sociability int // 0 to 100
}

const (
	AnimalFoodValueMultiplier = 3.0
	FruitFoodValueMultiplier  = 1.0
	WaterValueMultiplier      = 2.0
)

func (h *Human) EvaluateSurroundings(hex *Hexagone) float64 {
	var score = 0.0

	switch hex.Resource {
	case ANIMAL:
		score = (float64(h.Hungriness) / 100) * AnimalFoodValueMultiplier
	case FRUIT:
		score = (float64(h.Hungriness) / 100) * FruitFoodValueMultiplier
	}

	if hex.Biome.BiomeType == WATER {
		score = (float64(h.Thirstiness) / 100) * WaterValueMultiplier
	}

	return score
}
