package main

import "github.com/yarcat/playground/geom/app/intersect"

type intersector struct {
	circles []*circle
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
