// Copyright (c) 2021 Wireleap

// Package h2conn implements a net.Conn reads and writes over which are
// actually directed towards a h/2 stream.
package h2conn

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

// T is the type of h/2 overlay connections.
type T struct {
	net.Conn
	io.ReadCloser
	io.WriteCloser

	e  chan error
	er error

	once   sync.Once
	cancel context.CancelFunc
	dl     *time.Timer
	mu     sync.Mutex
}

// New creates a new T given an underlying net.Conn to fall back for non-r/w
// methods and a remote URL string to connect to via h/2.
func New(t http.RoundTripper, remote string, headers map[string]string) (c *T, err error) {
	c = &T{}

	var (
		pr  *io.PipeReader
		req *http.Request
	)

	pr, c.WriteCloser = io.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	req, err = http.NewRequestWithContext(ctx, http.MethodPut, remote, pr)
	if err != nil {
		cancel()
		return
	}
	c.cancel = cancel
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	c.e = make(chan error)

	// the request will block until something is written to the pipe writer end

	go func() {
		res, err := t.RoundTrip(req)
		if err == nil {
			c.ReadCloser = res.Body
		}
		c.e <- err
	}()

	return
}

func (c *T) check() {
	c.once.Do(func() {
		c.er = <-c.e
		close(c.e)
		c.e = nil
	})
}

// Write writes the given data to a pipe the other end of which is read as the
// sent request body.
func (c *T) Write(p []byte) (int, error) {
	if c.er != nil {
		return 0, c.er
	}
	return c.WriteCloser.Write(p)
}

// Read makes the initial request if needed, after which it reads from the
// response body if no error was encountered.
func (c *T) Read(p []byte) (int, error) {
	c.check()
	if c.er != nil {
		return 0, c.er
	}
	return c.ReadCloser.Read(p)
}

func (c *T) SetDeadline(t time.Time) error {
	c.mu.Lock()
	if c.dl != nil {
		c.dl.Stop()
	}
	if t.IsZero() {
		c.mu.Unlock()
		return nil
	}
	c.dl = time.AfterFunc(time.Until(t), func() { c.Close() })
	c.mu.Unlock()
	return nil
}

// Close closes both the sent request body and the response body.
func (c *T) Close() error {
	c.cancel()
	if c.ReadCloser != nil {
		c.ReadCloser.Close()
	}
	if c.WriteCloser != nil {
		c.WriteCloser.Close()
	}
	return nil
}
