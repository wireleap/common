// Copyright (c) 2021 Wireleap

// The ststore package provides a concurrent in-memory sharetoken store which
// is synced to disk after modifications.
package ststore

func adaptedPath(rootFolder string, ps ...string) []string {
	if len(ps) != 2 {
		return nil
	}

	return append([]string{rootFolder}, ps...)
}

func MalformedPath(ps ...string) []string {
	return adaptedPath("malformed", ps...)
}

func ExpiredPath(ps ...string) []string {
	return adaptedPath("expired", ps...)
}
