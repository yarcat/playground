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
	mm.sendEventToElementUnderCursor(&MouseButtonPressedEvent{})
	mm.leftButtonState = mm.waitButtonReleased
}

func (mm *mouseManager) waitButtonReleased(pressed bool) {
	if pressed {
		return
	}
	mm.sendEventToElementUnderCursor(&MouseButtonReleasedEvent{})
	mm.leftButtonState = mm.waitButtonPressed
}

func (mm *mouseManager) sendEventToElementUnderCursor(event Event) {
	element := ElementAt(mm.ui, image.Pt(ebiten.CursorPosition()))
	SendEvent(element, event)
	mm.leftButtonState = mm.waitButtonReleased
}
