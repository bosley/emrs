package core

import (
	"sync"
)

// --

type Set map[string]bool

func NewSet() Set {
	return make(map[string]bool)
}

func (s Set) Contains(x string) bool {
	_, o := s[x]
	return o
}

func (s Set) Insert(x string) {
	s[x] = true
}

func (s Set) Remove(x string) {
	delete(s, x)
}

// --

type Stack[V comparable] struct {
	lock     sync.Mutex
	s        []V
	emptyVal V
}

func NewStack[V comparable](emptyVal V) *Stack[V] {
	return &Stack[V]{
		lock:     sync.Mutex{},
		s:        make([]V, 0),
		emptyVal: emptyVal,
	}
}

func (s *Stack[V]) Push(v V) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.s = append(s.s, v)
}

func (s *Stack[V]) Pop() (V, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	l := len(s.s)
	if l == 0 {
		return s.emptyVal, NErr("empty stack")
	}
	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

// --

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Iter[V any] struct {
  Idx   int
  Value V
}

func iterate[V any](list []V, fn func(it Iter[V]) error) error {
	for i, v := range list {
    if err := fn(Iter[V]{ Idx: i, Value: v}); err != nil {
			return err
		}
	}
	return nil
}

func forEach[V any](list []V, fn func(i int, x V) error) error {
	for i, v := range list {
		if err := fn(i, v); err != nil {
			return err
		}
	}
	return nil
}

func contains[V any](in []V, lookFor V, isSame func(i int, l V, r V) bool) bool {
	for i, v := range in {
		if isSame(i, lookFor, v) {
			return true
		}
	}
	return false
}

func isUnique[V any](target []V, isSame func(i int, l V, r V) bool) bool {
	for value_idx, value := range target {
		if contains(target, value, func(scan_idx int, l V, r V) bool {
			return isSame(scan_idx, l, r) && scan_idx != value_idx
		}) {
			return false
		}
	}
	return true
}
