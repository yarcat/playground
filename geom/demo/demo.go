package main

import (
	"image/color"
	"log"

	"github.com/yarcat/playground/geom/body"
	"github.com/yarcat/playground/geom/shape"
	"github.com/yarcat/playground/geom/simulation"
	vec "github.com/yarcat/playground/geom/vector"
)

const (
	screenWidth, screenHeight = 800, 600
)

func main() {
	sim := simulation.New(screenWidth, screenHeight)
	body := &body.Body{
		Shape: &shape.Circle{R: 100, Color: color.White},
		Pos:   vec.New(screenWidth/2, screenHeight/2),
	}
	simulation.AddBody(sim, body)
	if err := simulation.Run(sim); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
