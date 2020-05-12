package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/shapes"
	"github.com/yarcat/playground/geom/vector"
)

func main() {
	const screenWidth, screenHeight = 800, 600
	app := application.New(screenWidth, screenHeight)
	updateStatus := newHUD(app)

	addButtons(app)

	var r, c, t *canvas.Canvas

	c = canvas.New(func(img *canvas.Image) {
		if img.Invalidated() {
			w, h := img.Size()
			d := w
			if h < d {
				d = h
			}
			shapes.DrawCircle(img.Image, w/2, h/2, d/2, color.White)
		}
	})
	c.SetBounds(image.Rect(50, 150, 150, 250))
	d := drag.EnableFor(c)
	d.AddDragListener(func() { updateStatus(c.Bounds()) })
	app.AddComponent(d)

	r = canvas.New(func(img *canvas.Image) {
		if img.Invalidated() {
			w, h := img.Size()
			shapes.DrawRectangle(img.Image, w, h, color.White)
		}
	})
	r.SetBounds(image.Rect(200, 150, 300, 250))
	d = drag.EnableFor(r)
	d.AddDragListener(func() { updateStatus(r.Bounds()) })
	app.AddComponent(d)

	t = canvas.New(func(img *canvas.Image) {
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
	d = drag.EnableFor(t)
	d.AddDragListener(func() { updateStatus(t.Bounds()) })
	app.AddComponent(d)

	if err := application.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
