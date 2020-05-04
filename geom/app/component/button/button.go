// Package button implements Button that can be clicked.
package button

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/yarcat/playground/geom/app/component"
	"golang.org/x/image/font"
)

const borderWidth = 2

var (
	pt               = image.Pt(borderWidth, borderWidth)
	defaultBgColor   = color.RGBA{0x80, 0x80, 0x80, 0xff}
	defaultTextColor = color.White
)

// State defines button state.
type State int

const (
	// Released is a released button state.
	Released State = iota
	// Pressed is a pressed button state.
	Pressed
)

// Button is a clickable element.
type Button struct {
	rect      image.Rectangle
	bgColor   color.Color
	textColor color.Color
	images    map[State]*ebiten.Image
	state     State
	font      font.Face
	text      string
}

// New returns new Button instance.
func New(text string) *Button {
	b := &Button{
		bgColor:   defaultBgColor,
		textColor: defaultTextColor,
		text:      text,
		images:    make(map[State]*ebiten.Image),
		font:      component.DefaultFont(),
	}
	return b
}

// Bounds returns a rectangle in logical screen coordinates.
func (b Button) Bounds() image.Rectangle {
	return b.rect
}

// SetBounds sets a rectangle in logical screen coordinates.
func (b *Button) SetBounds(rect image.Rectangle) {
	b.rect = rect
}

// SetFont sets the font.
func (b *Button) SetFont(f font.Face) {
	b.font = f
}

// Draw presents the button on the screen.
func (b *Button) Draw(screen *ebiten.Image) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		b.state = Pressed
	} else {
		b.state = Released
	}
	img := b.getImage()
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.rect.Min.X), float64(b.rect.Min.Y))
	screen.DrawImage(img, op)
}

func (b *Button) getImage() *ebiten.Image {
	if image, ok := b.images[b.state]; ok {
		return image
	}

	image, _ := ebiten.NewImage(b.rect.Dx(), b.rect.Dy(), ebiten.FilterDefault)
	image.Fill(b.bgColor)

	w, h := bounds(b.font, b.text)
	x, y := (b.rect.Dx()-w)/2, (b.rect.Dy()+h)/2
	var tx, ty int
	if b.state == Pressed {
		tx, ty = 1, 1
	}
	text.Draw(image, b.text, b.font, x+tx, y+ty, b.textColor)

	x1, y1 := float64(b.rect.Dx()), float64(b.rect.Dy())
	if b.state == Released {
		// TODO(yarcat): Replace white with smth configurable.
		ebitenutil.DrawLine(image, 0, 0, x1, 0, color.White)
		ebitenutil.DrawLine(image, 1, 0, 1, y1, color.White)
	}

	b.images[b.state] = image
	return image
}

func bounds(f font.Face, t string) (width, height int) {
	b, _ := font.BoundString(f, t)
	b = b.Sub(b.Min)
	return b.Max.X.Ceil(), b.Max.Y.Ceil()
}
