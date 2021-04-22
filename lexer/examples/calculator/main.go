package main

import (
	"fmt"
)

func main() {
	for _, expr := range []string{
		"-10",
		"1 + 2*3",
		"2*3 + 4",
		"(2+2)*2",
		"-2+2*-2",
		"(2+2*2",
		"2+2)*2",
		"2*10-",
	} {
		x, err := eval(newScanner(expr))
		if err == nil {
			fmt.Printf("   OK: %13s = %d\n", expr, x)
		} else {
			fmt.Printf("ERROR: %13s   %v\n", expr, err)
		}
	}
}
