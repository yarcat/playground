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
	"github.com/yarcat/playground/geom/shapes"
)

type rect struct {
	intersect.R
	other *rect
	xi    intersect.I
	xInfo func(intersect.I)
}

func (r *rect) intersected(other *rect, xi intersect.I) {
	r.other = other
	r.xi = xi
}

func (r rect) hasIntersection() bool {
	return r.other != nil
}

func (r *rect) draw(img *canvas.Image) {
	img.Clear()
	img.Fill((color.RGBA{0xf0, 0xf0, 0xf0, 0xa0}))
	w, h := img.Size()
	if r.hasIntersection() {
		r.xInfo(r.xi)
		shapes.DrawRectangle(img.Image, 0, 0, w, h, color.RGBA{0xff, 0, 0, 0xff})
		r.drawX(img.Image)
	}
	shapes.DrawRectangle(img.Image, 0, 0, w, h, color.White)
	img.Op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	img.Op.GeoM.Rotate(math.Pi / 100)
	img.Op.GeoM.Translate(float64(w)/2, float64(h)/2)
}

func (r rect) drawX(img *ebiten.Image) {
	var x1, y1, x2, y2 float64
	if r.xi.N.X != 0 {
		dy := (r.other.H+r.H)/2 - r.other.Y + r.Y
		y1 = r.H - dy/2
		y2 = y1
		if r.xi.N.X > 0 {
			x1 = r.W
			x2 = r.W - r.xi.P
		} else {
			x1 = 0
			x2 = r.xi.P
		}
	} else {
		dx := (r.other.W+r.W)/2 - r.other.X + r.X
		x1 = r.W - dx/2
		x2 = x1
		if r.xi.N.Y > 0 {
			y1 = r.H
			y2 = r.H - r.xi.P
		} else {
			y1 = 0
			y2 = r.xi.P
		}
	}
	ebitenutil.DrawLine(img, x1, y1, x2, y2, color.RGBA{0xff, 0, 0, 0xff})
}

func addRectangle(app *application.App, x, y, w, h int, hud *hud, is *intersector) {
	xr := &rect{
		R:     intersect.R{X: float64(x), Y: float64(y), W: float64(w), H: float64(h)},
		xInfo: hud.crossInfo,
	}
	r := canvas.New(xr.draw)
	r.SetBounds(image.Rect(x-w/2, y-h/2, x+w/2, y+h/2))
	d := drag.EnableFor(r)
	d.AddDragListener(func() {
		b := r.Bounds()
		hud.shapeInfo(b)
		is.computeR(xr)
		xr.MoveTo(
			math.Round(float64(b.Min.X+b.Max.X)/2),
			math.Round(float64(b.Min.Y+b.Max.Y)/2),
		)
	})
	is.addR(xr)
	app.AddComponent(d)
}
