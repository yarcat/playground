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
	motions              map[*features.Features]*states.MouseMotionState
}

func (m *gestureManagerImpl) update() {
	underCursor := &underCursor{app: m.app}

	if underCursor.features().ListensMouseMotion() {
		if _, ok := m.motions[underCursor.features()]; !ok {
			m.motions[underCursor.features()] = &states.MouseMotionState{
				Host:      (*motionHostAdapter)(m),
				Component: underCursor.component(),
				Features:  underCursor.features(),
			}
		}
	}
	pt := image.Pt(ebiten.CursorPosition())
	for _, state := range m.motions {
		state.Update(pt)
	}

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

type motionHostAdapter gestureManagerImpl

// RemoveMouseButtonState unregisters the state provided.
func (adap *motionHostAdapter) RemoveMouseMotionState(state *states.MouseMotionState) {
	delete((*gestureManagerImpl)(adap).motions, state.Features)
}

// MotionEvent returns a motion event instance.
func (*motionHostAdapter) MotionEvent() features.MotionEvent {
	return nil
}
