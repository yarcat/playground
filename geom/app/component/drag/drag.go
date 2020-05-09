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
		ftrs.ListenDrag(func(evt ftrs.DragEvent) {}),
	)
}
