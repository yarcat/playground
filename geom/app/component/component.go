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

// WithLifecycle represents a component, which can be added to the application.
type WithLifecycle interface {
	Component
	// HandleAdded is called right after the component was added to its parent.
	// A callee should use provided features object to register component functionality.
	HandleAdded(parent Component, features *Features)
}
