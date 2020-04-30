// Package simulation provides a simple framework which intention is to render
// shapes, provide physics simulation, etc.
package simulation

import (
	"github.com/hajimehoshi/ebiten"
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
	image  *ebiten.Image
}

// New returns new instance of Simulation, which will draw itself in the
// provided element.
func New(element ui.Element) *Simulation {
	size := element.Rect().Size()
	image, _ := ebiten.NewImage(size.X, size.Y, ebiten.FilterDefault)
	return &Simulation{
		Element: element,
		bodies:  orderedmap.New(),
		image:   image,
	}
}

// AddBody registers new body in the simulation.
func (s *Simulation) AddBody(b *body.Body) {
	s.bodies.Set(b, nil)
}

// OnDraw redraws the simulation.
func (s *Simulation) OnDraw(evt *ui.DrawEvent) {
	img, rect := ui.Image(s)
	s.image.Clear()
	body.Present(s.image, newBodyIterator(s.bodies))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	img.DrawImage(s.image, op)
}
