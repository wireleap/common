// Copyright (c) 2021 Wireleap

package synccounters

type Map interface {
	// Returns exisiting counter or initialises a new one
	GetOrInit(key string) (value *uint64)
	// Iterate over the entire itemlist
	// f should return false to abort iteration
	// Range(f) returns if iteration was completed
	Range(f func(key string, value *uint64) bool) bool
	// Reset values to 0
	// Returns a map[string]uint64 with the previous values, and a completion flag
	Reset() (map[string]uint64, bool)
}
