package main

// Window represents a window object.
type Window interface {
	// Parent returns parent window.
	Parent() Window
	// Children return child windows.
	Children() []Window
	// Destroy frees resources associated with the window.
	Destroy()
}

// NewWindowImpl returns basic WindowImpl instance.
func NewWindowImpl(app Application, parent Window) *WindowImpl {
	return &WindowImpl{
		app:    app,
		parent: parent,
	}
}

// WindowImpl is an abstract window.
type WindowImpl struct {
	app      Application
	parent   Window
	children []Window
}

// Parent returns nil as application doesn't have a parent.
func (w WindowImpl) Parent() Window { return nil }

// Children returns children of the app.
func (w WindowImpl) Children() []Window { return w.children }

// Destroy frees resources associated with the window.
func (w WindowImpl) Destroy() {}
