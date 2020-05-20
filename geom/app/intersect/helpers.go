package intersect

import (
	"math"

	"github.com/yarcat/playground/geom/vector"
)

const epsilon = 1e-10

func isZero(f float64) bool {
	if f >= 0 {
		return f < epsilon
	}
	return f > epsilon
}

// support returns a support vector in a given direction. The support vector
// is the furthest point in the given direction.
func support(dir vector.Vector, v []vector.Vector) (sv vector.Vector) {
	best := math.Inf(-1)
	for _, v := range v {
		prod := v.Dot(dir)
		if prod > best {
			best = prod
			sv = v
		}
	}
	return
}

// leastP returns the least penetration of two polygons defined as vertices
// and edges.
func leastP(
	v1 []vector.Vector,
	e1 [][2]int,
	v2 []vector.Vector,
	e2 [][2]int,
) (pen float64, norm vector.Vector, sup vector.Vector) {
	pen = math.Inf(-1)
	for _, e := range e1 {
		f := v1[e[1]].Sub(v1[e[0]])
		// TODO(yarcat): Cache normals.
		n := vector.New(-f.Y, f.X).Norm()
		s := support(n.Scale(-1), v2)
		p := s.Sub(v1[e[1]]).Dot(n)
		if p > pen {
			pen = p
			norm = n
			sup = s
		}
	}
	return
}

// toScene returns vertices in the "scene" coordinates. Currently polygons are
// defined relatively to their center.
func toScene(x, y, phi float64, v []vector.Vector) []vector.Vector {
	if len(v) == 0 {
		return nil
	}
	c := vector.New(x, y)
	vs := make([]vector.Vector, 0, len(v))
	for _, v := range v {
		vs = append(vs, v.Rotate(phi).Add(c))
	}
	return vs
}
