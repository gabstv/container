package container

import "sort"

type MapSliceItem[KT comparable, VT any] struct {
	Key   KT
	Value VT
}

func MapToSlice[KT comparable, VT any](m map[KT]VT) []MapSliceItem[KT, VT] {
	var ret []MapSliceItem[KT, VT]
	for k, v := range m {
		ret = append(ret, MapSliceItem[KT, VT]{k, v})
	}
	return ret
}

func SliceToMap[KT comparable, VT any](s []MapSliceItem[KT, VT]) map[KT]VT {
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
