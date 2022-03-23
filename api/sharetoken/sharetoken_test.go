// Copyright (c) 2022 Wireleap

package sharetoken

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/servicekey"
	"github.com/wireleap/common/api/signer"
)

func TestVerify(t *testing.T) {
	pk, sk, err := ed25519.GenerateKey(nil)

	if err != nil {
		t.Fatal(err)
	}

	skey := servicekey.New(sk)

	// mock contract sig
	skey.Contract = &servicekey.Contract{
		PublicKey:       jsonb.PK(pk),
		SettlementOpen:  9999999999,
		SettlementClose: 99999999999,
	}

	skey.Contract.Sign(signer.New(sk))

	sharetoken, err := New(skey, pk)

	if err != nil {
		t.Fatal(err)
	}

	err = sharetoken.Verify()

	if err != nil {
		t.Fatal(err)
	}

	if sharetoken.IsExpiredAt(time.Now().Unix()) {
		t.Fatal("sharetoken is expired")
	}
}
