// Package body provides various bodies that could be used by simulation.
package body

import (
	"github.com/yarcat/playground/geom/shape"
	vec "github.com/yarcat/playground/geom/vector"
)

// Body is a simple body that cannot move by itself, but has a position on the
// screen.
type Body struct {
	Shape shape.Shape
	Pos   vec.Vector
}

// Iterator helps to iterate over a seqence of bodies.
type Iterator interface {
	// Next advances iterator to the next value.
	Next() bool
	// Value returns current value of the iterator.
	Value() *Body
}
