package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/intersect"
	"github.com/yarcat/playground/geom/vector"
)

type poly struct {
	intersect.P
	other *poly
	xi    intersect.I
	xInfo func(intersect.I)
	img   *ebiten.Image
	v     []vector.Vector
}

func (p *poly) intersected(other *poly, xi intersect.I) {
	p.other = other
	p.xi = xi
}

func (p poly) hasIntersection() bool {
	return p.other != nil
}

func (p *poly) draw(img *canvas.Image) {
	img.Clear()
	img.Fill((color.RGBA{0xf0, 0xf0, 0xf0, 0xa0}))
	var c color.Color = color.White
	if p.hasIntersection() {
		p.xInfo(p.xi)
		c = color.RGBA{255, 0, 0, 255}
		p.drawX(img.Image)
	}
	w, h := img.Size()
	o := vector.New(float64(w)/2, float64(h)/2)
	for i, v := range p.V {
		p.v[i] = v.Rotate(p.Phi).Add(o)
	}
	for _, e := range p.E {
		ebitenutil.DrawLine(
			img.Image,
			p.v[e[0]].X, p.v[e[0]].Y,
			p.v[e[1]].X, p.v[e[1]].Y,
			c)

	}
}

func (p poly) drawX(img *ebiten.Image) {
}

func addPolygon(
	app *application.App,
	x, y, angle float64,
	v []vector.Vector,
	e [][2]int,
	hud *hud, is *intersector) {
	xp := &poly{
		P: intersect.P{
			X:   x,
			Y:   y,
			Phi: angle,
			E:   e,
			V:   v,
		},
		xInfo: hud.crossInfo,
		v:     make([]vector.Vector, len(v)),
	}
	p := canvas.New(xp.draw)
	// TODO(yarcat): Compute these based on the polygon vertices.
	w, h := 100, 100
	p.SetBounds(image.Rect(
		int(x)-w/2, int(y)-h/2,
		int(x)+w/2, int(y)+h/2),
	)
	d := drag.EnableFor(p)
	d.AddDragListener(func() {
		b := p.Bounds()
		hud.shapeInfo(b)
		is.computeP(xp)
		xp.MoveTo(
			math.Round(float64(b.Min.X+b.Max.X)/2),
			math.Round(float64(b.Min.Y+b.Max.Y)/2),
		)
	})
	is.addP(xp)
	app.AddComponent(d)
}

func newEquilateralTriangle(a float64) (v []vector.Vector, e [][2]int) {
	p1 := vector.Vector{}
	p2 := vector.New(a, 0)
	p3 := p2.Rotate(math.Pi / 3)
	o := vector.New(a/2, math.Sqrt(0.75)/3.0*a)
	v = []vector.Vector{p1.Sub(o), p2.Sub(o), p3.Sub(o)}
	e = [][2]int{{0, 1}, {1, 2}, {2, 0}}
	return
}
