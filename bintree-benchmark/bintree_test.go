package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkIterativeInsert(b *testing.B)          { benchmarkInsert(b, iterativeInsert) }
func BenchmarkIterativeInsertParentPtr(b *testing.B) { benchmarkInsert(b, iterativeInsertParentPtr) }

func BenchmarkRecursiveInsert(b *testing.B)          { benchmarkInsert(b, recursiveInsert) }
func BenchmarkRecursiveInsertParentPtr(b *testing.B) { benchmarkInsert(b, recursiveInsertParentPtr) }

func iterativeInsert(t *Tree, val int) {
	if t.root == nil {
		t.root = &node{val: val}
		return
	}
	current := t.root
	for {
		if val < current.val {
			if current.l == nil {
				current.l = &node{val: val}
				return
			} else {
				current = current.l
			}
		} else {
			if current.r == nil {
				current.r = &node{val: val}
				return
			} else {
				current = current.r
			}
		}
	}
}

func iterativeInsertParentPtr(t *Tree, val int) {
	p := &t.root
	for *p != nil {
		if val < (*p).val {
			p = &(*p).l
		} else {
			p = &(*p).r
		}
	}
	*p = &node{val: val}
}

func recursiveNodeInsert(n *node, val int) {
	if val < n.val {
		if n.l == nil {
			n.l = &node{val: val}
		} else {
			recursiveNodeInsert(n.l, val)
		}
	} else {
		if n.r == nil {
			n.r = &node{val: val}
		} else {
			recursiveNodeInsert(n.r, val)
		}
	}
}

func recursiveInsert(t *Tree, val int) {
	if t.root == nil {
		t.root = &node{val: val}
	} else {
		recursiveNodeInsert(t.root, val)
	}
}

func recursiveNodeInsertParentPtr(p **node, val int) {
	if *p == nil {
		*p = &node{val: val}
	} else if val < (*p).val {
		recursiveNodeInsertParentPtr(&(*p).l, val)
	} else {
		recursiveNodeInsertParentPtr(&(*p).r, val)
	}
}

func recursiveInsertParentPtr(t *Tree, val int) {
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

func benchmarkInsert(b *testing.B, insert func(*Tree, int)) {
	for _, data := range [][]int{
		data100,
		data1k,
		data10k,
		data100k,
	} {
		b.Run(fmt.Sprint(len(data)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				var t Tree
				for _, x := range data {
					insert(&t, x)
				}
			}
		})
	}
}
