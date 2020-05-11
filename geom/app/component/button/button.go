// Package button implements Button that can be clicked.
package button

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/yarcat/playground/geom/app/component"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
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
	rect            image.Rectangle
	bgColor         color.Color
	textColor       color.Color
	images          map[State]*ebiten.Image
	state           State
	entered         bool
	font            font.Face
	text            string
	actionListeners []func()
	dragged         bool
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

// HandleAdded is called after the button is added to its parent.
func (b *Button) HandleAdded(parent component.Component, features *ftrs.Features) {
	features.Add(
		ftrs.Draw(b.draw),
		ftrs.ListenMouseButtons(func(evt ftrs.GestureEvent) {
			if evt.Pressed() {
				b.state = Pressed
			} else {
				b.state = Released
			}
		}),
		ftrs.ListenAction(func() {
			for _, fn := range b.actionListeners {
				fn()
			}
		}),
		ftrs.ListenMouseEnter(func(evt ftrs.MotionEvent) {
			b.entered = true
		}),
		ftrs.ListenMouseLeave(func(evt ftrs.MotionEvent) {
			b.entered = false
		}),
		ftrs.ListenDrag(func(evt ftrs.DragEvent) {
			b.dragged = evt.State() != ftrs.DragStateEnded
		}),
	)
}

func (b *Button) invalidate() {
	for k, img := range b.images {
		if img != nil {
			img.Dispose()
		}
		delete(b.images, k)
	}
}

// SetText sets button text.
func (b *Button) SetText(text string) {
	b.text = text
	b.invalidate()
}

// AddActionListener registers action listener callback.
func (b *Button) AddActionListener(fn func()) {
	b.actionListeners = append(b.actionListeners, fn)
}

// SetBgColor sets background color.
func (b *Button) SetBgColor(c color.Color) {
	b.bgColor = c
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

// draw presents the button on the screen.
func (b *Button) draw(screen *ebiten.Image) {
	img := b.getImage()
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.rect.Min.X), float64(b.rect.Min.Y))
	screen.DrawImage(img, op)
}

func (b *Button) getImage() *ebiten.Image {
	if image, ok := b.images[b.drawState()]; ok {
		return image
	}

	image, _ := ebiten.NewImage(b.rect.Dx(), b.rect.Dy(), ebiten.FilterDefault)
	image.Fill(b.bgColor)

	w, h := bounds(b.font, b.text)
	x, y := (b.rect.Dx()-w)/2, (b.rect.Dy()+h)/2
	var tx, ty int
	if b.drawState() == Pressed {
		tx, ty = 1, 1
	}
	text.Draw(image, b.text, b.font, x+tx, y+ty, b.textColor)

	x1, y1 := float64(b.rect.Dx()), float64(b.rect.Dy())
	if b.drawState() == Released {
		// TODO(yarcat): Replace white with smth configurable.
		ebitenutil.DrawLine(image, 0, 0, x1, 0, color.White)
		ebitenutil.DrawLine(image, 1, 0, 1, y1, color.White)
	}

	b.images[b.state] = image
	return image
}

func (b *Button) drawState() State {
	if b.state == Pressed && b.entered && !b.dragged {
		return Pressed
	}
	return Released
}

func bounds(f font.Face, t string) (width, height int) {
	b, _ := font.BoundString(f, t)
	b = b.Sub(b.Min)
	return b.Max.X.Ceil(), b.Max.Y.Ceil()
}
