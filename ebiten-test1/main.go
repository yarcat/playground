package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Point struct {
	X, Y, Z float64
}

func (p Point) Add(a Point) Point { return Point{X: p.X + a.X, Y: p.Y + a.Y, Z: p.Z + a.Z} }
func (p Point) Sub(a Point) Point { return Point{X: p.X - a.X, Y: p.Y - a.Y, Z: p.Z - a.Z} }

func (p Point) CrossProduct(a Point) Point {
	return Point{
		X: p.Y*a.Z - p.Z*a.Y,
		Y: -p.X*a.Z + p.Z*a.X,
		Z: p.X*a.Y - p.Y*a.X,
	}
}

func (p Point) DotProduct(a Point) float64 {
	return p.X*a.X + p.Y*a.Y + p.Z*a.Z
}

func (p Point) Scale(k float64) Point {
	return Point{X: p.X * k, Y: p.Y * k, Z: p.Z * k}
}

func (p Point) Normalize() Point {
	l := p.Len()
	return Point{X: p.X / l, Y: p.Y / l, Z: p.Z / l}
}

func (p Point) Len() float64 { return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z) }

func RotateX(phi float64, pts []Point) {
	for i := range pts {
		y := pts[i].Y*math.Cos(phi) + pts[i].Z*math.Sin(phi)
		z := -pts[i].Y*math.Sin(phi) + pts[i].Z*math.Cos(phi)
		pts[i].Y = y
		pts[i].Z = z
	}
}

func RotateY(phi float64, pts []Point) {
	for i := range pts {
		x := pts[i].X*math.Cos(phi) - pts[i].Z*math.Sin(phi)
		z := pts[i].X*math.Sin(phi) + pts[i].Z*math.Cos(phi)
		pts[i].X = x
		pts[i].Z = z
	}
}

func RotateZ(phi float64, pts []Point) {
	for i := range pts {
		x := pts[i].X*math.Cos(phi) - pts[i].Y*math.Sin(phi)
		y := pts[i].X*math.Sin(phi) + pts[i].Y*math.Cos(phi)
		pts[i].X = x
		pts[i].Y = y
	}
}

func DrawCircle(img *ebiten.Image, x, y, r int, c color.Color) {
	img.Set(x, y+r, c)
	img.Set(x, y-r, c)
	img.Set(x+r, y, c)
	img.Set(x-r, y, c)
	xi, yi, d := 0, r, 3-2*r
	for xi <= yi {
		if d > 0 {
			d, yi = d+4*(xi-yi)+10, yi-1
		} else {
			d += 4*xi + 6
		}
		xi++
		img.Set(x+xi, y+yi, c)
		img.Set(x-xi, y+yi, c)
		img.Set(x+xi, y-yi, c)
		img.Set(x-xi, y-yi, c)
		img.Set(x+yi, y+xi, c)
		img.Set(x-yi, y+xi, c)
		img.Set(x+yi, y-xi, c)
		img.Set(x-yi, y-xi, c)
	}
}

type Game struct {
	width, height int
	verts         []Point
	facets        [][]int
	//edges         [][2]int
}

var white = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func DrawLine(img *ebiten.Image, k, dx, dy float64, a, b Point, c color.Color) {
	x1, y1 := k*a.X/a.Z, k*a.Y/a.Z
	x2, y2 := k*b.X/b.Z, k*b.Y/b.Z
	// k = 1
	// x1, y1 := a.X, a.Y
	// x2, y2 := b.X, b.Y
	ebitenutil.DrawLine(img, x1+dx, y1+dy, x2+dx, y2+dy, c)
}

var r, dr = 0, 1

func (g *Game) Draw(screen *ebiten.Image) {
	maxR := min(g.height, g.width)/2 - 1
	if r == 0 {
		r, dr = maxR, -1
	} else {
		r += dr
		if r <= maxR/4 {
			r, dr = maxR/4, 1
		} else if r >= maxR {
			r, dr = maxR, -1
		}
	}
	// DrawCircle(screen, g.width/2, g.height/2, r, white)

	// for _, e := range g.edges {

	// 	const r, k = 500, 500
	// 	z1, z2 := g.verts[e[0]].Z+r, g.verts[e[1]].Z+r
	// 	x1, y1 := g.verts[e[0]].X/z1, g.verts[e[0]].Y/z1
	// 	x2, y2 := g.verts[e[1]].X/z2, g.verts[e[1]].Y/z2

	// 	ebitenutil.DrawLine(screen,
	// 		x1*k+float64(g.width)/2, y1*k+float64(g.height)/2,
	// 		x2*k+float64(g.width)/2, y2*k+float64(g.height)/2,
	// 		color.White,
	// 	)
	// }
	r := Point{Z: 300}
	for _, facet := range g.facets {

		a := g.verts[facet[2]].Sub(g.verts[facet[1]])
		b := g.verts[facet[0]].Sub(g.verts[facet[1]])
		n := a.CrossProduct(b)

		// c := g.verts[facet[0]].Add(g.verts[facet[2]]).Scale(0.5)
		c := g.verts[facet[0]]
		// c := Point{0, 0, 1}
		col := color.RGBA{255, 255, 255, 255}

		if c.Add(r).DotProduct(n) >= 0 {
			continue
			// col = color.RGBA{50, 50, 50, 255}
		}

		dx, dy := float64(g.width/2), float64(g.height/2)
		const k = 500
		// DrawLine(screen, k, dx, dy, g.verts[facet[0]].Add(r), g.verts[facet[2]].Add(r))
		// DrawLine(screen, k, dx, dy, g.verts[facet[1]].Add(r), g.verts[facet[3]].Add(r))
		// DrawLine(screen, k, dx, dy, c.Add(r), c.Add(n.Scale(25)).Add(r))
		for n, i := range facet {
			j := facet[(n+1)%len(facet)]
			DrawLine(screen, k, dx, dy, g.verts[i].Add(r), g.verts[j].Add(r), col)
		}
	}

}

func (g *Game) Layout(outWidth, outHeight int) (screenWidth int, screenHeight int) {
	return g.width, g.height
}

var dxfi = rand.Float64()*math.Pi/5e4 + math.Pi/1e5
var dyfi = rand.Float64()*math.Pi/5e4 + math.Pi/1e5
var dzfi = rand.Float64()*math.Pi/5e4 + math.Pi/1e5
var xfi, yfi, zfi float64

func (g *Game) Update() error {
	xfi += dxfi
	if math.Abs(xfi) > math.Pi/100 {
		dxfi = -dxfi
	}
	yfi += dyfi
	if math.Abs(yfi) > math.Pi/100 {
		dyfi = -dyfi
	}
	zfi += dzfi
	if math.Abs(zfi) > math.Pi/100 {
		dzfi = -dzfi
	}
	RotateZ(zfi, g.verts)
	RotateX(xfi, g.verts)
	RotateY(yfi, g.verts)
	return nil
}

func main() {
	/*
		verts := []Point{
			{100, 100, 100},
			{-100, 100, 100},
			{-100, -100, 100},
			{100, -100, 100},
			{100, 100, -100},
			{-100, 100, -100},
			{-100, -100, -100},
			{100, -100, -100},
		}
		// edges := [][2]int{
		// 	{0, 1}, {1, 2}, {2, 3}, {3, 0},
		// 	{4, 5}, {5, 6}, {6, 7}, {7, 4},
		// 	{0, 4}, {1, 5}, {2, 6}, {3, 7},
		// }
		facets := [][]int{
			{0, 1, 2, 3},
			{7, 6, 5, 4},
			{2, 1, 5, 6},
			{3, 7, 4, 0},
			{4, 5, 1, 0},
			{6, 7, 3, 2},
		}
	*/
	verts := make([]Point, 10)
	verts[0] = Point{X: 100, Z: 50}
	{
		const fi = 2 * math.Pi / 10
		cosfi, sinfi := math.Cos(fi), math.Sin(fi)
		for i := 1; i < len(verts); i++ {
			verts[i] = Point{
				X: verts[i-1].X*cosfi - verts[i-1].Y*sinfi,
				Y: verts[i-1].X*sinfi + verts[i-1].Y*cosfi,
				Z: -verts[i-1].Z,
			}
		}
	}
	facets := make([][]int, 10, 12)
	for i := range verts {
		if i%2 == 0 {
			facets[i] = []int{i, (i + 1) % 10, (i + 2) % 10}
		} else {
			facets[i] = []int{(i + 2) % 10, (i + 1) % 10, i}
		}
	}
	verts = append(verts,
		Point{Z: (50*math.Sqrt(5)-1)/2 + 50},
		Point{Z: -(50*math.Sqrt(5)-1)/2 - 50},
	)
	for i := 0; i < 10; i += 2 {
		facets = append(facets, []int{10, i, (i + 2) % 10})
		facets = append(facets, []int{(i + 3) % 10, i + 1, 11})
		//facets = append(facets, []int{(i + 2 - 1) % 10, (i + 1) % 10, 11})
	}

	const width, height = 640, 480
	ebiten.SetWindowSize(width, height)
	g := &Game{width: width, height: height /*edges: edges,*/, facets: facets, verts: verts}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
