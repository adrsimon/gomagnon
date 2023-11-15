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

func DrawAgents(screen *ebiten.Image, b *typing.Board) {
	agents := make([]*typing.Human, 0)
	for _, ag := range b.AgentManager.Agents {
		agents = append(agents, ag)
	}

	for _, agent := range agents {
		hexSize := float32(b.HexSize)
		x := float32(agent.Position.Position.X)
		y := float32(agent.Position.Position.Y)

		var offsetX, offsetY float32
		offsetY = 0.75 * hexSize
		offsetX = 0

		if int(y)%2 == 0 {
			offsetX = hexSize / 2
			DrawAgent(screen, x*hexSize+offsetX, y*offsetY, hexSize)
		} else {
			DrawAgent(screen, x*hexSize+offsetX, y*offsetY, hexSize)
		}

		for _, neighbor := range agent.GetNeighborsWithin5() {
			if neighbor == nil {
				continue
			}
			xN := float32(neighbor.Position.X)
			yN := float32(neighbor.Position.Y)
			if int(yN)%2 == 0 {
				offsetX = hexSize / 2
				DrawAgentNeighbor(screen, xN*hexSize+offsetX, yN*offsetY, hexSize)
			} else {
				DrawAgentNeighbor(screen, xN*hexSize, yN*offsetY, hexSize)
			}
		}
	}
}
