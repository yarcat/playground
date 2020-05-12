package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/intersect"
	"github.com/yarcat/playground/geom/shapes"
	"github.com/yarcat/playground/geom/vector"
)

func addRectangle(app *application.App, updateStatus func(image.Rectangle)) {
	r := canvas.New(func(img *canvas.Image) {
		if img.Invalidated() {
			w, h := img.Size()
			shapes.DrawRectangle(img.Image, w, h, color.White)
		}
	})
	r.SetBounds(image.Rect(200, 150, 300, 250))
	d := drag.EnableFor(r)
	d.AddDragListener(func() { updateStatus(r.Bounds()) })
	app.AddComponent(d)
}

var circles = make(map[component.Component]*intersect.C)

func addCircle(app *application.App, x, y, r int, updateStatus func(image.Rectangle)) {
	xc := intersect.C{X: float64(x), Y: float64(y), R: float64(r)}
	var c *canvas.Canvas
	c = canvas.New(func(img *canvas.Image) {
		img.Clear()
		w, h := img.Size()
		var col color.Color = color.White
		for otherc, otherxc := range circles {
			if otherc == c {
				continue
			}
			if _, ok := intersect.Circles(xc, *otherxc); ok {
				col = color.RGBA{0xff, 0, 0, 0xff}
			}
		}
		shapes.DrawCircle(img.Image, w/2, h/2, r, col)
	})
	c.SetBounds(image.Rect(x-r, y-r, x+r, y+r))
	d := drag.EnableFor(c)
	d.AddDragListener(func() {
		updateStatus(c.Bounds())

	})
	circles[c] = &xc
	app.AddComponent(d)
}

func addTriangle(app *application.App, updateStatus func(image.Rectangle)) {
	t := canvas.New(func(img *canvas.Image) {
		if img.Invalidated() {
			w, h := img.Size()
			x0, y0 := float64(w)/2, float64(h)/2
			side := float64(w)
			if float64(h) < side {
				side = float64(h)
			}
			p0 := vector.New(-side/2, -y0+math.Sqrt(0.75*side*side)/6)
			p := vector.New(side, 0)
			v := []vector.Vector{
				p0,
				p.Add(p0),
				p.Rotate(math.Pi / 3).Add(p0),
			}
			edges := [...][2]int{{0, 1}, {1, 2}, {2, 0}}
			for _, e := range edges {
				ebitenutil.DrawLine(
					img.Image,
					x0+v[e[0]].X, y0+v[e[0]].Y,
					x0+v[e[1]].X, y0+v[e[1]].Y,
					color.White,
				)
			}
		}
	})
	t.SetBounds(image.Rect(350, 150, 450, 250))
	d := drag.EnableFor(t)
	d.AddDragListener(func() { updateStatus(t.Bounds()) })
	app.AddComponent(d)
}
