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
	app.Attach(sim, app.Root())
	circle := &body.Body{
		Image: shapes.Circle(50, color.White),
		Pos:   vector.New(300, 300),
	}
	rect := &body.Body{
		Image: shapes.Rectangle(50, 50, color.White),
		Pos:   vector.New(100, 100),
	}
	app.Attach(contrib.NewDragger(
		app,
		image.Rect(250, 250, 350, 350),
		func(r image.Rectangle) {
			circle.Pos = vector.New(float64(r.Min.X)+50, float64(r.Min.Y)+50)
		},
	), sim)
	app.Attach(contrib.NewDragger(
		app,
		image.Rect(75, 75, 125, 125),
		func(r image.Rectangle) {
			rect.Pos = vector.New(float64(r.Min.X)+25, float64(r.Min.Y)+25)
		},
	), sim)
	sim.AddBody(circle)
	sim.AddBody(rect)

	if err := ui.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
