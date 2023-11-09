package hexmap

import (
	"image/color"
)

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
	Agents  []*Human
}

type BiomesType int

const (
	DIRT BiomesType = iota
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
