package slicemerge

import "sort"

// MergeSorted merges sorted a and b into dst. Result is also sorted.
func MergeSorted(dst, a, b []int) (size int) {
	for len(a) > 0 && len(b) > 0 && size < len(dst) {
		if a[0] < b[0] {
			dst[size] = a[0]
			a = a[1:]
		} else if b[0] < a[0] {
			dst[size] = b[0]
			b = b[1:]
		} else { // a[0] == b[0]
			dst[size] = a[0]
			a, b = a[1:], b[1:]
		}
		size++
	}
	if size < len(dst) {
		if len(a) > 0 {
			size += copy(dst[size:], a)
		} else { // len(b) > 0.
			size += copy(dst[size:], b)
		}
	}
	return
}

// Merge merges a and b into dst.
func Merge(dst, a, b []int) (size int) {
	m := make(map[int]struct{}, len(dst))
	for i := 0; len(m) < len(dst) && i < len(a); i++ {
		m[a[i]] = struct{}{}
	}
	for i := 0; len(m) < len(dst) && i < len(b); i++ {
		m[b[i]] = struct{}{}
	}
	for x := range m {
		dst[size] = x
		size++
	}
	sort.Ints(dst[:size])
	return
}
