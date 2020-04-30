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
	sim.AddBody(&body.Body{
		Image: shapes.Circle(50, color.White),
		Pos:   vector.New(screenWidth/2, screenHeight/2),
	})
	sim.AddBody(&body.Body{
		Image: shapes.Rectangle(60, 60, color.White),
		Pos:   vector.New(100, 50),
	})

	if err := ui.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
