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
	fruits, animals, rocks, woods int,
) *GameMap {
	return &GameMap{
		Board:           NewBoard(xmax, ymax, hexSize, fruits, animals, rocks, woods),
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		BackgroundColor: backgroundColor,
	}
}

type Hexagone struct {
	Position *Point2D
	Resource ResourceType
	Biome    *Biome
	Agents   []*Human
}

func (h *Hexagone) OddRToAxial() (int, int) {
	q := h.Position.X - (h.Position.Y-(h.Position.Y&1))/2
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
