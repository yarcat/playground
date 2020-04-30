package shape

import (
	"image/color"

	"github.com/yarcat/playground/geom/image"
)

// Circle represents a circle that could be used in the simulation.
type Circle struct {
	// R is the radius of the circle.
	R float64
	// C is the color of the circle.
	Color color.Color
}

// Draw draws a complete circle with the center in origin.
func (c Circle) Draw(image image.Image) {
	drawCircle(image, int(c.R), c.Color)
}

func drawCircle(image image.Image, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)
	for x > y {
		image.Set(x, y, c)
		image.Set(y, x, c)
		image.Set(-y, x, c)
		image.Set(-x, y, c)
		image.Set(-x, -y, c)
		image.Set(-y, -x, c)
		image.Set(y, -x, c)
		image.Set(x, -y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}
