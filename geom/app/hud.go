package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/component/label"
	"github.com/yarcat/playground/geom/app/intersect"
)

type hud struct {
	shapeInfo func(image.Rectangle)
	crossInfo func(intersect.I)
}

func newHUD(app *application.App) *hud {
	var labels [3]*label.Label

	const s = "Hello, world!"
	for i, data := range []struct {
		color  color.Color
		halign label.HAlign
	}{
		{color.RGBA{0x80, 0x00, 0x00, 0xff}, label.HLeft},
		{color.RGBA{0x80, 0x80, 0x00, 0xff}, label.HCenter},
		{color.RGBA{0x00, 0x80, 0x00, 0xff}, label.HRight},
	} {
		l := label.New(s)
		labels[i] = l
		x, y := 265*i, 5
		l.SetBounds(image.Rect(x, y, x+len(s)*20, y+20))
		l.SetBgColor(data.color)
		l.SetHAlign(data.halign)
		app.AddComponent(drag.EnableFor(l))
	}

	return &hud{
		shapeInfo: func(r image.Rectangle) {
			status := fmt.Sprintf("(%d,%d) %dx%d", r.Min.X, r.Min.Y, r.Dx(), r.Dy())
			labels[1].SetText(status)
		},
		crossInfo: func(i intersect.I) {
			status := fmt.Sprintf("(%.2f,%.2f) %.6f", i.N.X, i.N.Y, i.P)
			labels[0].SetText(status)
		},
	}
}
