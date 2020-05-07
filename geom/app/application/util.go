package application

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/app/application/states"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

type underCursor struct {
	app *App
	f   *ftrs.Features
}

func (uc *underCursor) features() *ftrs.Features {
	if uc.f == nil {
		pt := image.Pt(ebiten.CursorPosition())
		comp := uc.app.ComponentAt(pt)
		uc.f = uc.app.features[comp]
	}
	return uc.f
}

func removeStates(states, remove []states.MouseButtonState) []states.MouseButtonState {
	n := 0
outer:
	for _, s := range states {
		for _, r := range remove {
			if r == s {
				continue outer
			}
		}
		states[n] = s
		n++
	}
	return states[:n]
}
