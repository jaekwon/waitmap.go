// Package waitmap implements a simple thread-safe map.
package waitmap

import (
	"sync"
	"sync/atomic"
)

// An entry in a WaitMap.
type entry struct {
	cond *sync.Cond
	data interface{}
	ok   uint32
}

type WaitMap struct {
	lock *sync.Mutex
	ents map[interface{}]*entry
}

func New() *WaitMap {
	return &WaitMap{
		lock: new(sync.Mutex),
		ents: make(map[interface{}]*entry),
	}
}

func NewCap(c int) *WaitMap {
	return &WaitMap{
		lock: new(sync.Mutex),
		ents: make(map[interface{}]*entry, c),
	}
}

// Retrieves the value mapped to by k. If no such value is yet available, waits
// until one is, and then returns that. Otherwise, it returns the value
// available at the time Get is called.
func (m *WaitMap) Get(k interface{}) interface{} {
	m.lock.Lock()
	e, ok := m.ents[k]
	if !ok {
		e = new(entry)
		e.cond = sync.NewCond(new(sync.Mutex))
		m.ents[k] = e
	}
	m.lock.Unlock()

	e.cond.L.Lock()
	for atomic.LoadUint32(&e.ok) != 1 {
		e.cond.Wait()
	}
	e.cond.L.Unlock()
	return e.data
}

// Maps the given key and value, waking any waiting calls to Get. Returns false
// (and changes nothing) if the key is already in the map.
func (m *WaitMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	e, ok := m.ents[k]
	if !ok {
		e = new(entry)
		e.cond = sync.NewCond(new(sync.Mutex))
		e.data = v
		e.ok = 1
		m.ents[k] = e
		m.lock.Unlock()
		return true
	}
	m.lock.Unlock()

	if atomic.CompareAndSwapUint32(&e.ok, 0, 1) {
		e.data = v
		e.cond.Broadcast()
		return true
	} else {
		return false
	}
}

// Returns true if k is a key in the map.
func (m *WaitMap) Check(k interface{}) bool {
	m.lock.Lock()
	e, ok := m.ents[k]
	m.lock.Unlock()
	if !ok {
		return false
	}
	return atomic.LoadUint32(&e.ok) == 1
}
