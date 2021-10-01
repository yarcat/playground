package tree

type Node interface {
	Left() Node
	Right() Node
	Empty() bool
	Visit()
}

func Traverse(n Node) {
	if n.Empty() {
		return
	}
	Traverse(n.Left())
	n.Visit()
	Traverse(n.Right())
}
