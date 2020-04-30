package image

import "image/color"

// Image represents a drawable object.
type Image interface {
	Set(x, y int, c color.Color)
}
