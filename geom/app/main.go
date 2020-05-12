package main

import (
	"log"

	"github.com/yarcat/playground/geom/app/application"
)

func main() {
	const screenWidth, screenHeight = 800, 600
	app := application.New(screenWidth, screenHeight)
	updateStatus := newHUD(app)

	addButtons(app) // No usefull buttons yet.

	addRectangle(app, updateStatus)
	addCircle(app, updateStatus)
	addTriangle(app, updateStatus)

	if err := application.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
