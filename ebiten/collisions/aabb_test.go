package main

import "testing"

func makeAabb(x1, y1, x2, y2 float64) aabb {
	return aabb{vector{x1, y1}, vector{x2, y2}}
}

func TestAabbIntersects(t *testing.T) {
	for _, tc := range []struct {
		name     string
		b1, b2   aabb
		expected bool
	}{
		{"the same", makeAabb(0, 0, 100, 100), makeAabb(0, 0, 100, 100), true},
		{"right touch", makeAabb(0, 0, 100, 100), makeAabb(100, 0, 200, 100), true},
		{"left touch", makeAabb(0, 0, 100, 100), makeAabb(-100, 0, 0, 100), true},
		{"bottom touch", makeAabb(0, 0, 100, 100), makeAabb(0, -100, 100, 0), true},
		{"top touch", makeAabb(0, 0, 100, 100), makeAabb(0, 100, 100, 200), true},
		{"inside", makeAabb(0, 0, 100, 100), makeAabb(10, 10, 90, 90), true},
		{"outside", makeAabb(0, 0, 100, 100), makeAabb(-10, -10, 110, 110), true},
		{"no intersection", makeAabb(0, 0, 100, 100), makeAabb(0, -100, 100, -10), false},
	} {
		actual1 := tc.b1.intersects(tc.b2)
		if actual1 != tc.expected {
			t.Errorf("%s: %v.intersects(%v) = %v, expected %v\n",
				tc.name, tc.b1, tc.b2, actual1, tc.expected)
		}
	}
}
