// Package features provides collection of component callbacks and settings.
package features

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Features represents an object that it used to configure component's behavior.
// It allows to enable or disable events, make the component drawable, etc.
// Main intention of this class is to ensure we avoid enormous amount of methods
// in the Component. Custom components should register their functionality.
type Features struct {
	drawFn                func(*ebiten.Image)
	mouseButtonListenerFn MouseButtonListener
	actionListenerFn      ActionListener
}

// FeatureOption is a function that can update concrete feature in the set of features.
type FeatureOption func(*Features)

// Add applies given features.
func (f *Features) Add(opts ...FeatureOption) {
	for _, optFn := range opts {
		optFn(f)
	}
}

// Draw registers a draw function to be called to render a frame.
func Draw(fn func(*ebiten.Image)) FeatureOption {
	return func(features *Features) {
		features.drawFn = fn
	}
}

// Draw calls registered draw function that renders a component every frame.
func (f *Features) Draw(screen *ebiten.Image) {
	if f != nil && f.drawFn != nil {
		f.drawFn(screen)
	}

}

// GestureEvent represents a generic mouse or touch event.
// TODO(yarcat): Move this away.
type GestureEvent interface {
	Pressed() bool
	Pos() image.Point
}

// MouseButtonListener is callback called upon mouse button presses and releases.
type MouseButtonListener func(GestureEvent)

// ListenMouseButtons registers a callback called upon mouse button presses and
// releases.
func ListenMouseButtons(fn MouseButtonListener) FeatureOption {
	return func(features *Features) {
		features.mouseButtonListenerFn = fn
	}
}

// ListensMouseButtons returns true if any mouse event listeners are registered.
func (f *Features) ListensMouseButtons() bool {
	return f != nil && (f.mouseButtonListenerFn != nil || f.actionListenerFn != nil)
}

// NotifyMouseButtons fires mouse button listening callback.
func (f *Features) NotifyMouseButtons(evt GestureEvent) {
	if f != nil && f.mouseButtonListenerFn != nil {
		f.mouseButtonListenerFn(evt)
	}
}

// ActionListener is a function called for actions e.g. clicks.
type ActionListener func()

// ListenAction registers a callback executed for actions e.g. clicks.
func ListenAction(fn ActionListener) FeatureOption {
	return func(features *Features) {
		features.actionListenerFn = fn
	}
}

// NotifyAction fires action callback.
func (f *Features) NotifyAction() {
	if f != nil && f.actionListenerFn != nil {
		f.actionListenerFn()
	}
}
