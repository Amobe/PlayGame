package utils

import "fmt"

// LinkedList contains the items linked to the next and sorted by the cond function.
// The item can be any type
type LinkedList[T any] struct {
	head *LinkedListItem[T]
	cond func(T, T) bool
	size int
}

// NewLinkedList creates a linked list with a specific type and compare condition function.
func NewLinkedList[T any](cond func(T, T) bool) *LinkedList[T] {
	return &LinkedList[T]{
		cond: cond,
	}
}

func (l *LinkedList[T]) Insert(val T) {
	l.size++
	tmp := &LinkedListItem[T]{
		Value: val,
	}
	if l.head == nil {
		l.head = tmp
		return
	}
	if l.cond(tmp.Value, l.head.Value) {
		tmp.Next = l.head
		l.head = tmp
		return
	}
	ptr := l.head
	for ptr.Next != nil {
		if l.cond(tmp.Value, ptr.Next.Value) {
			tmp.Next = ptr.Next
			ptr.Next = tmp
			return
		}
		ptr = ptr.Next
	}
	ptr.Next = tmp
}

func (l *LinkedList[T]) InsertMany(vals []T) {
	for _, val := range vals {
		l.Insert(val)
	}
}

func (l *LinkedList[T]) Iterator() *LinkedListIterator[T] {
	ptr := &LinkedListItem[T]{
		Next: l.head,
	}
	return &LinkedListIterator[T]{
		ptr: ptr,
	}
}

type LinkedListItem[T any] struct {
	Value T
	Next  *LinkedListItem[T]
}

type LinkedListIterator[T any] struct {
	ptr *LinkedListItem[T]
}

func (i *LinkedListIterator[T]) HasNext() bool {
	return i.ptr.Next != nil
}

func (i *LinkedListIterator[T]) Next() (t T, err error) {
	if i.ptr == nil {
		return t, fmt.Errorf("no element left")
	}
	item := i.ptr.Next
	if item == nil {
		return t, fmt.Errorf("next element is empty")
	}
	i.ptr = i.ptr.Next
	return item.Value, nil
}
