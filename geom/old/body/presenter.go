package body

import "github.com/hajimehoshi/ebiten"

// Present draws bodies on the screen.
func Present(image *ebiten.Image, it Iterator) {
	for it.Next() {
		b := it.Value()
		h, w := b.Image.Size()
		dx := b.Pos.X - float64(h)/2
		dy := b.Pos.Y - float64(w)/2
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(dx, dy)
		image.DrawImage(b.Image, op)
	}
}
