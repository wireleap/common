// Copyright (c) 2022 Wireleap

package fsdir

import (
	"os"
	"reflect"
	"testing"
)

type testStruct1 struct {
	I int `json:"i"`
}

type testStruct2 struct {
	*testStruct1
	F float64 `json:"f"`
}

type testStruct3 struct {
	*testStruct2
	S string `json:"s"`
}

func TestFileMap(t *testing.T) {
	ts1 := testStruct1{12345}
	ts2 := testStruct2{&ts1, 123.45}
	ts3 := testStruct3{&ts2, "foobar"}

	ts11 := ts1
	ts22 := ts2
	ts33 := ts3

	data := map[string]interface{}{
		"testStruct1": &ts11,
		"testStruct2": &ts22,
		"testStruct3": &ts33,
	}

	m, err := New("testdata")

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll("testdata")
	})

	for f, x := range data {
		// rescope
		f, x := f, x

		t.Run(f, func(t *testing.T) {
			t.Parallel()

			err := m.Set(x, f)

			if err != nil {
				t.Fatal(err)
			}

			err = m.Get(x, f)

			if err != nil {
				t.Fatal(err)
			}
		})
	}

	if !reflect.DeepEqual(ts11, ts1) {
		t.Fatalf("got: %s %+v, expected: %s %+v", reflect.TypeOf(ts11), ts11, reflect.TypeOf(ts1), ts1)
	}

	if !reflect.DeepEqual(ts22, ts2) {
		t.Fatalf("got: %s %+v, expected: %s %+v", reflect.TypeOf(ts22), ts22, reflect.TypeOf(ts2), ts2)
	}

	if !reflect.DeepEqual(ts33, ts3) {
		t.Fatalf("got: %s %+v, expected: %s %+v", reflect.TypeOf(ts33), ts33, reflect.TypeOf(ts3), ts3)
	}
}
