// Package simulation provides a simple framework which intention is to render
// shapes, provide physics simulation, etc.
package simulation

import (
	"github.com/yarcat/playground/geom/body"
	"github.com/yarcat/playground/geom/contrib/container/orderedmap"
	"github.com/yarcat/playground/geom/ui"
)

// Simulation is a container that should be executed with ebiten.RunGame. It
// represents a "world" -- it embeds rigid bodies and other objects that will be
// fired periodically to execute various simulations.
type Simulation struct {
	ui.Element
	bodies *orderedmap.OrderedMap // *body.Body: nil
}

// New returns new instance of Simulation, which will draw itself in the
// provided element.
func New(element ui.Element) *Simulation {
	return &Simulation{
		Element: element,
		bodies:  orderedmap.New(),
	}
}

// AddBody registers new body in the simulation.
func (s *Simulation) AddBody(b *body.Body) {
	s.bodies.Set(b, nil)
}

// OnDraw redraws the simulation.
func (s *Simulation) OnDraw(evt *ui.DrawEvent) {
	img := ui.Image(s)
	body.Present(img, newBodyIterator(s.bodies))
}
