package main

import (
	"github.com/adrsimon/gomagnon/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	sim := simulation.NewSimulation()

	sim.UI, sim.Selector, sim.AgentDesc = simulation.BuildUI(&sim)
	ebiten.SetWindowSize(sim.ScreenWidth, sim.ScreenHeight)
	ebiten.SetWindowTitle("Gomagnon - Neanderthal vs Sapiens")
	ebiten.SetTPS(20)

	if err := ebiten.RunGameWithOptions(&sim, &ebiten.RunGameOptions{
		GraphicsLibrary:   0,
		InitUnfocused:     false,
		ScreenTransparent: false,
		SkipTaskbar:       false,
	}); err != nil {
		log.Fatal(err)
	}
}
