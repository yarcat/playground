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
	Phi        float64
}

// MoveTo moves the rectangle center.
func (r *R) MoveTo(x, y float64) {
	r.X, r.Y = x, y
}

// P represents a polygon with the center in X, Y.
type P struct {
	X, Y float64
	Phi  float64
	// V is vertices.
	V []vector.Vector
	// E is edges defined as connection of two vertex points.
	E [][2]int
}

// MoveTo moves the polygon center.
func (p *P) MoveTo(x, y float64) {
	p.X, p.Y = x, y
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
	dx := b.X - a.X
	w2 := (a.W + b.W) / 2
	px := w2 - math.Abs(dx)
	if px <= 0 { // No X-axis intersection.
		return
	}

	dy := b.Y - a.Y
	h2 := (a.H + b.H) / 2
	py := h2 - math.Abs(dy)
	if py <= 0 { // No Y-axis intersection.
		return
	}

	ok = true

	if px <= py {
		intersection.P = px
		if dx > 0 {
			intersection.N = vector.New(1, 0)
		} else {
			intersection.N = vector.New(-1, 0)
		}
	} else {
		intersection.P = py
		if dy > 0 {
			intersection.N = vector.New(0, 1)
		} else {
			intersection.N = vector.New(0, -1)
		}
	}

	return intersection, true
}

// Polygons intersects two polygons and returns intersection info.
func Polygons(a, b P) (intersection I, ok bool) {
	// TODO(yarcat): Switch to the scene coordinates. Ensure we don't commpute
	// this every time.
	av := toScene(a.X, a.Y, a.Phi, a.V)
	bv := toScene(b.X, b.Y, b.Phi, b.V)

	p1 := leastP(av, a.E, bv, b.E)
	if p1 > 0 {
		return
	}
	p2 := leastP(bv, b.E, av, a.E)
	if p2 > 0 {
		return
	}
	intersection.P = max(p1, p2)
	// TODO(yarcat): Fill in the normal vector.
	return intersection, true
}
