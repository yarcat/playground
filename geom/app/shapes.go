package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/vector"
)

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
