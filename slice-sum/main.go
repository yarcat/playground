package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	const sliceLen = 2e8
	numCPU := runtime.NumCPU()
	//numCPU = 1

	s := make([]int64, sliceLen)
	for i := range s {
		s[i] = 2
	}

	for range 10 {
		start := time.Now()

		chunkLen := sliceLen / numCPU
		out := make(chan int64, numCPU)
		// TODO: Check if we are missing any pieces.
		for chunkStart := 0; chunkStart < sliceLen; chunkStart += chunkLen {
			job := s[chunkStart : chunkStart+chunkLen]
			go func() {
				var sum int64
				for _, x := range job {
					sum += x
				}
				out <- sum
			}()
		}
		var sum int64
		for range cap(out) {
			sum += <-out
		}

		fmt.Println(time.Since(start), sum)
	}
}
