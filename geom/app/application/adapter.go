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
	for _, c := range App(ga).drawable {
		c.Draw(screen)
	}
}

// Layout accepts window dimensions and returns logical dimensions.
func (ga gameAdapter) Layout(windowWidth, windowHeight int) (screenWidth, screenHeight int) {
	return App(ga).width, App(ga).height
}
