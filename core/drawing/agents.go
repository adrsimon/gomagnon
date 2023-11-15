package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

func DrawAgent(screen *ebiten.Image, x, y, size float32) {
	vector.DrawFilledCircle(screen, x, y, size/4, colornames.Black, false)
}

func DrawAgentNeighbor(background *ebiten.Image, xCenter float32, yCenter float32, hexSize float32) {
	x0 := xCenter
	x1 := xCenter - hexSize/2
	x2 := xCenter + hexSize/2
	y1 := yCenter - hexSize/2
	y2 := yCenter + hexSize/2
	y3 := yCenter - hexSize/4
	y4 := yCenter + hexSize/4
	vector.StrokeLine(background, x1, y3, x1, y4, 2, colornames.Red, false)
	vector.StrokeLine(background, x1, y4, x0, y2, 2, colornames.Red, false)
	vector.StrokeLine(background, x0, y2, x2, y4, 2, colornames.Red, false)
	vector.StrokeLine(background, x2, y4, x2, y3, 2, colornames.Red, false)
	vector.StrokeLine(background, x2, y3, x0, y1, 2, colornames.Red, false)
	vector.StrokeLine(background, x0, y1, x1, y3, 2, colornames.Red, false)
}
