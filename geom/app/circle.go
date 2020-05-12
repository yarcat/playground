package main

import (
	"image"
	"image/color"
	"math"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/intersect"
	"github.com/yarcat/playground/geom/shapes"
)

var circles = make(map[component.Component]*intersect.C)

func addCircle(app *application.App, x, y, r int, updateStatus func(image.Rectangle)) {
	xc := &intersect.C{X: float64(x), Y: float64(y), R: float64(r)}
	var c *canvas.Canvas
	c = canvas.New(func(img *canvas.Image) {
		img.Clear()
		img.Fill(color.RGBA{0xff, 0xff, 0xff, 0xa0})
		w, h := img.Size()
		var col color.Color = color.White
		for otherc, otherxc := range circles {
			if otherc == c {
				continue
			}
			if _, ok := intersect.Circles(*xc, *otherxc); ok {
				col = color.RGBA{0xff, 0, 0, 0xff}
			}
		}
		shapes.DrawCircle(img.Image, w/2, h/2, r, col)
	})
	c.SetBounds(image.Rect(x-r, y-r, x+r, y+r))
	d := drag.EnableFor(c)
	d.AddDragListener(func() {
		updateStatus(c.Bounds())
		b := c.Bounds()
		xc.MoveTo(
			math.Round(float64(b.Min.X+b.Max.X)/2),
			math.Round(float64(b.Min.Y+b.Max.Y)/2),
		)
	})
	circles[c] = xc
	app.AddComponent(d)
}
