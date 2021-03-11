package main

import (
	"sync"
)


type (
	SkipListInt struct {
		next []*skipListNodeInt
	}
	skipListNodeInt struct {
		next  []*skipListNodeInt
		value int
		a [10]*skipListNodeInt
	}
)

func (l *SkipListInt) Init(levels int) *SkipListInt {
	l.next = make([]*skipListNodeInt, levels)
	return l
}

var pool = sync.Pool{
	New: func() interface{} {
		return make([]**skipListNodeInt, 100)
	},
}

func (l *SkipListInt) Insert(value int) {
	// prev := pool.Get().([]**skipListNodeInt)
	prev := make([]**skipListNodeInt, 10)
	for next, level := l.next, l.Levels()-1; level >= 0; level-- {
		for next[level] != nil && next[level].value < value {
			next = next[level].next
		}
		prev[level] = &next[level]
	}
	n := &skipListNodeInt{value: value}
	n.next = n.a[:randLevel(l.Levels())]
	for level := len(n.next) - 1; level >= 0; level-- {
		*prev[level], n.next[level] = n, *prev[level]
	}
	// pool.Put(prev)
}

func (l *SkipListInt) Levels() int { return len(l.next) }
