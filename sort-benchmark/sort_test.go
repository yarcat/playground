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
	ordered = count(1_000)
	random  = shuffled(ordered)
	rev     = reversed(ordered)
)

func BenchmarkQuickOrdered(b *testing.B) { benchmark(b, ordered, qsorted) }
func BenchmarkQuickRandom(b *testing.B)  { benchmark(b, random, qsorted) }
func BenchmarkQuickRev(b *testing.B)     { benchmark(b, rev, qsorted) }

func BenchmarkMergeOrdered(b *testing.B) { benchmark(b, ordered, msorted) }
func BenchmarkMergeRandom(b *testing.B)  { benchmark(b, random, msorted) }
func BenchmarkMergeRev(b *testing.B)     { benchmark(b, rev, msorted) }

func benchmark(b *testing.B, s []int, sort func(list) list) {
	for i := 0; i < b.N; i++ {
		var l list
		for j := 0; j < len(s); j++ {
			l.pushBack(s[j])
		}
		sort(l)
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
