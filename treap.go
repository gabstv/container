// MIT License
//
// Copyright (c) 2018 perdata
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package container

type CompareFn[T any] func(left, right T) int

// // Comparer compares two values. The return value is zero if the
// // values are equal, negative if the first is smaller and positive
// // otherwise.
// type Comparer[T any] interface {
// 	Compare(left, right T) int
// }

// Treap is the basic recursive treap data structure
// Package treap implements a persistent treap (tree/heap combination).
//
// https://en.wikipedia.org/wiki/Treap
//
// A treap is a binary search tree for storing ordered distinct values
// (duplicates not allowed). In addition, each node actually a random
// priority field which is stored in heap order (i.e. all children
// have lower priority than the parent)
//
// This provides the basis for efficient immutable ordered Set
// operations.  See the ordered map example for how this can be used
// as an ordered map
//
// Much of this is based on "Fast Set Operations Using Treaps"
// by Guy E Blelloch and Margaret Reid-Miller:
// https://www.cs.cmu.edu/~scandal/papers/treaps-spaa98.pdf
//
// Benchmark
//
// The most interesting benchmark is the performance of insert where a
// single random key is inserted into a 5k sized map.  As the example
// shows, the treap structure does well here as opposed to a regular
// persistent map (which involves full copying).  This benchmark does not
// take into account the fact that the regular maps are not sorted unlike
// treaps.
//
// The intersection benchmark compares the case where two 10k sets with
// 5k in common being interesected. The regular persistent array is about
// 30% faster but this is still respectable showing for treaps.
//
//    $ go test --bench=. -benchmem
//    goos: darwin
//    goarch: amd64
//    cpu: Intel(R) Core(TM) i5-8210Y CPU @ 1.60GHz
//    BenchmarkInsert-4                         650380              1925 ns/op            1146 B/op         35 allocs/op
//    BenchmarkInsertRegularMap-4                 1627            755162 ns/op          336394 B/op          9 allocs/op
//    BenchmarkIntersection-4                      526           2428991 ns/op         1144594 B/op      35734 allocs/op
//    BenchmarkIntersectionRegularMap-4            475           2484496 ns/op          626639 B/op        124 allocs/op
//    BenchmarkUnion-4                             970           1439862 ns/op          626906 B/op      19579 allocs/op
//    BenchmarkDiff-4                              391           2749617 ns/op         1172926 B/op      36608 allocs/op
//    PASS
//
type Treap[T any] struct {
	Value       T
	Priority    int
	Left, Right *Treap[T]
}

// ForEach does inorder traversal of the treap
func (n *Treap[T]) ForEach(fn func(v T)) {
	if n != nil {
		n.Left.ForEach(fn)
		fn(n.Value)
		n.Right.ForEach(fn)
	}
}

// Find finds the node in the treap with matching value
func (n *Treap[T]) Find(v T, c CompareFn[T]) *Treap[T] {
	for {
		if n == nil {
			return nil
		}
		diff := c(n.Value, v)
		switch {
		case diff == 0:
			return n
		case diff < 0:
			n = n.Right
		case diff > 0:
			n = n.Left
		}
	}
}

// Union combines any two treaps. In case of duplicates, the overwrite
// field controls whether the union keeps the original value or
// whether it is updated based on value in the "other" arg
func (n *Treap[T]) Union(other *Treap[T], c CompareFn[T], overwrite bool) *Treap[T] {
	if n == nil {
		return other
	}
	if other == nil {
		return n
	}

	if n.Priority < other.Priority {
		other, n, overwrite = n, other, !overwrite
	}

	left, dupe, right := other.Split(n.Value, c)
	value := n.Value
	if overwrite && dupe != nil {
		value = dupe.Value
	}
	left = n.Left.Union(left, c, overwrite)
	right = n.Right.Union(right, c, overwrite)
	return &Treap[T]{value, n.Priority, left, right}
}

// Split splits the treap into all nodes that compare less-than, equal
// and greater-than the provided value.  The resulting values are
// properly formed treaps or nil if they contain no values.
func (n *Treap[T]) Split(v T, c CompareFn[T]) (left, mid, right *Treap[T]) {
	leftp, rightp := &left, &right
	for {
		if n == nil {
			*leftp = nil
			*rightp = nil
			return left, nil, right
		}

		root := &Treap[T]{n.Value, n.Priority, nil, nil}
		diff := c(n.Value, v)
		switch {
		case diff < 0:
			*leftp = root
			root.Left = n.Left
			leftp = &root.Right
			n = n.Right
		case diff > 0:
			*rightp = root
			root.Right = n.Right
			rightp = &root.Left
			n = n.Left
		default:
			*leftp = n.Left
			*rightp = n.Right
			return left, root, right
		}
	}
}

// Intersection returns a new treap with all the common values in the
// two treaps.
//
// see https://www.cs.cmu.edu/~scandal/papers/treaps-spaa98.pdf
// "Fast Set Operations Using Treaps"
//   by Guy E Blelloch and Margaret Reid-Miller.
//
// The algorithm is a very slight variation on that.
func (n *Treap[T]) Intersection(other *Treap[T], c CompareFn[T]) *Treap[T] {
	if n == nil || other == nil {
		return nil
	}

	if n.Priority < other.Priority {
		n, other = other, n
	}

	left, found, right := other.Split(n.Value, c)
	left = n.Left.Intersection(left, c)
	right = n.Right.Intersection(right, c)

	if found == nil {
		// TODO: use a destructive join as both left/right are copies
		return left.join(right)
	}

	return &Treap[T]{n.Value, n.Priority, left, right}
}

// Delete removes a node if it exists.
func (n *Treap[T]) Delete(v T, c CompareFn[T]) *Treap[T] {
	left, _, right := n.Split(v, c)
	return left.join(right)
}

// Diff finds all elements of current treap which aren't present in
// the other heap
func (n *Treap[T]) Diff(other *Treap[T], c CompareFn[T]) *Treap[T] {
	if n == nil || other == nil {
		return n
	}

	// TODO -- use  count
	if n.Priority >= other.Priority {
		left, dupe, right := other.Split(n.Value, c)
		left, right = n.Left.Diff(left, c), n.Right.Diff(right, c)
		if dupe != nil {
			return left.join(right)
		}
		return &Treap[T]{n.Value, n.Priority, left, right}
	}

	left, _, right := n.Split(other.Value, c)
	left = left.Diff(other.Left, c)
	right = right.Diff(other.Right, c)
	return left.join(right)
}

// see https://www.cs.cmu.edu/~scandal/papers/treaps-spaa98.pdf
// "Fast Set Operations Using Treaps"
//   by Guy E Blelloch and Margaret Reid-Miller.
//
// The algorithm is a very slight variation on that provided there.
//
// Note that all nodes in n have priority <= that of "other" for
// this call to work correctly.  It traverses  the right spine of n
// and left-spine of other, merging things along the way
//
// The algorithm is not that  different from zipping up a spine
func (n *Treap[T]) join(other *Treap[T]) *Treap[T] {
	var result *Treap[T]
	resultp := &result
	for {
		if n == nil {
			*resultp = other
			return result
		}
		if other == nil {
			*resultp = n
			return result
		}

		if n.Priority <= other.Priority {
			root := &Treap[T]{n.Value, n.Priority, n.Left, nil}
			*resultp = root
			resultp = &root.Right
			n = n.Right
		} else {
			root := &Treap[T]{other.Value, other.Priority, nil, other.Right}
			*resultp = root
			resultp = &root.Left
			other = other.Left
		}
	}
}

// Union combines two or more treaps. In case of duplicates, the overwrite field
// controls whether the union keeps the original value or whether it is updated
// based on value in the items
//
func TreapUnion[T any](comparer CompareFn[T], priority int, items ...T) *Treap[T] {
	var t *Treap[T]
	for _, elt := range items {
		t = t.Union(&Treap[T]{elt, priority, nil, nil}, comparer, false)
	}
	return t
}
