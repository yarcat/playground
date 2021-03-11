package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type (
	LessFn   func(a, b interface{}) bool
	SkipList struct {
		next []*skipListNode
		less LessFn
	}
	skipListNode struct {
		next  []*skipListNode
		value interface{}
	}
)

func NewSkipList(levels int, less LessFn) *SkipList {
	return new(SkipList).Init(levels, less)
}

func (l *SkipList) Init(levels int, less LessFn) *SkipList {
	l.next = make([]*skipListNode, levels)
	l.less = less
	return l
}

func (l *SkipList) Insert(value interface{}) (hops int) {
	prev := make([]**skipListNode, l.Levels())
	for next, level := l.next, l.Levels()-1; level >= 0; level-- {
		for next[level] != nil && l.less(next[level].value, value) {
			next = next[level].next
			hops++
		}
		prev[level] = &next[level]
	}
	n := &skipListNode{
		next:  make([]*skipListNode, randLevel(l.Levels())),
		value: value,
	}
	for level := len(n.next) - 1; level >= 0; level-- {
		*prev[level], n.next[level] = n, *prev[level]
	}
	return
}

func (l *SkipList) Lookup(value interface{}) (interface{}, bool, int) {
	hops := 1
	for next, level := l.next, l.Levels()-1; level >= 0; level-- {
		for next[level] != nil && l.less(next[level].value, value) {
			next = next[level].next
			hops++
		}
		if next[level] != nil && !l.less(value, next[level].value) {
			return next[level].value, true, hops
		}
	}
	return nil, false, hops
}

func (l *SkipList) Levels() int { return len(l.next) }

func main() {
	rand.Seed(time.Now().UnixNano())
	var l SkipList
	l.Init(10, func(a, b interface{}) bool { return a.(int) < b.(int) })
	for i := 0; i < 20; i++ {
		x := rand.Intn(100)
		hops := l.Insert(x)
		fmt.Println(x, hops)
	}
	fmt.Println(debug(l))
	for _, n := range []int{
		-1, 0, 1, 24, 28, 32, 35, 39, 45, 62, 87, 88, 89, 90, 95,
	} {
		_, ok, hops := l.Lookup(n)
		fmt.Printf("Lookup(%d) = %v (hops = %d)\n", n, ok, hops)
	}

}

func randLevel(max int) (level int) {
	for level = 1; level < max; level++ {
		if rand.Float32() > 0.5 {
			break
		}
	}
	return
}

func debug(l SkipList) string {
	lines := make([]strings.Builder, l.Levels())
	for i := range lines {
		fmt.Fprintf(&lines[i], "%-2d", l.Levels()-1-i)
	}
	prev := make([]*skipListNode, l.Levels())
	copy(prev, l.next)
	for node := l.next[0]; node != nil; node = node.next[0] {
		for i := 0; i < l.Levels(); i++ {
			level := len(prev) - 1 - i
			k := &prev[level]
			if *k == node {
				lines[i].WriteString("->")
				lines[i].WriteString(fmt.Sprintf("[%-2v]", node.value))
				*k = node.next[level]
			} else {
				lines[i].WriteString("------")
			}
		}
	}

	ll := make([]string, len(lines))
	for i := range lines {
		lines[i].WriteString("->[nil]")
		ll[i] = lines[i].String()
	}
	return strings.Join(ll, "\n")
}
