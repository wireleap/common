// Copyright (c) 2021 Wireleap

// Package h2rwc implements a io.ReadWriteCloser composed from an io.Writer and
// io.ReadCloser (most probably, http.ResponseWriter and server-side
// http.Request.Body).
package h2rwc

import "io"

// T is a composite io.ReadWriteCloser from io.Writer and io.ReadCloser.
type T struct {
	io.Writer
	io.ReadCloser
}
