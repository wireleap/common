// Copyright (c) 2021 Wireleap

// The ststore package provides a concurrent in-memory sharetoken store which
// is synced to disk after modifications.
package ststore

func adaptedPath(subfolder string, ps ...string) []string {
	if len(ps) != 2 {
		return nil
	}

	return []string{ps[0], subfolder, ps[1]}
}

func MalformedPath(ps ...string) []string {
	return adaptedPath("malformed", ps...)
}

func ExpiredPath(ps ...string) []string {
	return adaptedPath("expired", ps...)
}
