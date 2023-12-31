package drawing

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

func DrawAgent(screen *ebiten.Image, x, y, size float32, ag *typing.Agent) {
	if ag.Race == typing.SAPIENS {
		if ag.Body.Age < 10 {
			drawImage(screen, x, y, size, imgBabySapiens)
		} else {
			drawImage(screen, x, y, size/1.25, imgSapiens)
		}
	} else {
		if ag.Body.Age < 10 {
			drawImage(screen, x, y, size, imgBabyNeanderthal)
		} else {
			drawImage(screen, x, y, size/1.25, imgNeanderthal)
		}
	}
}

func DrawAgentNeighbor(background *ebiten.Image, xCenter float32, yCenter float32, hexSize float32, color color.Color) {
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
}

func DrawAgentPath(background *ebiten.Image, xa, ya, xb, yb float32, color color.Color) {
	vector.StrokeLine(background, xa, ya, xb, yb, 2, color, false)
}
