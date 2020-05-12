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

var circles []*circle

type circle struct {
	intersect.C
	xInfo func(intersect.I)
}

func (c *circle) draw(img *canvas.Image) {
	img.Clear()
	img.Fill(color.RGBA{0xff, 0xff, 0xff, 0xa0})
	w, h := img.Size()
	var col color.Color = color.White
	for _, other := range circles {
		if other == c {
			continue
		}
		if xi, ok := intersect.Circles(c.C, other.C); ok {
			c.xInfo(xi)
			col = color.RGBA{0xff, 0, 0, 0xff}
		}
	}
	shapes.DrawCircle(img.Image, w/2, h/2, int(c.R), col)
}

func addCircle(app *application.App, x, y, r int, hud *hud) {
	xc := circle{
		C:     intersect.C{X: float64(x), Y: float64(y), R: float64(r)},
		xInfo: hud.crossInfo,
	}
	c := canvas.New(xc.draw)
	c.SetBounds(image.Rect(x-r, y-r, x+r, y+r))
	d := drag.EnableFor(c)
	d.AddDragListener(func() {
		hud.shapeInfo(c.Bounds())
		b := c.Bounds()
		xc.MoveTo(
			math.Round(float64(b.Min.X+b.Max.X)/2),
			math.Round(float64(b.Min.Y+b.Max.Y)/2),
		)
	})
	circles = append(circles, &xc)
	app.AddComponent(d)
}
