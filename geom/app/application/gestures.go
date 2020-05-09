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
	states, removeStates []*states.MouseButtonState
	motions              map[*features.Features]*states.MouseMotionState
	drags, removeDrags   []*states.MouseDragState
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
	for _, state := range m.drags {
		state.Update(pt)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		f := underCursor.features()
		var action *states.Callback
		if f.ListensMouseButtons() || f.ListensDrag() {
			action = states.NewCallback(f.NotifyAction)
		}
		if f.ListensDrag() {
			state := states.NewMouseDragState(
				(*removerAdapter)(m), underCursor.features(), action, pt)
			m.drags = append(m.drags, state)
		}
		if f.ListensMouseButtons() {
			underCursor.features().NotifyMouseButtons(gestureEvent{})
			state := states.NewMouseButtonState(
				(*removerAdapter)(m), underCursor.features(), action)
			m.states = append(m.states, state)
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, mstate := range m.states {
			mstate.Released(underCursor.features())
		}
		for _, state := range m.drags {
			state.Released(underCursor.features())
		}
	}

	m.states = removeStates(m.states, m.removeStates)
	m.removeStates = m.removeStates[:0]

	m.drags = removeDrags(m.drags, m.removeDrags)
	m.removeDrags = m.removeDrags[:0]
}

type removerAdapter gestureManagerImpl

func (ra *removerAdapter) RemoveMouseButtonState(state *states.MouseButtonState) {
	(*gestureManagerImpl)(ra).removeStates = append(
		(*gestureManagerImpl)(ra).removeStates, state)
}

func (ra *removerAdapter) RemoveMouseDragState(state *states.MouseDragState) {
	(*gestureManagerImpl)(ra).removeDrags = append(
		(*gestureManagerImpl)(ra).removeDrags, state)
}

func (ra removerAdapter) GestureEvent() features.GestureEvent {
	return gestureEvent{}
}

func (ra removerAdapter) DragEvent(delta image.Point, state features.DragState) features.DragEvent {
	return dragEvent{delta: delta, state: state}
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

// TODO(yarcat): Move this to a better place.
type dragEvent struct {
	delta image.Point
	state features.DragState
}

// D returns a relative cursor motion since the last time the event was triggered.
func (evt dragEvent) D() image.Point {
	return evt.delta
}

func (dragEvent) Pos() image.Point {
	return image.Pt(ebiten.CursorPosition())
}

func (evt dragEvent) State() features.DragState {
	return evt.state
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
