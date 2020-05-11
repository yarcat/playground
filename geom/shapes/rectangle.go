package shapes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Rectangle returns an image with a rectangle in it.
func Rectangle(width, height int, c color.Color) *ebiten.Image {
	img, _ := ebiten.NewImage(width, height, ebiten.FilterDefault)
	DrawRectangle(img, width, height, c)
	return img
}

// DrawRectangle draws rectangle in the provided image.
func DrawRectangle(img *ebiten.Image, width, height int, c color.Color) {
	w, h := float64(width), float64(height)
	ebitenutil.DrawLine(img, 5, 5, w-5, 5, c)
	ebitenutil.DrawLine(img, 5, 5, 5, h-5, c)
	ebitenutil.DrawLine(img, w-5, 5, w-5, h-5, c)
	ebitenutil.DrawLine(img, 5, h-5, w-5, h-5, c)
}
