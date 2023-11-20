package simulation

import (
	"github.com/adrsimon/gomagnon/core/typing"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type Widget struct {
	X, Y           float32
	Width, Height  float32
	displayedAgent *typing.Human
}

func (w *Widget) Update() error {
	return nil
}

func (w *Widget) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, w.X, w.Y, w.Width, w.Height, color.White, false)
	vector.DrawFilledRect(screen, w.X+3, w.Y+3, w.Width-6, w.Height-6, color.Black, false)

	// TODO : write agent stats here
}

func (w *Widget) Layout() (int, int) {
	return int(w.Width), int(w.Height)
}

func NewWidget(x, y, width, height float32, ag *typing.Human) *Widget {
	return &Widget{
		X:              x,
		Y:              y,
		Width:          width,
		Height:         height,
		displayedAgent: ag,
	}
}
