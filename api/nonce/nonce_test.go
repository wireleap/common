// Copyright (c) 2022 Wireleap

package nonce

import (
	"testing"
)

// Nonces up to size 32 are tested -- this is arbitrary.
func TestNonce(t *testing.T) {
	for i := 0; i <= 32; i++ {
		t.Log("generating nonce of size", i)

		x, err := New(i)

		t.Log("generated nonce:", x)

		if err != nil {
			t.Error(err)
		}

		lx := len(x)

		if lx != i {
			t.Errorf("invalid size %d for nonce of requested size %d", lx, i)
		}
	}
}
