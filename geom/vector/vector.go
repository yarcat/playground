// Package vector provides vectors and varios operations on vector e.g. dot and
// cross products.
package vector

// Vector represents a geometric object that has magnitude and direction.
type Vector struct {
	X, Y float64
}

// New returns new vector pointing to the provided point.
func New(X, Y float64) Vector {
	return Vector{X, Y}
}
