// Package intersect provides utility functions for intersecting various
// geometrical objects.
package intersect

import (
	"math"

	"github.com/yarcat/playground/geom/vector"
)

// C represents a circle on a plane with the given center and radius.
type C struct {
	X, Y, R float64
}

// MoveTo moves the circle center.
func (c *C) MoveTo(x, y float64) {
	c.X, c.Y = x, y
}

// R represents an Axis-Aligned Bounding Box.
type R struct {
	X, Y, W, H float64
}

// MoveTo moves the rectangle center.
func (r *R) MoveTo(x, y float64) {
	r.X, r.Y = x, y
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
	dx, dy := b.X-a.X, b.Y-a.Y
	len2 := dx*dx + dy*dy
	rr := a.R + b.R
	if len2 > rr*rr {
		return // No intersection
	}

	ok = true

	len := math.Sqrt(len2)
	intersection.P = rr - len

	if isZero(len2) {
		intersection.N = vector.ZN
	} else {
		intersection.N = vector.New(dx/len, dy/len)
	}
	return
}

// Rectangles intersects two AABB rectangles and returns intersection
// information.
func Rectangles(a, b R) (intersection I, ok bool) {
	return intersection, true
}
