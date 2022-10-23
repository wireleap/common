// Copyright (c) 2022 Wireleap

package ststore

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	_, err = New(tmpd, ContractKeyFunc)

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestLoad(t *testing.T) {
	tmpd, err := ioutil.TempDir("", "wltest.*")

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpd)
	})

	// malformed JSON file
	data := []byte("{\"key\": \"value\"")

	// Write unreachable file in root folder
	fp := filepath.Join(tmpd, "somedata.json")

	if err := os.WriteFile(fp, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Write unreachable file in depth=2 folder
	path := filepath.Join(tmpd, "somecontract", "somefolder")

	if err = os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}

	fp = filepath.Join(path, "somedata.json")

	if err := os.WriteFile(fp, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Write reachable (depth=1) malformed file
	path = filepath.Join(tmpd, "somecontract")
	fp = filepath.Join(path, "somedata.json")

	if err := os.WriteFile(fp, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Initialise STStore, should handle errored files
	if _, err := New(tmpd, ContractKeyFunc); err != nil {
		t.Fatalf("unmarshal shouldn't have failed")
	}

	// Load malformed file and compare with original data
	fp = filepath.Join(path, "malformed", "somedata.json")

	var data2 []byte
	if data2, err = os.ReadFile(fp); err != nil {
		t.Fatal(err)
	}

	if comp := bytes.Compare(data, data2); comp != 0 {
		t.Errorf("Strings do not match %v", comp)
	}
}

func BenchmarkSTStore(b *testing.B) {
	tmpd, err := ioutil.TempDir("", "wltest.*")

	if err != nil {
		b.Fatal(err)
	}

	b.Cleanup(func() {
		os.RemoveAll(tmpd)
	})

	s, err := New(tmpd, ContractKeyFunc)

	if err != nil {
		b.Fatal(err)
	}

	pub, priv, err := ed25519.GenerateKey(nil)

	if err != nil {
		b.Fatal(err)
	}

	sk := servicekey.New(priv)

	b.Run("FillSTStore", func(b *testing.B) {
		// generate b.N sharetokens
		for n := 0; n < b.N; n++ {
			st, err := sharetoken.New(sk, pub)

			if err != nil {
				b.Fatal(err)
			}

			err = s.Add(st)

			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("LoadSTStore", func(b *testing.B) {
		// load b.N sharetokens
		if _, err := New(tmpd, ContractKeyFunc); err != nil {
			fmt.Println(err)
			b.Fatal(err)
		}
	})
}
