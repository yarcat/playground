package features

// DragListener is a callback executed for drag events.
type DragListener func(DragEvent)

// ListenDrag registers a listener.
func ListenDrag(fn DragListener) FeatureOption {
	return func(features *Features) {
		features.dragFn = append(features.dragFn, fn)
	}
}

// ListensDrag returns true if drag listener is set.
func (f *Features) ListensDrag() bool {
	return f != nil && len(f.dragFn) > 0
}

// NotifyDrag fires drag event callback.
func (f *Features) NotifyDrag(evt DragEvent) {
	if f == nil {
		return
	}
	for _, fn := range f.dragFn {
		fn(evt)
	}
}
