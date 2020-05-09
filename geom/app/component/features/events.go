package features

import "image"

// GestureEvent represents a generic mouse or touch event.
// TODO(yarcat): Move this away.
type GestureEvent interface {
	Pressed() bool
	Pos() image.Point
}

// MotionEvent represents an event sent with Leave/Enter
// notifications.
// TODO(yarcat): Move this away.
type MotionEvent interface{}
