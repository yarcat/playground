package features

// MotionListener represents mouse leave/enter callback function.
type MotionListener func(MotionEvent)

// ListenMouseEnter registers a callback executed when mouse enters an element.
func ListenMouseEnter(fn MotionListener) FeatureOption {
	return func(features *Features) {
		features.mouseEnterFn = fn
	}
}

// ListensMouseEnter returns true if mouse enter callback is set.
func (f *Features) ListensMouseEnter() bool {
	return f != nil && f.mouseEnterFn != nil
}

// NotifyMouseEnter executes mouse enter callback.
func (f *Features) NotifyMouseEnter(evt MotionEvent) {
	if f != nil && f.mouseEnterFn != nil {
		f.mouseEnterFn(evt)
	}
}

// ListenMouseLeave registers a callback executed when mouse pointer leaves an
// element boundary.
func ListenMouseLeave(fn MotionListener) FeatureOption {
	return func(features *Features) {
		features.mouseLeaveFn = fn
	}
}

// ListensMouseLeave returns true if mouse leave callback is set.
func (f *Features) ListensMouseLeave() bool {
	return f != nil && f.mouseLeaveFn != nil
}

// NotifyMouseLeave executes mouse leave callback.
func (f *Features) NotifyMouseLeave(evt MotionEvent) {
	if f != nil && f.mouseLeaveFn != nil {
		f.mouseLeaveFn(evt)
	}
}

// ListensMouseMotion returns true if either mouse enter or leave callback
// functions are set.
func (f *Features) ListensMouseMotion() bool {
	return f != nil && (f.mouseEnterFn != nil || f.mouseLeaveFn != nil)
}
