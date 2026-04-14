package utils

type Node[T any] struct {
	id   int
	val  T
	next *Node[T]
}

type LinkedList[T any] struct {
	head *Node[T]
	tail *Node[T]
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{head: nil, tail: nil}
}

func (l *LinkedList[T]) Push(val T) int {
	if l.tail == nil {
		new := Node[T]{id: 1, val: val, next: nil}
		l.head = &new
		l.tail = &new
	} else {
		new := Node[T]{id: l.tail.id + 1, val: val, next: nil}
		l.tail.next = &new
		l.tail = &new
	}
	return l.tail.id
}

func (l *LinkedList[T]) Pop() T {
	pop := l.head.val
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	}
	return pop
}

func (l *LinkedList[T]) Del(id int) {
	var prev *Node[T] = nil
	node := l.head
	for node.id < id {
		prev = node
		node = node.next
	}
	if prev == nil {
		l.head = node.next
	} else {
		prev.next = node.next
	}
	if l.tail.id == id {
		l.tail = prev
	}
	if l.head == nil {
		l.tail = nil
	}
}

func (l *LinkedList[T]) IsEmpty() bool {
	return l.head == nil
}
