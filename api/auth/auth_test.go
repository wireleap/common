// Copyright (c) 2021 Wireleap

package auth

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/blang/semver"
)

func TestSignedReqBody(t *testing.T) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	msg := ioutil.NopCloser(bytes.NewReader([]byte("foo")))
	sig := ed25519.Sign(sk, []byte("foo"))

	pkb64 := base64.RawURLEncoding.EncodeToString(pk)
	sigb64 := base64.RawURLEncoding.EncodeToString(sig)

	r := new(http.Request)
	r.Body = msg

	// Should pass
	r.Header = map[string][]string{
		join(Relay, Pubkey):    {pkb64},
		join(Relay, Signature): {sigb64},
	}

	if _, err := SignedReqBody(r, Relay); err != nil {
		t.Fatal(err)
	}

	// Should fail
	r.Header = map[string][]string{
		"foo": {"bar"},
	}

	if _, err := SignedReqBody(r, Relay); err == nil {
		t.Fatal("Unsigned body passed SignedReqBody function")
	}

	// Should fail
	sig = ed25519.Sign(sk, []byte("bar"))
	sigb64 = base64.RawURLEncoding.EncodeToString(sig)

	r.Header = map[string][]string{
		join(Relay, Pubkey):    {pkb64},
		join(Relay, Signature): {sigb64},
	}

	if _, err := SignedReqBody(r, Relay); err == nil {
		t.Fatal("Invalid signature passed as valid")
	}
}

func TestVersionCheck(t *testing.T) {
	r := new(http.Response)
	r.Header = map[string][]string{
		join(Directory, Version): {"1.0.0"},
	}

	v, err := semver.Make("1.0.0")
	if err != nil {
		t.Fatal(err)
	}

	// Should pass
	if err := VersionCheck(r.Header, Directory, &v); err != nil {
		t.Fatal(err)
	}

	// Should fail
	r.Header = map[string][]string{
		join(Directory, Version): {"foobar"},
	}

	if err := VersionCheck(r.Header, Directory, &v); err == nil {
		t.Fatal("Invalid version passed as valid")
	}
}
