package main

import (
	"image"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/yarcat/playground/geom/app/application"
	"github.com/yarcat/playground/geom/app/component/label"
	"golang.org/x/image/font"
)

func main() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	font := truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

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
		l.SetFont(font)
		l.SetBgColor(data.color)
		l.SetHAlign(data.halign)
		app.AddDrawable(l)
	}

	if err := application.Run(app); err != nil {
		log.Fatalf("RunGame failed: %v", err)
	}
}
