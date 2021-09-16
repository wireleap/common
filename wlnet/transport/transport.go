// Copyright (c) 2021 Wireleap

package transport

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/wireleap/common/wlnet"
	"github.com/wireleap/common/wlnet/h2conn"
)

// T is a complete Wireleap network transport which can dial to other
// wireleap-relays via H/2 over TCP and targets via TCP or UDP.
type T struct{ *http.Transport }

// Options is a struct which contains options for initializing a T.
type Options struct {
	// TLSVerify is the same as !InsecureSkipVerify in tls.Config
	TLSVerify bool
	// Certs is a list of TLS certificates to use
	Certs []tls.Certificate
	// Timeout is the maximum time for new connections
	Timeout time.Duration
}

// New creates a default T with the supplied options.
func New(opts Options) *T {
	var (
		tc = &tls.Config{
			Certificates:       opts.Certs,
			InsecureSkipVerify: !opts.TLSVerify,
			MinVersion:         tls.VersionTLS13,
			NextProtos:         []string{"h2"}, // H/2 only
		}
		nd = &net.Dialer{Timeout: opts.Timeout}
		td = &tls.Dialer{NetDialer: nd, Config: tc}
		t  = &T{
			Transport: &http.Transport{
				Dial:                  nd.Dial,
				DialContext:           nd.DialContext,
				DialTLS:               td.Dial,
				DialTLSContext:        td.DialContext,
				TLSClientConfig:       tc,
				ResponseHeaderTimeout: opts.Timeout,
				ForceAttemptHTTP2:     true,
				MaxConnsPerHost:       0,
				MaxIdleConnsPerHost:   0,
				MaxIdleConns:          4096,
				IdleConnTimeout:       5 * time.Minute,
			},
		}
	)
	return t
}

// DialWL creates a new connection to relay or target.
func (t *T) DialWL(c0 net.Conn, protocol string, remote *url.URL, payload *wlnet.Init) (c net.Conn, err error) {
	switch remote.Scheme {
	case "target":
		// c0/payload unused, could both be nil
		c, err = t.Transport.Dial(protocol, remote.Host)
	case "wireleap":
		tt := t.Transport
		if c0 != nil {
			// if previous connection supplied, use it to tunnel
			t2 := t.Transport.Clone()
			t2.DialTLS = func(_ string, _ string) (net.Conn, error) { return c0, nil }
			tt = t2
		}
		// convert to a stdlib-known scheme
		u2 := *remote
		u2.Scheme = "https"
		// payload used for headers
		c, err = h2conn.New(tt, u2.String(), payload.Headers())
	default:
		err = fmt.Errorf("unsupported dial scheme '%s' in %s", remote.Scheme, remote)
	}
	return
}
