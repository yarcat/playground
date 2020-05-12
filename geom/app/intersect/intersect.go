// Package intersect provides utility functions for intersecting various
// geometrical objects.
package intersect

import (
	"github.com/yarcat/playground/geom/vector"
)

// C represents a circle on a plane with the given center and radius.
type C struct {
	X, Y, R float64
}

// I represents an intersection information.
type I struct {
	// P is a penetration value.
	P float64
	// N is an intersection normal vector.
	N vector.Vector
}

// Circles intersects two circles and returns an intersection information.
func Circles(a, b C) (intersection I, ok bool) {
	ok = true
	return
}
