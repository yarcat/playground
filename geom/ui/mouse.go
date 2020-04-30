package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

type buttonStateFn func(pressed bool)

type mouseManager struct {
	ui              *UI
	leftButtonState buttonStateFn
}

func newMouseManager(ui *UI) *mouseManager {
	manager := &mouseManager{ui: ui}
	manager.leftButtonState = manager.waitButtonPressed
	return manager
}

func (mm *mouseManager) update() {
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	mm.leftButtonState(pressed)
}

func (mm *mouseManager) waitButtonPressed(pressed bool) {
	if !pressed {
		return
	}
	element, point := underCursor(mm.ui)
	SendEvent(element, &MouseButtonPressedEvent{Cursor: point})
	mm.leftButtonState = mm.waitButtonReleased
}

func (mm *mouseManager) waitButtonReleased(pressed bool) {
	if pressed {
		return
	}
	element, point := underCursor(mm.ui)
	SendEvent(element, &MouseButtonReleasedEvent{Cursor: point})
	mm.leftButtonState = mm.waitButtonPressed
}

// underCursor returns an element under cursor and cursor position
// in the element coordinates.
func underCursor(ui *UI) (element Element, cursor image.Point) {
	cursor = image.Pt(ebiten.CursorPosition())
	var rect image.Rectangle
	element, rect = ElementAt(ui, cursor)
	return element, cursor.Sub(rect.Min)
}
