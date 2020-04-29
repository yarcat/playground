package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type boxBldr struct {
	w, h float64
	r    vector
}

func (bb *boxBldr) move(v vector) *boxBldr {
	bb.r = bb.r.add(v)
	return bb
}

func (bb *boxBldr) moveXY(dx, dy float64) *boxBldr {
	return bb.move(vector{dx, dy})
}

func (bb *boxBldr) build() shape {
	halfW, halfH := bb.w/2, bb.h/2
	return &box{
		pts:   aabb{vector{-halfW, -halfH}, vector{halfW, halfH}},
		r:     bb.r,
		theta: 0,
	}
}

type cirBldr struct {
	r, x, y float64
}

func (cb *cirBldr) build() shape {
	return &circle{rad: cb.r, r: vector{cb.x, cb.y}}
}

type shape = translatingGeomWithAabb

type shapeBuilder interface {
	build() shape
}

func main() {
	// seed := time.Now().Unix()
	// fmt.Println(seed)
	rand.Seed(1587571564)

	wrld := newWorld()

	midscreen := vector{screenWidth / 2, screenHeight / 2}
	for _, d := range []struct {
		builder shapeBuilder
		vx, vy  float64
		m       float64
	}{
		{builder: (&boxBldr{w: 50, h: 200}).move(midscreen), m: 20},
		{builder: (&boxBldr{w: 20, h: 150}).move(midscreen).moveXY(-100, 100), m: 10},
		{builder: (&boxBldr{w: 20, h: 150}).move(midscreen).moveXY(150, -30), vx: -1, m: 10},
		{builder: &cirBldr{r: 20, x: 30, y: 100}, vx: 3, vy: 1, m: 2},
		{builder: &cirBldr{r: 20, x: 30, y: 300}, vx: 3, m: 2},
		{builder: &cirBldr{r: 50, x: 120, y: 250}, m: 6},
		// Screen bounding box.
		{builder: (&boxBldr{w: screenWidth, h: 10}).moveXY(screenWidth/2, 5)},
		{builder: (&boxBldr{w: screenWidth, h: 10}).moveXY(screenWidth/2, screenHeight-5)},
		{builder: (&boxBldr{w: 10, h: screenHeight - 30}).moveXY(5, screenHeight/2)},
		{builder: (&boxBldr{w: 10, h: screenHeight - 30}).moveXY(screenWidth-5, screenHeight/2)},
	} {
		g := d.builder.build()
		wrld.addSprite(&drawable{g, color.RGBA{0xff, 0xff, 0xff, 0xff}})
		rb := &rigidBody{g, vector{d.vx, d.vy}, d.m}
		wrld.addRigidBody(rb)
	}

	err := ebiten.Run(wrld.update, screenWidth, screenHeight, 1, "Test")
	if err != nil {
		log.Fatal(err)
	}
}

func findClosest(min, max, curr float64) float64 {
	if curr < min {
		return min
	}
	if curr > max {
		return max
	}
	return curr
}

func collideBoxWithCircle(b *box, c *circle) (n vector, pen float64, ok bool) {
	diag := b.pts[1].sub(b.pts[0])
	halfW, halfH := diag[0]/2, diag[1]/2
	bpos, _ := b.pos()
	cpos, _ := c.pos()
	dist := cpos.sub(bpos)
	closest := vector{
		findClosest(-halfW, halfW, dist[0]),
		findClosest(-halfH, halfH, dist[1]),
	}
	// TODO(yarcat): Look into the case when c's center is inside of b.
	norm := dist.sub(closest)
	if norm.len2() > c.rad*c.rad {
		// No collision.
		return
	}
	pen = math.Abs(c.rad - cpos.sub(closest.add(bpos)).len())
	return norm.norm(), pen, true
}

type body interface {
	// TODO(yarcat): We must come up with a better interface that would allow us
	// to choose the vector (e.g. speed, position, etc) we want to operate on.
	translate(d vector)
	v() vector
	translatev(d vector)
	mass() float64
}

type bodyState struct {
	dv   vector
	dpos vector
}

type bodyStates map[body]*bodyState

type collisionState struct {
	n   vector
	pen float64
}

type collisionStates map[body]map[body]*collisionState

type collisionResolver struct {
	bodies     bodyStates
	collisions collisionStates
}

func newCollisionResolver() *collisionResolver {
	return &collisionResolver{make(bodyStates), make(collisionStates)}
}

func (cr *collisionResolver) update(w *world) {
	for b1, collisions := range cr.collisions {
		for b2, collisionState := range collisions {
			cr.resolveCollision(b1, b2, collisionState)
		}
	}
	for b, state := range cr.bodies {
		// Push bodies slightly away.
		b.translate(state.dpos)
		b.translatev(state.dv)
	}
	cr.bodies = make(bodyStates)
	cr.collisions = make(collisionStates)
}

func (cr *collisionResolver) getBodyState(b body) *bodyState {
	state, ok := cr.bodies[b]
	if !ok {
		state = &bodyState{}
		cr.bodies[b] = state
	}
	return state
}

func (cr *collisionResolver) addCollision(b1, b2 body, n vector, pen float64) {
	collisions, ok := cr.collisions[b1]
	if !ok {
		// b1 doesn't exist, try b2.
		if collisions, ok = cr.collisions[b2]; ok {
			// For convenience, arrange bodies so that keys are [b1][b2].
			b1, b2 = b2, b1
		} else {
			collisions = make(map[body]*collisionState)
			cr.collisions[b1] = collisions
		}
	}
	// This is where the convenience comes handy. We know that we have to check for b2 here.
	if _, ok := collisions[b2]; !ok {
		collisions[b2] = &collisionState{n: n, pen: pen}
	}
}

func invertMass(m float64) float64 {
	if m == 0 {
		return 0
	}
	return 1 / m
}
func (cr *collisionResolver) resolveCollision(b1, b2 body, collision *collisionState) {
	const (
		Cr = 1
	)
	invm1, invm2 := invertMass(b1.mass()), invertMass(b2.mass())

	u1 := b1.v().dotProduct(collision.n)
	u2 := b2.v().dotProduct(collision.n)
	du := u2 - u1

	if du > 0 {
		// Two bodies are moving in the same direction. They will split.
		return
	}

	j := du * (1 + Cr) / (invm1 + invm2)

	if !iszero(j) {
		const (
			percent = 0.2
			minpen  = 0.02
		)
		_, k := minmax(0, collision.pen-minpen)
		k = k / (invm1 + invm2) * percent
		correction := collision.n.scale(k)

		state1 := cr.getBodyState(b1)
		state1.dv = state1.dv.add(collision.n.scale(j * invm1))
		state1.dpos = state1.dpos.sub(correction.scale(invm2))

		state2 := cr.getBodyState(b2)
		state2.dv = state2.dv.sub(collision.n.scale(j * invm2))
		state2.dpos = state2.dpos.add(correction.scale(invm2))
	}
}

type bodyAdapter struct {
	rb *rigidBody
}

func (ba bodyAdapter) mass() float64 {
	return ba.rb.m
}

func (ba bodyAdapter) v() vector {
	return ba.rb.v
}

func (ba bodyAdapter) translate(d vector) {
	ba.rb.shape.translate(d)
}

func (ba bodyAdapter) translatev(d vector) {
	ba.rb.v = ba.rb.v.add(d)
}

type translatingGeomWithAabb interface {
	geomWithAabb
	translate(vector)
}

type rigidBody struct {
	shape translatingGeomWithAabb
	v     vector
	m     float64
}

func (s rigidBody) update(w *world) {
	s.shape.translate(s.v)
}

type geomWithAabb interface {
	geom
	aabb() aabb
}

type aabbOf struct {
	geomWithAabb
}

func (ao *aabbOf) accept(a geomAcceptor) {
	r, _ := ao.pos()
	b := &box{
		pts:   ao.aabb(),
		r:     r,
		theta: 0,
	}
	b.accept(a)
}

type geomCollider struct {
	g      geomWithAabb
	fn     func(g1, g2 geom, n vector, pen float64)
	fnAabb func()
}

type circleCollider struct {
	*geomCollider
}

func (c *circleCollider) check(sc specificCollider) {
	sc.checkCircle(c.g.(*circle))
}

func (c *circleCollider) checkCircle(cl *circle) {
	// Checking AABB intersections is cheap.
	bb1, bb2 := c.g.aabb(), cl.aabb()
	if !bb1.intersects(bb2) {
		return
	}
	c.fnAabb()
	p1, _ := c.g.(*circle).pos()
	p2, _ := cl.pos()
	r1r2 := c.g.(*circle).rad + cl.rad
	n := p2.sub(p1)
	dist2 := n.len2()
	if r1r2*r1r2 >= dist2 {
		c.fn(c.g, cl, n.norm(), r1r2-math.Sqrt(dist2))
	}
}

func (c *circleCollider) checkBox(b *box) {
	// Checking AABB intersections is cheap.
	bb1, bb2 := c.g.aabb(), b.aabb()
	if !bb1.intersects(bb2) {
		return
	}
	c.geomCollider.fnAabb()

	cir := c.g.(*circle)
	if n, pen, ok := collideBoxWithCircle(b, cir); ok {
		c.fn(cir, b, n.scale(-1), pen)
	}
}

type boxCollider struct {
	*geomCollider
}

func (c *boxCollider) check(sc specificCollider) {
	sc.checkBox(c.g.(*box))
}

func (c *boxCollider) checkBox(b *box) {
	// Checking AABB intersections is cheap.
	bb1, bb2 := c.g.aabb(), b.aabb()
	if !bb1.intersects(bb2) {
		return
	}
	c.geomCollider.fnAabb()

	if n, pen, ok := collideAabbAndAabb(c.g.(*box), b); ok {
		c.fn(c.g.(*box), b, n, pen)
	}
}

func (c *boxCollider) checkCircle(cr *circle) {
	// Checking AABB intersections is cheap.
	bb1, bb2 := c.g.aabb(), cr.aabb()
	if !bb1.intersects(bb2) {
		return
	}
	c.geomCollider.fnAabb()

	b := c.g.(*box)
	if n, pen, ok := collideBoxWithCircle(b, cr); ok {
		c.fn(b, cr, n, pen)
	}
}

func collideAabbAndAabb(b1 *box, b2 *box) (n vector, pen float64, ok bool) {
	aabb1 := b1.aabb()
	aabb2 := b2.aabb()

	if !aabb1.intersects(aabb2) {
		return
	}

	r1, _ := b1.pos()
	r2, _ := b2.pos()
	n = r2.sub(r1)

	ext1 := b1.pts[1].sub(b1.pts[0])
	ext2 := b2.pts[1].sub(b2.pts[0])

	xoverlap := (ext1[0]+ext2[0])/2 - math.Abs(n[0])
	if xoverlap < 0 {
		return
	}
	yoverlap := (ext1[1]+ext2[1])/2 - math.Abs(n[1])
	if yoverlap < 0 {
		return
	}

	if yoverlap < xoverlap {
		pen = yoverlap
		if n[1] < 0 {
			n = vector{0, -1}
		} else {
			n = vector{0, 1}
		}
	} else {
		pen = xoverlap
		if n[0] < 0 {
			n = vector{-1, 0}
		} else {
			n = vector{1, 0}
		}
	}
	return n, pen, true
}

func (b *box) changeOrigin(r vector, theta float64) *box {
	return &box{
		pts:   b.pts,
		r:     b.r.sub(r).rotate(-theta),
		theta: b.theta - theta,
	}
}

type rotator struct {
	g     geom
	theta float64 // per tick
}

func (r *rotator) update(w *world) {
	r.g.rotate(r.theta)
}

const epsilon0 = 1e-10

func iszero(f float64) bool {
	if f > 0 {
		return f < epsilon0
	}
	return f > -epsilon0
}

type vector [2]float64

func (v vector) x() float64 {
	return v[0]
}

func (v vector) y() float64 {
	return v[1]
}

func (v vector) dotProduct(v2 vector) float64 {
	return v[0]*v2[0] + v[1]*v2[1]
}

func (v vector) sub(v2 vector) vector {
	return vector{v[0] - v2[0], v[1] - v2[1]}
}

func (v vector) norm() vector {
	return v.scale(1 / v.len())
}

func (v vector) len() float64 {
	return math.Sqrt(v.len2())
}

func (v vector) len2() float64 {
	return v[0]*v[0] + v[1]*v[1]
}

func (v vector) add(v2 vector) vector {
	return vector{v[0] + v2[0], v[1] + v2[1]}
}

func (v vector) scale(k float64) vector {
	for i := range v {
		v[i] *= k
	}
	return v
}

func (v vector) apply(fn func(float64) float64) vector {
	for i, val := range v {
		v[i] = fn(val)
	}
	return v
}

func (v vector) rotate(theta float64) vector {
	cos, sin := math.Cos(theta), math.Sin(theta)
	return vector{
		v[0]*cos - v[1]*sin,
		v[0]*sin + v[1]*cos,
	}
}

type geom interface {
	pos() (r vector, theta float64)
	rotate(theta float64)
	accept(geomAcceptor)
}

type geomAcceptor interface {
	box(*box)
	circle(*circle)
}

type withCollision struct {
	normal, collision   *drawable
	ticksSinceCollision int
}

// TODO(yarcat): Probably need enter/exit collision callbacks.
func (d *withCollision) onCollision() {
	d.ticksSinceCollision = -1
}

func (d *withCollision) update(w *world) {
	d.ticksSinceCollision++
}

func (d *withCollision) draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if d.ticksSinceCollision == 0 {
		d.collision.draw(screen, op)
	} else {
		d.normal.draw(screen, op)
	}
}

type drawable struct {
	geom
	color color.RGBA
}

func (d *drawable) draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	disp := drawDispatcher{screen, op, d.color}
	d.geom.accept(&disp)
}

type drawDispatcher struct {
	screen *ebiten.Image
	op     *ebiten.DrawImageOptions
	color  color.RGBA
}

func (d *drawDispatcher) box(b *box) {
	diag := b.pts[1].sub(b.pts[0]).apply(math.Abs)
	img, _ := ebiten.NewImage(
		/* width */ int(diag[0]),
		/* height */ int(diag[1]),
		0)

	const thickness = 2

	ebitenutil.DrawRect(img, 0, 0, diag[0], thickness, d.color)
	ebitenutil.DrawRect(img, 0, diag[1]-thickness, diag[0], diag[1], d.color)
	ebitenutil.DrawRect(img, 0, 0, thickness, diag[1], d.color)
	ebitenutil.DrawRect(img, diag[0]-thickness, 0, diag[0], diag[1], d.color)

	halfDiag := diag.scale(0.5)
	d.op.GeoM.Translate(-halfDiag[0], -halfDiag[1])
	d.op.GeoM.Rotate(b.theta)
	d.op.GeoM.Translate(b.r[0], b.r[1])
	d.screen.DrawImage(img, d.op)
	img.Dispose()
}

type pixel struct {
	*ebiten.Image
	op ebiten.DrawImageOptions
}

func (p *pixel) at(img *ebiten.Image, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	img.DrawImage(p.Image, op)
}

func (d *drawDispatcher) circle(c *circle) {
	D := int(c.rad*2 + 0.5)
	r := D / 2
	img, _ := ebiten.NewImage(D, D, 0)

	pixelImg, _ := ebiten.NewImage(2, 2, 0)
	pixelImg.Fill(d.color)
	pixel := &pixel{Image: pixelImg}

	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	x0, y0 := r, r

	for x > y {
		pixel.at(img, x0+x, y0+y)
		pixel.at(img, x0+y, y0+x)
		pixel.at(img, x0-y, y0+x)
		pixel.at(img, x0-x, y0+y)
		pixel.at(img, x0-x, y0-y)
		pixel.at(img, x0-y, y0-x)
		pixel.at(img, x0+y, y0-x)
		pixel.at(img, x0+x, y0-y)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
	pixel.Dispose()

	d.op.GeoM.Translate(c.r[0]-c.rad, c.r[1]-c.rad)
	d.screen.DrawImage(img, d.op)

	img.Dispose()
}

func minmax(a, b float64) (min, max float64) {
	if a <= b {
		return a, b
	}
	return b, a
}

// aabb stands for Axis-Alligned Bounding Box. The first element must contain
// min and the second max coordinates.
//       [1] max x; y
// +-------+
// |       |
// +-------+
// [0] min x; y
type aabb [2]vector

func (b1 aabb) intersects(b2 aabb) bool {
	const (
		min, max = 0, 1
		x, y     = 0, 1
	)
	return b1[min][x] <= b2[max][x] && b2[min][x] <= b1[max][x] &&
		b1[min][y] <= b2[max][y] && b2[min][y] <= b1[max][y]
}

type circle struct {
	rad float64
	r   vector
}

func (c *circle) translate(r vector) {
	c.r = c.r.add(r)
}

func (c circle) pos() (r vector, theta float64) {
	return c.r, theta
}

func (c circle) rotate(theta float64) {}

func (c *circle) accept(a geomAcceptor) {
	a.circle(c)
}

func (c *circle) aabb() aabb {
	return aabb{
		vector{-c.rad, -c.rad}.add(c.r),
		vector{c.rad, c.rad}.add(c.r),
	}
}

type box struct {
	pts   aabb
	r     vector
	theta float64
}

func (b *box) translate(v vector) {
	b.r = b.r.add(v)
}

func (b *box) aabb() aabb {
	// [1]      [2]
	// x0;y1  x1;y1
	//   +------+
	//   |      |
	//   +------+
	// x0;y0  x1;y0
	// [0]      [3]
	vertices := [4]vector{
		b.pts[0],
		vector{b.pts[0][0], b.pts[1][1]},
		b.pts[1],
		vector{b.pts[1][0], b.pts[0][1]},
	}

	maxx, maxy := math.Inf(-1), math.Inf(-1)
	minx, miny := math.Inf(1), math.Inf(1)

	for _, p := range vertices {
		p = p.rotate(b.theta)
		if p[0] > maxx {
			maxx = p[0]
		}
		if p[1] > maxy {
			maxy = p[1]
		}
		if p[0] < minx {
			minx = p[0]
		}
		if p[1] < miny {
			miny = p[1]
		}
	}

	return aabb{
		vector{minx, miny}.add(b.r),
		vector{maxx, maxy}.add(b.r),
	}
}

func (b *box) accept(g geomAcceptor) {
	g.box(b)
}

func (b *box) pos() (r vector, theta float64) {
	return b.r, b.theta
}

func (b *box) rotate(theta float64) {
	b.theta += theta
}

type sprite interface {
	draw(*ebiten.Image, *ebiten.DrawImageOptions)
}

type behavior interface {
	update(w *world)
}

type specificCollider interface {
	checkBox(*box)
	checkCircle(*circle)
}

type collider interface {
	check(specificCollider)
	specificCollider
}

type world struct {
	sprites   map[sprite]interface{}
	behaviors []behavior
	colliders map[collider]interface{}
	bodies    map[geom]body
	cr        *collisionResolver
}

func newWorld() *world {
	return &world{
		sprites:   make(map[sprite]interface{}),
		colliders: make(map[collider]interface{}),
		bodies:    make(map[geom]body),
		cr:        newCollisionResolver(),
	}
}

func (w *world) update(screen *ebiten.Image) error {
	for c1 := range w.colliders {
		for c2 := range w.colliders {
			if c1 != c2 {
				c2.check(c1)
			}
		}
	}
	for _, b := range w.behaviors {
		b.update(w)
	}
	w.cr.update(w)
	if ebiten.IsRunningSlowly() {
		return nil
	}
	for s := range w.sprites {
		op := &ebiten.DrawImageOptions{}
		s.draw(screen, op)
	}
	w.debugInfo(screen)
	return nil
}

func (w *world) addSprite(s sprite) {
	w.sprites[s] = nil
}

func (w *world) addBehavior(b behavior) {
	w.behaviors = append(w.behaviors, b)
}

func (w *world) addCollider(c collider) {
	w.colliders[c] = nil
}

func (w *world) addRigidBody(rb *rigidBody) {
	var g geom = rb.shape
	w.bodies[g] = &bodyAdapter{rb}
	w.addBehavior(rb)
	cb := &colliderBuilder{
		gc: &geomCollider{
			rb.shape,
			func(g1, g2 geom, n vector, pen float64) {
				w.cr.addCollision(w.bodies[g1], w.bodies[g2], n, pen)
			},
			func() {},
		}}
	rb.shape.accept(cb)
	w.addCollider(cb.collider)
}

type colliderBuilder struct {
	gc       *geomCollider
	collider collider
}

func (cb *colliderBuilder) box(b *box) {
	cb.collider = &boxCollider{cb.gc}
}

func (cb *colliderBuilder) circle(c *circle) {
	cb.collider = &circleCollider{cb.gc}
}

func (w *world) debugInfo(screen *ebiten.Image) {
	const template = `
FPS       : %.2f
Sprites   : %d
Behaviors : %d
Colliders : %d
`
	msg := fmt.Sprintf(
		template,
		ebiten.CurrentFPS(),
		len(w.sprites),
		len(w.behaviors),
		len(w.colliders),
	)
	ebitenutil.DebugPrint(screen, strings.TrimSpace(msg))
}
