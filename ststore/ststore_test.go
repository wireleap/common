// Copyright (c) 2021 Wireleap

package ststore

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/wireleap/common/api/servicekey"
	"github.com/wireleap/common/api/sharetoken"
)

func Test(t *testing.T) {
	tmpd, err := ioutil.TempDir("", "wltest.*")

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpd)
	})

	s, err := New(tmpd, ContractKeyFunc)

	if err != nil {
		t.Fatal(err)
	}

	pub, priv, err := ed25519.GenerateKey(nil)

	if err != nil {
		t.Fatal(err)
	}

	sk := servicekey.New(priv)

	sk.Contract.SettlementOpen = time.Now().Unix()
	sk.Contract.SettlementClose = time.Now().Unix() + 100

	st, err := sharetoken.New(sk, pub)

	if err != nil {
		t.Fatal(err)
	}

	err = s.Add(st)

	if err != nil {
		t.Fatal(err)
	}

	err = s.Add(st)

	if !errors.Is(err, DuplicateSTError) {
		t.Fatalf("returned error does not match expected value: %s", err)
	}

	pk := base64.RawURLEncoding.EncodeToString(pub)

	set := s.SettlingAt(pk, sk.Contract.SettlementOpen)

	if len(set) != 1 {
		t.Fatalf("invalid settling map length: %d, map: %+v", len(set), set)
	}
}
