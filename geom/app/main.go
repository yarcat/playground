package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
)

func main() {
	const screenWidth, screenHeight = 800, 600
	app := NewApp(screenWidth, screenHeight)

	label := NewLabel(app, app)
	label.SetText("label")

	if err := ebiten.RunGame(app); err != nil {
		log.Fatalf("RunGame failed: %v", err)
	}
}
