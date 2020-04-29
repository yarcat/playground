package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func randomShootTicks() int {
	return 200 + rand.Int()%200

}

type world struct {
	explosiveBullet *ebiten.Image
	normalBullet    *ebiten.Image
	ticks           uint64
	shootTicks      int
	sprites         map[sprite]interface{}
}

func newWorld() (*world, error) {
	eb, _ := ebiten.NewImage(2, 2, 0 /* filter=defaultFilter */)
	eb.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	nb, _ := ebiten.NewImage(1, 1, 0 /* filter=defaultFilter */)
	nb.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	return &world{
		explosiveBullet: eb,
		normalBullet:    nb,
		sprites:         make(map[sprite]interface{}, 5000),
		shootTicks:      0,
	}, nil
}

func (w *world) add(s sprite) {
	w.sprites[s] = nil
}

var dt = time.Duration(16 * time.Millisecond)

func (w *world) Update(screen *ebiten.Image) error {
	deltaSec := dt.Seconds()
	for s := range w.sprites {
		if s.active() {
			s.update(w, deltaSec)
		}
		if !s.active() {
			delete(w.sprites, s)
		}
	}

	w.shootTicks--
	if w.shootTicks <= 0 {
		addExplosiveBullets(w)
		w.shootTicks = randomShootTicks()
	}
	w.ticks++
	return nil
}

func (w *world) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (w *world) Draw(screen *ebiten.Image) {
	for s := range w.sprites {
		draw(screen, s)
	}

	w.debug(screen)
}

func (w *world) debug(screen *ebiten.Image) {
	const template = `
FPS            : %.2f
Ticks          : %d
Sprites        : %d
Shoot in       : %d
`
	msg := fmt.Sprintf(
		template,
		ebiten.CurrentFPS(),
		w.ticks,
		len(w.sprites),
		w.shootTicks,
	)
	ebitenutil.DebugPrint(screen, strings.TrimSpace(msg))
}

type sprite interface {
	active() bool
	image() *ebiten.Image
	update(*world, float64)
	pos() (x, y float64)
	color() (r, g, b, a float64)
}

type bullet struct {
	point
	drawing
	enabled bool
	vx, vy  float64
}

const gravY = -200

func (s *bullet) active() bool {
	return s.y >= 0 && s.enabled
}

func (s *bullet) image() *ebiten.Image {
	return s.img
}

func (s *bullet) pos() (x, y float64) {
	return s.x, s.y
}

func (s *bullet) color() (r, g, b, a float64) {
	return s.r, s.g, s.b, s.a
}

func (s *bullet) update(w *world, ds float64) {
	if s.y < 0 {
		return
	}
	s.vy += gravY * ds
	s.x += s.vx * ds
	s.y += s.vy * ds
}

func draw(screen *ebiten.Image, s sprite) {
	opt := &ebiten.DrawImageOptions{}
	x, y := s.pos()
	opt.GeoM.Translate(x, screenHeight-y)
	opt.ColorM.Scale(s.color())
	screen.DrawImage(s.image(), opt)
}

type explosive struct {
	bullet
	ticks int
}

func (s *explosive) active() bool {
	return s.ticks >= 0
}

func (s *explosive) update(w *world, ds float64) {
	s.ticks--
	if s.ticks <= 0 {
		s.explode(w)
		return
	}
	s.bullet.update(w, ds)
}

func (s *explosive) explode(w *world) {
	sprites := 5 + rand.Int()%50
	for i := 0; i < sprites; i++ {
		strength := 200 + rand.Float64()*100
		angle := rand.Float64() * 2 * math.Pi

		var b bullet = s.bullet
		b.img = w.normalBullet
		b.vx += strength * math.Cos(angle)
		b.vy += strength * math.Sin(angle)

		w.add(&withShadow{&b, 10})
	}
}

type point struct {
	x, y float64
}

func (p *point) pos() (x, y float64) {
	return p.x, p.y
}

type drawing struct {
	img        *ebiten.Image
	r, g, b, a float64
}

func (d *drawing) image() *ebiten.Image {
	return d.img
}

func (d *drawing) color() (r, g, b, a float64) {
	return d.r, d.g, d.b, d.a
}

type withShadow struct {
	sprite
	ticks int
}

func (s *withShadow) update(w *world, dt float64) {
	if s.ticks <= 0 {
		return
	}
	x, y := s.pos()
	r, g, b, a := s.color()
	w.add(&shadow{
		point{x, y},
		drawing{s.image(), r, g, b, a},
		s.ticks,
		a / float64(s.ticks+1),
	})
	s.sprite.update(w, dt)
}

type shadow struct {
	point
	drawing
	ticks int
	da    float64
}

func (s *shadow) update(w *world, dt float64) {
	s.ticks--
	s.a -= s.da
}

func (s *shadow) active() bool {
	return s.ticks > 0
}

func addExplosiveBullets(w *world) {
	sprites := 2 + rand.Int()%10
	shotX := screenWidth * rand.Float64()
	for i := 0; i < sprites; i++ {
		s := &withShadow{
			&explosive{
				bullet{
					point: point{x: shotX + rand.Float64()*10},
					drawing: drawing{
						img: w.explosiveBullet,
						r:   rand.Float64(),
						g:   rand.Float64(),
						b:   rand.Float64(),
						a:   1,
					},
					enabled: true,
					vx:      -50 + rand.Float64()*100,
					vy:      400 + rand.Float64()*50,
				},
				rand.Int()%100 + 100,
			},
			/* ticks */ 20,
		}

		w.add(s)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	w, err := newWorld()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Firework!")
	if err == nil {
		err = ebiten.RunGame(w)
	}
	if err != nil {
		panic(err)
	}
}
