package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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
