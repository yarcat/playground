package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// NewLabel return new label instance.
func NewLabel(app Application, parent Window) *Label {
	l := &Label{Window: NewWindowImpl(app, parent)}
	parent.AppendChild(l)
	return l
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

// Draw draws label on screen.
func (l Label) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, l.text)
}
