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

// DragState represents current drag state.
type DragState int

const (
	// DragStateDragged is sent when there was mouse motion detected while in drag.
	DragStateDragged DragState = iota
	// DragStateEnded is sent when a drag state ends.
	DragStateEnded
)

// DragEvent represents a mouse motion event which disables action event.
// This allows to move components on the screen with a mouse.
type DragEvent interface {
	// D returns a relative cursor motion since the last time the event was triggered.
	D() image.Point
	Pos() image.Point
	State() DragState
}
