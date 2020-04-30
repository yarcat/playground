package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Element represents a user interface element. It will receive mouse, keyboard
// and other events.
type Element interface {
	// Rect returns a rectangle that defines dimensions of this element.
	Rect() image.Rectangle
	// AddChild adds new child to this element.
	AddChild(Element)
	// Image returns an image area representing this element.
	Image() *ebiten.Image
	// OnDraw handles draw events. This method shouldn't be executed directly.
	OnDraw(*DrawEvent)
}
