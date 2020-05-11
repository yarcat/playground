// Package vector provides vectors and varios operations on vector e.g. dot and
// cross products.
package vector

import "math"

// Vector represents a geometric object that has magnitude and direction.
type Vector struct {
	X, Y float64
}

// New returns new vector pointing to the provided point.
func New(X, Y float64) Vector {
	return Vector{X, Y}
}

// Rotate returns new vector rotated by the specified rad.
func (v Vector) Rotate(a float64) Vector {
	cos, sin := math.Cos(a), math.Sin(a)
	return New(v.X*cos-v.Y*sin, v.X*sin+v.Y*cos)
}

// Add adds a to v.
func (v Vector) Add(a Vector) Vector {
	return New(v.X+a.X, v.Y+a.Y)
}

// Sub subtracts a from v.
func (v Vector) Sub(a Vector) Vector {
	return New(v.X-a.X, v.Y-a.Y)
}
