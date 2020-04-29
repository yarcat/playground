package body

import "github.com/hajimehoshi/ebiten"

// Present draws bodies on the screen.
func Present(screen *ebiten.Image, it Iterator) {
	op := &ebiten.DrawImageOptions{}
	img := &originImage{
		img: screen,
		op:  op,
	}
	for it.Next() {
		b := it.Value()
		op.GeoM.Translate(b.Pos.X, b.Pos.Y)
		b.Shape.Draw(img)
		op.GeoM.Translate(-b.Pos.X, -b.Pos.Y)
	}
}
