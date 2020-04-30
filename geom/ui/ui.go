package ui

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/image"
)

// UI represents an abstract user interface manager. The UI manager is reponsible
// for bookkeeping the world e.g. parenting info, etc.
type UI struct {
	screenWidth, screenHeight int
	// screen is currently used image. It changes every time draw() callback is invoked.
	screen *ebiten.Image
	img    image.Image
	root   Element
	// elements maps elements to their parents.
	elements map[Element]Element
}

// NewUI returns new UI instance ready to be executed with Run().
func NewUI(screenWidth, screenHeight int) *UI {
	ui := &UI{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		elements:     make(map[Element]Element),
	}
	screenRect := Rect{
		Left: 0, Bottom: 0,
		Right: screenWidth, Top: screenHeight,
	}
	ui.root = NewElement(ui, screenRect)
	return ui
}

// Attach creates a parent-child association.
func (ui *UI) Attach(child, parent Element) {
	ui.elements[child] = parent
}

// Root returns top-most element, which represents the screen.
func (ui *UI) Root() Element {
	return ui.root
}

func (ui *UI) draw(screen *ebiten.Image) {
	ui.screen = screen
	ui.img = image.NewImage(ui.screen, &ebiten.DrawImageOptions{})
	drawEvent := &DrawEvent{}
	for element := range ui.elements {
		SendEvent(element, drawEvent)
	}
}

func (ui *UI) image(element Element) Image {
	return ui.img
}
