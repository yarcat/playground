package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/component/label"
)

func newHUD(app *application.App) func(image.Rectangle) {
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

	return func(r image.Rectangle) {
		status := fmt.Sprintf("(%d,%d) %dx%d", r.Min.X, r.Min.Y, r.Dx(), r.Dy())
		labels[1].SetText(status)
	}
}
