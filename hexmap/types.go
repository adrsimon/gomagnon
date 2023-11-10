package hexmap

import (
	"image/color"
)

type GameMap struct {
	Board           *Board
	ScreenWidth     int
	ScreenHeight    int
	BackgroundColor color.RGBA
}

type Board struct {
	Cases           map[string]*Hexagone
	XMax            int
	YMax            int
	HexSize         int
	Biomes          []*Biome
	ResourceManager *ResourceManager
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

type ResourceType int

const (
	NONE ResourceType = iota
	FRUIT
	ANIMAL
	ROCK
	WOOD
)

type ResourceManager struct {
	Resources         []ResourceType
	MaxFruitQuantity  int
	MaxAnimalQuantity int
	MaxRockQuantity   int
	MaxWoodQuantity   int
	FruitQuantity     int
	AnimalQuantity    int
	RockQuantity      int
	WoodQuantity      int
}

type Hexagone struct {
	Position *Point2D
	Resource ResourceType
}

type Point2D struct {
	X int
	Y int
}
