package component

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Component represents an abstract rectangular UI element.
type Component interface {
	// Bounds returns a rectangle in logical screen coordinates.
	Bounds() image.Rectangle
	// SetBounds sets a rectangle in logical screen coordinates.
	SetBounds(image.Rectangle)
}

// DrawableComponent represents an abstract UI element that can draw itself.
type DrawableComponent interface {
	Component
	Draw(*ebiten.Image)
}
