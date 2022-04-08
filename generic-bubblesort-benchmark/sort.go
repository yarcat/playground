package main

import (
	"fmt"
	"sort"
)

func SortInt(x sort.Interface) {
	l := x.Len() // Do not call x.Len() too often.
	for i := l - 1; i > 0; i-- {
		swapped := false
		for j := 0; j < i; j++ {
			if !x.Less(j, j+1) {
				x.Swap(j, j+1)
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func SortGen[T sort.Interface](x T) {
	l := x.Len() // Do not call x.Len() too often.
	for i := l - 1; i > 0; i-- {
		swapped := false
		for j := 0; j < i; j++ {
			if !x.Less(j, j+1) {
				x.Swap(j, j+1)
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func main() {
	x := []int{10, 5, 1, 2, 8, 4, 3, 7, 6, 9}

	y := make([]int, len(x))
	copy(y, x)
	SortInt(sort.IntSlice(y))
	fmt.Println(y)

	z := make([]int, len(x))
	copy(z, x)
	SortGen(sort.IntSlice(z))
	fmt.Println(z)
}
