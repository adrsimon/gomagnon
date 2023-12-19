package drawing

import (
	"image/color"

	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func DrawBoard(screen *ebiten.Image, b *typing.Board, cameraX, cameraY, zoomFactor float32, agent *typing.Agent) {
	for _, line := range b.Cases {
		for _, hex := range line {
			if agent != nil && agent.MapVision[hex.Position.X][hex.Position.Y].Position.X == -1 && agent.MapVision[hex.Position.X][hex.Position.Y].Position.Y == -1 {
				continue
			}

			hexSize := b.HexSize
			x := hex.Position.X
			y := hex.Position.Y

			xc, yc := GetHexGraphicalCenter(x, y, hexSize)
			xc, yc = xc-cameraX, yc-cameraY
			xc, yc = xc*zoomFactor, yc*zoomFactor

			if &hex.Hut == nil {
				DrawHex(screen, xc, yc, hex.Biome, hexSize*zoomFactor, hex.Resource, nil)
			} else {
				DrawHex(screen, xc, yc, hex.Biome, hexSize*zoomFactor, hex.Resource, hex.Hut)
			}
		}
	}
}

func DrawAgents(screen *ebiten.Image, agents []*typing.Agent, cameraX, cameraY, zoomFactor, hexSize float32, debug bool) {
	for _, agent := range agents {
		var col color.Color
		if agent.Race == typing.SAPIENS {
			col = colornames.Blue
		} else {
			col = colornames.Red
		}

		x := agent.Position.Position.X
		y := agent.Position.Position.Y

		xA, yA := GetHexGraphicalCenter(x, y, hexSize)
		xA, yA = xA-cameraX, yA-cameraY
		xA, yA = xA*zoomFactor, yA*zoomFactor
		DrawAgent(screen, xA, yA, hexSize*zoomFactor, agent)

		if !debug {
			continue
		}

		for _, neighbor := range agent.Behavior.GetNeighboursWithinAcuity() {
			if neighbor == nil {
				continue
			}
			xN, yN := GetHexGraphicalCenter(neighbor.Position.X, neighbor.Position.Y, hexSize)
			xN, yN = xN-cameraX, yN-cameraY
			xN, yN = xN*zoomFactor, yN*zoomFactor
			DrawAgentNeighbor(screen, xN, yN, hexSize*zoomFactor, col)
		}

		if agent.CurrentPath != nil && len(agent.CurrentPath) > 0 {
			x0, y0 := GetHexGraphicalCenter(agent.Position.Position.X, agent.Position.Position.Y, hexSize)
			x1, y1 := GetHexGraphicalCenter(agent.CurrentPath[len(agent.CurrentPath)-1].X, agent.CurrentPath[len(agent.CurrentPath)-1].Y, hexSize)
			x0, y0 = x0-cameraX, y0-cameraY
			x1, y1 = x1-cameraX, y1-cameraY
			x0, y0 = x0*zoomFactor, y0*zoomFactor
			x1, y1 = x1*zoomFactor, y1*zoomFactor
			DrawAgentPath(screen, x0, y0, x1, y1, col)
			for i := 0; i < len(agent.CurrentPath)-1; i++ {
				xa, ya := GetHexGraphicalCenter(agent.CurrentPath[i].X, agent.CurrentPath[i].Y, hexSize)
				xb, yb := GetHexGraphicalCenter(agent.CurrentPath[i+1].X, agent.CurrentPath[i+1].Y, hexSize)
				xa, ya = xa-cameraX, ya-cameraY
				xb, yb = xb-cameraX, yb-cameraY
				xa, ya = xa*zoomFactor, ya*zoomFactor
				xb, yb = xb*zoomFactor, yb*zoomFactor
				DrawAgentPath(screen, xa, ya, xb, yb, col)
			}
		}
	}
}

func GetHexGraphicalCenter(x, y int, hexSize float32) (float32, float32) {
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
