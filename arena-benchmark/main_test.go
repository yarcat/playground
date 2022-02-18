package main

import (
	"math/rand"
	"testing"
	"time"
)

var items []int

func init() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1_000; i++ {
		items = append(items, rand.Int())
	}
}

func benchmark(b *testing.B, newTree func() *tree) {
	for n := 0; n < b.N; n++ {
		t := newTree()
		for _, x := range items {
			t.Insert(x)
		}
	}
}

func BenchmarkNew(b *testing.B) {
	benchmark(b, func() *tree {
		return &tree{New: func(x int) *node { return &node{data: x} }}
	})
}

func BenchmarkArena(b *testing.B) {
	benchmark(b, func() *tree {
		a := arena{
			arena: make([]node, len(items)),
			p:     -1,
		}
		return &tree{
			New: func(x int) *node {
				n := a.New()
				n.data = x
				return n
			}}
	})
}
