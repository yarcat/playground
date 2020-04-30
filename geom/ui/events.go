package ui

import "image"

// Event is an abstract event which can be sent to elements.
type Event interface {
	DispatchEvent(EventHandler)
}

// EventHandler dispatches events by calling corresponding methods.
type EventHandler interface {
	OnDraw(*DrawEvent)
	OnMouseButtonPressed(*MouseButtonPressedEvent)
	OnMouseButtonReleased(*MouseButtonReleasedEvent)
}

// SendEvent sends the event to the element.
func SendEvent(element Element, e Event) {
	e.DispatchEvent(element)
}

// MouseButtonPressedEvent is sent to the element on which mouse button was pressed.
type MouseButtonPressedEvent struct {
	Cursor image.Point
}

// DispatchEvent calls EventHandler.OnMouseButtonPressed.
func (evt *MouseButtonPressedEvent) DispatchEvent(d EventHandler) {
	d.OnMouseButtonPressed(evt)
}

// MouseButtonReleasedEvent is sent to the element on which mouse button was released.
type MouseButtonReleasedEvent struct {
	Cursor image.Point
}

// DispatchEvent calls EventHandler.OnMouseButtonReleased.
func (evt *MouseButtonReleasedEvent) DispatchEvent(d EventHandler) {
	d.OnMouseButtonReleased(evt)
}

// DrawEvent is sent to elements when they need to draw themselves.
type DrawEvent struct{}

// DispatchEvent calls EventHandler.OnDraw.
func (evt *DrawEvent) DispatchEvent(d EventHandler) {
	d.OnDraw(evt)
}
