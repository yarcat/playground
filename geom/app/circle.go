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
	"github.com/yarcat/playground/geom/vector"
)

type circle struct {
	intersect.C
	other *circle
	xi    intersect.I
	xInfo func(intersect.I)
}

func (c *circle) intersected(other *circle, xi intersect.I) {
	c.other = other
	c.xi = xi
}

func (c circle) hasIntersection() bool {
	return c.other != nil
}

func (c *circle) draw(img *canvas.Image) {
	img.Clear()
	img.Fill(color.RGBA{0xff, 0xff, 0xff, 0xa0})
	w, h := img.Size()
	if c.hasIntersection() {
		center := vector.New(float64(w/2), float64(h/2))
		c.xInfo(c.xi)
		// TODO(yarcat): Remove duplicates.
		c.drawX(img.Image, center)
	}
	shapes.DrawCircle(img.Image, w/2, h/2, int(c.R), color.White)
}

func (c circle) drawX(img *ebiten.Image, center vector.Vector) {
	p1 := center.Add(c.xi.N.Scale(c.R))
	p2 := p1.Sub(c.xi.N.Scale(c.xi.P))
	ebitenutil.DrawLine(
		img,
		p1.X, p1.Y,
		p2.X, p2.Y,
		color.RGBA{0xff, 0, 0, 0xff},
	)
}

func addCircle(app *application.App, x, y, r int, hud *hud, is *intersector) {
	xc := &circle{
		C:     intersect.C{X: float64(x), Y: float64(y), R: float64(r)},
		xInfo: hud.crossInfo,
	}
	c := canvas.New(xc.draw)
	c.SetBounds(image.Rect(x-r, y-r, x+r, y+r))
	d := drag.EnableFor(c)
	d.AddDragListener(func() {
		hud.shapeInfo(c.Bounds())
		is.computeC(xc)
		b := c.Bounds()
		xc.MoveTo(
			math.Round(float64(b.Min.X+b.Max.X)/2),
			math.Round(float64(b.Min.Y+b.Max.Y)/2),
		)
	})
	is.addC(xc)
	app.AddComponent(d)
}
