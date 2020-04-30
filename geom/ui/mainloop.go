package ui

import (
	"github.com/hajimehoshi/ebiten"
)

// Run executes main simulation loop.
func Run(ui *UI) error {
	return ebiten.RunGame((*gameAdapter)(ui))
}

type gameAdapter UI

// Update updates a game by one tick. The screen provided is ignored.
func (ga *gameAdapter) Update(screen *ebiten.Image) error {
	(*UI)(ga).update()
	return nil
}

// Draw draws the game screen by one frame. The given argument represents
// a screen image.
func (ga *gameAdapter) Draw(screen *ebiten.Image) {
	(*UI)(ga).draw(screen)
	debugPrint(screen)
}

// Layout returns the desired screen dimensions. Current implementation ignores
// the window dimensions provided.
func (ga *gameAdapter) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return ga.screenWidth, ga.screenHeight
}
