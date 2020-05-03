package main

import (
	"github.com/hajimehoshi/ebiten"
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
	app.Window = NewWindowImpl(
		app,
		nil, /* parent */
	)
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

// Run executes main application loop. It dispatches events, updates
// the world, etc.
func Run(a *App) error {
	return ebiten.RunGame((*gameAdapter)(a))

}

type drawable interface {
	Draw(*ebiten.Image)
}

type gameAdapter App

// Update updates the world by one frame.
func (ga *gameAdapter) Update(screen *ebiten.Image) error { return nil }

// Draw draws the world by one frame.
func (ga gameAdapter) Draw(screen *ebiten.Image) {
	var drawRecursively func(win Window)
	drawRecursively = func(win Window) {
		if d, ok := win.(drawable); ok {
			d.Draw(screen)
		}
		for _, childWin := range win.Children() {
			drawRecursively(childWin)
		}
	}
	drawRecursively(App(ga))
}

// Layout accepts window dimensions and returns logical dimensions.
func (ga gameAdapter) Layout(windowWidth, windowHeight int) (screenWidth, screenHeight int) {
	return App(ga).width, App(ga).height
}
