// Package application provides an App class that is a main container of elements,
// event dispatcher, etc.
package application

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/app/component"
)

// New returns new application instance with initialized logical window width
// and height.
func New(width, height int) *App {
	// Setting defaults.
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowResizable(true)

	app := &App{width: width, height: height}
	app.mouseManager.app = app
	return app
}

// App represents an application. It manages internal objects and present
// the on-screen.
type App struct {
	width, height int
	drawable      []component.DrawableComponent
	mouseManager  mouseManager
}

// AddDrawable adds a component to the top-most container.
func (app *App) AddDrawable(c component.DrawableComponent) {
	app.drawable = append(app.drawable, c)
}

// ComponentAt returns compotents under given window point.
// Coordinates are given in logical window units.
func (app App) ComponentAt(pt image.Point) []component.Component {
	var comps []component.Component
	for _, c := range app.drawable {
		if pt.In(c.Bounds()) {
			comps = append(comps, c)
		}
	}
	return comps
}
