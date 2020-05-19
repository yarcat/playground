package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type buttonStateFn func(pressed bool)

type mouseManager struct {
	ui              *UI
	leftButtonState buttonStateFn
	capture         map[Element]interface{}
}

func newMouseManager(ui *UI) *mouseManager {
	manager := &mouseManager{
		ui:      ui,
		capture: make(map[Element]interface{}),
	}
	manager.leftButtonState = manager.waitButtonPressed
	return manager
}

func (mm *mouseManager) update() {
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if len(mm.capture) > 0 {
		mm.sendMouseEvent(func(element Element, point image.Point) {
			// TODO(yarcat): Send those only if relative position has changed.
			SendEvent(element, &MousePositionEvent{Cursor: point})
		})
	}
	mm.leftButtonState(pressed)
}

func (mm *mouseManager) waitButtonPressed(pressed bool) {
	if !pressed {
		return
	}
	mm.sendMouseEvent(func(element Element, point image.Point) {
		SendEvent(element, &MouseButtonPressedEvent{Cursor: point})
	})
	mm.leftButtonState = mm.waitButtonReleased
}

func (mm *mouseManager) waitButtonReleased(pressed bool) {
	if pressed {
		return
	}
	mm.sendMouseEvent(func(element Element, point image.Point) {
		SendEvent(element, &MouseButtonReleasedEvent{Cursor: point})
	})
	mm.leftButtonState = mm.waitButtonPressed
}

func (mm *mouseManager) sendMouseEvent(sendEvent func(Element, image.Point)) {
	cursor := image.Pt(ebiten.CursorPosition())
	if len(mm.capture) == 0 {
		element, point := mm.underPoint(cursor)
		sendEvent(element, point)
		return
	}
	// This is a potential infinite recursion.
	for element := range mm.capture {
		rect := ScreenRect(element)
		sendEvent(element, cursor.Sub(rect.Min))
	}
}

func (mm *mouseManager) captureMouse(element Element) {
	mm.capture[element] = nil
}

func (mm *mouseManager) uncaptureMouse(element Element) {
	delete(mm.capture, element)
}

// underPoint returns an element under cursor and cursor position
// in the element coordinates.
func (mm *mouseManager) underPoint(point image.Point) (element Element, cursor image.Point) {
	var rect image.Rectangle
	element, rect = ElementAt(mm.ui, point)
	return element, point.Sub(rect.Min)
}
