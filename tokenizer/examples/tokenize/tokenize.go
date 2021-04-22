package main

import (
	"fmt"

	"github.com/yarcat/playground/tokenizer"
	"github.com/yarcat/playground/tokenizer/examples/contrib"
)

func main() {
	t := tokenizer.New("(2+2)*2", contrib.DefaultNumericExpression)
	for t.Next() {
		fmt.Println(t.Token())
	}
}
