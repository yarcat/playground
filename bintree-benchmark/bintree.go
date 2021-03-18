package main

import "fmt"

type (
	node struct {
		l, r *node
		val  int
	}
	Tree struct {
		root *node
	}
)

func (n *node) visit(f func(int)) {
	if n != nil {
		n.l.visit(f)
		f(n.val)
		n.r.visit(f)
	}
}

func (t *Tree) Insert(x int) {
	p := &t.root
	for *p != nil {
		if (*p).val > x {
			p = &(*p).l
		} else {
			p = &(*p).r
		}
	}
	*p = &node{val: x}
}

func (t Tree) ForEach(f func(int)) {
	t.root.visit(f)
}

func main() {
	var t Tree
	for _, x := range []int{5, 7, 2, 1, 9, 6} {
		t.Insert(x)
	}
	var res []int
	t.ForEach(func(x int) { res = append(res, x) })
	fmt.Println(res)
}
