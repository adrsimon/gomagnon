package typing

import "image/color"

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
	hexSize int,
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
	Agents   []*Agent
}

type Point2D struct {
	X int
	Y int
}
