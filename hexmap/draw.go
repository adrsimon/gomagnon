package hexmap

import (
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var hexImage *ebiten.Image

func init() {
	img, _, err := ebitenutil.NewImageFromFile("images/forest.png")
	if err != nil {
		log.Fatal(err)
	}
	hexImage = img
}

// func (g *GameMap) DrawHex(background *ebiten.Image, xCenter float32, yCenter float32, color color.Color, hexSize float32) {
// 	x0 := xCenter
// 	x1 := xCenter - hexSize/2
// 	x2 := xCenter + hexSize/2
// 	y1 := yCenter - hexSize/2
// 	y2 := yCenter + hexSize/2
// 	y3 := yCenter - hexSize/4
// 	y4 := yCenter + hexSize/4
// 	vector.StrokeLine(background, x1, y3, x1, y4, 1, color, false)
// 	vector.StrokeLine(background, x1, y4, x0, y2, 1, color, false)
// 	vector.StrokeLine(background, x0, y2, x2, y4, 1, color, false)
// 	vector.StrokeLine(background, x2, y4, x2, y3, 1, color, false)
// 	vector.StrokeLine(background, x2, y3, x0, y1, 1, color, false)
// 	vector.StrokeLine(background, x0, y1, x1, y3, 1, color, false)
// 	vector.DrawFilledCircle(background, xCenter, yCenter, hexSize/4, color, false)
// }

func (g *GameMap) DrawHex(background *ebiten.Image, xCenter float32, yCenter float32, color color.Color, hexSize float32) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(xCenter-hexSize), float64(yCenter-hexSize))

	// Redimensionner l'image pour correspondre à la taille de l'hexagone
	//op.GeoM.Scale((1.5)*(float64(hexSize)/float64(hexImage.Bounds().Dx())), (1.5)*(float64(hexSize)/float64(hexImage.Bounds().Dy())))

	// Dessiner l'image hexagonale sur le fond
	background.DrawImage(hexImage, op)
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
