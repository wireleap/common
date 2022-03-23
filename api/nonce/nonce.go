// Copyright (c) 2022 Wireleap

// Package nonce is used to create nonces.
package nonce

import (
	"crypto/rand"
	"encoding/hex"
)

// New generates a nonce (random hex sequence) of size n.
func New(n int) (string, error) {
	byten := hex.DecodedLen(n) + 1
	b := make([]byte, byten)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	h := hex.EncodeToString(b)
	return h[:n], nil
}
