// Copyright (c) 2021 Wireleap

package pofrequest

import "testing"

func TestPRValidate(t *testing.T) {
	x := &T{}
	err := x.Validate()

	if err == nil {
		t.Fatal("no error returned where validation should fail")
	}

	i := int64(10)
	x.Quantity = &i
	err = x.Validate()

	if err == nil {
		t.Fatal("no error returned where validation should fail")
	}

	y := "pofex"
	x.Type = &y
	err = x.Validate()

	if err == nil {
		t.Fatal("no error returned where validation should fail")
	}

	x.Duration = &i
	err = x.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
