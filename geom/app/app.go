package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// NewApp returns new application instance with initialized logical window width
// and height.
func NewApp(width, height int) *App {
	// Setting defaults.
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowResizable(true)
	app := &App{
		width:  width,
		height: height,
	}
	app.Window = NewWindowImpl(app, nil)
	return app
}

// Application defines and object that provides application functionality
// e.g. scheduling.
type Application interface{}

// App represents an application. It manages internal objects and present
// the on-screen.
type App struct {
	Window
	width, height int
}

// Update updates a by one tick. The image provided is ignored.
func (a App) Update(screen *ebiten.Image) error { return nil }

// Draw draws the app screen by one frame.
func (a App) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, world!")
}

// Layout accepts a native outside size in device-independent pixels and
// returns the game's logical screen size.
func (a App) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return a.width, a.height
}
