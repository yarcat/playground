package main

import (
	"fmt"
	"runtime"
	"simd/archsimd"
	"slices"
	"time"
	"unsafe"
)

type TT interface {
	int32 | int64
}

func simpleLoop[T TT](a []T) (sum T) {
	for i := range a {
		sum += a[i]
	}
	return sum
}

func indepSums2[T TT](a []T) T {
	aa := a
	if len(a)%2 != 0 {
		aa = a[:len(a)-1]
	}

	var a0, a1 T
	for i := 0; i < len(aa); i += 2 {
		b := aa[i : i+2 : i+2]
		a0 += b[0]
		a1 += b[1]
	}

	if len(a)%2 == 0 {
		return a0 + a1
	}
	return a0 + a1 + a[len(a)-1]
}

func indepSums4[T TT](a []T) T {
	n := len(a)
	end := n - n%4
	aa := a[:end]

	var a0, a1, a2, a3 T
	for i := 0; i < len(aa); i += 4 {
		b := aa[i : i+4 : i+4]
		a0 += b[0]
		a1 += b[1]
		a2 += b[2]
		a3 += b[3]
	}

	for _, x := range a[end:] {
		a0 += x
	}

	return a0 + a1 + a2 + a3
}

func indepSums8[T TT](a []T) T {
	n := len(a)
	end := n - n%8
	aa := a[:end]

	var a0, a1, a2, a3, a4, a5, a6, a7 T
	for i := 0; i < len(aa); i += 8 {
		b := aa[i : i+8 : i+8]
		a0 += b[0]
		a1 += b[1]
		a2 += b[2]
		a3 += b[3]
		a4 += b[4]
		a5 += b[5]
		a6 += b[6]
		a7 += b[7]
	}

	for _, x := range a[end:] {
		a0 += x
	}

	return a0 + a1 + a2 + a3 + a4 + a5 + a6 + a7
}

func indepSums8p[T TT](a []T) T {
	n := uintptr(len(a))
	end := n - n%8
	
	basePtr := unsafe.Pointer(&a[0])
	size := unsafe.Sizeof(a[0])

	var a0, a1, a2, a3, a4, a5, a6, a7 T
	for i := uintptr(0); i < end; i += 8 {
		block := (*[8]T)(unsafe.Add(basePtr, uintptr(i)*size))

		a0 += block[0]
		a1 += block[1]
		a2 += block[2]
		a3 += block[3]
		a4 += block[4]
		a5 += block[5]
		a6 += block[6]
		a7 += block[7]
	}

	for i := end; i < n; i++ {
		val := (*T)(unsafe.Add(basePtr, uintptr(i)*size))
		a0 += *val
	}

	return a0 + a1 + a2 + a3 + a4 + a5 + a6 + a7
}

func reduceAdd(acc archsimd.Int32x8) int32 {
	s1 := acc.AddPairs(acc)
	s2 := s1.AddPairs(s1)
	var tmp [8]int32
	s2.Store(&tmp)
	return tmp[0] + tmp[4]
}

func indepSums8simdInt32(a []int32) int32 {
	n := len(a)
	const laneSize = 8
	end := n - n%laneSize

	var acc archsimd.Int32x8

	for i := 0; i < end; i += laneSize {
		ptr := (*[8]int32)(unsafe.Pointer(&a[i]))
		v := archsimd.LoadInt32x8(ptr)
		acc = acc.Add(v)
	}

	total := reduceAdd(acc)

	for _, x := range a[end:] {
		total += x
	}

	return total
}

func indepSums16simdInt32(a []int32) int32 {
	n := len(a)
	const laneSize = 16
	end := n - n%laneSize

	var acc0, acc1 archsimd.Int32x8

	for i := 0; i < end; i += laneSize {
		ptr0 := (*[8]int32)(unsafe.Pointer(&a[i]))
		v0 := archsimd.LoadInt32x8(ptr0)
		acc0 = acc0.Add(v0)

		ptr1 := (*[8]int32)(unsafe.Pointer(&a[i+8]))
		v1 := archsimd.LoadInt32x8(ptr1)
		acc1 = acc1.Add(v1)
	}

	total := reduceAdd(acc0.Add(acc1))

	for _, x := range a[end:] {
		total += x
	}

	return total
}

func paral[T TT](a []T) (sum T) {
	n := runtime.NumCPU() * 2
	seg := (len(a) + n - 1) / n
	out := make(chan T, n)

	for i := range n {
		go func() {
			// out <- simpleLoop(a[i*seg : min(i*seg+seg, len(a))])
			// out <- indepSums2(a[i*seg : min(i*seg+seg, len(a))])
			// out <- indepSums4(a[i*seg : min(i*seg+seg, len(a))])
			// out <- indepSums8(a[i*seg : min(i*seg+seg, len(a))])
			out <- indepSums8p(a[i*seg : min(i*seg+seg, len(a))])
		}()
	}

	for range n {
		sum += <-out
	}

	return sum
}

func report[T TT](name string, f func([]T) T, a []T) {
	const iters = 101

	d := make([]time.Duration, iters)

	for i := range d {
		t := time.Now()
		_ = f(a)
		d[i] = time.Since(t)
	}

	slices.Sort(d)

	fmt.Printf("%15s %8.2f %8.2f %8.2f\n", name,
		d[0].Seconds()*1000,
		d[iters/2].Seconds()*1000,
		d[iters-1].Seconds()*1000,
	)
}

func main() {
	a := make([]int32, 2e8)
	for i := range a {
		a[i] = 2
	}

	fmt.Printf("%15s %8s %8s %8s\n", "",
		"min(ms)", "50(ms)", "max(ms)",
	)
	report("simpleLoop", simpleLoop, a)
	report("indepSums2", indepSums2, a)
	report("indepSums4", indepSums4, a)
	report("indepSums8", indepSums8, a)
	report("indepSums8p", indepSums8p, a)
	report("8simdInt32", indepSums8simdInt32, a)
	report("16simdInt32", indepSums16simdInt32, a)
	report("paral", paral, a)
}