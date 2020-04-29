package simulation

import "github.com/hajimehoshi/ebiten"

// Run executes main simulation loop.
func Run(s *Simulation) error {
	return ebiten.RunGame((*gameAdapter)(s))
}

type gameAdapter Simulation

// Update updates a game by one tick. The screen provided is ignored.
func (ga *gameAdapter) Update(screen *ebiten.Image) error {
	(*Simulation)(ga).update()
	return nil
}

// Draw draws the game screen by one frame. The given argument represents
// a screen image.
func (ga *gameAdapter) Draw(screen *ebiten.Image) {
	(*Simulation)(ga).draw(screen)
}

// Layout returns the desired screen dimensions. Current implementation ignores
// the window dimensions provided.
func (ga *gameAdapter) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return ga.screenWidth, ga.screenHeight
}
