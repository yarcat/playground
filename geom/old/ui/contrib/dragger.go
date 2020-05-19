package contrib

import (
	"image"

	"github.com/yarcat/playground/geom/old/ui"
)

// Dragger is an element that could be dragged with a mouse.
type Dragger struct {
	ui.Element
	lastCursor image.Point
	fn         func(image.Rectangle)
}

// NewDragger creates new Dragger element.
func NewDragger(app *ui.UI, rect image.Rectangle, fn func(image.Rectangle)) *Dragger {
	element := ui.NewElement(app, rect)
	return &Dragger{Element: element, fn: fn}
}

// OnMouseButtonPressed handles mouse button presses.
func (dragger *Dragger) OnMouseButtonPressed(event *ui.MouseButtonPressedEvent) {
	dragger.lastCursor = ui.ScreenRect(dragger).Min.Add(event.Cursor)
	ui.CaptureMouse(dragger)
}

// OnMouseButtonReleased handles mouse button releases.
func (dragger *Dragger) OnMouseButtonReleased(event *ui.MouseButtonReleasedEvent) {
	ui.UncaptureMouse(dragger)
}

// OnMousePosition handles mouse position updates.
func (dragger *Dragger) OnMousePosition(event *ui.MousePositionEvent) {
	cursor := ui.ScreenRect(dragger).Min.Add(event.Cursor)
	d := cursor.Sub(dragger.lastCursor)
	rect := dragger.Rect().Add(d)
	dragger.lastCursor = cursor
	dragger.SetRect(rect)
	dragger.fn(rect)
}
