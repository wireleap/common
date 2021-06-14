// Copyright (c) 2021 Wireleap

package relay

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/wireleap/common/api/sharetoken"
	"github.com/wireleap/common/api/status"
	"github.com/wireleap/common/wlnet"
	"github.com/wireleap/common/wlnet/flushwriter"
	"github.com/wireleap/common/wlnet/h2rwc"
	"github.com/wireleap/common/wlnet/transport"
)

type T struct {
	*transport.T
	Options
}

type Options struct {
	// BufSize is the size in bytes of the send/receive buffers of a relay.
	BufSize int
	// MaxTime is the maximum time for a single connection.
	MaxTime time.Duration
	// HandleST is a generic function which is called on incoming sharetokens.
	HandleST func(*sharetoken.T) error
	// ErrorOrigin is an optional string to use when signaling the origin of
	// errors downstream.
	ErrorOrigin string
	// AllowLoopback sets whether to allow dialing loopback addresses. While
	// useful for testing, it presents a security risk in production.
	AllowLoopback bool
}

func New(tt *transport.T, o Options) *T { return &T{T: tt, Options: o} }

// isLoopback determines whether the presented address is a loopback interface
// address.
func isLoopback(addr string) bool {
	if addr == "localhost" {
		return true
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		// probably a fqdn
		return false
	}
	// unspecified ips (0.0.0.0/::) can be used to access loopback too
	return ip.IsLoopback() || ip.IsUnspecified()
}

// ServeTLS is the handler function for listening and relaying incoming data.
// It handles the initial init payload and brokers the subsequent tunnel
// connections or an exit connection if needed.
func (t *T) ServeTLS(c io.ReadWriteCloser) {
	defer c.Close()

	origin := t.ErrorOrigin
	p, err := wlnet.ReadInit(c)

	if err != nil {
		wlnet.WriteStatus(c, &status.T{
			Code:   http.StatusBadRequest,
			Desc:   err.Error(),
			Origin: origin,
		})
		return
	}

	if p.Command == "PING" {
		// raw, not in wlnet wire format
		(&status.T{
			Code:   http.StatusOK,
			Desc:   "PONG",
			Origin: origin,
		}).WriteTo(c)
		return
	}

	if t.HandleST != nil {
		err = t.HandleST(p.Token)

		if err != nil {
			wlnet.WriteStatus(c, &status.T{
				Code:   http.StatusBadRequest,
				Desc:   err.Error(),
				Origin: origin,
			})

			return
		}
	}

	// signal target errors differently
	if p.Remote.Scheme == "target" {
		origin = "target"
	}

	// no dials to localhost (this relay's host)
	if !t.AllowLoopback && isLoopback(p.Remote.Hostname()) {
		wlnet.WriteStatus(c, &status.T{
			Code: http.StatusBadRequest,
			Desc: fmt.Sprintf(
				"loopback address '%s' requested, refusing to dial",
				p.Remote.Hostname(),
			),
			Origin: origin,
		})
		return
	}

	// hide requested target for privacy
	shown := "(target)"

	// ok to show relays though
	if p.Remote.Scheme != "target" {
		shown = p.Remote.String()
	}

	log.Printf("Dialing %s connection to %s", p.Protocol, shown)
	c2, err := t.DialWL(p.Protocol, &p.Remote.URL)

	if err != nil {
		// TODO more granular errors

		if os.IsTimeout(err) {
			wlnet.WriteStatus(c, &status.T{
				Code:   http.StatusRequestTimeout,
				Desc:   err.Error(),
				Origin: origin,
			})
		} else {
			wlnet.WriteStatus(c, &status.T{
				Code:   http.StatusBadGateway,
				Desc:   err.Error(),
				Origin: origin,
			})
		}

		return
	}

	if p.Remote.Scheme == "target" {
		c = wlnet.FragWriteCloser{ReadWriteCloser: c}
	}

	err = wlnet.Splice(c, c2, t.MaxTime, t.BufSize)

	if err != nil {
		// TODO more granular errors

		if os.IsTimeout(err) {
			wlnet.WriteStatus(c, &status.T{
				Code:   http.StatusRequestTimeout,
				Desc:   err.Error(),
				Origin: origin,
			})
		} else {
			wlnet.WriteStatus(c, &status.T{
				Code:   http.StatusGone,
				Desc:   err.Error(),
				Origin: origin,
			})
		}
	}
}

// ServeHTTP is the handler function for H2. It being named ServeHTTP allows
// T to expose the http.Handler interface.
func (t *T) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		status.ErrMethod.WriteTo(w)
		return
	}
	// TODO process the h/2 connection in a more seamless way
	t.ServeTLS(h2rwc.T{Writer: flushwriter.T{Writer: w}, ReadCloser: r.Body})
}

// ListenAndServeHTTP listens on the specified address and passes the
// connections to ServeHTTP.
func (t *T) ListenAndServeHTTP(addr string) error {
	l, err := tls.Listen("tcp", addr, t.Transport.TLSClientConfig)
	if err != nil {
		return err
	}
	s := http.Server{
		Addr:      addr,
		Handler:   t,
		TLSConfig: t.Transport.TLSClientConfig,
	}
	go s.Serve(l)
	return nil
}
