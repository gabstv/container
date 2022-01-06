package container

type node[T comparable] struct {
	data   T
	prev   *node[T]
	next   *node[T]
	parent *LinkedList[T]
}

func (n *node[T]) Data() T {
	if n == nil {
		var d T
		return d
	}
	return n.data
}

func (n *node[T]) Prev() Node[T] {
	if n == nil {
		return nil
	}
	return n.prev
}

func (n *node[T]) Next() Node[T] {
	if n == nil {
		return nil
	}
	return n.next
}

func (n *node[T]) Remove() bool {
	if n.parent == nil {
		return false
	}
	var prev *node[T]
	if n.prev != nil {
		prev = n.prev
	}
	var next *node[T]
	if n.next != nil {
		next = n.next
	}
	if prev != nil {
		prev.next = next
	}
	if next != nil {
		next.prev = prev
	}
	if n.parent.tail == n {
		n.parent.tail = prev
	}
	if n.parent.head == n {
		n.parent.head = next
	}
	n.parent.length--
	n.next = nil
	n.prev = nil
	n.parent = nil
	return true
}

type Node[T comparable] interface {
	Data() T
	Prev() Node[T]
	Next() Node[T]
	Remove() bool
}

type LinkedList[T comparable] struct {
	length int
	head   *node[T]
	tail   *node[T]
}

func (ll *LinkedList[T]) Len() int {
	return ll.length
}

func (ll *LinkedList[T]) First() Node[T] {
	return ll.head
}

func (ll *LinkedList[T]) Last() Node[T] {
	return ll.tail
}

func (ll *LinkedList[T]) Contains(d T) bool {
	for n := ll.First(); n != nil; n = n.Next() {
		if n != nil && n.Data() == d {
			return true
		}
	}
	return false
}

func (ll *LinkedList[T]) Remove(d T) bool {
	for n := ll.First(); n != nil; n = n.Next() {
		if n.Data() == d {
			n.(*node[T]).Remove()
			return true
		}
	}
	return false
}

// Pop removes the last item of the list and returns its data.
func (ll *LinkedList[T]) Pop() T {
	var d T
	if ll.length == 0 || ll.tail == nil {
		return d
	}
	d = ll.tail.data
	_ = ll.tail.Remove()
	return d
}

// Shift removes the first item of the list and returns its data.
func (ll *LinkedList[T]) Shift() T {
	var d T
	if ll.length == 0 || ll.head == nil {
		return d
	}
	d = ll.head.data
	_ = ll.head.Remove()
	return d
}

// Push adds an item to the end of the list.
func (ll *LinkedList[T]) Push(data T) Node[T] {
	n := &node[T]{
		data:   data,
		parent: ll,
	}
	if ll.length == 0 {
		ll.length = 1
		ll.head = n
		ll.tail = n
		return n
	}
	prevItem := ll.tail
	n.prev = prevItem
	prevItem.next = n
	ll.tail = n
	ll.length++
	return n
}

// Unshift adds an item to the beginning of the list.
func (ll *LinkedList[T]) Unshift(data T) Node[T] {
	n := &node[T]{
		data:   data,
		parent: ll,
	}
	if ll.length == 0 {
		ll.length = 1
		ll.head = n
		ll.tail = n
		return n
	}
	nextItem := ll.head
	n.next = nextItem
	nextItem.prev = n
	ll.head = n
	ll.length++
	return n
}

// Clear removes all items from the list.
func (ll *LinkedList[T]) Clear() {
	if ll.length == 0 {
		return
	}
	items := make([]Node[T], 0, ll.length)
	for n := ll.First(); n != nil; n = n.Next() {
		if n != nil {
			items = append(items, n)
		}
	}
	for _, item := range items {
		if item != nil {
			n := item.(*node[T])
			if n != nil {
				n.parent = nil
				n.prev = nil
				n.next = nil
			}
		}
	}
	ll.length = 0
	ll.head = nil
	ll.tail = nil
}
