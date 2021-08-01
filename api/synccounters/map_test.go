// Copyright (c) 2021 Wireleap

package synccounters

import (
	"testing"
)

func testMap(t *testing.T, newMap func() Map) {
	m := newMap()
	x := m.GetOrInit("key1")
	x_ := m.GetOrInit("key1")

	if x != x_ {
		t.Error("Incorrect address")
	}

	m.GetOrInit("key2")

	var counter int

	f := func(key string, value *uint64) bool {
		counter += 1
		if counter == 2 {
			return false
		}
		return true
	}

	if m.Range(f) {
		t.Error("Iteration should have been aborted")
	}

	f = func(key string, value *uint64) bool {
		return true
	}

	if !m.Range(f) {
		t.Error("Iteration shouldn't have been aborted")
	}
}

func testMapReset(t *testing.T, newMap func() Map) {
	test_map := map[string]uint64{
		"key1": uint64(5),
		"key2": uint64(5),
	}

	m := newMap()

	// Case1: Initialised synccounters.Map.
	for k, v := range test_map {
		x := m.GetOrInit(k)
		*x = v
	}

	// Case2: This counter has been pushed in a RWC pipe.
	x := m.GetOrInit("key1")

	if *x != 5 {
		t.Error("Item has worng value")
	}

	ms, ok := m.Reset()

	// Case3: Poped reset result matches original synccounters.Map.
	if !ok {
		t.Error("Iteration shouldn't have been aborted")
	} else if len(ms) != len(test_map) {
		t.Error("Map length should match")
	}

	// Case3: Poped reset result matches original synccounters.Map.
	for k, v := range ms {
		if test_map[k] != v {
			t.Error("Values should match")
		}
	}

	// Case1: Reset synccounters.Map.
	m.Range(func(_ string, value *uint64) bool {
		if *value != 0 {
			t.Error("Item in map wasn't reset")
		}
		return true
	})

	// Case2: Delegated counter has also been reset.
	if *x != 0 {
		t.Error("Item has worng value")
	}
}

func TestCMap(t *testing.T) {
	testMap(t, NewCMap)
	testMapReset(t, NewCMap)
}

func TestAtomicList(t *testing.T) {
	testMap(t, NewList)
	testMapReset(t, NewCMap)
}
