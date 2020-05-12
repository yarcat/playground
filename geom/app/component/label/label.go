// Package label implements a Label that can display text strings.
package label

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/yarcat/playground/geom/app/component"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
	"golang.org/x/image/font"
)

// HAlign represents horizontal aligment of the text.
type HAlign int

const (
	// HLeft makes the text align left in the boundary box.
	HLeft HAlign = iota
	// HCenter makes the text align center in the boundary box.
	HCenter
	// HRight makes the text align right in the boundary box.
	HRight
)

// Label can display text strings.
type Label struct {
	rect image.Rectangle

	text      string
	image     *ebiten.Image
	font      font.Face
	halign    HAlign
	bgColor   color.Color
	textColor color.Color
}

// New returns new Label instance.
func New(text string) *Label {
	return &Label{
		text:      text,
		font:      component.DefaultFont(),
		bgColor:   component.DefaultBgColor(),
		textColor: component.DefaultTextColor(),
	}
}

// SetText updates label text.
func (l *Label) SetText(text string) {
	if l.text != text {
		l.text = text
		l.image = nil
	}
}

// HandleAdded is called right after this component is added to its parent.
func (l *Label) HandleAdded(parent component.Component, features *ftrs.Features) {
	features.Add(ftrs.Draw(l.draw))
}

// SetBgColor sets background color.
func (l *Label) SetBgColor(c color.Color) {
	l.bgColor = c
}

// SetTextColor sets text color.
func (l *Label) SetTextColor(c color.Color) {
	l.textColor = c
}

// SetFont sets font used by the label.
func (l *Label) SetFont(f font.Face) {
	l.font = f
}

// SetHAlign sets horizontal alignment of the text.
func (l *Label) SetHAlign(a HAlign) {
	l.halign = a
}

// Bounds returns a rectangle of this element in logical screen coordinates.
func (l Label) Bounds() image.Rectangle {
	return l.rect
}

// SetBounds sets a rectangle for this elemet in logical screen coordinates.
func (l *Label) SetBounds(rect image.Rectangle) {
	l.rect = rect
}

// draw draws the label.
func (l *Label) draw(screen *ebiten.Image) {
	image := l.getImage()
	if image == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(l.rect.Min.X), float64(l.rect.Min.Y))
	screen.DrawImage(l.image, op)
}

func (l *Label) getImage() *ebiten.Image {
	if l.image != nil {
		return l.image
	}
	if l.text == "" || l.rect.Dy() == 0 || l.font == nil {
		return nil
	}
	l.image, _ = ebiten.NewImage(l.rect.Dx(), l.rect.Dy(), ebiten.FilterDefault)
	if l.bgColor != nil {
		l.image.Fill(l.bgColor)
	}
	w, h := bounds(l.font, l.text)
	x := 0
	if l.halign == HCenter {
		x = (l.rect.Dx() - w) / 2
	} else if l.halign == HRight {
		x = l.rect.Dx() - w
	}
	text.Draw(l.image, l.text, l.font, x, (l.rect.Dy()+h)/2, l.textColor)
	return l.image
}

func bounds(f font.Face, t string) (width, height int) {
	b, _ := font.BoundString(f, t)
	b = b.Sub(b.Min)
	return b.Max.X.Ceil(), b.Max.Y.Ceil()
}
