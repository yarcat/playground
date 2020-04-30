package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
)

// UI represents an abstract user interface manager. The UI manager is reponsible
// for bookkeeping the world e.g. parenting info, etc.
type UI struct {
	screenWidth, screenHeight int
	// screen is currently used image. It changes every time draw() callback is invoked.
	screen *ebiten.Image

	// elements maps elements to their parents.
	elements map[Element]Element

	root  Element
	mouse *mouseManager
}

// NewUI returns new UI instance ready to be executed with Run().
func NewUI(screenWidth, screenHeight int) *UI {
	ui := &UI{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		elements:     make(map[Element]Element),
	}
	ui.mouse = newMouseManager(ui)
	screenRect := image.Rect(0, 0, screenWidth, screenHeight)
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
	drawEvent := &DrawEvent{}
	for element := range ui.elements {
		SendEvent(element, drawEvent)
	}
}

func (ui *UI) update() {
	ui.mouse.update()
}

// elementImage returns current used screen image and the element's rectangle
// in screen coordinates.
func elementImage(ui *UI, element Element) (*ebiten.Image, image.Rectangle) {
	if element == ui.root {
		return ui.screen, ui.root.Rect()
	}
	return ui.screen, screenRect(ui, element)
}

// screenRect returns a rectangle for the given element in screen coordinates.
func screenRect(ui *UI, element Element) image.Rectangle {
	if element == ui.root {
		return ui.root.Rect()
	}
	rect := element.Rect()
	parent := ui.elements[element]
	parentRect := screenRect(ui, parent)
	min := parentRect.Min.Add(rect.Min)
	max := min.Add(rect.Size())
	return image.Rectangle{Max: max, Min: min}
}

// elementAt returns an element under the point in logical screen coordinates
// (this includes scaling).
func elementAt(ui *UI, point image.Point) Element {
	// TODO(yarcat): Return a real element under the point.
	return ui.root
}
