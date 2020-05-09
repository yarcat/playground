package features

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
