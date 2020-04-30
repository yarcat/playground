package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/body"
	"github.com/yarcat/playground/geom/shape"
	"github.com/yarcat/playground/geom/simulation"
	"github.com/yarcat/playground/geom/ui"
	vec "github.com/yarcat/playground/geom/vector"
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
		Shape: &shape.Circle{R: 10, Color: color.White},
		Pos:   vec.New(screenWidth-10, 10),
	})

	if err := ui.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
