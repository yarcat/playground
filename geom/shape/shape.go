package shape

import "github.com/yarcat/playground/geom/image"

// Shape is an abstract shape that could draw itself in the image.
type Shape interface {
	Draw(image.Image)
}
