package states

import (
	"image"

	"github.com/yarcat/playground/geom/app/component"
	"github.com/yarcat/playground/geom/app/component/features"
)

// MouseMotionStateHost is the element that owns this state.
type MouseMotionStateHost interface {
	RemoveMouseMotionState(*MouseMotionState)
	MotionEvent() features.MotionEvent
}

// MouseMotionState represents a state machine that sends mouse enter/leave
// notifications.
type MouseMotionState struct {
	Host      MouseMotionStateHost
	Features  *features.Features
	Component component.Component
	inside    bool
}

// Update desides whether leave/enter notifications should be called, and call
// them (actually only one can be sent at once) if they should.
func (state *MouseMotionState) Update(pt image.Point) {
	inside := pt.In(state.Component.Bounds())
	if inside == state.inside {
		return
	}
	state.inside = inside
	if inside {
		if state.Features.ListensMouseEnter() {
			state.Features.NotifyMouseEnter(state.Host.MotionEvent())
		}
		return
	}
	if state.Features.ListensMouseLeave() {
		state.Features.NotifyMouseLeave(state.Host.MotionEvent())
	}
	// We don't really have to remove it until the component is removed, or
	// it's reconfigured to stop receiving motion notifications. But I'm doing
	// this to ensure there will be no memory leaks in the future, if I decide
	// to add/change some functionality.
	state.Host.RemoveMouseMotionState(state)
}
