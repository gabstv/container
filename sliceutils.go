package container

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// SliceInsert only copies elements
// in a[i:] once and allocates at most once.
// But, as of Go toolchain 1.16, due to lacking of
// optimizations to avoid elements clearing in the
// "make" call, the verbose way is not always faster.
//
// Future compiler optimizations might implement
// both in the most efficient ways.
func SliceInsert[T any](s []T, k int, vs ...T) []T {
	if k >= len(s) || k == -1 {
		return append(s, vs...)
	}
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]T, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}

// SliceBinarySearch returns the index of the first element. If the element is not found,
// it returns where the element would be inserted.
func SliceBinarySearch[T any, TV constraints.Ordered](s []T, getv func(T) TV, v TV) (pos int, exists bool) {
	x := sort.Search(len(s), func(i int) bool {
		return getv(s[i]) >= v
	})
	if x < len(s) && getv(s[x]) == v {
		// x is present at data[i]
		return x, true
	}
	// x is not present in data,
	// but i is the index where it would be inserted.
	return x, false
}
