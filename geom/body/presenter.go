package body

import (
	"github.com/yarcat/playground/geom/image"
)

// Present draws bodies on the screen.
func Present(image image.Image, it Iterator) {
	for it.Next() {
		b := it.Value()
		b.Shape.Draw(image)
	}
}
