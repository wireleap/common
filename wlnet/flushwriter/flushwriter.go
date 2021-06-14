// Copyright (c) 2021 Wireleap

package flushwriter

import (
	"io"
	"net/http"
)

// T is an io.Writer which flushes after every write.
type T struct{ io.Writer }

// Write flushes after calling the underlying io.Writer's Write.
func (t T) Write(p []byte) (n int, err error) {
	n, err = t.Writer.Write(p)
	if f, ok := t.Writer.(http.Flusher); ok {
		f.Flush()
	}
	return
}
