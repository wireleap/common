// Copyright (c) 2021 Wireleap

package relayentry

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/texturl"

	"github.com/blang/semver"
)

func TestValidate(t *testing.T) {
	pk, _, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		t.Fatal(err)
	}

	v, err := semver.Make("1.0.0")
	if err != nil {
		t.Fatal(err)
	}

	// Should pass
	r := T{
		Role:    "fronting",
		Addr:    texturl.URLMustParse("wireleap://wireleap.com"),
		Pubkey:  jsonb.PK(pk),
		Version: &v,
	}

	if err = r.Validate(); err != nil {
		t.Fatal(err)
	}

	// Should fail with invalid URL scheme
	r = T{
		Role:    "fronting",
		Addr:    texturl.URLMustParse("gopher://foo.bar"),
		Pubkey:  jsonb.PK(pk),
		Version: &v,
	}

	if err = r.Validate(); err == nil {
		t.Fatal(err)
	}

	// Should fail with invalid relay role
	r = T{
		Role:    "foobar",
		Addr:    texturl.URLMustParse("wireleap://wireleap.com"),
		Pubkey:  jsonb.PK(pk),
		Version: &v,
	}

	if err = r.Validate(); err == nil {
		t.Fatal(err)
	}
}
