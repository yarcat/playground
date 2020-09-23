package insert

import (
	"container/list"
	"testing"
)

var (
	slice10   = slice(10)
	slice100  = slice(100)
	slice1k   = slice(1_000)
	slice10k  = slice(10_000)
	slice100k = slice(100_000)
	slice1m   = slice(1_000_000)
)

func slice(size int) []byte { return make([]byte, size) }

func BenchmarkInsertHeadSlice10(b *testing.B)      { benchmarkInsertSlice(b, slice10) }
func BenchmarkInsertHeadSlice100(b *testing.B)     { benchmarkInsertSlice(b, slice100) }
func BenchmarkInsertHeadSlice1000(b *testing.B)    { benchmarkInsertSlice(b, slice1k) }
func BenchmarkInsertHeadSlice10000(b *testing.B)   { benchmarkInsertSlice(b, slice10k) }
func BenchmarkInsertHeadSlice100000(b *testing.B)  { benchmarkInsertSlice(b, slice100k) }
func BenchmarkInsertHeadSlice1000000(b *testing.B) { benchmarkInsertSlice(b, slice1m) }

func benchmarkInsertSlice(b *testing.B, a []byte) {
	c := cap(a)
	for n := 0; n < b.N; n++ {
		// insert at i trick:
		// a = append(a[:i], append([]T{x}, a[i:]...)...)
		// could be simplified as we insert to head
		a = append([]byte{0}, a...)
		// reset a
		a = a[:c:c]
	}
}

var (
	linkedList10   = linkedList(10)
	linkedList100  = linkedList(100)
	linkedList1k   = linkedList(1_000)
	linkedList10k  = linkedList(10_000)
	linkedList100k = linkedList(100_000)
	linkedList1m   = linkedList(1_000_000)
)

func linkedList(size int) *list.List {
	l := list.New()
	for i := 0; i < size; i++ {
		l.PushBack(0)
	}
	return l
}

func BenchmarkInsertHeadList10(b *testing.B)      { benchmarkInsertList(b, linkedList10) }
func BenchmarkInsertHeadList100(b *testing.B)     { benchmarkInsertList(b, linkedList100) }
func BenchmarkInsertHeadList1000(b *testing.B)    { benchmarkInsertList(b, linkedList1k) }
func BenchmarkInsertHeadList10000(b *testing.B)   { benchmarkInsertList(b, linkedList10k) }
func BenchmarkInsertHeadList100000(b *testing.B)  { benchmarkInsertList(b, linkedList100k) }
func BenchmarkInsertHeadList1000000(b *testing.B) { benchmarkInsertList(b, linkedList1m) }

func benchmarkInsertList(b *testing.B, l *list.List) {
	for n := 0; n < b.N; n++ {
		l.PushFront(0)
		// reset l
		l.Remove(l.Front())
	}
}
