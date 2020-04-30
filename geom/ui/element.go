package ui

import (
	"image"
	"log"
)

// elementImpl represents a rectangular UI area. The area could be used to receive
// user interface events e.g. DrawEvent.
type elementImpl struct {
	// ui is a main user interface manager.
	ui *UI
	// rect is a rectangle that defines this element.
	rect image.Rectangle
}

// NewElement returns rectangular element. The element will not display itself
// or do anything useful until attached to the screen.
//
// Provided rectangle represents a rectangular area in the coordinates of the
// element's parent. These will be appliation window coordinates if this is a
// top-level element.
func NewElement(ui *UI, rect image.Rectangle) Element {
	return &elementImpl{
		ui:   ui,
		rect: rect,
	}
}

// UI returns app intance that hosts this element.
func (e *elementImpl) UI() *UI {
	return e.ui
}

// Rect returns a rectangle that defines dimensions of this element.
func (e *elementImpl) Rect() image.Rectangle {
	return e.rect
}

// SetRect allows to move and/or resize this element relatively to its parent.
func (e *elementImpl) SetRect(rect image.Rectangle) {
	log.Printf("elementImplSetRect: old=%#v, new=%#v", e.rect, rect)
	e.rect = rect
}

// AddChild adds new child to this element.
func (e *elementImpl) AddChild(child Element) {
	e.ui.Attach(child, e)
}

// OnDraw handles draw events.
// This method shouldn't be executed directly.
func (e *elementImpl) OnDraw(evt *DrawEvent) {}

// OnMouseButtonPressed handles mouse button press notifications.
// This method shouldn't be executed directly.
func (e *elementImpl) OnMouseButtonPressed(*MouseButtonPressedEvent) {}

// OnMouseButtonPressed handles mouse button release notifications.
// This method shouldn't be executed directly.
func (e *elementImpl) OnMouseButtonReleased(*MouseButtonReleasedEvent) {}

// OmMousePosition handles mouse position updates.
// This method shouldn't be executed directly.
func (e *elementImpl) OnMousePosition(*MousePositionEvent) {}
