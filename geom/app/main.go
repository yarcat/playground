package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/button"
	"github.com/yarcat/playground/geom/app/component/canvas"
	"github.com/yarcat/playground/geom/app/component/drag"
	"github.com/yarcat/playground/geom/app/component/features"
	"github.com/yarcat/playground/geom/app/component/label"
	"github.com/yarcat/playground/geom/shapes"
	"github.com/yarcat/playground/geom/vector"
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
	t.SetBounds(image.Rect(500, 500, 600, 600))
	app.AddComponent(drag.EnableFor(t))

	if err := application.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
