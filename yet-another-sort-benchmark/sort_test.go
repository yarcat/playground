package sort_test

import (
	"math/rand"
	"reflect"
	"testing"

	"golang.org/x/exp/constraints"
)

var sortFns = []struct {
	name string
	f    func([]int)
}{
	{"merge", SortMerge[int]},
	{"bubble", SortBubble[int]},
	{"insert", SortInsert[int]},
	{"insert2", SortInsert2[int]},
	{"insert bisect", SortInsertBisect[int]},
	{"quick", SortQuick[int]},
}

const minRng, maxRng = -10_000, 10_000

var s []int

func init() {
	// rand.Seed(time.Now().UnixMicro())
	s = make([]int, 10_000)
	for i := range s {
		s[i] = rand.Intn(maxRng-minRng+1) + minRng
	}
}

func benchmarkSort(b *testing.B, sort func([]int)) {
	s2 := make([]int, len(s))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		copy(s2, s)
		b.StartTimer()
		sort(s2)
	}
}

func BenchmarkSort(b *testing.B) {
	for _, bc := range sortFns {
		b.Run(bc.name, func(b *testing.B) { benchmarkSort(b, bc.f) })
	}
}

func testSort(t *testing.T, sort func([]int)) {
	for _, tc := range []struct {
		name        string
		input, want []int
	}{
		{"nil", nil, nil},
		{"empty", []int{}, []int{}},
		{"single", []int{123}, []int{123}},
		{"two_sorted", []int{1, 2}, []int{1, 2}},
		{"two_reversed", []int{2, 1}, []int{1, 2}},
		{"three_random", []int{2, 3, 1}, []int{1, 2, 3}},
		{"three_sorted", []int{1, 2, 3}, []int{1, 2, 3}},
		{"three_reversed", []int{3, 2, 1}, []int{1, 2, 3}},
		{"four_reversed", []int{4, 3, 2, 1}, []int{1, 2, 3, 4}},
		{"ten_random", []int{4, 3, 7, 1, 0, 5, 2, 6, 9, 8}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"ten_sorted", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"ten_reversed", []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"same", []int{123, 123, 123}, []int{123, 123, 123}},
		{"repeated", []int{123, 1, 123, 1, 123, 1}, []int{1, 1, 1, 123, 123, 123}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var got []int
			if tc.input != nil {
				got = make([]int, len(tc.input))
				copy(got, tc.input)
			}
			sort(got)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got = %v, want = %v", got, tc.want)
			}
		})
	}
}

func TestSort(t *testing.T) {
	for _, tc := range sortFns {
		t.Run(tc.name, func(t *testing.T) { testSort(t, tc.f) })
	}
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func SortMerge[T constraints.Ordered](s []T) {
	from, to, fromS := s, make([]T, len(s)), true
	merge := func(segLen, tailIdx int) int {
		a, b := tailIdx, tailIdx+segLen
		A, B := min(segLen, len(s)-a), min(segLen, len(s)-b)
		for ; A > 0 && B > 0; tailIdx++ {
			if from[a] <= from[b] {
				to[tailIdx], a, A = from[a], a+1, A-1
			} else {
				to[tailIdx], b, B = from[b], b+1, B-1
			}
		}
		if A > 0 {
			return tailIdx + copy(to[tailIdx:], from[a:a+A])
		}
		return tailIdx + copy(to[tailIdx:], from[b:b+B])
	}
	for segLen := 1; segLen < len(s); segLen *= 2 {
		for tailIdx := 0; tailIdx < len(s); {
			tailIdx = merge(segLen, tailIdx)
		}
		from, to, fromS = to, from, !fromS
	}
	if !fromS {
		copy(s, from)
	}
}

func SortBubble[T constraints.Ordered](s []T) {
	for i := len(s); i > 0; i-- {
		swapped := false
		for j := 1; j < i; j++ {
			if s[j-1] > s[j] {
				s[j-1], s[j] = s[j], s[j-1]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func SortInsert[T constraints.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		j := 0
		for s[j] < s[i] {
			j++
		}
		for j < i && s[j] == s[i] {
			j++
		}
		// for k := i; k > j; k-- {
		// 	s[k], s[k-1] = s[k-1], s[k]
		// }
		x := s[i]
		copy(s[j+1:], s[j:i])
		s[j] = x
	}
}

func SortInsertBisect[T constraints.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		x, j, k := s[i], 0, i
		for j != k {
			if m := (j + k) / 2; s[m] <= x {
				j = m + 1
			} else {
				k = m
			}
		}
		copy(s[j+1:], s[j:i])
		s[j] = x
	}
}

func SortInsert2[T constraints.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		n, j := s[i], i-1
		for ; j >= 0 && s[j] > n; j-- {
			s[j+1] = s[j]
		}
		s[j+1] = n
	}
}

func partitionForQuick[T constraints.Ordered](s []T, piv T) (left, right []T) {
	i, j := 0, len(s)-1
	for i < j {
		for s[i] < piv {
			i++
		}
		for s[j] > piv {
			j--
		}
		if s[j] == piv {
			if s[i] > piv {
				s[i], s[j] = s[j], s[i]
			}
			j--
		} else {
			s[i], s[j] = s[j], s[i]
			i, j = i+1, j-1
		}
	}
	for i < len(s) && s[i] == piv {
		i++
	}
	for j >= 0 && s[j] == piv {
		j--
	}
	return s[:j+1], s[i:]
}

func SortQuick[T constraints.Ordered](s []T) {
	if len(s) <= 1 {
		return
	}
	left, right := partitionForQuick(s, s[len(s)/2])
	SortQuick(left)
	SortQuick(right)
}
