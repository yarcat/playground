package body

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
)

type originImage struct {
	img *ebiten.Image
	op  *ebiten.DrawImageOptions
}

func (oi *originImage) Set(x, y int, c color.Color) {
	newx, newy := oi.op.GeoM.Apply(float64(x), float64(y))
	c = oi.op.ColorM.Apply(c)
	oi.img.Set(int(newx), int(newy), c)
}
