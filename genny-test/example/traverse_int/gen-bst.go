// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package main

import "genny-test/tree"

type (
	LessFn func(int, int) bool

	Node struct {
		left, right *Node
		value       int
	}
	BST struct {
		root *Node
		less LessFn
	}
)

func (n *Node) Left() *Node  { return n.left }
func (n *Node) Right() *Node { return n.right }
func (n Node) Value() int    { return n.value }

func New(less LessFn) *BST { return &BST{less: less} }

func (bst *BST) Insert(v int) {
	n := &bst.root
	for *n != nil {
		if bst.less(v, (*n).value) {
			n = &(*n).left
		} else {
			n = &(*n).right
		}
	}
	*n = &Node{value: v}
}

func (bst *BST) Traverse(visit func(*Node)) { tree.Traverse(asNode{bst.root, visit}) }

type asNode struct {
	*Node
	visit func(*Node)
}

func (n asNode) Left() tree.Node  { return asNode{n.left, n.visit} }
func (n asNode) Right() tree.Node { return asNode{n.right, n.visit} }
func (n asNode) Empty() bool      { return n.Node == nil }
func (n asNode) Visit()           { n.visit(n.Node) }
