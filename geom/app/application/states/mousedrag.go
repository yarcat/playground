package states

import (
	"image"

	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

// MouseDragStateHost is an object that can remove/unregister mouse drag states.
type MouseDragStateHost interface {
	RemoveMouseDragState(*MouseDragState)
	DragEvent(delta image.Point, state ftrs.DragState) ftrs.DragEvent
}

// MouseDragState detects mouse drags and notifies a component. The state disables
// the action if mouse drag was detected.
type MouseDragState struct {
	host     MouseDragStateHost
	features *ftrs.Features
	action   *Callback
	cursor   image.Point
	drag     bool
}

// NewMouseDragState returns new mouse drag state instance. The instance resets
// (disables) action callback if mouse drag was detected.
func NewMouseDragState(host MouseDragStateHost, features *ftrs.Features, action *Callback, cursor image.Point) *MouseDragState {
	return &MouseDragState{
		host:     host,
		features: features,
		action:   action,
		cursor:   cursor,
	}
}

// Released handles mouse button release, which notifies DragStateEnded if the
// drag was detected. The state also removes itself from the host.
func (state *MouseDragState) Released(features *ftrs.Features) {
	state.host.RemoveMouseDragState(state)
	if !state.drag {
		return
	}
	evt := state.host.DragEvent(image.Point{}, ftrs.DragStateEnded)
	state.features.NotifyDrag(evt)
}

// Update handles mouse motion changes and decides whether DragStateDragged
// should be sent out.
func (state *MouseDragState) Update(pt image.Point) {
	if pt.Eq(state.cursor) {
		return
	}
	delta := pt.Sub(state.cursor)
	if !state.drag {
		// TODO(yarcat): Make this configurable.
		const minDragToDetectSqr = 100
		len2 := delta.X*delta.X + delta.Y*delta.Y
		if len2 < minDragToDetectSqr {
			return
		}
		state.drag = true
		state.action.Reset()
	}
	evt := state.host.DragEvent(delta, ftrs.DragStateDragged)
	state.cursor = pt
	state.features.NotifyDrag(evt)
}
