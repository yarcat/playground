package main

import (
	"image/color"
	"log"
	"test1/object"
	"test1/raycasting"
	"test1/vec"
	"test1/world"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type game struct {
	width, height int
	world         world.World
	object        object.Object
	t             time.Time
	img           *ebiten.Image
}

func (g *game) Update() error {
	t := time.Now()
	var d time.Duration
	if !g.t.IsZero() {
		d = t.Sub(g.t)
	}
	g.t = t
	g.object = object.HandleInput(g.object, &g.world, d,
		ebiten.IsKeyPressed(ebiten.KeyUp),
		ebiten.IsKeyPressed(ebiten.KeyDown),
		ebiten.IsKeyPressed(ebiten.KeyRight),
		ebiten.IsKeyPressed(ebiten.KeyLeft),
	)
	return nil
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.width, g.height
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.img == nil {
		g.img = ebiten.NewImage(g.width, g.height)
	}
	g.img.Clear()
	side := world.Draw(g.img, g.world, color.RGBA{150, 150, 150, 255})
	c := color.RGBA{255, 255, 255, 255}
	raycasting.FOV(&g.world, g.object.P(), g.object.V(), g.img.Bounds().Dx(),
		func(v1 vec.Vec, screenX int, dist float64, vertSide bool) {
			p1, p2 := g.object.P(), g.object.P().Add(v1.Scale(dist))
			vector.StrokeLine(g.img, float32(side*p1.X), float32(side*p1.Y), float32(side*p2.X), float32(side*p2.Y), 1, c, true)
			vector.StrokeCircle(g.img, float32(side*p1.X), float32(side*p1.Y), 0.1*float32(side), 1, c, true)
		})
	object.Draw(g.img, g.object, side)

	var opts ebiten.DrawImageOptions
	opts.GeoM.Translate(-g.object.P().X*side, -g.object.P().Y*side)
	var f_1 ebiten.GeoM
	/*
		00 10
		a  b
		01 11
		c  d
	*/
	//[x;y] -> [-y;x]
	f_1.SetElement(0, 0, -g.object.V().Y) // a
	f_1.SetElement(0, 1, g.object.V().X)  // b
	f_1.SetElement(1, 0, g.object.V().X)  // c
	f_1.SetElement(1, 1, g.object.V().Y)  // d
	f_1.Invert()
	opts.GeoM.Concat(f_1)
	opts.GeoM.Scale(-1, -1)
	opts.GeoM.Translate(float64(screen.Bounds().Dx())/2, float64(screen.Bounds().Dy())/2)
	screen.DrawImage(g.img, &opts)
}

func main() {
	const screenWidth, screenHeight = 640, 480
	ebiten.SetWindowSize(screenWidth, screenHeight)
	p, v := vec.Vec{X: 1.5, Y: 1.5}, vec.Vec{X: 0, Y: 1}
	g := game{
		width:  screenWidth,
		height: screenHeight,
		world:  world.MustGet(),
		object: object.New(p, v),
	}
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatalf("run game failed: %v", err)
	}
}
