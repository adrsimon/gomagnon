package typing

type Human struct {
	id          int
	Position    Hexagone
	Type        rune // can be cromagnon or neandertal
	Hungriness  int  // 0 to 100
	Thirstiness int  // 0 to 100
	Age         int
	Gender      rune
	Strength    int // 0 to 100
	Sociability int // 0 to 100
}

func (h *Human) EvaluateSurroundings(hex *Hexagone) float64 {
	var foodValue, waterValue float64

	switch hex.Resource {
	case ANIMAL:
		foodValue = (float64(h.Hungriness) / 100) * 3
	case FRUIT:
		foodValue = (float64(h.Hungriness) / 100)
	}

	if hex.Biome.BiomeType == WATER {
		waterValue = (float64(h.Thirstiness) / 100)
	}

	return (foodValue + waterValue) / 2
}
