package shift

func shiftLeft(s []int, n int) {
	if len(s) == 0 || len(s) == 1 {
		return
	}
	if n %= len(s); n == 0 {
		return
	}
	// positioned contains amount of elements positioned.
	// left is the left-most index of an element that we start from. All
	// elements before are already positioned and must not be touched anymore,
	// hense the invariant pos > left (inner loop).
	for left, positioned := 0, 0; positioned < len(s); left++ {
		// Initial pos is the position the left element should be positioned to.
		// Left and pos elements are already at their places. This always puts
		// two elements at their places, giving one extra positioned element at
		// every iteration.
		positioned++
		for pos := len(s) - n + left; pos > left; pos = (len(s) + pos - n) % len(s) {
			// Using left-th element as an intermediate location. Eventually it
			// will get the element that should be in that position.
			s[left], s[pos] = s[pos], s[left]
			positioned++
		}
	}
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func shiftLeftRev(s []int, n int) {
	if len(s) == 0 || len(s) == 1 {
		return
	}
	if n %= len(s); n == 0 {
		return
	}
	reverse(s)
	reverse(s[:len(s)-n])
	reverse(s[len(s)-n:])
}
