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
	// SetRect allows to move and/or resize this element relatively to its
	// parent.
	SetRect(rect image.Rectangle)
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
	// OnMousePosition handles mouse position updates.
	// This method shouldn't be executed directly.
	// Note: Mouse position updates aren't sent until mouse capturing is set
	// for the element.
	OnMousePosition(*MousePositionEvent)
}

// Image returns an image area representing this element and its rectangle in
// screen coordinates.
func Image(element Element) (*ebiten.Image, image.Rectangle) {
	return elementImage(element.UI(), element)
}

// ElementAt returns an element under the point. The top-most child get
// returned if there are multiple elements at this point. The point is in
// logical screen coordinates (this includes screen scaling).
func ElementAt(ui *UI, point image.Point) (Element, image.Rectangle) {
	return elementAt(ui, point)
}

// CaptureMouse captures mouse inputs to the element specified.
func CaptureMouse(element Element) {
	captureMouse(element.UI(), element)
}

// UncaptureMouse stops sending mouse inputs to this element set with
// CaptureMouse. The function does nothing if this isn't the same element as
// the one passed to CaptureMouse call.
func UncaptureMouse(element Element) {
	uncaptureMouse(element.UI(), element)
}

// ScreenRect returns the element's rectangle in screen coordinates.
func ScreenRect(element Element) image.Rectangle {
	return screenRect(element.UI(), element)
}
