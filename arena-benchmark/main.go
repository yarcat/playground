package main

import "fmt"

type node struct {
	l, r *node
	data int
}

type tree struct {
	root *node
	New  func(int) *node
}

func (t *tree) Insert(x int) {
	r := &t.root
	for *r != nil {
		if x < (*r).data {
			r = &(*r).l
		} else {
			r = &(*r).r
		}
	}
	*r = t.New(x)
}

func (t *tree) Traverse(f func(int)) {
	var traverse func(*node)
	traverse = func(r *node) {
		if r == nil {
			return
		}
		traverse(r.l)
		f(r.data)
		traverse(r.r)
	}
	traverse(t.root)
}

type arena struct {
	arena []node
	p     int
}

func (a *arena) New() *node {
	a.p++
	return &a.arena[a.p]
}

func main() {
	t1 := tree{
		New: func(x int) *node { return &node{data: x} },
	}
	a := arena{arena: make([]node, 100), p: -1}
	t2 := tree{
		New: func(x int) *node {
			n := a.New()
			n.data = x
			return n
		},
	}
	for _, x := range []int{10, 5, 2, 7, 1, 3, 6, 9} {
		t1.Insert(x)
		t2.Insert(x)
	}
	t1.Traverse(func(x int) { fmt.Println(x) })
	t2.Traverse(func(x int) { fmt.Println(x) })
}
