package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

type Vec struct {
	X, Y, Z float64
}

func Add(a, b Vec) Vec {
	return Vec{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

func Sub(a, b Vec) Vec {
	return Vec{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

func Divide(v Vec, a float64) Vec {
	return Vec{v.X / a, v.Y / a, v.Z / a}
}

func Multiply(v Vec, a float64) Vec {
	return Vec{v.X * a, v.Y * a, v.Z * a}
}

func Mod(a Vec) float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

func Normalize(v Vec) Vec {
	return Vec{
		v.X / Mod(v),
		v.Y / Mod(v),
		v.Z / Mod(v),
	}
}

func Cross(a, b Vec) Vec {
	return Vec{
		a.Y*b.Z - b.Y*a.Z,
		a.Z*b.X - b.Z*a.X,
		a.X*b.Y - b.X*a.Y,
	}
}

func Dot(a, b Vec) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

type Rotator struct {
	X, Y, Z float64
}

func (v *Vec) RotateZ(rad float64) {
	x := v.X*math.Cos(rad) - v.Y*math.Sin(rad)
	y := v.X*math.Sin(rad) + v.Y*math.Cos(rad)
	v.X, v.Y = x, y
}

func (v *Vec) RotateY(rad float64) {
	x := v.X*math.Cos(rad) + v.Z*math.Sin(rad)
	z := v.X*-math.Sin(rad) + v.Z*math.Cos(rad)
	v.X, v.Z = x, z
}

func (v *Vec) RotateX(rad float64) {
	y := v.Y*math.Cos(rad) - v.Z*math.Sin(rad)
	z := v.Y*math.Sin(rad) + v.Z*math.Cos(rad)
	v.Y, v.Z = y, z
}

func (v *Vec) Rotate(rad Rotator) {
	v.RotateZ(rad.Z)
	v.RotateY(rad.Y)
	v.RotateX(rad.X)
}

func CentralProjection(v Vec, k float64) Vec {
	return Vec{
		-(v.X / v.Z) * k,
		-(v.Y / v.Z) * k,
		-1,
	}
}

func DrawLine(screen *ebiten.Image, a, b Vec, clr color.Color) {
	halfWidth, halfHeight := float64(screenWidth/2), float64(screenHeight/2)
	k := float64(250)
	a = CentralProjection(a, k)
	b = CentralProjection(b, k)
	ebitenutil.DrawLine(screen, a.X+halfWidth, -a.Y+halfHeight, b.X+halfWidth, -b.Y+halfHeight, clr)
}

type Rect struct {
	A, B, C, D Vec
}

func (r *Rect) Draw(screen *ebiten.Image, clr color.Color) {
	DrawLine(screen, r.A, r.B, clr)
	DrawLine(screen, r.B, r.C, clr)
	DrawLine(screen, r.C, r.D, clr)
	DrawLine(screen, r.D, r.A, clr)
}

type Cube struct {
	p [8]Vec
}

func (c *Cube) Rotate(screen *ebiten.Image, r Rotator) {
	ctr := Add(Divide(Sub(c.p[6], c.p[0]), 2), c.p[0])
	for i := range c.p {
		c.p[i] = Sub(c.p[i], ctr)
		c.p[i].RotateX(math.Pi / 300) //r)
		c.p[i].RotateY(math.Pi / 500) //r)
		c.p[i].RotateZ(math.Pi / 700) //r)
		c.p[i] = Add(c.p[i], ctr)
	}
}

func (c *Cube) Draw(screen *ebiten.Image, clr color.Color) {

	// da, db stands for diagonal a and b(diagonal starting and ending points)
	DrawNormal := func(da, db, v, w Vec, col color.Color) {
		return
		ctr := Add(Divide(Sub(da, db), 2), db)
		DrawLine(screen, ctr, Add(Multiply(Normalize(Cross(v, w)), 200), ctr), col)
	}

	for _, f := range [][4]int{
		{0, 1, 2, 3},
		{7, 6, 5, 4},
		{0, 4, 5, 1},
		{1, 5, 6, 2},
		{3, 2, 6, 7},
		{4, 0, 3, 7},
	} {
		// Near plane
		DrawNormal(c.p[f[2]], c.p[f[0]], Sub(c.p[f[1]], c.p[f[2]]), Sub(c.p[f[2]], c.p[f[3]]), color.RGBA{255, 0, 0, 255})
		cr := Cross(Sub(c.p[f[1]], c.p[f[2]]), Sub(c.p[f[2]], c.p[f[3]]))
		center := Multiply(Add(c.p[f[0]], c.p[f[2]]), 0.5)
		if Dot(center, cr) < 0 {
			for i := range f {
				j := (i + 1) % len(f)
				DrawLine(screen, c.p[f[i]], c.p[f[j]], clr)
				DrawLine(screen, c.p[f[i]], c.p[f[j]], clr)
				DrawLine(screen, c.p[f[i]], c.p[f[j]], clr)
				DrawLine(screen, c.p[f[i]], c.p[f[j]], clr)
			}
		}
	}
	// // Far plane
	// DrawNormal(c.p[5], c.p[7], Sub(c.p[6], c.p[5]), Sub(c.p[5], c.p[4]), color.RGBA{0, 255, 0, 255})
	// if cr := Cross(Sub(c.p[6], c.p[5]), Sub(c.p[5], c.p[4])); Dot(Vec{0, 0, 1}, cr) < 0 {
	// 	DrawLine(screen, c.p[4], c.p[5], clr)
	// 	DrawLine(screen, c.p[5], c.p[6], clr)
	// 	DrawLine(screen, c.p[6], c.p[7], clr)
	// 	DrawLine(screen, c.p[7], c.p[4], clr)
	// }

	// //Left plane
	// DrawNormal(c.p[0], c.p[5], Sub(c.p[5], c.p[1]), Sub(c.p[1], c.p[0]), color.RGBA{0, 255, 255, 255})
	// if cr := Cross(Sub(c.p[5], c.p[1]), Sub(c.p[1], c.p[0])); Dot(Vec{0, 0, 1}, cr) < 0 {
	// 	DrawLine(screen, c.p[4], c.p[5], clr)
	// 	DrawLine(screen, c.p[5], c.p[1], clr)
	// 	DrawLine(screen, c.p[1], c.p[0], clr)
	// 	DrawLine(screen, c.p[0], c.p[4], clr)
	// }

	// // Top plane
	// DrawNormal(c.p[6], c.p[1], Sub(c.p[5], c.p[6]), Sub(c.p[6], c.p[2]), color.RGBA{255, 0, 255, 255})
	// if cr := Cross(Sub(c.p[5], c.p[6]), Sub(c.p[6], c.p[2])); Dot(Vec{0, 0, 1}, cr) < 0 {
	// 	DrawLine(screen, c.p[1], c.p[5], clr)
	// 	DrawLine(screen, c.p[5], c.p[6], clr)
	// 	DrawLine(screen, c.p[6], c.p[2], clr)
	// 	DrawLine(screen, c.p[2], c.p[1], clr)
	// }

	// // Right plane
	// DrawNormal(c.p[6], c.p[3], Sub(c.p[2], c.p[6]), Sub(c.p[6], c.p[7]), color.White)
	// if cr := Cross(Sub(c.p[2], c.p[6]), Sub(c.p[6], c.p[7])); Dot(Vec{0, 0, 1}, cr) < 0 {
	// 	DrawLine(screen, c.p[3], c.p[2], clr)
	// 	DrawLine(screen, c.p[2], c.p[6], clr)
	// 	DrawLine(screen, c.p[6], c.p[7], clr)
	// 	DrawLine(screen, c.p[7], c.p[3], clr)
	// }

	// // Bottom plane
	// DrawNormal(c.p[0], c.p[7], Sub(c.p[0], c.p[3]), Sub(c.p[3], c.p[7]), color.RGBA{255, 165, 0, 255})
	// if cr := Cross(Sub(c.p[0], c.p[3]), Sub(c.p[3], c.p[7])); Dot(Vec{0, 0, 1}, cr) < 0 {
	// 	DrawLine(screen, c.p[4], c.p[0], clr)
	// 	DrawLine(screen, c.p[0], c.p[3], clr)
	// 	DrawLine(screen, c.p[3], c.p[7], clr)
	// 	DrawLine(screen, c.p[4], c.p[4], clr)
	// }
}

type game struct {
	c            []Cube
	screenBuffer *ebiten.Image
}

func NewGame() *game {
	return &game{
		[]Cube{
			{
				[8]Vec{
					{-200, -200, 400}, // NearBottomLeft
					{-200, 200, 400},  // NearTopLeft
					{200, 200, 400},   // NearTopRight
					{200, -200, 400},  // NearBottomRight

					{-200, -200, 800}, // FarBottomLeft
					{-200, 200, 800},  // FarTopLeft
					{200, 200, 800},   // FarTopRight
					{200, -200, 800},  // FarBottomRight
				},
			},
		},
		ebiten.NewImage(screenWidth, screenHeight),
	}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }
func (g *game) Update() error {
	for i := range g.c {
		g.c[i].Rotate(g.screenBuffer, Rotator{math.Pi / 180 / 2, math.Pi / 180 / 2, math.Pi / 180 / 2})
	}
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	for i := range g.c {
		g.c[i].Draw(screen, color.RGBA{255, 102, 204, 255})
	}
	screen.DrawImage(g.screenBuffer, &ebiten.DrawImageOptions{})
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	g := NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
