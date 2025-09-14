package safe

import (
	"fmt"
	"sync"
)

// Please only use Map for primitive type (excluding slice and other few types).
// Map do not guard against the integrity of the data itself. It only guard 
// against the integrity of map itself. 
// Thus, it's strongly against to use it with `struct`, pointer type, interface, 
// string, and so on. The only exception of this is immutability is enforced, 
// which update is done by overwriting the entire value of a given key.

type Map[K comparable, T any] struct {
	Mutex sync.RWMutex
	Map   map[K]T
}

func AllocMap[K comparable, T any](size uint32) *Map[K, T] {
	if size <= 0 {
		return &Map[K, T]{
			Map: make(map[K]T),
		}
	}

	return &Map[K, T]{
		Map: make(map[K]T, size),
	}
}

func (m *Map[K, T]) RLock() {
	m.Mutex.RLock()
}

func (m *Map[K, T]) RUnlock() {
	m.Mutex.RUnlock()
}

func (m *Map[K, T]) Lock() {
	m.Mutex.Lock()
}

func (m *Map[K, T]) Unlock() {
	m.Mutex.Unlock()
}

// Thread safe
func (m *Map[K, T]) Get(k K) (v T, in bool) {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()

	v, in = m.Map[k]
	return v, in
}

// Thread safe
// Use this to enforce key collision free
func (m *Map[K, T]) Add(k K, v T) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, in := m.Map[k]; in {
		return fmt.Errorf("%v key already exists.", k)
	}

	m.Map[k] = v

	return nil
}

// Thread safe
// Set is used when caller knows a key `k` is in the map. Otherwise, it will 
// return error.
func (m *Map[K, T]) Set(k K, v T) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if _, in := m.Map[k]; !in {
		return fmt.Errorf("%v key does not exist.", k)
	}

	m.Map[k] = v

	return nil
}

// Thread safe 
func (m *Map[K, T]) Del(k K) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	delete(m.Map, k)
}

// Thread safe
func (m *Map[K, T]) Len() int {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	return len(m.Map)
}

// Thread safe
func (m* Map[K, T]) ForR(f func (k K, v T) bool) {
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
	for k, v := range m.Map {
		c := f(k, v)
		if !c {
			break
		}
	}
}

// Thread safe
func (m *Map[K, T]) ForW(f func (k K, v T) bool) {
	m.Mutex.Lock()
	defer m.Mutex.Lock()
	for k, v := range m.Map {
		c := f(k, v)
		if !c {
			break
		}
	}
}
