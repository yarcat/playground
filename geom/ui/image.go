package ui

import "image/color"

// Image represents a drawable area.
type Image interface {
	// Set sets the color value of the given pixel.
	Set(x, y int, c color.Color)
}
