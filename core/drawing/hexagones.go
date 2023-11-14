package drawing

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"golang.org/x/image/colornames"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var imgPlains, imgForest, imgWater, imgCaves *ebiten.Image

func init() {
	plains, _, err := ebitenutil.NewImageFromFile("assets/images/plains.png")
	if err != nil {
		log.Fatal(err)
	}
	imgPlains = plains

	forest, _, err := ebitenutil.NewImageFromFile("assets/images/forest.png")
	if err != nil {
		log.Fatal(err)
	}
	imgForest = forest

	water, _, err := ebitenutil.NewImageFromFile("assets/images/water.png")
	if err != nil {
		log.Fatal(err)
	}
	imgWater = water

	caves, _, err := ebitenutil.NewImageFromFile("assets/images/caves.png")
	if err != nil {
		log.Fatal(err)
	}
	imgCaves = caves
}

func DrawHex(background *ebiten.Image, xCenter float32, yCenter float32, biome typing.BiomesType, hexSize float32, resource typing.ResourceType) {
	var hexImage *ebiten.Image
	switch biome {
	case typing.PLAINS:
		hexImage = imgPlains
	case typing.FOREST:
		hexImage = imgForest
	case typing.WATER:
		hexImage = imgWater
	case typing.CAVE:
		hexImage = imgCaves
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale((3/2.0)*float64(hexSize)/float64(hexImage.Bounds().Dx()), (4.0/3.0)*float64(hexSize)/float64(hexImage.Bounds().Dy()))
	op.GeoM.Translate(float64(xCenter-(3.0/4.0)*hexSize), float64(yCenter-(2.0/3.0)*hexSize))
	background.DrawImage(hexImage, op)

	switch resource {
	case typing.FRUIT:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Green, false)
	case typing.ANIMAL:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Red, false)
	case typing.ROCK:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Grey, false)
	case typing.WOOD:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Black, false)
	case typing.NONE:
	}
}
