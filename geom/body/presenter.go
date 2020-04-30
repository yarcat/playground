package body

import "github.com/hajimehoshi/ebiten"

// Present draws bodies on the screen.
func Present(image *ebiten.Image, it Iterator) {
	for it.Next() {
		b := it.Value()
		b.Shape.Draw(image)
	}
}
