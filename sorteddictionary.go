package container

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

type SortedDictionary[KT constraints.Ordered, VT any] struct {
	lock  sync.RWMutex
	items []sortedDictionaryItem[KT, VT]
}

type sortedDictionaryItem[KT constraints.Ordered, VT any] struct {
	key KT
	val VT
}

func (m *SortedDictionary[KT, VT]) sort() {
	if len(m.items) < 2 {
		return
	}
	sort.Slice(m.items, func(i, j int) bool {
		return m.items[i].key < m.items[j].key
	})
}

// binary search
func (m *SortedDictionary[KT, VT]) search(k KT) (int, bool) {
	x := sort.Search(len(m.items), func(i int) bool {
		return m.items[i].key >= k
	})
	if x < len(m.items) && m.items[x].key == k {
		// x is present at data[i]
		return x, true
	}
	// x is not present in data,
	// but i is the index where it would be inserted.
	return x, false
}

// Set sets a key=value pair in the map.
func (m *SortedDictionary[KT, VT]) Set(k KT, v VT) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.items == nil {
		m.items = make([]sortedDictionaryItem[KT, VT], 0)
	}
	if len(m.items) == 0 {
		m.items = append(m.items, sortedDictionaryItem[KT, VT]{k, v})
		return
	}
	x, ok := m.search(k)
	if ok {
		m.items[x].val = v
		return
	}
	// add at x
	if x == len(m.items) {
		m.items = append(m.items, sortedDictionaryItem[KT, VT]{k, v})
		return
	}
	m.items = append(m.items[:x], append([]sortedDictionaryItem[KT, VT]{{key: k, val: v}}, m.items[x:]...)...)
}

func (m *SortedDictionary[KT, VT]) Get(k KT) VT {
	m.lock.RLock()
	defer m.lock.RUnlock()
	i, ok := m.search(k)
	if !ok {
		var d VT
		return d
	}
	return m.items[i].val
}

// Contains returns true if the dictionary contains the key.
func (m *SortedDictionary[KT, VT]) Contains(k KT) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.search(k); ok {
		return true
	}
	return false
}

// Remove deletes the key from the dictionary.
func (m *SortedDictionary[KT, VT]) Remove(k KT) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if i, ok := m.search(k); ok {
		m.items = append(m.items[:i], m.items[i+1:]...)
		return true
	}
	return false
}

// Pop removes and returns the value of the first item.
func (m *SortedDictionary[KT, VT]) Pop() VT {
	m.lock.Lock()
	defer m.lock.Unlock()
	if len(m.items) == 0 {
		var zv VT
		return zv
	}
	v := m.items[0].val
	m.items = m.items[1:]
	return v
}

func (m *SortedDictionary[KT, VT]) Clear() {
	m.lock.Lock()
	m.items = make([]sortedDictionaryItem[KT, VT], 0)
	m.lock.Unlock()
}

// Each calls the given function for each key=value pair in the map.
// It creates a copy of the map to iterate, so setting a key inside
// the loop will not affect the iteration.
func (m *SortedDictionary[KT, VT]) Each(fn func(KT, VT) bool) {
	m.lock.RLock()
	if len(m.items) < 1 {
		m.lock.RUnlock()
		return
	}
	m2 := make([]sortedDictionaryItem[KT, VT], len(m.items))
	copy(m2, m.items)
	m.lock.RUnlock()
	for _, v := range m2 {
		if !fn(v.key, v.val) {
			return
		}
	}
}

// Values returns a slice copy of the values.
func (m *SortedDictionary[KT, VT]) Values() []VT {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.items) < 1 {
		return []VT{}
	}
	m2 := make([]VT, len(m.items))
	for k, v := range m.items {
		m2[k] = v.val
	}
	return m2
}

// Keys returns a slice copy of the keys.
func (m *SortedDictionary[KT, VT]) Keys() []KT {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if len(m.items) < 1 {
		return []KT{}
	}
	m2 := make([]KT, len(m.items))
	for k, v := range m.items {
		m2[k] = v.key
	}
	return m2
}

// Index returns the index of the key in the dictionary.
func (m *SortedDictionary[KT, VT]) Index(k KT) int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	i, ok := m.search(k)
	if !ok {
		return -1
	}
	return i
}
