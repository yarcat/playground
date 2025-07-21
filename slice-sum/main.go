package main

import (
	"fmt"
	"time"
)

func main() {
	s := make([]int64, 2e8)
	for i := range s {
		s[i] = 2
	}
	start := time.Now()
	var sum int64
	for _, x := range s {
		sum += x
	}
	fmt.Println(time.Since(start), sum)
}
