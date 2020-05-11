package main

import (
	"image"
	"image/color"
	"log"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/button"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/component/features"
	"github.com/yarcat/playground/geom/app/component/label"
	"github.com/yarcat/playground/geom/shapes"
)

func main() {
	const screenWidth, screenHeight = 800, 600
	app := application.New(screenWidth, screenHeight)

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
		x, y := 50*i, (20+5)*i
		l.SetBounds(image.Rect(x, y, x+len(s)*20, y+20))
		l.SetBgColor(data.color)
		l.SetHAlign(data.halign)
		app.AddComponent(drag.EnableFor(l))
	}

	b := button.New("Press me")
	b.SetBounds(image.Rect(100, 300, 200, 350))
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
	b.SetBounds(image.Rect(300, 320, 450, 400))
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

	c := canvas.New(func(img *canvas.Image) {
		if img.Invalidated() {
			w, h := img.Size()
			d := w
			if h < d {
				d = h
			}
			shapes.DrawCircle(img.Image, w/2, h/2, d/2, color.White)
		}
	})
	c.SetBounds(image.Rect(500, 500, 600, 600))
	app.AddComponent(drag.EnableFor(c))

	r := canvas.New(func(img *canvas.Image) {
		if img.Invalidated() {
			w, h := img.Size()
			shapes.DrawRectangle(img.Image, w, h, color.White)
		}
	})
	r.SetBounds(image.Rect(500, 500, 600, 600))
	app.AddComponent(drag.EnableFor(r))

	if err := application.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
