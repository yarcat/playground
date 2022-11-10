package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

type (
	node struct {
		next *node
		val  int
	}
	list struct {
		head, tail *node
	}
)

func (l *list) qsort() {
	if l.head == l.tail { // Single element or empty.
		return
	}
	less, other := split(l.head.next, func(val int) bool { return val < l.head.val })
	*l = concat(qsorted(less), single(l.head), qsorted(other))
}

func (l *list) msort() {
	if l.head == l.tail { // Single element or empty.
		return
	}
	var first bool
	a, b := split(l.head, func(int) bool { first = !first; return first })
	*l = merge(msorted(a), msorted(b))
}

func (l *list) pushBack(val int) {
	if l.head == nil {
		l.head = &node{val: val}
		l.tail = l.head
	} else {
		l.tail.next = &node{val: val}
		l.tail = l.tail.next
	}
}

func (l *list) foreach(f func(int)) {
	for p := l.head; p != nil; p = p.next {
		f(p.val)
	}
}

func (l *list) appendList(lst list) {
	if l.head == nil {
		l.head, l.tail = lst.head, lst.tail
	} else if lst.head != nil {
		l.tail.next, l.tail = lst.head, lst.tail
	}
	if l.tail != nil {
		l.tail.next = nil
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var l list
	for i := 0; i < 10; i++ {
		l.pushBack(rand.Intn(1000))
	}
	fmt.Println(stringify(l))
	// fmt.Println(stringify(msorted(l)))
	fmt.Println(stringify(qsorted(l)))

	//zz := []int{71, 1, 87, 58, 15, 55, 41, 56, 6, 27, 75, 26, 63, 35, 16, 64, 80, 36, 18, 79, 94, 54, 81, 91, 93, 90, 95, 3, 66, 53, 68, 9, 19, 62, 23, 86, 31, 5, 2, 40, 82, 98, 77, 76, 89, 44, 67, 28, 60, 99, 74, 83, 73, 0, 57, 39, 59, 17, 24, 25, 96, 14, 32, 22, 52, 38, 42, 85, 21, 7, 45, 37, 33, 12, 69, 49, 97, 84, 65, 4, 92, 46, 50, 88, 48, 8, 72, 11, 30, 51, 61, 70, 13, 29, 47, 20, 34, 78, 43, 10}
	zz := []int{6, 5, 0, 2, 7, 3, 1, 8, 9}
	l1, _ := genqsort(genlist(zz...), genless[int])
	genprint(l1)
	l2, _ := genqsort(genlist("hello", "world", "!!!"), genless[string])
	genprint(l2)
	l3 := genmsort(genlist(zz...), genless[int])
	genprint(l3)
}

func stringify(l list) string {
	var s []string
	l.foreach(func(val int) {
		s = append(s, fmt.Sprint(val))
	})
	return fmt.Sprintf("[%s]", strings.Join(s, " "))
}

func next(n *node) *node {
	if n == nil {
		return nil
	}
	return n.next
}

func split(n *node, f func(int) bool) (tru, fal list) {
	for nex := next(n); n != nil; n, nex = nex, next(nex) {
		head, tail := &fal.head, &fal.tail
		if f(n.val) {
			head, tail = &tru.head, &tru.tail
		}
		n.next = nil
		if *head == nil {
			*head, *tail = n, n
		} else {
			(*tail).next, *tail = n, n
		}
	}
	return
}

func concat(lists ...list) (l list) {
	for _, lst := range lists {
		l.appendList(lst)
	}
	return
}

func merge(a, b list) (l list) {
	ah, bh := a.head, b.head
	for ah != nil && bh != nil {
		if ah.val < bh.val {
			n := next(ah)
			l.appendList(list{head: ah, tail: ah})
			ah = n
		} else {
			n := next(bh)
			l.appendList(list{head: bh, tail: bh})
			bh = n
		}
	}
	if ah != nil {
		l.appendList(list{head: ah, tail: a.tail})
	} else if bh != nil {
		l.appendList(list{head: bh, tail: b.tail})
	}
	return l
}

func single(n *node) list { return list{head: n, tail: n} }
func qsorted(l list) list { l.qsort(); return l }
func msorted(l list) list { l.msort(); return l }

type gennode[T any] struct {
	v    T
	next *gennode[T]
}

func genlist[T any](vals ...T) (l *gennode[T]) {
	for i := len(vals) - 1; i >= 0; i-- {
		l = &gennode[T]{v: vals[i], next: l}
	}
	return l
}

func genprint[T any](n *gennode[T]) {
	for ; n != nil; n = n.next {
		fmt.Print(" -> ", n.v)
	}
	fmt.Println(" -> nil ")
}

func genconcat[T any](l, lt, n, r, rt *gennode[T]) (head, tail *gennode[T]) {
	n.next = r
	if l == nil {
		return n, rt
	}
	lt.next = n
	if r == nil {
		return l, n
	}
	return l, rt
}

func genpartition[T any](n *gennode[T], v T, less func(a, b T) bool) (l, r *gennode[T]) {
	pl, pr := &l, &r
	for ; n != nil; n = n.next {
		if less(n.v, v) {
			*pl = n
			pl = &(*pl).next
		} else {
			*pr = n
			pr = &(*pr).next
		}
	}
	*pl, *pr = nil, nil
	return l, r
}

func genqsort[T any](n *gennode[T], less func(T, T) bool) (head, tail *gennode[T]) {
	if n == nil { // Empty.
		return nil, nil
	} else if n.next == nil { // One element.
		return n, n
	}
	l, r := genpartition(n.next, n.v, less)
	l, lt := genqsort(l, less)
	r, rt := genqsort(r, less)
	return genconcat(l, lt, n, r, rt)
}

// genskip skips not more than x elements. The function expects the list not
// to be empty, and it guarantees returned value is not nil (it could stop
// earlier, if next element is nil).
func genskip[T any](n *gennode[T], x int) *gennode[T] {
	for x > 0 && n.next != nil {
		x, n = x-1, n.next
	}
	return n
}

// merge accepts a segment of a list like a...bprev...cprev..., merges all these
// elements in the order defined by less. The function returns new head and tail elements.
func genmerge[T any](a, bprev, cprev *gennode[T], less func(T, T) bool) (a2, cprev2 *gennode[T]) {
	p, b, c := &a2, bprev.next, cprev.next
	for ; a != bprev.next && b != c; p = &(*p).next {
		if less(a.v, b.v) {
			*p, a = a, a.next
		} else {
			*p, b = b, b.next
		}
	}
	if a != bprev.next {
		*p, bprev.next = a, c
		return a2, bprev
	}
	*p = b
	return a2, cprev
}

// msort implements non-recursive merge sort for lists.
func genmsort[T any](n *gennode[T], less func(T, T) bool) *gennode[T] {
	for i := 1; ; i *= 2 {
		merged := false
		// n is split into segments a...bprev...cprev....
		for a := &n; *a != nil; {
			bprev := genskip(*a, i-1)
			if bprev.next == nil { // There is nothing to merge with.
				break
			}
			cprev := genskip(bprev.next, i-1)
			*a, cprev = genmerge(*a, bprev, cprev, less)
			a, merged = &cprev.next, true
		}
		if !merged {
			return n
		}
	}
}

func genless[T constraints.Ordered](a, b T) bool { return a < b }
