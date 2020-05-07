package application

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/app/application/states"
	"github.com/yarcat/playground/geom/app/component"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

type underCursor struct {
	app *App
	c   component.Component
	f   *ftrs.Features
}

func (uc *underCursor) features() *ftrs.Features {
	if uc.f == nil {
		comp := uc.component()
		uc.f = uc.app.features[comp]
	}
	return uc.f
}

func (uc *underCursor) component() component.Component {
	if uc.c == nil {
		pt := image.Pt(ebiten.CursorPosition())
		uc.c = uc.app.ComponentAt(pt)
	}
	return uc.c
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
