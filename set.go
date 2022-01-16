package container

import (
	"errors"
	"sync"
)

var (
	ErrItemExists = errors.New("item already exists")
)

type Set[T comparable] struct {
	lock sync.RWMutex
	m    map[T]struct{}
}

func (s *Set[T]) Add(item T) error {
	s.lock.RLock()
	if s.m == nil {
		s.lock.RUnlock()
		s.lock.Lock()
		s.m = make(map[T]struct{})
		s.lock.Unlock()
		s.lock.RLock()
	}
	_, ok := s.m[item]
	s.lock.RUnlock()
	if ok {
		return ErrItemExists
	}
	s.lock.Lock()
	s.m[item] = struct{}{}
	s.lock.Unlock()
	return nil
}

func (s *Set[T]) Remove(item T) {
	s.lock.Lock()
	if s.m != nil {
		delete(s.m, item)
	}
	s.lock.Unlock()
}

func (s *Set[T]) Each(fn func(item T)) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.m != nil {
		for item := range s.m {
			fn(item)
		}
	}
}

func (s *Set[T]) Contains(item T) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if s.m == nil {
		return false
	}
	_, ok := s.m[item]
	return ok
}
