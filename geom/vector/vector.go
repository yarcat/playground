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

// Dot returns a dot product of two vectors.
func (v Vector) Dot(a Vector) float64 {
	return v.X*a.X + v.Y*a.Y
}

// Norm returns a normalized vector in a given direction.
func (v Vector) Norm() Vector {
	len2 := v.Len2()
	if len2 == 0 {
		return ZN
	}
	return v.Scale(1 / math.Sqrt(len2))
}

// Len2 returns a squared length.
func (v Vector) Len2() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Len returns vector length.
func (v Vector) Len() float64 {
	return math.Sqrt(v.Len2())
}
