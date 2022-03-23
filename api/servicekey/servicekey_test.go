// Copyright (c) 2022 Wireleap

package servicekey

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/wireleap/common/api/signer"
)

func TestServicekey(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(nil)

	if err != nil {
		t.Fatal(err)
	}

	s := signer.New(priv)
	sk := New(priv)

	sk.Contract = &Contract{
		SettlementOpen:  time.Now().Unix(),
		SettlementClose: time.Now().Unix() + 100,
	}

	sk.Contract.Sign(s)

	err = sk.Contract.Verify()

	if err != nil {
		t.Fatal(err)
	}
}
