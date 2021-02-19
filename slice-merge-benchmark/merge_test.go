package slicemerge

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

const elements = 10_000_000 // <- Change me!

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	dst    = make([]int, elements*2)
	sliceA = makeRandInts(elements)
	sliceB = makeRandInts(elements)
)

func makeRandInts(n int) []int {
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = rand.Int()
	}
	return res
}

func BenchmarkMergeSorted(b *testing.B) {
	benchmarkMerge(b, MergeSorted)
}

func BenchmarkMerge(b *testing.B) {
	benchmarkMerge(b, Merge)
}

func benchmarkMerge(b *testing.B, merge func(dst, a, b []int) int) {
	for i := 0; i < b.N; i++ {
		merge(dst, sliceA, sliceB)
	}
}

type testCase struct {
	name     string
	a, b     []int
	dst      []int
	wantDst  []int
	wantSize int
}

var commonCases = []testCase{
	{name: "nil"},
	{name: "nil dst", a: []int{1}, b: []int{2}},
	{name: "empty dst", a: []int{1}, b: []int{2}, dst: []int{}, wantDst: []int{}},
	{name: "large dst", a: []int{1, 3}, b: []int{1, 2, 3}, dst: make([]int, 4), wantDst: []int{1, 2, 3, 0}, wantSize: 3},
	{name: "only a", a: []int{1, 2, 3}, dst: make([]int, 3), wantDst: []int{1, 2, 3}, wantSize: 3},
	{name: "only b", b: []int{1, 2, 3}, dst: make([]int, 3), wantDst: []int{1, 2, 3}, wantSize: 3},
	{name: "same", a: []int{1, 2, 3}, b: []int{1, 2, 3}, dst: make([]int, 3), wantDst: []int{1, 2, 3}, wantSize: 3},
	// "Small dst" would be merge implementation specific.
}

func TestMergeSorted(t *testing.T) {
	testMerge(t, append(commonCases,
		testCase{name: "short dst", a: []int{1, 3, 5}, b: []int{2, 4, 6},
			dst: []int{0, 0, 0, 0}, wantDst: []int{1, 2, 3, 4}, wantSize: 4},
	), MergeSorted)
}

func TestMerge(t *testing.T) {
	testMerge(t, append(commonCases,
		testCase{name: "short dst", a: []int{10, 3, 15}, b: []int{2, 4, 6},
			dst: []int{0, 0, 0, 0}, wantDst: []int{2, 3, 10, 15}, wantSize: 4},
	), Merge)
}

func testMerge(t *testing.T, cases []testCase, merge func(dst, a, b []int) int) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotSize := merge(tc.dst, tc.a, tc.b)
			if !reflect.DeepEqual(tc.wantDst, tc.dst) || gotSize != tc.wantSize {
				t.Errorf("got dst = %v, size = %v; want dst = %v, size = %v",
					tc.dst, gotSize, tc.wantDst, tc.wantSize)
			}
		})
	}
}
