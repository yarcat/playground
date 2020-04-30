package ui

import (
	"image"
	"log"

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
	if parent == nil {
		return ui.root.Rect()
	}
	parentRect := screenRect(ui, parent)
	min := parentRect.Min.Add(rect.Min)
	max := min.Add(rect.Size())
	return image.Rectangle{Max: max, Min: min}
}

// elementAt returns an element under the point in logical screen coordinates
// (this includes scaling).
func elementAt(ui *UI, point image.Point) (Element, image.Rectangle) {
	log.Printf("elementAt: point=%#v, ui.elements=%#v", point, ui.elements)
	allElementsUnderPoint := func() (elements []Element) {
		// TODO(yarcat): Optimize this, cache rectangles calculated, etc.
		// This works while we don't have many elements thought.
		for el := range ui.elements {
			rect := screenRect(ui, el)
			log.Printf("elementAt: el=%#v, rect=%#v", el, rect)
			if point.In(rect) {
				elements = append(elements, el)
			}
		}
		return
	}

	removeElements := func(registry, toRemove []Element) (elements []Element) {
	registryLoop:
		for _, el := range registry {
			for _, remEl := range toRemove {
				if el == remEl {
					continue registryLoop
				}
			}
			elements = append(elements, el)
		}
		return
	}

	underPoint := allElementsUnderPoint()
	log.Printf("elementAt: underPoint=%#v", underPoint)

	// Now we need to precise the element. We'll try to find the first element,
	// which doesn't have a child in the elements registry.
	ancestors := make([]Element, 0, 10)
	for _, testElem := range underPoint {
		log.Printf("elementAt: testElem = %#v", testElem)
		ancestors = ancestors[:0]
		for el, ok := ui.elements[testElem]; ok; el, ok = ui.elements[el] {
			ancestors = append(ancestors, el)
		}
		log.Printf("elementAt: ancestors = %#v", ancestors)
		underPoint = removeElements(underPoint, ancestors)
	}
	// TODO(yarcat): Implement a way to return the same element even if there
	// are few of them under this point.
	if len(underPoint) != 1 {
		log.Printf("elementAt: len(underPoint) = %d, want 1. Returning root.",
			len(underPoint))
		return ui.root, ui.root.Rect()
	}
	return underPoint[0], screenRect(ui, underPoint[0])
}

func captureMouse(ui *UI, element Element) {
	ui.mouse.captureMouse(element)
}

func uncaptureMouse(ui *UI, element Element) {
	ui.mouse.uncaptureMouse(element)
}
