package main

import (
	"image"
	"image/color"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/button"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/component/features"
)

func addButtons(app *application.App) {
	b := button.New("Press me")
	b.SetBounds(image.Rect(100, 500, 250, 550))
	b.AddActionListener(func(b *button.Button) features.ActionListener {
		labels := [...]string{"Press me", "Drag me"}
		n := 0
		return func() {
			n = (n + 1) % len(labels)
			b.SetText(labels[n])
		}
	}(b))
	app.AddComponent(drag.EnableFor(b))

	b = button.New("Press me")
	b.SetBounds(image.Rect(300, 500, 450, 550))
	b.SetBgColor(color.RGBA{0x00, 0xf0, 0x00, 0xff})
	b.AddActionListener(func(b *button.Button) features.ActionListener {
		labels := [...]string{"Press me", "Yeah!", "Do it again!"}
		n := 0
		return func() {
			n = (n + 1) % len(labels)
			b.SetText(labels[n])
		}
	}(b))
	app.AddComponent(b)
}
