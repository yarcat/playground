// Package drag implements a component that can be dragged.
package drag

import (
	"github.com/yarcat/playground/geom/app/component"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

type dragFeatures struct {
	ftrs.Features
}

// Drag makes an underlying component mouse-draggable.
type Drag struct {
	component.WithLifecycle
	pressed bool
}

// EnableFor wraps a component with lifecycle and enable its dragging.
func EnableFor(c component.WithLifecycle) *Drag {
	return &Drag{WithLifecycle: c}
}

// HandleAdded is called after the button is added to its parent.
func (d *Drag) HandleAdded(parent component.Component, features *ftrs.Features) {
	d.WithLifecycle.HandleAdded(parent, features)
	features.Add(
		ftrs.ListenDrag(func(evt ftrs.DragEvent) {
			// TODO(yarcat): Make this a public interface.
			type container interface {
				CloserToUser(component.Component)
			}
			if parent, ok := parent.(container); ok {
				parent.CloserToUser(d)
			} else {
				// Panic here to ensure we don't forget fix this when container
				// interface becomes public.
				panic("parent must implement container")
			}
			b := d.Bounds().Add(evt.D())
			d.SetBounds(b)
		}),
	)
}
