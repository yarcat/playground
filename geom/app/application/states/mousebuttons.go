package states

import (
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

// MouseButtonStateHost is an object that can remove/unregister mouse button
// states.
type MouseButtonStateHost interface {
	RemoveMouseButtonState(*MouseButtonState)
	GestureEvent() ftrs.GestureEvent
}

// MouseButtonState sends mouse button notifications.
type MouseButtonState struct {
	host     MouseButtonStateHost
	features *ftrs.Features
	action   *Callback
}

// NewMouseButtonState returns new mouse button state instance. It usually makes
// sense to call its Pressed or Released method after creating.
func NewMouseButtonState(host MouseButtonStateHost, features *ftrs.Features, action *Callback) *MouseButtonState {
	return &MouseButtonState{
		host:     host,
		features: features,
		action:   action,
	}
}

// Released notifies the state that a mouse button was released. It notifies
// mouse button state using features and attempts to remove itself.
func (state *MouseButtonState) Released(features *ftrs.Features) {
	state.features.NotifyMouseButtons(state.host.GestureEvent())
	if state.features == features {
		state.action.Run()
	}
	state.host.RemoveMouseButtonState(state)
}
