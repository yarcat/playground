// Package application provides an App class that is a main container of elements,
// event dispatcher, etc.
package application

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/yarcat/playground/geom/app/application/states"
	"github.com/yarcat/playground/geom/app/component"
	ftrs "github.com/yarcat/playground/geom/app/component/features"
)

// New returns new application instance with initialized logical window width
// and height.
func New(width, height int) *App {
	// Setting defaults.
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowResizable(true)

	app := &App{
		width:    width,
		height:   height,
		features: make(map[component.Component]*ftrs.Features),
	}
	app.gestureManager = gestureManagerImpl{
		app:     app,
		motions: make(map[*ftrs.Features]*states.MouseMotionState),
	}
	return app
}

// App represents an application. It manages internal objects and present
// the on-screen.
type App struct {
	width, height  int
	components     []component.Component
	gestureManager gestureManagerImpl
	features       map[component.Component]*ftrs.Features
}

// AddComponent adds a component and fires life-cycle events.
func (app *App) AddComponent(c component.WithLifecycle) {
	app.components = append(app.components, c)
	features := &ftrs.Features{}
	app.features[c] = features
	// TODO(yarcat): Parent must not be nil.
	c.HandleAdded(nil /* parent */, features)
}

// ComponentAt returns a component under a window point.
func (app App) ComponentAt(pt image.Point) component.Component {
	for _, c := range app.components {
		if pt.In(c.Bounds()) {
			return c
		}
	}
	return nil
}
