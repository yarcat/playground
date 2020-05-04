// Package application provides an App class that is a main container of elements,
// event dispatcher, etc.
package application

import (
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

	return &App{width: width, height: height}
}

// App represents an application. It manages internal objects and present
// the on-screen.
type App struct {
	width, height int
	drawable      []component.DrawableComponent
}

// AddDrawable adds a component to the top-most container.
func (app *App) AddDrawable(c component.DrawableComponent) {
	app.drawable = append(app.drawable, c)
}
