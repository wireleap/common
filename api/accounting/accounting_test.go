// Copyright (c) 2021 Wireleap

package accounting

import "testing"

func TestValidate(t *testing.T) {
	x := &T{}
	err := x.Validate()

	if err == nil {
		t.Fatal("no error returned where validation should fail")
	}

	p := int64(1)
	x.Price = &p
	err = x.Validate()

	if err == nil {
		t.Fatal("no error returned where validation should fail")
	}

	c := "usd"
	x.Currency = &c
	err = x.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
