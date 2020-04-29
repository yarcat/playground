package simulation

import "github.com/yarcat/playground/geom/body"

// AddBody registers new body in the simulation.
func AddBody(s *Simulation, b *body.Body) {
	s.addBody(b)
}
