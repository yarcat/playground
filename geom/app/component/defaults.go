package component

import (
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
)

var defaultFont font.Face

// DefaultFont is a default font used for components that contain text e.g. labels and buttons.
func DefaultFont() font.Face {
	if defaultFont == nil {
		tt, err := truetype.Parse(fonts.ArcadeN_ttf)
		if err != nil {
			log.Fatal(err)
		}
		const dpi = 72
		defaultFont = truetype.NewFace(tt, &truetype.Options{
			Size:    12,
			DPI:     dpi,
			Hinting: font.HintingFull,
		})

	}
	return defaultFont
}

// DefaultBgColor returns default background color.
func DefaultBgColor() color.Color {
	return color.Black
}

// DefaultTextColor returns default text color.
func DefaultTextColor() color.Color {
	return color.White
}
