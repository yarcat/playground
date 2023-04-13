package object

import (
	"test1/vec"
)

type Object struct {
	pos vec.Vec
	v   vec.Vec
}

func New(p, v vec.Vec) Object {
	return Object{pos: p, v: v}
}

func (o Object) P() vec.Vec { return o.pos }
func (o Object) V() vec.Vec { return o.v }

func (o Object) Move(dtSec float64) Object {
	return Object{
		pos: o.pos.Add(o.v.Scale(dtSec)),
		v:   o.v,
	}
}

func (o Object) Rotate(a float64) Object {
	return Object{pos: o.pos, v: o.v.Rotate(a)}
}
