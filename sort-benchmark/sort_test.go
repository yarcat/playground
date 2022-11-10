package main

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	ordered = count(10_000)
	random  = shuffled(ordered)
	rev     = reversed(ordered)
)

func oldAlgo(sort func(list) list) func(*testing.B, []int) {
	return func(b *testing.B, s []int) {
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			var l list
			for _, x := range s {
				l.pushBack(x)
			}
			b.StartTimer()
			sort(l)
		}
	}
}

func newAlgo(sort func(*gennode[int], func(a, b int) bool) *gennode[int]) func(*testing.B, []int) {
	return func(b *testing.B, ints []int) {
		for n := 0; n < b.N; n++ {
			b.StopTimer()
			l := genlist(ints...)
			b.StartTimer()
			sort(l, genless[int])
		}
	}
}
func BenchmarkSort(b *testing.B) {
	genqsort := func(n *gennode[int], less func(a, b int) bool) *gennode[int] {
		n, _ = genqsort(n, less)
		return n
	}
	for _, data := range []struct {
		name string
		in   []int
	}{
		{"ord", ordered},
		{"rng", random},
		{"rev", rev},
	} {
		b.Run(data.name, func(b *testing.B) {
			for _, bc := range []struct {
				name string
				f    func(*testing.B, []int)
			}{
				{"old qsort", oldAlgo(qsorted)},
				{"new qsort", newAlgo(genqsort)},
				{"old msort", oldAlgo(msorted)},
				{"new msort", newAlgo(genmsort[int])},
			} {
				b.Run(bc.name, func(b *testing.B) { bc.f(b, data.in) })
			}
		})
	}
}

func count(n int) (s []int) {
	s = make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	return
}

func shuffled(s []int) (ss []int) {
	ss = make([]int, len(s))
	copy(ss, s)
	rand.Shuffle(len(ss), func(i, j int) {
		ss[i], ss[j] = ss[j], ss[i]
	})
	return
}

func reversed(s []int) (ss []int) {
	ss = make([]int, len(s))
	for i, j := 0, len(s)-1; i < len(s); i, j = i+1, j-1 {
		ss[j] = s[i]
	}
	return
}
