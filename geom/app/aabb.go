package main

import (
	"image"
	"image/color"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/shapes"
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
