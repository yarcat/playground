package application

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

// gestureManagerImpl provides functionality that helps to handle mouse and touch screen.
type gestureManagerImpl struct {
	app    *App
	states []*mouseState
}

type mouseState struct {
	features *ftrs.Features
	// stateFn is a state handler function executed in update.
	stateFn func() (keep bool)
}

func newMouseState(features *ftrs.Features) *mouseState {
	s := &mouseState{features: features}
	s.stateFn = s.sendPressed
	return s
}

func (state *mouseState) sendPressed() (keep bool) {
	state.features.NotifyMouseButtons(state)
	state.stateFn = state.sendReleased
	return true
}

func (state *mouseState) sendReleased() (keep bool) {
	if !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return true
	}
	state.features.NotifyMouseButtons(state)
	state.features.NotifyAction()
	return false
}

// update returns whether this state should be kept.
func (state *mouseState) update() (keep bool) {
	return state.stateFn()
}

// Pos returns current cursor position.
func (state mouseState) Pos() image.Point {
	return image.Pt(ebiten.CursorPosition())
}

// Pressed returns whether the left button is pressed.
func (state mouseState) Pressed() bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

type gestureEventImpl struct {
	pressed bool
	pos     image.Point
}

// Pressed returns true if the input is in a pressed state.
func (evt gestureEventImpl) Pressed() bool {
	return evt.pressed
}

// Pos returns cursor position.
func (evt gestureEventImpl) Pos() image.Point {
	return evt.pos
}

func (m *gestureManagerImpl) update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		comp := m.app.ComponentAt(image.Pt(ebiten.CursorPosition()))
		f := m.app.features[comp]
		if f.ListensMouseButtons() {
			m.states = append(m.states, newMouseState(f))
		}
	}

	for i := 0; i < len(m.states); {
		if m.states[i].update() {
			i++
		} else {
			m.states = append(m.states[:i], m.states[i+1:]...)
		}
	}
}
