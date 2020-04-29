package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func main() {
	rand.Seed(time.Now().Unix())

	var firework Firework
	firework.Init()

	var trigger *DoAfterTicksComponent
	trigger = NewDoAfterTicksComponent(0, func() {
		firework.Shoot()
		trigger.Activated = false
		trigger.TicksUntilActivation = 100 + uint(rand.Int()%200)
	})

	shooter := NewElement(Vector{0, 0})
	shooter.AddComponents(trigger)
	firework.AddElement(shooter)

	if err := ebiten.Run(firework.Update, screenWidth, screenHeight, 2, "Test"); err != nil {
		panic(err)
	}
}

func explode(element *Element, pieces uint, addElement func(*Element)) func() {
	return func() {
		element.Active = false
		var ecv Vector
		element.ForEachComponent(func(c Component) {
			switch c := c.(type) {
			case *VelocityComponent:
				ecv.AddI(&c.Velocity)
			}
		})
		for i := uint(0); i < pieces; i++ {
			alpha := rand.Float64() * 2 * math.Pi
			v := rand.Float64()*50 + 20
			element := NewElement(element.Pos)
			ev := Vector{
				ecv.X + v*math.Cos(alpha),
				ecv.Y + v*math.Sin(alpha),
			}
			element.AddComponents(
				NewDeactivateInvisibleComponent(),
				NewSquareComponent(2 /* size */, RandomColor, DecayAlpha(1)),
				NewVelocityComponent(ev),
				NewDoEveryTickComponent(addTrail(element, 10, addElement,
					FixedColor(255, 255, 255, 255))),
			)
			addGravity(element)
			addElement(element)
		}
	}
}

func addTrail(element *Element, count uint, addElement func(*Element), colorOpt func(*SquareComponent)) func() {
	return func() {
		trail := NewElement(element.Pos)
		trail.AddComponents(
			NewSquareComponent(1 /* size */, colorOpt, DecayAlpha(uint8(255.0/count))),
			NewDoAfterTicksComponent(count, DeactivateFn(trail)),
		)
		addGravity(trail)
		addElement(trail)
	}
}

type Firework struct {
	elems     map[*Element]interface{}
	scheduled []*Element
	last      time.Time
}

func (f *Firework) Init() {
	f.elems = make(map[*Element]interface{}, 1000)
	f.last = time.Now()
}

func (f *Firework) Shoot() {
	initialShells := rand.Int()%20 + 10

	shootPosition := Vector{rand.Float64() * screenWidth, 0}
	shootAngle := math.Atan(screenHeight / (screenWidth/2 - shootPosition.X))
	if shootPosition.X > screenWidth/2 {
		shootAngle += math.Pi
	}

	for i := 0; i < initialShells; i++ {
		// Randomize initial angle and velocity.
		alpha := shootAngle + math.Pi/16 - rand.Float64()*math.Pi/8
		v := rand.Float64()*20 + 170
		velocity := Vector{v * math.Cos(alpha), v * math.Sin(alpha)}
		velocityComponent := NewVelocityComponent(velocity)

		pos := shootPosition
		pos.AddXY(rand.Float64()*10, 0)
		element := NewElement(pos)
		colorOpt := RandomFixedColor()
		element.AddComponents(
			NewDeactivateInvisibleComponent(),
			NewSquareComponent(2 /* size */, colorOpt),
			NewDoAfterTicksComponent(
				80+uint(rand.Int()%80), /* ticks */
				explode(
					element,
					10+uint(rand.Int()%5), /* pieces */
					f.ScheduleAddElement),
			),
			velocityComponent,
			NewDoEveryTickComponent(addTrail(element, 40, f.AddElement,
				colorOpt)),
		)
		addGravity(element)

		f.AddElement(element)
	}
}

func (f *Firework) AddElement(e *Element) {
	f.elems[e] = nil
}

func (f *Firework) ScheduleAddElement(e *Element) {
	f.scheduled = append(f.scheduled, e)
}

const debugTemplate = `Current FPS    : %.2f
Scheduled      : %d
Active Elements: %d
`

func (f *Firework) Update(screen *ebiten.Image) error {
	debugMsg := fmt.Sprintf(debugTemplate,
		ebiten.CurrentFPS(), len(f.scheduled), len(f.elems))
	ebitenutil.DebugPrint(screen, debugMsg)
	for _, element := range f.scheduled {
		f.AddElement(element)
	}
	f.scheduled = f.scheduled[:0]

	last := f.last
	f.last = time.Now()
	dt := f.last.Sub(last)
	for element := range f.elems {
		if element.Active {
			element.Update(dt)
		}
		if !element.Active {
			delete(f.elems, element)
		}
	}

	if ebiten.IsRunningSlowly() {
		return nil
	}

	// Image buffer that we'll flip later to move the origin to the
	// buttom left corner.
	img, err := ebiten.NewImage(screenWidth, screenHeight, 0)
	if err != nil {
		return err
	}

	for element := range f.elems {
		element.Draw(img)
	}

	// Move the origin to the bottom left corner by flipping the image.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, screenHeight)
	return screen.DrawImage(img, op)
}

func addGravity(e *Element) {
	e.AddComponents(NewAccellerationComponent(Vector{0, -80}))
}

type BaseComponent struct {
	Element *Element
}

func (bc *BaseComponent) SetElement(e *Element) {
	bc.Element = e
}

type DoEveryTickComponent struct {
	BaseComponent
	Fn func()
}

func NewDoEveryTickComponent(fn func()) *DoEveryTickComponent {
	return &DoEveryTickComponent{BaseComponent{}, fn}
}

func (c DoEveryTickComponent) Update(d time.Duration) {
	c.Fn()
}

func (tc DoEveryTickComponent) Draw(screen *ebiten.Image) {}

type DeactivateInvisibleComponent struct {
	BaseComponent
}

func NewDeactivateInvisibleComponent() *DeactivateInvisibleComponent {
	return &DeactivateInvisibleComponent{BaseComponent{}}
}

func (dic DeactivateInvisibleComponent) Update(d time.Duration) {
	if dic.Element.Pos.Y < 0 {
		dic.Element.Active = false
	}
}

func (dic DeactivateInvisibleComponent) Draw(screen *ebiten.Image) {}

type DoAfterTicksComponent struct {
	BaseComponent
	TicksUntilActivation uint
	Fn                   func()
	Activated            bool
}

func NewDoAfterTicksComponent(ticks uint, fn func()) *DoAfterTicksComponent {
	return &DoAfterTicksComponent{BaseComponent{}, ticks, fn, false}
}

func DeactivateFn(e *Element) func() {
	return func() {
		e.Active = false
	}
}

func (c *DoAfterTicksComponent) Update(d time.Duration) {
	if c.Activated {
		return
	}
	if c.TicksUntilActivation > 0 {
		c.TicksUntilActivation--
		return
	}
	c.Activated = true // Allow to reactivate.
	c.Fn()
}

func (c *DoAfterTicksComponent) Draw(screen *ebiten.Image) {}

type AccellerationComponent struct {
	BaseComponent
	Accelleration Vector
}

func NewAccellerationComponent(accelleration Vector) *AccellerationComponent {
	return &AccellerationComponent{
		BaseComponent{},
		accelleration,
	}
}

func (ac *AccellerationComponent) Update(d time.Duration) {
	dt := d.Seconds()
	dv := ac.Accelleration
	dv.Scale(dt)
	ac.Element.ForEachComponent(func(c Component) {
		switch c := c.(type) {
		case *VelocityComponent:
			c.Velocity.AddI(&dv)
		}
	})
}

func (ac *AccellerationComponent) Draw(screen *ebiten.Image) {}

type VelocityComponent struct {
	BaseComponent
	Velocity Vector
}

func NewVelocityComponent(velocity Vector) *VelocityComponent {
	return &VelocityComponent{
		BaseComponent{},
		velocity,
	}
}

func (vc *VelocityComponent) Update(d time.Duration) {
	dt := d.Seconds()
	dv := vc.Velocity
	dv.Scale(dt)
	vc.Element.Pos.AddI(&dv)
}

func (vc *VelocityComponent) Draw(screen *ebiten.Image) {}

type SquareComponent struct {
	BaseComponent
	Color color.RGBA
	Size  float64
	Fns   []func(*color.RGBA)
}

func DecayAlpha(step uint8) func(*SquareComponent) {
	return func(c *SquareComponent) {
		c.Fns = append(c.Fns, func(c *color.RGBA) {
			if c.A < step {
				step = c.A
			}
			c.A -= step
		})
	}
}

var RandomColor = func(c *SquareComponent) {
	c.Fns = append(c.Fns, func(c *color.RGBA) {
		c.R = uint8(rand.Int() % 256)
		c.G = uint8(rand.Int() % 256)
		c.B = uint8(rand.Int() % 256)
	})
}

func FixedColor(r, g, b, a uint8) func(*SquareComponent) {
	return func(c *SquareComponent) {
		c.Color = color.RGBA{r, g, b, a}
	}
}

func RandomFixedColor() func(*SquareComponent) {
	return FixedColor(
		uint8(128+rand.Int()%128),
		uint8(128+rand.Int()%128),
		uint8(128+rand.Int()%128),
		255)
}

func NewSquareComponent(size float64, opts ...func(*SquareComponent)) *SquareComponent {
	c := &SquareComponent{
		BaseComponent{},
		color.RGBA{255, 255, 255, 255},
		size,
		nil, /* Fns */
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *SquareComponent) Update(d time.Duration) {
	for _, fn := range c.Fns {
		fn(&c.Color)
	}
}

func (c *SquareComponent) Draw(screen *ebiten.Image) {
	pos := c.Element.Pos
	x := pos.X
	y := pos.Y
	ebitenutil.DrawRect(screen, x, y,
		c.Size /* width */, c.Size /* height */, c.Color)
}

type Vector struct {
	X, Y float64
}

func (v *Vector) Scale(k float64) {
	v.X *= k
	v.Y *= k
}

func (v *Vector) AddI(v2 *Vector) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector) AddXY(x, y float64) {
	v.X += x
	v.Y += y
}

type Element struct {
	Pos        Vector
	Active     bool
	components []Component
}

func NewElement(pos Vector) *Element {
	return &Element{
		Pos:    pos,
		Active: true,
	}
}

func (e *Element) AddComponents(comps ...ComponentWithElement) {
	for _, comp := range comps {
		comp.SetElement(e)
		e.components = append(e.components, comp)
	}
}

func (e *Element) ForEachComponent(f func(Component)) {
	for _, c := range e.components {
		f(c)
	}
}

func (e *Element) Update(d time.Duration) {
	for _, c := range e.components {
		c.Update(d)
	}
}

func (e *Element) Draw(screen *ebiten.Image) {
	for _, c := range e.components {
		c.Draw(screen)
	}
}

type Component interface {
	Update(time.Duration)
	Draw(screen *ebiten.Image)
}

type ComponentWithElement interface {
	Component
	SetElement(e *Element)
}
