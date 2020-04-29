// Package simulation provides a simple framework which intention is to render
// shapes, provide physics simulation, etc.
package simulation

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/body"
	"github.com/yarcat/playground/geom/contrib/container/orderedmap"
)

// Simulation is a container that should be executed with ebiten.RunGame. It
// represents a "world" -- it embeds rigid bodies and other objects that will be
// fired periodically to execute various simulations.
type Simulation struct {
	// Desired screen dimensions.
	screenWidth, screenHeight int
	bodies                    *orderedmap.OrderedMap // *body.Body: nil
}

// New returns new instance of Simulation ready to be executed with
// ebiten.RunGame.
func New(screenWidth, screenHeight int) *Simulation {
	return &Simulation{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		bodies:       orderedmap.New(),
	}
}

func (s *Simulation) update() {}

func (s *Simulation) draw(screen *ebiten.Image) {
	body.Present(screen, newBodyIterator(s.bodies))
	debugPrint(screen, s)
}

func (s *Simulation) addBody(b *body.Body) {
	s.bodies.Set(b, nil)
}
