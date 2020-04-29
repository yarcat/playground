package shape

import "image/color"

// OriginImage is an abstraction on top of ebiten.Image that ensures automagical
// translations and rotations.
type OriginImage interface {
	Set(x, y int, c color.Color)
}
