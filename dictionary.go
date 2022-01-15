package container

import "sync"

// Dictionary is a thread-safe dictionary.
type Dictionary[KT comparable, VT any] struct {
	lock sync.RWMutex
	m    map[KT]VT
}

// Set sets a key=value pair in the map.
func (m *Dictionary[KT, VT]) Set(k KT, v VT) {
	m.lock.Lock()
	if m.m == nil {
		m.m = make(map[KT]VT)
	}
	m.m[k] = v
	m.lock.Unlock()
}

func (m *Dictionary[KT, VT]) Get(k KT) VT {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.m == nil {
		var d VT
		return d
	}
	return m.m[k]
}

// Contains returns true if the dictionary contains the key.
func (m *Dictionary[KT, VT]) Contains(k KT) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.m == nil {
		return false
	}
	_, ok := m.m[k]
	return ok
}

// Remove deletes the key from the dictionary.
func (m *Dictionary[KT, VT]) Remove(k KT) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.m[k]
	if !ok {
		return false
	}
	delete(m.m, k)
	return true
}

// Pop removes and returns the value of the key.
func (m *Dictionary[KT, VT]) Pop(k KT) (VT, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.m[k]
	if !ok {
		var zv VT
		return zv, false
	}
	delete(m.m, k)
	return v, true
}

func (m *Dictionary[KT, VT]) Clear() {
	m.lock.Lock()
	m.m = make(map[KT]VT)
	m.lock.Unlock()
}

// Each calls the given function for each key=value pair in the map.
// It creates a copy of the map to iterate, so setting a key inside
// the loop will not affect the iteration.
func (m *Dictionary[KT, VT]) Each(fn func(KT, VT) bool) {
	m.lock.RLock()
	if m.m == nil {
		m.lock.RUnlock()
		return
	}
	m2 := make(map[KT]VT)
	for k, v := range m.m {
		m2[k] = v
	}
	m.lock.RUnlock()
	for k, v := range m2 {
		if !fn(k, v) {
			return
		}
	}
}

// Map returns a map copy of the dictionary.
func (m *Dictionary[KT, VT]) Map() map[KT]VT {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if m.m == nil {
		return nil
	}
	m2 := make(map[KT]VT)
	for k, v := range m.m {
		m2[k] = v
	}
	return m2
}
