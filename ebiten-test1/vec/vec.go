package vec

import "math"

type Vec struct {
	X, Y float64
}

func (a Vec) Add(b Vec) Vec       { return Vec{X: a.X + b.X, Y: a.Y + b.Y} }
func (a Vec) Sub(b Vec) Vec       { return Vec{X: a.X - b.X, Y: a.Y - b.Y} }
func (a Vec) Scale(k float64) Vec { return Vec{X: a.X * k, Y: a.Y * k} }

func (a Vec) Rotate(phi float64) Vec {
	return Vec{
		X: a.X*math.Cos(phi) - a.Y*math.Sin(phi),
		Y: a.X*math.Sin(phi) + a.Y*math.Cos(phi),
	}
}

func (a Vec) Len() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

func (a Vec) Norm() Vec {
	l := a.Len()
	return Vec{X: a.X / l, Y: a.Y / l}
}
