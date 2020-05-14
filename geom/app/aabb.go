package main

import (
	"image"
	"image/color"
	"math"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/intersect"
	"github.com/yarcat/playground/geom/shapes"
)

type aabb struct {
	intersect.R
	other *aabb
	xi    intersect.I
	xInfo func(intersect.I)
}

func (r *aabb) intersected(other *aabb, xi intersect.I) {
	r.other = other
	r.xi = xi
}

func (r aabb) hasIntersection() bool {
	return r.other != nil
}

func (r aabb) draw(img *canvas.Image) {
	img.Clear()
	img.Fill((color.RGBA{0xf0, 0xf0, 0xf0, 0xa0}))
	w, h := img.Size()
	if r.hasIntersection() {
		shapes.DrawRectangle(img.Image, 0, 0, w, h, color.RGBA{0xff, 0, 0, 0xff})
	} else {
		shapes.DrawRectangle(img.Image, 0, 0, w, h, color.White)
	}
}

func addRectangle(app *application.App, x, y, w, h int, hud *hud, is *intersector) {
	xr := &aabb{
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
