package shift

import (
	"reflect"
	"testing"
)

func TestShiftLeft2(t *testing.T)   { testShiftLeft(t, shiftLeft2) }
func TestShiftLeft(t *testing.T)    { testShiftLeft(t, shiftLeft) }
func TestShiftLeftRev(t *testing.T) { testShiftLeft(t, shiftLeftRev) }

func TestAllShifts(t *testing.T) {
	for i := 1; i <= 50; i++ {
		in1, in2, in3 := s(1, i), s(1, i), s(1, i)
		for j := 0; j <= i; j++ {
			shiftLeft(in1, j)
			shiftLeftRev(in2, j)
			shiftLeft2(in3, j)
			if !reflect.DeepEqual(in1, in2) || !reflect.DeepEqual(in1, in3) {
				t.Fatalf("i = %v, j = %v", i, j)
			}
		}
	}
}

func testShiftLeft(t *testing.T, fn func([]int, int)) {
	for _, tc := range []struct {
		name string
		in   []int
		n    int
		want []int
	}{
		{name: "nil"},
		{name: "nil-10", n: 10},
		{name: "2-0", in: s(1, 2), n: 0, want: s(1, 2)},
		{name: "2-1", in: s(1, 2), n: 1, want: s(2, 2, 1, 1)},
		{name: "2-2", in: s(1, 2), n: 2, want: s(1, 2)},
		{name: "3-0", in: s(1, 3), n: 0, want: s(1, 3)},
		{name: "3-1", in: s(1, 3), n: 1, want: s(2, 3, 1, 1)},
		{name: "3-2", in: s(1, 3), n: 2, want: s(3, 3, 1, 2)},
		{name: "3-3", in: s(1, 3), n: 3, want: s(1, 3)},
		{name: "6-1", in: s(1, 6), n: 1, want: s(2, 6, 1, 1)},
		{name: "6-2", in: s(1, 6), n: 2, want: s(3, 6, 1, 2)},
		{name: "6-3", in: s(1, 6), n: 3, want: s(4, 6, 1, 3)},
		{name: "6-4", in: s(1, 6), n: 4, want: s(5, 6, 1, 4)},
		{name: "6-5", in: s(1, 6), n: 5, want: s(6, 6, 1, 5)},
		{name: "7-1", in: s(1, 7), n: 1, want: s(2, 7, 1, 1)},
		{name: "7-2", in: s(1, 7), n: 2, want: s(3, 7, 1, 2)},
		{name: "7-3", in: s(1, 7), n: 3, want: s(4, 7, 1, 3)},
		{name: "7-4", in: s(1, 7), n: 4, want: s(5, 7, 1, 4)},
		{name: "7-5", in: s(1, 7), n: 5, want: s(6, 7, 1, 5)},
		{name: "7-6", in: s(1, 7), n: 6, want: s(7, 7, 1, 6)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fn(tc.in, tc.n)
			if !reflect.DeepEqual(tc.want, tc.in) {
				t.Errorf("got = %v, want = %v\n", tc.in, tc.want)
			}
		})
	}
}

func s(rr ...int) (res []int) {
	if len(rr)%2 != 0 {
		panic("expected [min; max] ranges")
	}
	for i := 0; i < len(rr); i += 2 {
		start, end := rr[i], rr[i+1]
		for v := start; v <= end; v++ {
			res = append(res, v)
		}
	}
	return
}

func benchmarkShift(b *testing.B, fn func([]int, int), s []int, n int) {
	for i := 0; i < b.N; i++ {
		fn(s, n)
	}
}

var (
	small = s(1, 100)
	large = s(1, 1_000_000)
)

func BenchmarkShiftSmall(b *testing.B)    { benchmarkShift(b, shiftLeft, small, len(small)/2) }
func BenchmarkShiftLarge(b *testing.B)    { benchmarkShift(b, shiftLeft, large, len(large)/2) }
func BenchmarkShift2Small(b *testing.B)   { benchmarkShift(b, shiftLeft2, small, len(small)/2) }
func BenchmarkShift2Large(b *testing.B)   { benchmarkShift(b, shiftLeft2, large, len(large)/2) }
func BenchmarkShiftRevSmall(b *testing.B) { benchmarkShift(b, shiftLeftRev, small, len(small)/2) }
func BenchmarkShiftRevLarge(b *testing.B) { benchmarkShift(b, shiftLeftRev, large, len(large)/2) }
