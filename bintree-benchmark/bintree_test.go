package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkIntTreeIterativeInsert(b *testing.B) { benchmarkIntTreeInsert(b, iterativeInsert) }
func BenchmarkIntTreeIterativeInsertParentPtr(b *testing.B) {
	benchmarkIntTreeInsert(b, iterativeInsertParentPtr)
}

func BenchmarkIntTreeRecursiveInsert(b *testing.B) { benchmarkIntTreeInsert(b, recursiveInsert) }
func BenchmarkIntTreeRecursiveInsertParentPtr(b *testing.B) {
	benchmarkIntTreeInsert(b, recursiveInsertParentPtr)
}

func BenchmarkTreeIterativeInsertParentPtr(b *testing.B) { benchmarkTreeInsert(b, (*Tree).Insert) }

func iterativeInsert(t *IntTree, val int) {
	if t.root == nil {
		t.root = &intNode{val: val}
		return
	}
	current := t.root
	for {
		if val < current.val {
			if current.l == nil {
				current.l = &intNode{val: val}
				return
			} else {
				current = current.l
			}
		} else {
			if current.r == nil {
				current.r = &intNode{val: val}
				return
			} else {
				current = current.r
			}
		}
	}
}

func iterativeInsertParentPtr(t *IntTree, val int) {
	p := &t.root
	for *p != nil {
		if val < (*p).val {
			p = &(*p).l
		} else {
			p = &(*p).r
		}
	}
	*p = &intNode{val: val}
}

func recursiveNodeInsert(n *intNode, val int) {
	if val < n.val {
		if n.l == nil {
			n.l = &intNode{val: val}
		} else {
			recursiveNodeInsert(n.l, val)
		}
	} else {
		if n.r == nil {
			n.r = &intNode{val: val}
		} else {
			recursiveNodeInsert(n.r, val)
		}
	}
}

func recursiveInsert(t *IntTree, val int) {
	if t.root == nil {
		t.root = &intNode{val: val}
	} else {
		recursiveNodeInsert(t.root, val)
	}
}

func recursiveNodeInsertParentPtr(p **intNode, val int) {
	if *p == nil {
		*p = &intNode{val: val}
	} else if val < (*p).val {
		recursiveNodeInsertParentPtr(&(*p).l, val)
	} else {
		recursiveNodeInsertParentPtr(&(*p).r, val)
	}
}

func recursiveInsertParentPtr(t *IntTree, val int) {
	recursiveNodeInsertParentPtr(&t.root, val)
}

func init() { rand.Seed(time.Now().UnixNano()) }

var (
	data100  = genData(100)
	data1k   = genData(1_000)
	data10k  = genData(10_000)
	data100k = genData(100_000)
)

func genData(n int) (out []int) {
	out = make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = rand.Intn(n)
	}
	return
}

func benchmarkIntTreeInsert(b *testing.B, insert func(*IntTree, int)) {
	benchmarkInsert(b, func(data []int) {
		var t IntTree
		for _, x := range data {
			insert(&t, x)
		}
	})
}

func benchmarkTreeInsert(b *testing.B, insert func(*Tree, interface{})) {
	benchmarkInsert(b, func(data []int) {
		t := NewTree(func(a, b interface{}) bool { return a.(int) < b.(int) })
		for _, x := range data {
			insert(&t, x)
		}
	})
}

func benchmarkInsert(b *testing.B, test func([]int)) {
	for _, data := range [][]int{
		data100,
		data1k,
		data10k,
		data100k,
	} {
		b.Run(fmt.Sprint(len(data)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				test(data)
			}
		})
	}
}
