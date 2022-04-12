// Copyright (c) 2022 Wireleap

// Package h2conn implements a net.Conn reads and writes over which are
// actually directed towards a h/2 stream.
package h2conn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/wireleap/common/api/status"
)

// T is the type of h/2 overlay connections.
type T struct {
	net.Conn
	io.ReadCloser
	io.WriteCloser

	e  chan error
	er error

	resp   *http.Response
	once   sync.Once
	cancel context.CancelFunc
	dl     *time.Timer
	mu     sync.Mutex
}

// New creates a new T given a http.Roundtripper and a remote URL string to
// connect to via h/2 as well as any headers that are needed.
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
		c.mu.Lock()
		defer c.mu.Unlock()
		if err == nil {
			c.resp = res
			c.ReadCloser = res.Body
		} else {
			c.cancel()
		}
		c.e <- err
		close(c.e)
	}()

	return
}

func (c *T) check() error {
	if err, ok := <-c.e; ok {
		if err == nil {
			return nil
		} else {
			c.mu.Lock()
			c.er = err
			c.mu.Unlock()
			return err
		}
	} else {
		return c.er
	}
}

// Write writes the given data to a pipe the other end of which is read as the
// sent request body.
func (c *T) Write(p []byte) (int, error) {
	return c.WriteCloser.Write(p)
}

// Read makes the initial request if needed, after which it reads from the
// response body if no error was encountered.
func (c *T) Read(p []byte) (int, error) {
	if err := c.check(); err != nil {
		return 0, err
	}
	n, err := c.ReadCloser.Read(p)
	if err != nil && c.resp != nil && c.resp.Trailer != nil {
		sth := c.resp.Trailer.Get(status.Header)
		if sth != "" {
			var st status.T
			if err = json.Unmarshal([]byte(sth), &st); err != nil {
				return 0, fmt.Errorf("error while unmarshaling status trailer: %w", err)
			}
			if st.Is(status.OK) {
				return n, nil
			}
			return n, &st
		}
	}
	return n, err
}

func (c *T) SetDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dl != nil {
		c.dl.Stop()
	}
	if t.IsZero() {
		return nil
	}
	c.dl = time.AfterFunc(time.Until(t), func() { c.Close() })
	return nil
}

// Close closes both the sent request body and the response body.
func (c *T) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cancel()
	if c.ReadCloser != nil {
		c.ReadCloser.Close()
	}
	if c.WriteCloser != nil {
		c.WriteCloser.Close()
	}
	return nil
}

// TODO?
func (c *T) SetReadDeadline(t time.Time) error  { return c.SetDeadline(t) }
func (c *T) SetWriteDeadline(t time.Time) error { return c.SetDeadline(t) }
func (c *T) RemoteAddr() net.Addr               { return nil }
func (c *T) LocalAddr() net.Addr                { return nil }
