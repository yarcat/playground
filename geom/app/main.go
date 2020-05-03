package main

import (
	"log"
)

func main() {
	const screenWidth, screenHeight = 800, 600
	app := NewApp(screenWidth, screenHeight)

	NewLabel(app, app).SetText("my label 1")
	NewLabel(app, app).SetText("my label 2")

	if err := Run(app); err != nil {
		log.Fatalf("RunGame failed: %v", err)
	}
}
