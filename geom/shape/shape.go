package shape

import "github.com/hajimehoshi/ebiten"

// Shape is an abstract shape that could draw itself in the image.
type Shape interface {
	Draw(*ebiten.Image)
}
