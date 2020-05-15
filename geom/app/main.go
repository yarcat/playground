package main

import (
	"log"
	"math"

	"github.com/yarcat/playground/geom/app/application"
)

func main() {
	const screenWidth, screenHeight = 800, 600
	app := application.New(screenWidth, screenHeight)
	hud := newHUD(app)
	is := &intersector{}

	addButtons(app) // No usefull buttons yet.

	addRectangle(app, 300, 100, 100, 100, 0, hud, is)
	addRectangle(app, 300, 250, 100, 100, 0, hud, is)
	addRectangle(app, 300, 400, 100, 100, math.Pi/4, hud, is)
	addCircle(app, 100, 100, 50, hud, is)
	addCircle(app, 100, 250, 50, hud, is)
	// addTriangle(app, hud.shapeInfo)

	if err := application.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
