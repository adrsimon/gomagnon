package hexmap

import (
	"image/color"
)

type GameMap struct {
	Board           *Board
	ScreenWidth     int
	ScreenHeight    int
	BackgroundColor color.RGBA
	DirtColor       color.RGBA
	ForestColor     color.RGBA
	WaterColor      color.RGBA
	CaveColor       color.RGBA
}

type Board struct {
	Cases   map[string]*Hexagone
	XMax    int
	YMax    int
	HexSize int
	Biomes  []*Biome
}

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

type Hexagone struct {
	Position *Point2D
}

type Point2D struct {
	X int
	Y int
}
