package ui

// Event is an abstract event which can be sent to elements.
type Event interface {
	DispatchEvent(EventHandler)
}

// EventHandler dispatches events by calling corresponding methods.
type EventHandler interface {
	OnDraw(*DrawEvent)
}

// SendEvent sends the event to the element.
func SendEvent(element Element, e Event) {
	e.DispatchEvent(element)
}

// DrawEvent is sent to elements when they need to draw themselves.
type DrawEvent struct{}

// DispatchEvent calls EventHandler.OnDraw.
func (evt *DrawEvent) DispatchEvent(d EventHandler) {
	d.OnDraw(evt)
}
