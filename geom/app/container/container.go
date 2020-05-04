// Package container provides interface for components that may contain other components.
package container

import (
	"github.com/yarcat/playground/geom/app/component"
)

// Container allows to add other components.
type Container interface {
	AddComponent(component.Component)
}
