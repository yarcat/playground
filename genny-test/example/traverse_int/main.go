package main

import "fmt"

//go:generate go run github.com/cheekybits/genny -pkg=main -in=../../tree/bst/bst.go -out=gen-bst.go gen "Type=int"

func main() {
	t := New(func(x, y int) bool { return x < y })
	t.Insert(10)
	t.Insert(5)
	t.Insert(15)
	t.Traverse(func(n *Node) { fmt.Println(n.Value()) })
}
