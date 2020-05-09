package features

// DragListener is a callback executed for drag events.
type DragListener func(DragEvent)

// ListenDrag registers a listener.
func ListenDrag(fn DragListener) FeatureOption {
	return func(features *Features) {
		features.dragFn = fn
	}
}

// ListensDrag returns true if drag listener is set.
func (f *Features) ListensDrag() bool {
	return f != nil && f.dragFn != nil
}

// NotifyDrag fires drag event callback.
func (f *Features) NotifyDrag(evt DragEvent) {
	if f != nil && f.dragFn != nil {
		f.dragFn(evt)
	}
}
