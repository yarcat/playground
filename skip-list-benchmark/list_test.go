package main

// Запуск из командной строки:
//         go test -bench . benchmark_sample_test.go
// В данном случае benchmark_sample_test.go - это имя файла,
// в котором находится данная программа. Имя тестируемого
// файла обязательно должно заканиваться на _test

import (
	"math"
	"math/rand"
	"testing"
)

const (
	PlusInfinity  = math.MaxInt64
	MinusInfinity = math.MinInt64
)

type (
	list struct {
		head *lmnt
	}
	lmnt struct {
		x    int
		next *lmnt
	}
)

var d [2048]int

func init() {
	for i, _ := range d {
		d[i] = rand.Intn(1000000000)
	}
}

func SortLinkedList(data []int) {
	l := NewList()
	for _, x := range data {
		l.Insert2(x)
	}
}

func SortSkipList(data []int) {
	var l SkipList
	l.Init(10, func(a, b interface{}) bool { return a.(int) < b.(int) })
	for _, x := range data {
		l.Insert(x)
	}
}

func SortSkipListInt(data []int) {
	var l SkipListInt
	l.Init(10)
	for _, x := range data {
		l.Insert(x)
	}
}

func benchmarkSort(b *testing.B, sort func([]int)) {
	for _, tc := range []struct {
		name string
		data []int
	}{
		{"256", d[:256]},
		{"512", d[:512]},
		{"1024", d[:1024]},
		{"2048", d[:2048]},
	} {
		b.Run(tc.name, func(b *testing.B) {
			d := make([]int, len(tc.data))
			for i := 0; i < b.N; i++ {
				copy(d, tc.data)
				sort(d)
			}
		})
	}
}

func BenchmarkSortLinkedList(b *testing.B)  { benchmarkSort(b, SortLinkedList) }
func BenchmarkSortSkipList(b *testing.B)    { benchmarkSort(b, SortSkipList) }
func BenchmarkSortSkipListInt(b *testing.B) { benchmarkSort(b, SortSkipListInt) }

func NewList() list {
	return list{&lmnt{MinusInfinity, &lmnt{PlusInfinity, nil}}}
}

func (s *list) Insert2(num int) {
	runner := (*s).head
	runner2 := (*runner).next
	for (*runner2).x < num {
		runner, runner2 = runner2, (*runner2).next
	}
	(*runner).next = &lmnt{num, runner2}
}
