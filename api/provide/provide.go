// Copyright (c) 2021 Wireleap

package provide

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"sync"
	"time"

	"github.com/wireleap/common/api/apiversion"
	"github.com/wireleap/common/api/auth"
	"github.com/wireleap/common/api/canned"
	"github.com/wireleap/common/api/status"

	"github.com/blang/semver"
)

type Routes map[string]http.Handler

func LogRequestGate(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got JSON API request: %s: %s %s %s\n", r.RemoteAddr, r.Method, r.Proto, r.URL)
		targetMux.ServeHTTP(w, r)
	})
}

func IdempotencyKeyGate(targetMux http.Handler) http.Handler {
	var (
		m  = map[string]canned.T{}
		t  = time.NewTicker(24 * time.Hour) // TODO unhardcode?
		mu sync.RWMutex
	)

	go func() {
		for _ = range t.C {
			log.Printf("cleaning up cached idempotency keys/responses...")

			mu.Lock()
			for k := range m {
				delete(m, k)
			}
			mu.Unlock()

			log.Printf("done cleaning up cached idempotency keys/responses")
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			targetMux.ServeHTTP(w, r)
			return
		}

		ik := r.Header.Get("Idempotency-Key")

		if ik != "" {
			var (
				can canned.T
				ok  bool
				err error
			)

			mu.RLock()
			can, ok = m[ik]
			mu.RUnlock()

			if !ok {
				rr := httptest.NewRecorder()
				targetMux.ServeHTTP(rr, r)

				res := rr.Result()
				can, err = canned.Can(res)

				if err != nil {
					log.Printf("could not put following http response in a can, this is weird...")
					b, err := httputil.DumpResponse(res, true)

					if err == nil {
						log.Print(string(b))
					} else {
						log.Printf("additionally, error while trying to dump response: %s", err)
					}

					status.ErrInternal.WriteTo(w)
					return
				}

				mu.Lock()
				m[ik] = can
				mu.Unlock()
			}

			can.Uncan(w)
			return
		}

		targetMux.ServeHTTP(w, r)
	})
}

func MethodGate(m Routes) http.Handler {
	supported := []string{}

	for k, _ := range m {
		supported = append(supported, k)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m == nil {
			// kinda weird but ok
			status.ErrMethod.Wrap(status.Cause(
				"no HTTP methods are supported for this endpoint",
			)).WriteTo(w)
			return
		}

		h := m[r.Method]

		if h == nil {
			cause := fmt.Sprintf(
				"request uses unsupported HTTP method '%s', supported methods are: %v",
				r.Method,
				supported,
			)

			status.ErrMethod.Wrap(status.Cause(cause)).WriteTo(w)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func AuthGate(targetMux http.Handler, components ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.SignedReqBody(r, components...)

		if err != nil {
			status.ErrForbidden.Wrap(err).WriteTo(w)
			return
		}

		targetMux.ServeHTTP(w, r)
	})
}

func VersionGate(targetMux http.Handler, v *semver.Version) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(auth.GetHeader(r.Header, auth.API, auth.Version)) > 0 {
			// only check version header if provided

			err := auth.VersionCheck(r.Header, auth.API, &apiversion.VERSION)

			if err != nil {
				status.ErrRequest.Wrap(fmt.Errorf(
					"invalid client API version: %w", err,
				)).WriteTo(w)
				return
			}
		}

		auth.SetHeader(w.Header(), auth.API, auth.Version, v.String())
		targetMux.ServeHTTP(w, r)
	})
}

func NewMux(routes ...Routes) *http.ServeMux {
	mux := http.NewServeMux()

	for _, fs := range routes {
		for path, f := range fs {
			mux.Handle(path, f)
		}
	}

	return mux
}

func DefaultServer(addr string, mux http.Handler) *http.Server {
	return &http.Server{
		Addr: addr,
		Handler: LogRequestGate(
			IdempotencyKeyGate(
				mux,
			),
		),
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
}

func UnversionedServer(addr string, mux http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           LogRequestGate(IdempotencyKeyGate(mux)),
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
}
