package main

import "github.com/yarcat/playground/geom/app/intersect"

type intersector struct {
	circles []*circle
	aabbs   []*rect
}

func (is *intersector) addC(c *circle) {
	is.circles = append(is.circles, c)
}

func (is *intersector) computeC(c *circle) {
	for _, other := range is.circles {
		if c == other {
			continue
		}
		if xi, ok := intersect.Circles(c.C, other.C); ok {
			c.intersected(other, xi)
		} else {
			c.intersected(nil, intersect.I{})
		}
	}
}

func (is *intersector) addR(r *rect) {
	is.aabbs = append(is.aabbs, r)
}

func (is *intersector) computeR(r *rect) {
	for _, other := range is.aabbs {
		if r == other {
			continue
		}
		if xi, ok := intersect.Rectangles(r.R, other.R); ok {
			r.intersected(other, xi)
		} else {
			r.intersected(nil, intersect.I{})
		}
	}
}
