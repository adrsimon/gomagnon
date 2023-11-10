package hexmap

import (
	"fmt"
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

func (g *GameMap) DrawHex(background *ebiten.Image, xCenter float32, yCenter float32, biome BiomesType, hexSize float32) {
	var hexImage *ebiten.Image
	switch biome {
	case PLAINS:
		hexImage = imgPlains
	case FOREST:
		hexImage = imgForest
	case WATER:
		hexImage = imgWater
	case CAVE:
		hexImage = imgCaves
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale((3/2.0)*float64(hexSize)/float64(hexImage.Bounds().Dx()), (4.0/3.0)*float64(hexSize)/float64(hexImage.Bounds().Dy()))
	op.GeoM.Translate(float64(xCenter-(3.0/4.0)*hexSize), float64(yCenter-(2.0/3.0)*hexSize))
	fmt.Println(biome, hexImage.Bounds().Dx(), hexImage.Bounds().Dy())
	background.DrawImage(hexImage, op)

	x0 := xCenter
	x1 := xCenter - hexSize/2
	x2 := xCenter + hexSize/2
	y1 := yCenter - hexSize/2
	y2 := yCenter + hexSize/2
	y3 := yCenter - hexSize/4
	y4 := yCenter + hexSize/4
	vector.StrokeLine(background, x1, y3, x1, y4, 1, colornames.Black, false)
	vector.StrokeLine(background, x1, y4, x0, y2, 1, colornames.Black, false)
	vector.StrokeLine(background, x0, y2, x2, y4, 1, colornames.Black, false)
	vector.StrokeLine(background, x2, y4, x2, y3, 1, colornames.Black, false)
	vector.StrokeLine(background, x2, y3, x0, y1, 1, colornames.Black, false)
	vector.StrokeLine(background, x0, y1, x1, y3, 1, colornames.Black, false)
}
