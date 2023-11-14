package main

import (
	"log"
	"time"

	"github.com/adrsimon/gomagnon/simulation"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	sim := simulation.NewSimulation()
	ebiten.SetWindowSize(sim.ScreenWidth, sim.ScreenHeight)
	ebiten.SetWindowTitle("Map Generated")

	if err := ebiten.RunGame(&sim); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			err := sim.Update()
			if err != nil {
				// Handle error
			}
		}
	}
}
