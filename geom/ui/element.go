package ui

import "github.com/hajimehoshi/ebiten"

// elementImpl represents a rectangular UI area. The area could be used to receive
// user interface events e.g. DrawEvent.
type elementImpl struct {
	// ui is a main user interface manager.
	ui *UI
	// rect is a rectangle that defines this element.
	rect Rect
}

// NewElement returns rectangular element. The element will not display itself
// or do anything useful until attached to the screen.
//
// Provided rectangle represents a rectangular area in the coordinates of the
// element's parent. These will be appliation window coordinates if this is a
// top-level element.
func NewElement(ui *UI, rect Rect) Element {
	return &elementImpl{
		ui:   ui,
		rect: rect,
	}
}

// Rect returns a rectangle that defines dimensions of this element.
func (e *elementImpl) Rect() Rect {
	return e.rect
}

// AddChild adds new child to this element.
func (e *elementImpl) AddChild(child Element) {
	e.ui.Attach(child, e)
}

// Image returns an image area representing this element.
func (e *elementImpl) Image() *ebiten.Image {
	return e.ui.image(e)
}

// OnDraw handles draw events. This method shouldn't be executed directly.
func (e *elementImpl) OnDraw(evt *DrawEvent) {}
