package typing

import (
	"fmt"
)

type Hexagone struct {
	Position *Point2D
	Resource ResourceType
	Biome    BiomeType
	Agents   []*Agent
	Hut      *Hut
}

func (h *Hexagone) EvenRToAxial() (int, int) {
	q := h.Position.X - ((h.Position.Y + (h.Position.Y & 1)) / 2)
	r := h.Position.Y
	return q, r
}

func (h *Hexagone) ToString() string {
	return fmt.Sprintf("%d:%d", h.Position.X, h.Position.Y)
}

type Point2D struct {
	X int
	Y int
}
