package main

// Window represents a window object.
type Window interface {
	// App returns an associated application instance.
	App() Application
	// Parent returns parent window.
	Parent() Window
	// Children return child windows.
	Children() []Window
	// Destroy frees resources associated with the window.
	Destroy()
	// AppendChild appends a child window.
	AppendChild(Window)
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

// App returns the application instance.
func (w WindowImpl) App() Application {
	return w.app
}

// Parent returns nil as application doesn't have a parent.
func (w WindowImpl) Parent() Window { return nil }

// Children returns children of the app.
func (w WindowImpl) Children() []Window { return w.children }

// AppendChild appends a child.
func (w *WindowImpl) AppendChild(child Window) {
	// TODO(yarcat): Need to ensure child/parent relation is correct.
	w.children = append(w.children, child)
}

// Destroy frees resources associated with the window.
func (w WindowImpl) Destroy() {}
