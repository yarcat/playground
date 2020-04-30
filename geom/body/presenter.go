package body

import "github.com/hajimehoshi/ebiten"

// Present draws bodies on the screen.
func Present(image *ebiten.Image, it Iterator) {
	op := &ebiten.DrawImageOptions{}
	for it.Next() {
		b := it.Value()
		op.GeoM.Translate(b.Pos.X, b.Pos.Y)
		image.DrawImage(b.Image, op)
		op.GeoM.Translate(-b.Pos.X, -b.Pos.Y)
	}
}
