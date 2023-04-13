package raycasting

import (
	"math"
	"test1/vec"
)

type world interface {
	EmptyAt(row, col int) bool
}

func ddaParams(p, v vec.Vec) (dirDist, dirLen vec.Vec, dx, dy int) {
	dirLen.X = math.Abs(1 / v.X)
	if v.X < 0 {
		dx = -1
		dirDist.X = (p.X - math.Trunc(p.X)) * dirLen.X
	} else {
		dx = 1
		dirDist.X = (math.Trunc(p.X+1) - p.X) * dirLen.X
	}
	dirLen.Y = math.Abs(1 / v.Y)
	if v.Y < 0 {
		dy = -1
		dirDist.Y = (p.Y - math.Trunc(p.Y)) * dirLen.Y
	} else {
		dy = 1
		dirDist.Y = (math.Trunc(p.Y+1) - p.Y) * dirLen.Y
	}
	return dirDist, dirLen, dx, dy
}

func Solve(w world, p, v vec.Vec) (dist float64, vertSide bool) {
	dirDist, dirLen, dcol, drow := ddaParams(p, v)
	for col, row := int(p.X), int(p.Y); ; {
		if dirDist.X < dirDist.Y {
			dist = dirDist.X
			dirDist.X += dirLen.X
			col += dcol
			vertSide = false
		} else {
			dist = dirDist.Y
			dirDist.Y += dirLen.Y
			row += drow
			vertSide = true
		}
		if !w.EmptyAt(row, col) {
			return dist, vertSide
		}
	}
}

type RayFunc func(v1 vec.Vec, screenX int, dist float64, vertSide bool)

func FOV(w world, p, v vec.Vec, screenWidth int, f RayFunc) {
	for x := 0; x < screenWidth; x++ {
		projX := float64(2*x)/float64(screenWidth) - 1
		v1 := v.Add(vec.Vec{X: v.Y * projX, Y: -v.X * projX})
		dist, vertSide := Solve(w, p, v1)
		f(v1, x, dist, vertSide)
	}
}
