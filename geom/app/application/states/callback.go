package states

// Callback stores a reference to a callback to invoke. Nothing is invoked,
// if underlying callback is set to nil.
type Callback struct {
	fn func()
}

// NewCallback returns an object which can invoke underlying callback.
func NewCallback(fn func()) *Callback {
	return &Callback{fn}
}

// Reset sets the underlying funtion to nil.
func (cb *Callback) Reset() {
	cb.fn = nil
}

// Run executes a callback if it isn't nil.
func (cb *Callback) Run() {
	if cb.fn != nil {
		cb.fn()
	}
}
