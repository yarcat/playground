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

	{
		// TODO(yarcat): Currently our bounding rect for polygons has fixed
		// boundaries of 100. a=86.6 is a min equilateral triangle that fits.
		v, e := newEquilateralTriangle(85)
		addPolygon(app, 500, 100, 0, v, e, hud, is)
		addPolygon(app, 500, 250, math.Pi/3, v, e, hud, is)
	}
	{
		v, e := newSquare(50)
		addPolygon(app, 500, 100, 0, v, e, hud, is)
		addPolygon(app, 500, 250, math.Pi/4, v, e, hud, is)
	}
	addRectangle(app, 300, 100, 100, 100, 0, hud, is)
	addRectangle(app, 300, 250, 100, 100, 0, hud, is)
	addCircle(app, 100, 100, 50, hud, is)
	addCircle(app, 100, 250, 50, hud, is)

	if err := application.Run(app); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
