// Copyright (c) 2021 Wireleap

package synccounters

import (
	"sync/atomic"
)

func resetMapFunction() (func(string, *uint64) bool, map[string]uint64) {
	m := make(map[string]uint64, 0)
	return func(key string, value *uint64) bool {
		if _, ok := m[key]; ok {
			// In theory we're only iterating over each element once
		} else if value == nil {
			// Safety check, no address to override ==> ABORT!
			return false
		} else {
			// Save netstat status if not null
			if old_value := atomic.SwapUint64(value, uint64(0)); old_value != uint64(0) {
				m[key] = old_value
			}
		}
		return true
	}, m
}

func resetMap(m Map) (map[string]uint64, bool) {
	f, ms := resetMapFunction()
	return ms, m.Range(f)
}
