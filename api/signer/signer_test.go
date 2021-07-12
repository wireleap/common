// Copyright (c) 2021 Wireleap

package signer

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

func TestSigner(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)

	if err != nil {
		t.Fatal(err)
	}

	s := New(priv)

	if !bytes.Equal(s.Public(), pub) {
		t.Error("signer pubkey and original pubkey are not identical")
	}

	msg := []byte("fnord")

	sig := s.Sign(msg)

	if !ed25519.Verify(pub, msg, sig) {
		t.Error("signer-signed message does not verify")
	}
}
