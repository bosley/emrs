package core

import (
	"encoding/json"
	"os"
	"sync"
)

// --

type Set map[string]bool

func NewSet() Set {
	return make(map[string]bool)
}

func SetFrom(options []string) Set {
	s := NewSet()
	iterate[string](options,
		func(it Iter[string]) error {
			s.Insert(it.Value)
			return nil
		})
	return s
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

func LoadJSON[T any](filename string) (T, error) {
	var data T
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}
	return data, json.Unmarshal(fileData, &data)
}

func PathExists(path string) bool {
  exists, err := safeCheckExists(path)
  if err != nil {
    panic(err.Error())
  }
  return exists
}

func safeCheckExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func deleteIf[V any](s []V, check func(V) bool) []V {
	for i := 0; i < len(s); i++ {
		if check(s[i]) {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func mapContains[K comparable, V any](m map[K]V, k K) bool {
	_, o := m[k]
	return o
}

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
		if err := fn(Iter[V]{Idx: i, Value: v}); err != nil {
			return err
		}
	}
	return nil
}

func itermap[K comparable, V any](m map[K]V, fn func(K, V)) {
	for k, v := range m {
		fn(k, v)
	}
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

func isIn[V comparable](in []V, x V) bool {
	if len(in) == 0 {
		return false
	}
	if in[0] == x {
		return true
	}
	return isIn(in[1:], x)
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
