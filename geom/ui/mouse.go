package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten"
)

type buttonStateFn func(pressed bool)

type mouseManager struct {
	leftButtonState buttonStateFn
}

func newMouseManager() *mouseManager {
	manager := &mouseManager{}
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
	log.Println("released -> pressed")
	mm.leftButtonState = mm.waitButtonReleased
}

func (mm *mouseManager) waitButtonReleased(pressed bool) {
	if pressed {
		return
	}
	log.Println("pressed -> released")
	mm.leftButtonState = mm.waitButtonPressed
}
