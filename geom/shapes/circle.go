package shapes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

// Circle returns an image with a circle in in. The image has width and
// height set to double radius.
func Circle(r int, c color.Color) *ebiten.Image {
	d := r * 2
	img, _ := ebiten.NewImage(d, d, ebiten.FilterDefault)
	DrawCircle(img, r, r, r, c)
	return img
}

// DrawCircle draws a circle with the provided center, radious and color.
func DrawCircle(image *ebiten.Image, x, y, r int, c color.Color) {
	x0, y0, x, y, dx, dy := x, y, r-1, 0, 1, 1
	err := dx - (r * 2)
	for x >= y {
		image.Set(x0+x, y0+y, c)
		image.Set(x0+y, y0+x, c)
		image.Set(x0-y, y0+x, c)
		image.Set(x0-x, y0+y, c)
		image.Set(x0-x, y0-y, c)
		image.Set(x0-y, y0-x, c)
		image.Set(x0+y, y0-x, c)
		image.Set(x0+x, y0-y, c)

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
