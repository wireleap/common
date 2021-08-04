// Copyright (c) 2021 Wireleap

package wlnet

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/wireleap/common/api/status"
)

// retransmit(src, dst, ec, bufsize) reads from src and writes to dst using a
// buffer of size bufsize while reporting any errors to channel ec.
func retransmit(src io.Reader, dst io.Writer, ec chan error, bufsize int) {
	buf := make([]byte, bufsize)
	_, err := io.CopyBuffer(dst, src, buf)
	ec <- err
}

// splice(ctx, src, dst, maxtime, bufsize) splices src and dst together
// end-to-end by performing a retransmit() in both directions with buffer size
// bufsize. If maxtime is not zero, connections are limited to this
// time-to-live. Can be cancelled through ctx.
func Splice(ctx context.Context, src, dst io.ReadWriteCloser, maxtime time.Duration, bufsize int) (err error) {
	if maxtime != time.Second*0 {
		dl := time.Now().Add(maxtime)

		if c, ok := src.(net.Conn); ok {
			c.SetDeadline(dl)
		} else {
			time.AfterFunc(maxtime, func() { src.Close() })
		}

		if c, ok := dst.(net.Conn); ok {
			c.SetDeadline(dl)
		} else {
			time.AfterFunc(maxtime, func() { dst.Close() })
		}
	}

	ec := make(chan error)

	go retransmit(src, dst, ec, bufsize)
	go retransmit(dst, src, ec, bufsize)

	cancelled := false
	select {
	// Regular flow
	case err = <-ec:
		st := &status.T{}
		if err != nil && errors.As(err, &st) {
			log.Printf("splice error: %s", err)
		}
	// Cancel flow
	case <- ctx.Done():
		err = nil
		cancelled = true
	}

	// interrupt the connection on errors or EOFs from either side
	// this is strict but does not let connections linger
	dst.Close()
	src.Close()

	// wait for stream termination
	<-ec
	if cancelled {
		<-ec
	}

	return
}
