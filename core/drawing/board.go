package drawing

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawBoard(screen *ebiten.Image, b *typing.Board) {
	for _, biome := range b.Biomes {
		for _, hex := range biome.Hexs {
			hexSize := float32(b.HexSize)
			x := hex.Position.X
			y := hex.Position.Y

			xc, yc := getHexGraphicalCenter(x, y, hexSize)
			DrawHex(screen, xc, yc, biome.BiomeType, hexSize, hex.Resource)
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
		x := agent.Position.Position.X
		y := agent.Position.Position.Y

		xA, yA := getHexGraphicalCenter(x, y, hexSize)
		DrawAgent(screen, xA, yA, hexSize)

		for _, neighbor := range agent.GetNeighborsWithin5() {
			if neighbor == nil {
				continue
			}
			xN, yN := getHexGraphicalCenter(neighbor.Position.X, neighbor.Position.Y, hexSize)
			DrawAgentNeighbor(screen, xN, yN, hexSize)
		}

		if agent.CurrentPath != nil && len(agent.CurrentPath) > 0 {
			x0, y0 := getHexGraphicalCenter(agent.Position.Position.X, agent.Position.Position.Y, hexSize)
			x1, y1 := getHexGraphicalCenter(agent.CurrentPath[len(agent.CurrentPath)-1].Position.X, agent.CurrentPath[len(agent.CurrentPath)-1].Position.Y, hexSize)
			DrawAgentPath(screen, x0, y0, x1, y1)
			for i := 0; i < len(agent.CurrentPath)-1; i++ {
				xa, ya := getHexGraphicalCenter(agent.CurrentPath[i].Position.X, agent.CurrentPath[i].Position.Y, hexSize)
				xb, yb := getHexGraphicalCenter(agent.CurrentPath[i+1].Position.X, agent.CurrentPath[i+1].Position.Y, hexSize)
				DrawAgentPath(screen, xa, ya, xb, yb)
			}
		}
	}
}

func getHexGraphicalCenter(x, y int, hexSize float32) (float32, float32) {
	var offsetX, offsetY float32
	offsetY = 0.75 * hexSize
	offsetX = 0

	if int(y)%2 == 0 {
		offsetX = hexSize / 2
		return float32(x)*hexSize + offsetX, float32(y) * offsetY
	} else {
		return float32(x)*hexSize + offsetX, float32(y) * offsetY
	}
}
