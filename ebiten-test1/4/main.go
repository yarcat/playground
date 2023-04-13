package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"test1/object"
	vec2d "test1/vec"
	"test1/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type vec = vec2d.Vec

type app struct {
	width, height int
	world         world.World
	object        object.Object
	t             time.Time
	mapImg        *ebiten.Image
}

func (a *app) Update() error {
	t := time.Now()
	dt := t.Sub(a.t).Seconds()
	a.t = t
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		a.object = a.object.Rotate(-math.Pi / 2 * dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		a.object = a.object.Rotate(math.Pi / 2 * dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			dt = -dt
		}
		o := a.object.Move(dt)
		p := o.P()
		if a.world.EmptyAt(int(p.Y), int(p.X)) {
			a.object = o
		}
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	mapImg := a.mapImg
	mapImg.Clear()
	mapW, mapH := mapImg.Size()
	size := world.Draw(mapImg, a.world, color.RGBA{255, 255, 255, 255})
	a.drawObjectAndMap(screen, mapImg, size)
	var opt ebiten.DrawImageOptions
	p, v := a.object.P(), a.object.V().Scale(-1)
	opt.GeoM.Translate(p.X*size, p.Y*size)
	opt.GeoM.SetElement(0, 0, v.Y)
	opt.GeoM.SetElement(1, 0, -v.X)
	opt.GeoM.SetElement(0, 1, v.X)
	opt.GeoM.SetElement(1, 1, v.Y)
	opt.GeoM.Invert()
	opt.GeoM.Translate(float64(mapW)/2, float64(mapH)/2)
	opt.CompositeMode = ebiten.CompositeModeCopy
	opt.Filter = ebiten.FilterLinear
	screen.DrawImage(mapImg, &opt)
}

func (a *app) Layout(int, int) (int, int) { return a.width, a.height }

func (a *app) drawObjectAndMap(img, mapImg *ebiten.Image, side float64) {
	p, v := a.object.P(), a.object.V()
	ebitenutil.DrawRect(mapImg, p.X*side-1, p.Y*side-1,
		3, 3, color.White)
	for x := 0; x < a.width; x++ {
		camX := 2*float64(x)/float64(a.width) - 1
		v := v.Add(vec{X: v.Y * camX * 2 / 3, Y: -v.X * camX * 2 / 3})
		dirDist, dirLen, dcol, drow := ddaParams(p, v)
		dist := 0.0
		var vertSide bool
		for col, row := int(p.X), int(p.Y); ; {
			if dirDist.X < dirDist.Y {
				dist = dirDist.X
				dirDist.X += dirLen.X
				col += dcol
				vertSide = false
			} else {
				dist = dirDist.Y
				dirDist.Y += dirLen.Y
				row += drow
				vertSide = true
			}
			if !a.world.EmptyAt(row, col) {
				break
			}
		}
		{
			h := p.Add(v.Scale(dist))
			ebitenutil.DrawLine(mapImg, p.X*side, p.Y*side,
				h.X*side, h.Y*side, color.RGBA{100, 100, 100, 255})
		}
		{
			// ebitenutil.DrawRect(img, h.X*side-1, h.Y*side-1, 3, 3, color.RGBA{R: 255, A: 255})
			maxH := float64(a.height)
			h := maxH / dist
			if h > maxH {
				h = maxH
			}
			y := (maxH - h) / 2
			kf := 255 / dist
			if kf > 255 {
				kf = 255
			}
			k := uint8(kf)
			var c color.RGBA
			if vertSide {
				c = color.RGBA{k, k, k, 255}
			} else {
				c = color.RGBA{k / 2, k / 2, k / 2, 255}
			}
			for row := int(y); row <= int(y+h); row++ {
				img.Set(a.width-x, row, c)
			}
			// ebitenutil.DrawLine(img,
			// 	float64(a.width-x), y, float64(a.width-x), y+h, c)
		}
	}
}

func ddaParams(p, v vec) (dirDist, dirLen vec, dx, dy int) {
	dirLen.X = math.Abs(1 / v.X)
	if v.X < 0 {
		dx = -1
		dirDist.X = (p.X - math.Trunc(p.X)) * dirLen.X
	} else {
		dx = 1
		dirDist.X = (math.Trunc(p.X+1) - p.X) * dirLen.X
	}
	dirLen.Y = math.Abs(1 / v.Y)
	if v.Y < 0 {
		dy = -1
		dirDist.Y = (p.Y - math.Trunc(p.Y)) * dirLen.Y
	} else {
		dy = 1
		dirDist.Y = (math.Trunc(p.Y+1) - p.Y) * dirLen.Y
	}
	return dirDist, dirLen, dx, dy
}

func main() {
	const screenWidth, screenHeight = 640, 480
	ebiten.SetWindowSize(screenWidth, screenHeight)
	d := screenHeight / 4
	if dd := screenWidth / 4; dd < d {
		d = dd
	}
	app := &app{
		width:  screenWidth,
		height: screenHeight,
		world:  world.MustGet(),
		object: object.New(
			vec{X: 5.5, Y: 1.5}, // P
			vec{Y: 1.0},         // V
		),
		mapImg: ebiten.NewImage(d, d),
	}
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
