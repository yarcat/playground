package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func Draw(img *ebiten.Image, w World, clr color.RGBA) (side float64) {
	rows, cols := w.Dims()
	cells := rows
	if cols > cells {
		cells = cols
	}

	imgSize, imgH := img.Bounds().Dx(), img.Bounds().Dy()
	if imgH < imgSize {
		imgSize = imgH
	}

	side = float64(imgSize) / float64(cells)

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if w.EmptyAt(r, c) {
				continue
			}
			vector.DrawFilledRect(img, float32(float64(c)*side), float32(float64(r)*side),
				float32(side), float32(side), clr, true)
		}
	}

	return side
}
