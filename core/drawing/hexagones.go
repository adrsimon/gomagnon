package drawing

import (
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	_ "image/png"

	"github.com/adrsimon/gomagnon/core/typing"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawHex(background *ebiten.Image, xCenter float32, yCenter float32, biome typing.BiomeType, hexSize float32, resource typing.ResourceType, hut *typing.Hut) {
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
		drawImage(background, xCenter, yCenter, hexSize/1.5, imgMushroom)
	case typing.ANIMAL:
		drawImage(background, xCenter, yCenter, hexSize/1.5, imgCow)
	case typing.ROCK:
		drawImage(background, xCenter, yCenter, hexSize/1.5, imgRock)
	case typing.WOOD:
		drawImage(background, xCenter, yCenter, hexSize/1.5, imgWood)
	case typing.NONE:
	}

	if hut != nil {
		if hut.Owner == nil {
			vector.DrawFilledRect(background, xCenter-hexSize/4, yCenter-hexSize/4, hexSize/2, hexSize/2, colornames.Black, false)
		} else if hut.Owner.Race == typing.SAPIENS {
			vector.DrawFilledRect(background, xCenter-hexSize/4, yCenter-hexSize/4, hexSize/2, hexSize/2, colornames.Blue, false)
		} else if hut.Owner.Race == typing.NEANDERTHAL {
			vector.DrawFilledRect(background, xCenter-hexSize/4, yCenter-hexSize/4, hexSize/2, hexSize/2, colornames.Red, false)
		}
	}
}

func drawImage(background *ebiten.Image, x, y, size float32, img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(size)/float64(img.Bounds().Dx()), float64(size)/float64(img.Bounds().Dy()))
	op.GeoM.Translate(float64(x-size/2), float64(y-size/2))
	background.DrawImage(img, op)
}
