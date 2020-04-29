package simulation

import (
	"github.com/yarcat/playground/geom/body"
	"github.com/yarcat/playground/geom/contrib/container/orderedmap"
)

type bodyIterator struct {
	*orderedmap.Iterator
}

func newBodyIterator(bodies *orderedmap.OrderedMap) *bodyIterator {
	return &bodyIterator{bodies.Iter()}
}

func (it bodyIterator) Value() *body.Body {
	b, _ := it.Iterator.Value()
	return b.(*body.Body)
}
