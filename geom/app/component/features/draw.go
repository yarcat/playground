package features

import "github.com/hajimehoshi/ebiten"

// Draw registers a draw function to be called to render a frame.
func Draw(fn func(*ebiten.Image)) FeatureOption {
	return func(features *Features) {
		features.drawFn = fn
	}
}

// Draw calls registered draw function that renders a component every frame.
func (f *Features) Draw(screen *ebiten.Image) {
	if f != nil && f.drawFn != nil {
		f.drawFn(screen)
	}

}
