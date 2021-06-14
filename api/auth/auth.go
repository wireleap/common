// Copyright (c) 2021 Wireleap

package auth

import (
	"bytes"
	"crypto/ed25519"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/blang/semver"
	"github.com/wireleap/common/api/jsonb"
)

const (
	Prefix string = "Wireleap"

	API       string = "Api"
	Relay     string = "Relay"
	Contract  string = "Contract"
	Directory string = "Directory"
	Client    string = "Client"

	Version   string = "Version"
	Pubkey    string = "Pubkey"
	Signature string = "Signature"

	Challenge string = "Challenge"
	Response  string = "Response"
)

func SetHeader(h http.Header, c, f, v string) {
	h.Set(join(c, f), v)
}

func GetHeader(h http.Header, c, f string) string {
	return h.Get(join(c, f))
}

func DelHeader(h http.Header, c, f string) {
	h.Del(join(c, f))
}

func join(c, f string) string {
	return strings.Join([]string{Prefix, c, f}, "-")
}

func SignedRead(r *io.ReadCloser, h http.Header, cs ...string) ([]byte, error) {
	body0, err := ioutil.ReadAll(*r)

	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		var (
			pks  = []byte(h.Get(join(c, Pubkey)))
			sigs = []byte(h.Get(join(c, Signature)))

			pk  jsonb.PK
			sig jsonb.B
		)

		if len(pks) == 0 || len(sigs) == 0 {
			return nil, fmt.Errorf("auth headers for %s missing", c)
		}

		err = (&pk).UnmarshalText(pks)

		if err != nil {
			return nil, err
		}

		err = (&sig).UnmarshalText(sigs)

		if err != nil {
			return nil, err
		}

		if !ed25519.Verify(pk.T(), body0, sig.T()) {
			return nil, fmt.Errorf("auth signature for %s does not verify", c)
		}
	}

	// reset body & deep-copy body contents
	body := make([]byte, len(body0))
	copy(body, body0)
	*r = ioutil.NopCloser(bytes.NewReader(body0))
	return body, nil
}

func SignedReqBody(r *http.Request, cs ...string) ([]byte, error) {
	return SignedRead(&r.Body, r.Header, cs...)
}

func SignedResBody(r *http.Response, cs ...string) ([]byte, error) {
	return SignedRead(&r.Body, r.Header, cs...)
}

func VersionCheck(h http.Header, component string, want *semver.Version) error {
	vers := h.Get(join(component, Version))
	have, err := semver.Parse(vers)

	if err != nil {
		return fmt.Errorf("error when parsing version: %w", err)
	}

	if have.Minor != want.Minor {
		return fmt.Errorf("expecting version 0.%d.x, got %s", want.Minor, have.String())
	}

	return nil
}
