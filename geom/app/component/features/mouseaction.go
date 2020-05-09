package features

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
