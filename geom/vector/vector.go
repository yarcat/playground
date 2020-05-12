// Package vector provides vectors and varios operations on vector e.g. dot and
// cross products.
package vector

import "math"

var (
	// ZN is a "zero-normal" vector. This vector should be used while trying to
	// produce a normal vector of length 0. It doesn't matter too much what
	// values to use here - it's just important to be consistent.
	ZN = Vector{X: 0, Y: 1}
)

// Vector represents a geometric object that has magnitude and direction.
type Vector struct {
	X, Y float64
}

// New returns new vector pointing to the provided point.
func New(x, y float64) Vector {
	return Vector{X: x, Y: y}
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

// Scale returns new vector with X and Y scaled equally.
func (v Vector) Scale(k float64) Vector {
	return New(v.X*k, v.Y*k)
}
