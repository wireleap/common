// Copyright (c) 2022 Wireleap

// Package fsdir provides an abstract interface to a directory on disk.
package fsdir

// Compare two paths and return if they're equal
func PathEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
