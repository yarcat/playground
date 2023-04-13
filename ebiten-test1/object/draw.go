package object

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func Draw(img *ebiten.Image, o Object, side float64) {
	vector.StrokeLine(img, float32((o.pos.X+o.v.X*0.1)*side), float32((o.pos.Y+o.v.Y*0.1)*side),
		float32((o.pos.X+o.v.X)*side), float32((o.pos.Y+o.v.Y)*side), 1, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeCircle(img, float32(o.pos.X*side), float32(o.pos.Y*side), float32(side*0.1), 1, color.RGBA{0, 255, 0, 255}, true)
}
