package main

import (
	"github.com/adrsimon/gomagnon/simulation"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	sim := simulation.NewSimulation()
	ebiten.SetWindowSize(sim.ScreenWidth, sim.ScreenHeight)
	ebiten.SetWindowTitle("Map Generated")
	ebiten.SetTPS(1)

	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
