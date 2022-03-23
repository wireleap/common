// Copyright (c) 2022 Wireleap

package tlscert

import (
	"crypto/ed25519"
	"crypto/tls"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	defer os.RemoveAll("testdata")
	err := os.MkdirAll("testdata", 0755)

	if err != nil {
		t.Fatal(err)
	}

	_, pk, err := ed25519.GenerateKey(nil)

	if err != nil {
		t.Fatal(err)
	}

	cert := filepath.Join("testdata", "cert.pem")
	key := filepath.Join("testdata", "key.pem")

	err = Generate(cert, key, pk)

	if err != nil {
		t.Fatal(err)
	}

	_, err = tls.LoadX509KeyPair(cert, key)

	if err != nil {
		t.Fatal(err)
	}
}
