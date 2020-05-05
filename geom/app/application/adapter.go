package application

import "github.com/hajimehoshi/ebiten"

type gameAdapter App

// Update updates the world by one frame.
func (ga *gameAdapter) Update(screen *ebiten.Image) error {
	(*App)(ga).gestureManager.update()
	return nil
}

// Draw draws the world by one frame.
func (ga gameAdapter) Draw(screen *ebiten.Image) {
	// TODO(yarcat): Get rid of this dynamic dispatching. Components know whether they
	// can be displayed, and they could register themselves using life-cycle hooks.
	// This is kinda similar to what Java AWT does with the default component by
	// utilizing a gazillion of methods, flags and event filters.
	type drawer interface {
		Draw(*ebiten.Image)
	}
	for _, c := range App(ga).components {
		if c, ok := c.(drawer); ok {
			c.Draw(screen)
		}
	}
}

// Layout accepts window dimensions and returns logical dimensions.
func (ga gameAdapter) Layout(windowWidth, windowHeight int) (screenWidth, screenHeight int) {
	return App(ga).width, App(ga).height
}
