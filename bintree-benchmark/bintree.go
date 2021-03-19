package main

import "fmt"

type (
	intNode struct {
		l, r *intNode
		val  int
	}
	IntTree struct {
		root *intNode
	}

	node struct {
		l, r *node
		val  interface{}
	}
	Tree struct {
		root *node
		less func(a, b interface{}) bool
	}
)

func NewTree(less func(a, b interface{}) bool) Tree { return Tree{less: less} }

func (n *intNode) visit(f func(int)) {
	if n != nil {
		n.l.visit(f)
		f(n.val)
		n.r.visit(f)
	}
}

func (n *node) visit(f func(interface{})) {
	if n != nil {
		n.l.visit(f)
		f(n.val)
		n.r.visit(f)
	}
}

func (t *IntTree) Insert(x int) {
	p := &t.root
	for *p != nil {
		if x < (*p).val {
			p = &(*p).l
		} else {
			p = &(*p).r
		}
	}
	*p = &intNode{val: x}
}

func (t *Tree) Insert(val interface{}) {
	p := &t.root
	for *p != nil {
		if t.less(val, (*p).val) {
			p = &(*p).l
		} else {
			p = &(*p).r
		}
	}
	*p = &node{val: val}
}

func (t IntTree) ForEach(f func(int))      { t.root.visit(f) }
func (t Tree) ForEach(f func(interface{})) { t.root.visit(f) }

func main() {
	var it IntTree
	t := NewTree(func(a, b interface{}) bool { return a.(int) < b.(int) })

	for _, x := range []int{5, 7, 2, 1, 9, 6} {
		it.Insert(x)
		t.Insert(x)
	}
	var ires, res []int
	it.ForEach(func(x int) { ires = append(ires, x) })
	t.ForEach(func(x interface{}) { res = append(res, x.(int)) })
	fmt.Println(ires, res)
}
