package drawing

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawBoard(screen *ebiten.Image, b *typing.Board) {

	for _, biome := range b.Biomes {
		for _, hex := range biome.Hexs {
			hexSize := float32(b.HexSize)
			x := float32(hex.Position.X)
			y := float32(hex.Position.Y)

			var offsetX, offsetY float32
			offsetY = 0.75 * hexSize
			offsetX = 0

			if int(y)%2 == 0 {
				offsetX = hexSize / 2
				DrawHex(screen, x*hexSize+offsetX, y*offsetY, biome.BiomeType, hexSize, hex.Resource)
			} else {
				DrawHex(screen, x*hexSize+offsetX, y*offsetY, biome.BiomeType, hexSize, hex.Resource)
			}
		}
	}
}
