package container

import "sort"

// MI is the type returned when MapToSlice is called.
type MI[KT comparable, VT any] struct {
	Key   KT
	Value VT
}

func MapToSlice[KT comparable, VT any](m map[KT]VT) []MI[KT, VT] {
	var ret []MI[KT, VT]
	for k, v := range m {
		ret = append(ret, MI[KT, VT]{k, v})
	}
	return ret
}

func SliceToMap[KT comparable, VT any](s []MI[KT, VT]) map[KT]VT {
	ret := make(map[KT]VT)
	for _, v := range s {
		ret[v.Key] = v.Value
	}
	return ret
}

func Sort[T any](items []T, lessfn func(a, b T) bool) {
	sort.Slice(items, func(i, j int) bool {
		return lessfn(items[i], items[j])
	})
}
