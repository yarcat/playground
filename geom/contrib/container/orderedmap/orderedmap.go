// Package orderedmap implements a map that remembers its key insertion order.
package orderedmap

import (
	"container/list"
)

// Keeping key and value types for potential generation of specialized ordered
// maps.
type (
	// KeyType is the type of keys stored in the OrderedMap.
	KeyType interface{}
	// ValueType is the type of values stored in the OrderedMap.
	ValueType interface{}
)

// OrderedMap is a map that remembers its key insertion order. It provides only
// minimal set of functionality required for my use-case.
type OrderedMap struct {
	// items is a list of (key, value) pairs. The list preserves the order of the
	// elements added to the map.
	items *list.List
	// dict is a lookup table that helps to find elements in the items list.
	dict map[KeyType]*list.Element
}

// New returns new OrderedMap.
func New() *OrderedMap {
	return &OrderedMap{
		items: list.New(),
		dict:  make(map[KeyType]*list.Element),
	}
}

// Get returns an associated value and a test flag that denotes an existence of
// the key in the dictionary. The ok is true if the key exists, false otherwise.
func (om OrderedMap) Get(k KeyType) (v ValueType, ok bool) {
	e, ok := om.dict[k]
	if !ok {
		return
	}
	return asMapItem(e).value, true
}

// Set associates the value with the key. If this key already existed in the
// dictionary it will be overwritten and the order won't be changed.
func (om *OrderedMap) Set(k KeyType, v ValueType) {
	if e, ok := om.dict[k]; ok {
		asMapItem(e).value = v
		return
	}
	item := &mapItem{key: k, value: v}
	om.dict[k] = om.items.PushBack(item)
}

// Iter returns an iterator over the key, pair values stored in the map.
//
// Example:
//
//	m := orderedmap.New()
//	m.Set("hello", "world")
//	for it := m.Iter(); it.Next(); {
//		fmt.Println(it.Value())
//	}
func (om OrderedMap) Iter() *Iterator {
	return newIterator(om.items)
}

type mapItem struct {
	key   KeyType
	value ValueType
}

func asMapItem(e *list.Element) *mapItem {
	if e == nil {
		return nil
	}
	return e.Value.(*mapItem)
}

type nextItemFn func() (item *mapItem, ok bool)

// Iterator implements Next/Value type iterator helper that assists going
// though the sequence of key/value pairs stored in the OrderedMap.
type Iterator struct {
	item     *mapItem
	ok       bool
	nextItem nextItemFn
}

func newIterator(pairs *list.List) *Iterator {
	// I've implemented a simple state machine here:
	//	- Being called for the first time it returns the
	//		front element of the list.
	//	- All consequitive calls return next element.
	var (
		getter, getFront, getNext nextItemFn
		elem                      *list.Element
	)

	getFront = func() (*mapItem, bool) {
		elem = pairs.Front()
		getter = getNext
		return asMapItem(elem), elem != nil
	}

	getNext = func() (*mapItem, bool) {
		elem = elem.Next()
		return asMapItem(elem), elem != nil
	}

	getter = getFront

	return &Iterator{nextItem: func() (*mapItem, bool) {
		return getter()
	}}
}

// Next moves the iterator to the next key/value pair in the OrderedMap. It
// returns boolean indicating whether the iterator still has elements. If the
// return value is stop, the iterator is exposed and there are no more key/pairs
// left to iterate over.
func (it *Iterator) Next() bool {
	it.item, it.ok = it.nextItem()
	return it.ok
}

// Value returns current key/value pair.
func (it Iterator) Value() (KeyType, ValueType) {
	return it.item.key, it.item.value
}
