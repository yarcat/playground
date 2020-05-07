// Package states contains various internal states e.g. mouse, keyboard, etc.
package states

import (
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

// MouseButtonState represent left button state.
type MouseButtonState interface {
	// Released should be called right after a button was released. It accepts
	// component features under mouse cursor.
	Released(underCursor *ftrs.Features)
}

// MouseButtonStateHost is an object that can remove/unregister mouse button
// states.
type MouseButtonStateHost interface {
	RemoveMouseButtonState(MouseButtonState)
	GestureEvent() ftrs.GestureEvent
}

type mouseButtonState struct {
	host     MouseButtonStateHost
	features *ftrs.Features
}

// NewMouseButtonState returns new mouse button state instance. It usually makes
// sense to call its Pressed or Released method after creating.
func NewMouseButtonState(host MouseButtonStateHost, features *ftrs.Features) MouseButtonState {
	return &mouseButtonState{
		host:     host,
		features: features,
	}
}

// Released notifies the state that a mouse button was released. It notifies
// mouse button state using features and attempts to remove itself.
func (state *mouseButtonState) Released(features *ftrs.Features) {
	if state.features == features {
		state.features.NotifyMouseButtons(state.host.GestureEvent())
		state.features.NotifyAction()
	}
	state.host.RemoveMouseButtonState(state)
}
