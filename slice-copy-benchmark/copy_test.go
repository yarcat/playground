package main

import "testing"

var (
	empty  []int
	small  = make([]int, 1_000)
	medium = make([]int, 100_000)
	large  = make([]int, 1_000_000)
)

func BenchmarkSliceCopyWithAppendEmpty(b *testing.B) {
	benchmark(b, sliceCopyWithAppend, empty)
}
func BenchmarkSliceCopyWithAppendSmall(b *testing.B) {
	benchmark(b, sliceCopyWithAppend, small)
}
func BenchmarkSliceCopyWithAppendMedium(b *testing.B) {
	benchmark(b, sliceCopyWithAppend, medium)
}
func BenchmarkSliceCopyWithAppendLarge(b *testing.B) {
	benchmark(b, sliceCopyWithAppend, large)
}

func BenchmarkSliceCopyWithTransformAppendEmpty(b *testing.B) {
	benchmark(b, sliceCopyWithTransformAppend(transform), empty)
}
func BenchmarkSliceCopyWithTransformAppendSmall(b *testing.B) {
	benchmark(b, sliceCopyWithTransformAppend(transform), small)
}
func BenchmarkSliceCopyWithTransformAppendMedium(b *testing.B) {
	benchmark(b, sliceCopyWithTransformAppend(transform), medium)
}
func BenchmarkSliceCopyWithTransformAppendLarge(b *testing.B) {
	benchmark(b, sliceCopyWithTransformAppend(transform), large)
}

func BenchmarkSliceCopyWithPreallocEmpty(b *testing.B) {
	benchmark(b, sliceCopyWithPrealloc, empty)
}
func BenchmarkSliceCopyWithPreallocSmall(b *testing.B) {
	benchmark(b, sliceCopyWithPrealloc, small)
}
func BenchmarkSliceCopyWithPreallocMedium(b *testing.B) {
	benchmark(b, sliceCopyWithPrealloc, medium)
}
func BenchmarkSliceCopyWithPreallocLarge(b *testing.B) {
	benchmark(b, sliceCopyWithPrealloc, large)
}

func BenchmarkSliceCopyWithTransformPreallocEmpty(b *testing.B) {
	benchmark(b, sliceCopyWithTransformPrealloc(transform), empty)
}
func BenchmarkSliceCopyWithTransformPreallocSmall(b *testing.B) {
	benchmark(b, sliceCopyWithTransformPrealloc(transform), small)
}
func BenchmarkSliceCopyWithTransformPreallocMedium(b *testing.B) {
	benchmark(b, sliceCopyWithTransformPrealloc(transform), medium)
}
func BenchmarkSliceCopyWithTransformPreallocLarge(b *testing.B) {
	benchmark(b, sliceCopyWithTransformPrealloc(transform), large)
}

func BenchmarkSliceCopyWithIndexAndPreallocEmpty(b *testing.B) {
	benchmark(b, sliceCopyWithIndexAndPrealloc, empty)
}
func BenchmarkSliceCopyWithIndexAndPreallocSmall(b *testing.B) {
	benchmark(b, sliceCopyWithIndexAndPrealloc, small)
}
func BenchmarkSliceCopyWithIndexAndPreallocMedium(b *testing.B) {
	benchmark(b, sliceCopyWithIndexAndPrealloc, medium)
}
func BenchmarkSliceCopyWithIndexAndPreallocLarge(b *testing.B) {
	benchmark(b, sliceCopyWithIndexAndPrealloc, large)
}

func BenchmarkSliceCopyWithTransformIndexAndPreallocEmpty(b *testing.B) {
	benchmark(b, sliceCopyWithTransformIndexAndPrealloc(transform), empty)
}
func BenchmarkSliceCopyWithTransformIndexAndPreallocSmall(b *testing.B) {
	benchmark(b, sliceCopyWithTransformIndexAndPrealloc(transform), small)
}
func BenchmarkSliceCopyWithTransformIndexAndPreallocMedium(b *testing.B) {
	benchmark(b, sliceCopyWithTransformIndexAndPrealloc(transform), medium)
}
func BenchmarkSliceCopyWithTransformIndexAndPreallocLarge(b *testing.B) {
	benchmark(b, sliceCopyWithTransformIndexAndPrealloc(transform), large)
}

type copyFunc func(in []int) (out []int)

func benchmark(b *testing.B, fn copyFunc, in []int) {
	for n := 0; n < b.N; n++ {
		fn(in)
	}
}

func sliceCopyWithPrealloc(in []int) (out []int) {
	if len(in) == 0 {
		return nil
	}
	out = make([]int, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}
	return out
}

func transform(x int) int { return x << 10 }

func sliceCopyWithTransformPrealloc(fn func(int) int) copyFunc {
	return func(in []int) (out []int) {
		if len(in) == 0 {
			return nil
		}
		out = make([]int, 0, len(in))
		for i, v := range in {
			out = append(out, fn(i+v))
		}
		return out
	}
}

func sliceCopyWithAppend(in []int) (out []int) {
	for _, v := range in {
		out = append(out, v)
	}
	return
}

func sliceCopyWithTransformAppend(fn func(int) int) copyFunc {
	return func(in []int) (out []int) {
		for _, v := range in {
			out = append(out, v)
		}
		return
	}
}

func sliceCopyWithIndexAndPrealloc(in []int) (out []int) {
	if len(in) == 0 {
		return nil
	}
	out = make([]int, len(in))
	for i, v := range in {
		out[i] = v
	}
	return
}

func sliceCopyWithTransformIndexAndPrealloc(fn func(int) int) copyFunc {
	return func(in []int) (out []int) {
		if len(in) == 0 {
			return nil
		}
		out = make([]int, len(in))
		for i, v := range in {
			out[i] = v
		}
		return
	}
}
