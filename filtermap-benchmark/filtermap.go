package main

import (
	"fmt"
	"sync"
)

type T1 = int
type T2 = float32

func f(v T1) (T2, bool) {
	if v%2 == 0 {
		return 0, false
	}
	return float32(v), true
}

func ConcurrentFilterMap1(in []T1, f func(T1) (T2, bool), maxConcurrency int) []T2 {
	input := make(chan T1, len(in))
	for _, v := range in {
		input <- v
	}
	close(input)
	output := make(chan T2, len(in)) // We don't know an output size in advance.
	filterMap(input, output, f, maxConcurrency)
	close(output)
	var result []T2
	for v := range output {
		result = append(result, v)
	}
	return result
}

func ConcurrentFilterMap2(in []T1, f func(T1) (T2, bool), maxConcurrency int) []T2 {
	input := make(chan T1) // No buffering!
	go func() {            // First go-routine.
		for _, v := range in {
			input <- v
		}
		close(input)
	}()
	output := make(chan T2) // No buffering!
	go func() {             // Second extra go routine.
		filterMap(input, output, f, maxConcurrency)
		close(output)
	}()
	var result []T2
	for v := range output {
		result = append(result, v)
	}
	return result
}

func filterMap(input chan T1, output chan T2, f func(T1) (T2, bool), maxConcurrency int) {
	var wg sync.WaitGroup
	wg.Add(maxConcurrency)
	defer wg.Wait()
	for i := 0; i < maxConcurrency; i++ {
		go func() {
			defer wg.Done()
			for v := range input {
				if v, ok := f(v); ok {
					output <- v
				}
			}
		}()
	}
}

func ConcurrentFilterMap3(in []T1, f func(T1) (T2, bool), maxConcurrency int) []T2 {
	out := make(chan []T2, 1)
	out <- nil

	type token struct{}
	sem := make(chan token, maxConcurrency)
	for _, v1 := range in {
		sem <- token{}
		v1 := v1
		go func() {
			defer func() { <-sem }()
			if v2, ok := f(v1); ok {
				out <- append(<-out, v2)
			}
		}()
	}

	// Acquire all semaphore slots to wait for work to complete.
	for n := cap(sem); n > 0; n-- {
		sem <- token{}
	}

	return <-out
}

func main() {
	in := []T1{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	out2 := ConcurrentFilterMap1(in, f, 5)
	out1 := ConcurrentFilterMap2(in, f, 5)
	fmt.Println(out1, out2)
}
