// Copyright (c) 2021 Wireleap

package texturl

import (
	"bytes"
	"testing"
)

func TestMarshalText(t *testing.T) {
	u := URLMustParse("https://wireleap.com")
	r := []byte("https://wireleap.com")

	m, err := u.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(m, r) != 0 {
		t.Fatal("Marshalled text doesn't equal requested result")
	}
}

func TestUnmarshalText(t *testing.T) {
	var b []byte
	u := URLMustParse("https://wireleap.com")

	if err := u.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
}

func TestURLMustParse(t *testing.T) {
	u := URLMustParse("https://wireleap.com")
	if u.Scheme != "https" {
		t.Fatal("Did not parse URL properly")
	}
}
