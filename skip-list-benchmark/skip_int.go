package main

type (
	SkipListInt struct {
		next []*skipListNodeInt
	}
	skipListNodeInt struct {
		next  []*skipListNodeInt
		value int
	}
)

func (l *SkipListInt) Init(levels int) *SkipListInt {
	l.next = make([]*skipListNodeInt, levels)
	return l
}

func (l *SkipListInt) Insert(value int) {
	prev := make([]**skipListNodeInt, l.Levels())
	for next, level := l.next, l.Levels()-1; level >= 0; level-- {
		for next[level] != nil && next[level].value < value {
			next = next[level].next
		}
		prev[level] = &next[level]
	}
	n := &skipListNodeInt{
		next:  make([]*skipListNodeInt, randLevel(l.Levels())),
		value: value,
	}
	for level := len(n.next) - 1; level >= 0; level-- {
		*prev[level], n.next[level] = n, *prev[level]
	}
}

func (l *SkipListInt) Levels() int { return len(l.next) }
