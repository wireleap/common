// Copyright (c) 2021 Wireleap

package synccounters

import (
	"sync"
	"sync/atomic"
)

type atomicList struct {
	atomic.Value
	sync.Mutex
}

func NewList() Map {
	r := &atomicList{}
	r.Value.Store([]Tuple{})
	return r
}

func (m *atomicList) load(key string) (value *uint64, ok bool) {
	for _, t := range m.Value.Load().([]Tuple) {
		if t.Key == key {
			value, ok = t.Val, true
			break
		}
	}
	return
}

func (m *atomicList) append(key string, value *uint64) {
	m.Lock()
	l := m.Value.Load().([]Tuple)
	m.Value.Store(append(l, Tuple{Key: key, Val: value}))
	m.Unlock()
}

func (m *atomicList) GetOrInit(key string) *uint64 {
	v, ok := m.load(key)
	if ok {
		return v
	}
	var zero uint64 = 0
	m.append(key, &zero)
	return &zero
}

func (m *atomicList) Range(f func(key string, value *uint64) bool) bool {
	for _, t := range m.Value.Load().([]Tuple) {
		if !f(t.Key, t.Val) {
			return false
		}
	}
	return true
}

func (m *atomicList) Reset() (map[string]uint64, bool) {
	return resetMap(m)
}
