package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Element represents a user interface element. It will receive mouse, keyboard
// and other events.
type Element interface {
	// UI returns UI instance that hosts this element.
	UI() *UI
	// Rect returns a rectangle that defines dimensions of this element.
	Rect() image.Rectangle
	// AddChild adds new child to this element.
	AddChild(Element)

	// OnDraw handles draw events.
	// This method shouldn't be executed directly.
	OnDraw(*DrawEvent)
	// OnMouseButtonPressed handles mouse button press notifications.
	// This method shouldn't be executed directly.
	OnMouseButtonPressed(*MouseButtonPressedEvent)
	// OnMouseButtonPressed handles mouse button release notifications.
	// This method shouldn't be executed directly.
	OnMouseButtonReleased(*MouseButtonReleasedEvent)
}

// Image returns an image area representing this element and its rectangle in
// screen coordinates.
func Image(element Element) (*ebiten.Image, image.Rectangle) {
	return elementImage(element.UI(), element)
}

// ElementAt returns an element under the point. The top-most child get
// returned if there are multiple elements at this point. The point is in
// logical screen coordinates (this includes screen scaling).
func ElementAt(ui *UI, point image.Point) Element {
	return elementAt(ui, point)
}
