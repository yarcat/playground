package application

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/app/component"
)

type mouseManager struct {
	app         *App
	leftPressed bool
	components  []component.Component
}

// MouseEvent represents a generic mouse event.
type MouseEvent interface {
	Pressed() bool
	Pos() image.Point
}

type mouseEventImpl struct {
	pressed bool
	pos     image.Point
}

// Pressed returns true if left mouse button is pressed.
func (evt mouseEventImpl) Pressed() bool {
	return evt.pressed
}

// Pos returns cursor position.
func (evt mouseEventImpl) Pos() image.Point {
	return evt.pos
}

type mousePressedHandler interface {
	OnMousePressed(MouseEvent)
}

type mouseClickedHandler interface {
	OnMouseClicked(MouseEvent)
}

func (m *mouseManager) update() {
	pressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if pressed == m.leftPressed {
		return
	}
	m.leftPressed = pressed
	x, y := ebiten.CursorPosition()
	pos := image.Pt(x, y)
	evt := mouseEventImpl{pressed: pressed, pos: pos}
	if pressed {
		m.components = m.app.ComponentAt(pos)
		for _, c := range m.components {
			if c, ok := c.(mousePressedHandler); ok {
				c.OnMousePressed(evt)
			}
		}
	} else {
		comps := m.app.ComponentAt(pos)
		for _, c := range comps {
			if c, ok := c.(mousePressedHandler); ok {
				c.OnMousePressed(evt)
			}
			if componentIn(c, m.components) {
				if c, ok := c.(mouseClickedHandler); ok {
					c.OnMouseClicked(evt)
				}
			}
		}
	}
}

func componentIn(c component.Component, set []component.Component) bool {
	for _, comp := range set {
		if c == comp {
			return true
		}
	}
	return false
}
