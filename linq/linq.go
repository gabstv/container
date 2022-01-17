// Package linq provides methods for querying and manipulating slices.
package linq

import "sort"

// Query is the type returned from query functions.
type Query[T, TR any] struct {
	source []T
}

// From initializes a linq query with passed slice as the source.
// It performs a copy of the original slice so it can preserve the order
// of the original slice.
func From[T, TR any](slc []T) Query[T, TR] {
	slc2 := make([]T, len(slc))
	copy(slc2, slc)
	return Query[T, TR]{source: slc2}
}

// Where uses the predicate to filter the source slice.
func (q Query[T, TR]) Where(pred func(T) bool) Query[T, TR] {
	q2 := make([]T, 0, len(q.source))
	for _, item := range q.source {
		if pred(item) {
			q2 = append(q2, item)
		}
	}
	return Query[T, TR]{source: q2}
}

func (q Query[T, TR]) All() []T {
	return q.source
}

func (q Query[T, TR]) First() T {
	if len(q.source) < 1 {
		var zv T
		return zv
	}
	return q.source[0]
}

// Result is the type returned from query functions after Select() is called.
type Result[TR any] struct {
	value []TR
}

func (q Query[T, TR]) Select(fn func(T) TR) Result[TR] {
	q2 := make([]TR, 0, len(q.source))
	for _, item := range q.source {
		q2 = append(q2, fn(item))
	}
	return Result[TR]{value: q2}
}

func (r Query[T, TR]) Sort(lessfn func(T, T) bool) Query[T, TR] {
	sort.Slice(r.source, func(i, j int) bool {
		return lessfn(r.source[i], r.source[j])
	})
	return r
}

func (r Result[TR]) Sort(lessfn func(TR, TR) bool) Result[TR] {
	sort.Slice(r.value, func(i, j int) bool {
		return lessfn(r.value[i], r.value[j])
	})
	return r
}

func (r Result[TR]) All() []TR {
	return r.value
}

func (q Result[TR]) First() TR {
	if len(q.value) < 1 {
		var zv TR
		return zv
	}
	return q.value[0]
}
