package typing

import (
	"fmt"
	"image/color"
)

type GameMap struct {
	Board           *Board
	ScreenWidth     int
	ScreenHeight    int
	BackgroundColor color.RGBA
}

func NewGame(
	screenWidth, screenHeight int,
	backgroundColor color.RGBA,
	xmax, ymax int,
	hexSize float32,
	maxResources map[ResourceType]int,
) *GameMap {
	return &GameMap{
		Board:           NewBoard(xmax, ymax, hexSize, maxResources),
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		BackgroundColor: backgroundColor,
	}
}

type Hexagone struct {
	Position *Point2D
	Resource ResourceType
	Biome    BiomeType
	Agents   []*Human
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
