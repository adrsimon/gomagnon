package typing

type BiomesType int

const (
	PLAINS BiomesType = iota
	FOREST
	WATER
	CAVE
)

type Biome struct {
	BiomeType BiomesType
	Hexs      []*Hexagone
}
