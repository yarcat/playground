package application

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/yarcat/playground/geom/app/application/states"
	"github.com/yarcat/playground/geom/app/component/features"
)

// gestureManagerImpl provides functionality that helps to handle mouse and
// touch screen.
type gestureManagerImpl struct {
	app                  *App
	states, removeStates []states.MouseButtonState
}

func (m *gestureManagerImpl) update() {
	underCursor := &underCursor{app: m.app}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if underCursor.features().ListensMouseButtons() {
			underCursor.features().NotifyMouseButtons(gestureEvent{})
			state := states.NewMouseButtonState(
				(*removerAdapter)(m), underCursor.features())
			m.states = append(m.states, state)
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, mstate := range m.states {
			mstate.Released(underCursor.features())
		}
	}

	m.states = removeStates(m.states, m.removeStates)
	m.removeStates = m.removeStates[:0]
}

type removerAdapter gestureManagerImpl

func (ra *removerAdapter) RemoveMouseButtonState(state states.MouseButtonState) {
	(*gestureManagerImpl)(ra).removeStates = append(
		(*gestureManagerImpl)(ra).removeStates, state)
}

func (ra removerAdapter) GestureEvent() features.GestureEvent {
	return gestureEvent{}
}

// TODO(yarcat): Move this to a better place.
type gestureEvent struct{}

// Pos returns mouse cursor position in logical screen coordinates.
func (evt gestureEvent) Pos() image.Point {
	return image.Pt(ebiten.CursorPosition())
}

// Pressed returns whether left mouse button is pressed.
func (evt gestureEvent) Pressed() bool {
	return ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}
