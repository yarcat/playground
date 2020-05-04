package application

import "github.com/hajimehoshi/ebiten"

// Run executes main application loop. It dispatches events, updates
// the world, etc.
func Run(a *App) error {
	return ebiten.RunGame((*gameAdapter)(a))
}
