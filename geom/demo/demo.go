package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/body"
	"github.com/yarcat/playground/geom/shapes"
	"github.com/yarcat/playground/geom/simulation"
	"github.com/yarcat/playground/geom/ui"
	"github.com/yarcat/playground/geom/ui/contrib"
	"github.com/yarcat/playground/geom/vector"
)

const (
	screenWidth, screenHeight = 800, 600
)

func main() {
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Simulation Demo")

	app := ui.NewUI(screenWidth, screenHeight)

	simRect := image.Rect(0, 0, screenWidth, screenHeight)
	sim := simulation.New(ui.NewElement(app, simRect))
	app.Root().AddChild(sim)
	circle := &body.Body{
		Image: shapes.Circle(50, color.White),
		Pos:   vector.New(screenWidth/2, screenHeight/2),
	}
	rect := &body.Body{
		Image: shapes.Rectangle(60, 60, color.White),
		Pos:   vector.New(100, 50),
	}
	sim.AddChild(contrib.NewDragger(
		app,
		image.Rect(
			100,
			50,
			160,
			110,
		),
		func(r image.Rectangle) {
			circle.Pos = vector.New(float64(r.Min.X), float64(r.Min.Y))
		},
	))
	sim.AddChild(contrib.NewDragger(
		app,
		image.Rect(
			screenWidth/2,
			screenHeight/2,
			screenWidth/2+120,
			screenHeight/2+120,
		),
		func(r image.Rectangle) {
			rect.Pos = vector.New(float64(r.Min.X), float64(r.Min.Y))
		},
	))
	sim.AddBody(circle)
	sim.AddBody(rect)

	if err := ui.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
