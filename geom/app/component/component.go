package component

import (
	"image"
)

// Component represents an abstract rectangular UI element.
type Component interface {
	// Bounds returns a rectangle in logical screen coordinates.
	Bounds() image.Rectangle
	// SetBounds sets a rectangle in logical screen coordinates.
	SetBounds(image.Rectangle)
}
