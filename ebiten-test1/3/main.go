package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenWidth, screenHeight = 640, 480

const mapWidth, mapHeight = 24, 24

var worldMap = [mapWidth][mapHeight]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

type game struct{}

var posX, posY = 22.0, 12.0
var dirX, dirY = -1.0, 0.0
var planeX, planeY = 0.0, 0.66
var lastT time.Time

func (game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }
func (game) Update() error {
	now := time.Now()
	frameTime := now.Sub(lastT).Seconds()
	lastT = now
	moveSpeed := frameTime * 5
	rotSpeed := frameTime * 3
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if worldMap[int(posX+dirX*moveSpeed)][int(posY)] == 0 {
			posX += dirX * moveSpeed
		}
		if worldMap[int(posX)][int(posY+dirY*moveSpeed)] == 0 {
			posY += dirY * moveSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if worldMap[int(posX-dirX*moveSpeed)][int(posY)] == 0 {
			posX -= dirX * moveSpeed
		}
		if worldMap[int(posX)][int(posY-dirY*moveSpeed)] == 0 {
			posY -= dirY * moveSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		oldDirX := dirX
		dirX = dirX*math.Cos(-rotSpeed) - dirY*math.Sin(-rotSpeed)
		dirY = oldDirX*math.Sin(-rotSpeed) + dirY*math.Cos(-rotSpeed)
		oldPlaneX := planeX
		planeX = planeX*math.Cos(-rotSpeed) - planeY*math.Sin(-rotSpeed)
		planeY = oldPlaneX*math.Sin(-rotSpeed) + planeY*math.Cos(-rotSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		oldDirX := dirX
		dirX = dirX*math.Cos(rotSpeed) - dirY*math.Sin(rotSpeed)
		dirY = oldDirX*math.Sin(rotSpeed) + dirY*math.Cos(rotSpeed)
		oldPlaneX := planeX
		planeX = planeX*math.Cos(rotSpeed) - planeY*math.Sin(rotSpeed)
		planeY = oldPlaneX*math.Sin(rotSpeed) + planeY*math.Cos(rotSpeed)
	}
	return nil
}
func (game) Draw(screen *ebiten.Image) {
	const w, h = screenWidth, screenHeight
	for x := 0; x < w; x++ {
		cameraX := float64(2*x)/float64(w) - 1
		rayDirX := dirX + planeX*cameraX
		rayDirY := dirY + planeY*cameraX
		mapX := int(posX)
		mapY := int(posY)
		var sideDistX, sideDistY float64
		deltaDistX := math.Abs(1 / rayDirX)
		deltaDistY := math.Abs(1 / rayDirY)
		var perpWallDist float64
		var stepX, stepY int
		hit := false
		var side bool
		if rayDirX < 0 {
			stepX = -1
			sideDistX = (posX - float64(mapX)) * deltaDistX
		} else {
			stepX = 1
			sideDistX = (float64(mapX+1) - posX) * deltaDistX
		}
		if rayDirY < 0 {
			stepY = -1
			sideDistY = (posY - float64(mapY)) * deltaDistY
		} else {
			stepY = 1
			sideDistY = (float64(mapY+1) - posY) * deltaDistY
		}
		for !hit {
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				mapX += stepX
				side = false
			} else {
				sideDistY += deltaDistY
				mapY += stepY
				side = true
			}
			if worldMap[mapX][mapY] > 0 {
				hit = true
			}
		}
		if !side {
			perpWallDist = sideDistX - deltaDistX
		} else {
			perpWallDist = sideDistY - deltaDistY
		}

		lineHeight := int(h / perpWallDist)
		drawEnd := lineHeight/2 + h/2
		drawStart := drawEnd - lineHeight
		if drawStart < 0 {
			drawStart = 0
		}
		if drawEnd >= h {
			drawEnd = h
		}
		var c color.RGBA
		switch worldMap[mapX][mapY] {
		case 1:
			c = color.RGBA{R: 255, A: 255}
		case 2:
			c = color.RGBA{G: 255, A: 255}
		case 3:
			c = color.RGBA{B: 255, A: 255}
		case 4:
			c = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		default:
			c = color.RGBA{R: 255, G: 255, A: 255} // Yellow
		}
		if side {
			c = color.RGBA{R: c.R / 2, G: c.G / 2, B: c.B / 2, A: 255}
		}
		for y := drawStart; y <= drawEnd; y++ {
			screen.Set(x, y, c)
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(game{}); err != nil {
		log.Fatal(err)
	}
}
