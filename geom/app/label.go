package main

// NewLabel return new label instance.
func NewLabel(app Application, parent Window) *Label {
	return &Label{}
}

// Label is a simple text widget.
type Label struct {
	Window
	text string
}

// SetText sets label text.
func (l *Label) SetText(text string) {
	l.text = text
}
