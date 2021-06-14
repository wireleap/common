// Copyright (c) 2021 Wireleap

package transport

import (
	"net"
	"net/url"
	"testing"
	"time"
)

func TestWLTransport(t *testing.T) {
	// mock network dials
	o := Options{TLSVerify: false, Timeout: time.Second * 5}
	for _, scheme := range []string{"wireleap", "https", "target", "invalid"} {
		t.Run(scheme+" write", func(t *testing.T) {
			t.Parallel()

			tt := New(o)
			c1, c2 := net.Pipe()
			tt.Dial = func(_, _ string) (net.Conn, error) { return c1, nil }
			tt.DialTLS = tt.Dial
			u, err := url.Parse(scheme + "://test:1234")
			if err != nil {
				t.Fatal(err)
			}
			c, err := tt.DialWL("tcp", u)
			if err != nil {
				if scheme == "invalid" {
					return
				}
				t.Fatal(err)
			}
			p0 := []byte{'h', 'e', 'l', 'l', 'o', '!', '\r', '\n'}
			p1 := make([]byte, len(p0))
			t.Run(scheme+" read", func(t *testing.T) {
				t.Parallel()

				n, err := c2.Read(p1)
				if err != nil {
					t.Fatal(err)
				}
				if n != len(p1) {
					t.Fatal("partial read")
				}
			})
			n, err := c.Write(p0)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(p0) {
				t.Fatal("partial write")
			}

			c1.Close()
			c2.Close()
		})
	}
}
