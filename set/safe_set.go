package set

import (
	"fmt"
	"strings"
	"sync"
)

type SafeSet[K comparable] struct {
	m    map[K]bool
	lock sync.RWMutex
}

func NewSafe[K comparable]() *SafeSet[K] {
	return &SafeSet[K]{m: make(map[K]bool)}
}

func FromListSafe[K comparable](list []K) *SafeSet[K] {
	m := make(map[K]bool)
	for _, value := range list {
		m[value] = true
	}
	return &SafeSet[K]{m: m}
}

func (s *SafeSet[K]) Add(value K) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.m[value] = true
}

func (s *SafeSet[K]) Delete(value K) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.m, value)
}

func (s *SafeSet[K]) Has(value K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, found := s.m[value]
	return found
}

func (s *SafeSet[K]) List() []K {
	s.lock.RLock()
	defer s.lock.RUnlock()
	list := make([]K, len(s.m))
	i := 0
	for value := range s.m {
		list[i] = value
		i++
	}
	return list
}

func (s *SafeSet[K]) Iter(callback func(value K)) {
	values := s.List()
	for _, value := range values {
		callback(value)
	}
}

func (s *SafeSet[K]) IterIndex(callback func(value K, i int)) {
	values := s.List()
	for i, value := range values {
		callback(value, i)
	}
}

func (s *SafeSet[K]) Size() int {
	return len(s.m)
}

func (s *SafeSet[K]) Union(other *SafeSet[K]) *SafeSet[K] {
	res := NewSafe[K]()
	s.Iter(func(value K) { res.Add(value) })
	other.Iter(func(value K) { res.Add(value) })
	return res
}

func (s *SafeSet[K]) Intersection(other *SafeSet[K]) *SafeSet[K] {
	if s.Size() > other.Size() {
		s, other = other, s
	}

	res := NewSafe[K]()
	s.Iter(func(value K) {
		if other.Has(value) {
			res.Add(value)
		}
	})

	return res
}

func (s *SafeSet[K]) Copy() *SafeSet[K] {
	res := NewSafe[K]()
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k := range s.m {
		res.m[k] = true
	}
	return res
}

func (s *SafeSet[K]) Subtract(other *SafeSet[K]) *SafeSet[K] {
	res := NewSafe[K]()
	values := s.List()
	for _, k := range values {
		if !other.Has(k) {
			res.m[k] = true
		}
	}
	return res
}

func (s *SafeSet[K]) Difference(other *SafeSet[K]) *SafeSet[K] {
	return s.Subtract(other)
}

func (s *SafeSet[K]) Eq(other *SafeSet[K]) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if len(s.m) != len(other.m) {
		return false
	}
	for k := range s.m {
		if _, found := other.m[k]; !found {
			return false
		}
	}
	return true
}

func (s *SafeSet[K]) String() string {
	res := make([]string, len(s.m))
	i := 0

	s.lock.RLock()
	for k := range s.m {
		res[i] = fmt.Sprintf("%v", k)
		i++
	}
	s.lock.RUnlock()

	return fmt.Sprintf("*set.Set{%s}", strings.Join(res, ", "))
}
