package main

import (
	"image"
	"image/color"
	"log"

	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/button"
	"github.com/yarcat/playground/geom/app/component/label"
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
		app.AddComponent(l)
	}

	b := button.New("Press me")
	b.SetBounds(image.Rect(100, 300, 200, 350))
	app.AddComponent(b)

	lbls := [...]string{"Press me", "Yeah!", "Do it again!"}
	lblN := 0
	b = button.New(lbls[lblN])
	b.SetBounds(image.Rect(300, 320, 450, 400))
	b.SetBgColor(color.RGBA{0x00, 0xf0, 0x00, 0xff})
	b.AddActionListener(func() {
		lblN = (lblN + 1) % len(lbls)
		b.SetText(lbls[lblN])
	})
	app.AddComponent(b)

	if err := application.Run(app); err != nil {
		log.Fatalf("RunGame failed: %v", err)
	}
}
