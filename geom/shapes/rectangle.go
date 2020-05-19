package shapes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Rectangle returns an image with a rectangle in it.
func Rectangle(width, height int, c color.Color) *ebiten.Image {
	img, _ := ebiten.NewImage(width, height, ebiten.FilterDefault)
	DrawRectangle(img, 0, 0, width-1, height-1, c)
	return img
}

// DrawRectangle draws rectangle in the provided image.
func DrawRectangle(img *ebiten.Image, x1, y1, x2, y2 int, c color.Color) {
	fx1, fy1 := float64(x1+1), float64(y1)
	fx2, fy2 := float64(x2+1), float64(y2)
	ebitenutil.DrawLine(img, fx1, fy1, fx2, fy1, c)
	ebitenutil.DrawLine(img, fx1, fy2, fx2, fy2, c)
	ebitenutil.DrawLine(img, fx1, fy1, fx1, fy2, c)
	ebitenutil.DrawLine(img, fx2, fy1, fx2, fy2, c)
}
