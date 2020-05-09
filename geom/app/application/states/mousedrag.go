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
}

// Update handles mouse motion changes and decides whether DragStateDragged
// should be sent out.
func (state *MouseDragState) Update(pt image.Point) {
	if pt.Eq(state.cursor) {
		return
	}
	if !state.drag {
		ptMod2 := pt.X*pt.X + pt.Y*pt.Y
		curMod2 := state.cursor.X*state.cursor.X + state.cursor.Y*state.cursor.Y
		d4 := curMod2 - ptMod2
		d4 = d4 * d4
		if d4 < 10000 { // TODO(yarcat): Make this configurable.
			return
		}
		state.drag = true
		state.action.Reset()
	}
	evt := state.host.DragEvent(pt.Sub(state.cursor), ftrs.DragStateDragged)
	state.cursor = pt
	state.features.NotifyDrag(evt)
}
