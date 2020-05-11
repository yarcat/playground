// Package features provides collection of component callbacks and settings.
package features

import (
	"github.com/hajimehoshi/ebiten"
)

// Features represents an object that it used to configure component's behavior.
// It allows to enable or disable events, make the component drawable, etc.
// Main intention of this class is to ensure we avoid enormous amount of methods
// in the Component. Custom components should register their functionality.
type Features struct {
	drawFn                     func(*ebiten.Image)
	mouseButtonListenerFn      MouseButtonListener
	actionListenerFn           ActionListener
	mouseEnterFn, mouseLeaveFn MotionListener
	dragFn                     []DragListener
}

// FeatureOption is a function that can update concrete feature in the set of features.
type FeatureOption func(*Features)

// Add applies given features.
func (f *Features) Add(opts ...FeatureOption) {
	for _, optFn := range opts {
		optFn(f)
	}
}
