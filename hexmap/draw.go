package hexmap

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image/color"
)

func (g *GameMap) DrawHex(background *ebiten.Image, xCenter float32, yCenter float32, color color.Color, hexSize float32, resource ResourceType) {
	x0 := xCenter
	x1 := xCenter - hexSize/2
	x2 := xCenter + hexSize/2
	y1 := yCenter - hexSize/2
	y2 := yCenter + hexSize/2
	y3 := yCenter - hexSize/4
	y4 := yCenter + hexSize/4
	vector.StrokeLine(background, x1, y3, x1, y4, 1, color, false)
	vector.StrokeLine(background, x1, y4, x0, y2, 1, color, false)
	vector.StrokeLine(background, x0, y2, x2, y4, 1, color, false)
	vector.StrokeLine(background, x2, y4, x2, y3, 1, color, false)
	vector.StrokeLine(background, x2, y3, x0, y1, 1, color, false)
	vector.StrokeLine(background, x0, y1, x1, y3, 1, color, false)

	vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/2, color, false)

	switch resource {
	case FRUIT:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Green, false)
	case ANIMAL:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Red, false)
	case ROCK:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Grey, false)
	case WOOD:
		vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/8, colornames.Black, false)
	case NONE:
	}
}
