package main

import (
	"log"

	"github.com/yarcat/playground/geom/app/application"
)

func main() {
	const screenWidth, screenHeight = 800, 600
	app := application.New(screenWidth, screenHeight)
	hud := newHUD(app)

	addButtons(app) // No usefull buttons yet.

	addRectangle(app, hud.shapeInfo)
	addCircle(app, 100, 100, 50, hud)
	addCircle(app, 100, 250, 50, hud)
	addTriangle(app, hud.shapeInfo)

	if err := application.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
