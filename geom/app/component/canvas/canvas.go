// Package canvas implements Canavas component, which allows to draw shapes.
package canvas

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/app/component"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

// Image represents a drawable image buffer.
type Image struct {
	*ebiten.Image
	invalidated bool
}

func (img *Image) dispose() {
	if img.Image != nil {
		img.Image.Dispose()
		img.Image = nil
	}
}

func (img *Image) setImage(i *ebiten.Image) {
	img.Image = i
	img.invalidated = true
}

func (img Image) valid() bool {
	return img.Image != nil
}

// Invalidated returns true if the image requires to be redrawn.
func (img Image) Invalidated() bool {
	return img.invalidated
}

// Canvas allows to draw figures.
type Canvas struct {
	rect   image.Rectangle
	image  Image
	drawFn func(*Image)
}

// New returns new Canvas instance.
func New(fn func(*Image)) *Canvas {
	return &Canvas{drawFn: fn}
}

// HandleAdded is called right after this component is added to its parent.
func (c *Canvas) HandleAdded(parent component.Component, features *ftrs.Features) {
	features.Add(ftrs.Draw(c.draw))
}

// Bounds returns a rectangle of this element in logical screen coordinates.
func (c *Canvas) Bounds() image.Rectangle {
	return c.rect
}

// SetBounds sets a rectangle for this elemet in logical screen coordinates.
func (c *Canvas) SetBounds(rect image.Rectangle) {
	if c.rect.Size().Eq(rect.Size()) {
		c.rect = rect
		return
	}
	c.image.dispose()
	size := rect.Size()
	if size.Eq(image.Point{}) {
		c.image.setImage(nil)
	} else {
		image, _ := ebiten.NewImage(size.X, size.Y, ebiten.FilterDefault)
		c.image.setImage(image)
	}
	c.rect = rect
}

func (c *Canvas) draw(screen *ebiten.Image) {
	if !c.image.valid() {
		return
	}
	c.drawFn(&c.image)
	if c.image.invalidated {
		c.image.invalidated = false
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.rect.Min.X), float64(c.rect.Min.Y))
	screen.DrawImage(c.image.Image, op)
}
