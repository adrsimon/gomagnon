package main

import (
	"github.com/adrsimon/gomagnon/gui"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	sim := gui.NewSimulation()
	ebiten.SetWindowSize(sim.ScreenWidth, sim.ScreenHeight)
	ebiten.SetWindowTitle("Map Generated")

	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}
}
