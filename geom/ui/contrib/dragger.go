package contrib

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/ui"
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
	log.Println("Dragger.OnMouseButtonPressed: Mouse captured.")
}

// OnMouseButtonReleased handles mouse button releases.
func (dragger *Dragger) OnMouseButtonReleased(event *ui.MouseButtonReleasedEvent) {
	ui.UncaptureMouse(dragger)
	log.Println("Dragger.OnMouseButtonReleased: Mouse released.")
}

// OnMousePosition handles mouse position updates.
func (dragger *Dragger) OnMousePosition(event *ui.MousePositionEvent) {
	cursor := ui.ScreenRect(dragger).Min.Add(event.Cursor)
	d := cursor.Sub(dragger.lastCursor)
	log.Printf("Dragger.OnMousePosition: event.Cursor = %#v, d = %#v", event.Cursor, d)
	rect := dragger.Rect().Add(d)
	dragger.lastCursor = cursor
	dragger.SetRect(rect)
	dragger.fn(rect)
}

// OnDraw handles draw events.
func (dragger *Dragger) OnDraw(event *ui.DrawEvent) {
	img, rect := ui.Image(dragger)
	rectImg, _ := ebiten.NewImage(rect.Dx(), rect.Dy(), ebiten.FilterDefault)
	defer rectImg.Dispose()
	rectImg.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	img.DrawImage(rectImg, op)
}
