package object

import (
	"math"
	"test1/vec"
	"time"
)

type world interface {
	EmptyAt(row, col int) bool
}

func handleDir(p, v vec.Vec, w world, dtSec float64, fwd, bck bool) vec.Vec {
	if fwd && bck || !fwd && !bck {
		return p
	}
	dt := 0.0
	if fwd {
		dt += dtSec
	}
	if bck {
		dt -= dtSec
	}
	x, y := p.X, p.Y
	if w.EmptyAt(int(y+v.Y*dt), int(x)) {
		y += v.Y * dt
	}
	if w.EmptyAt(int(y), int(x+v.X*dt)) {
		x += v.X * dt
	}
	return vec.Vec{X: x, Y: y}
}

func handleRot(v vec.Vec, dtSec float64, left, right bool) vec.Vec {
	dt := 0.0
	if left {
		dt += dtSec
	}
	if right {
		dt -= dtSec
	}
	return v.Rotate(math.Pi * dt) // TODO: Do not hardcode the angle.
}

func HandleInput(o Object, w world, d time.Duration, fwd, bck, left, right bool) Object {
	return Object{
		pos: handleDir(o.pos, o.v, w, d.Seconds(), fwd, bck),
		v:   handleRot(o.v, d.Seconds(), left, right),
	}
}
