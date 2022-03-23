// Copyright (c) 2022 Wireleap

package duration

import (
	"encoding/json"
	"testing"
	"time"
)

var cases = map[string]time.Duration{
	`"1m42s"`:     1*time.Minute + 42*time.Second,
	`"1d"`:        24 * time.Hour,
	`"2d"`:        2 * 24 * time.Hour,
	`"3d"`:        3 * 24 * time.Hour,
	`"10d"`:       10 * 24 * time.Hour,
	`"1d2d3h"`:    3*24*time.Hour + 3*time.Hour,
	`"314d12s3m"`: 7536*time.Hour + 12*time.Second + 3*time.Minute,
	`"13m14d3s"`:  336*time.Hour + 13*time.Minute + 3*time.Second,
}

func TestDurationJSON(t *testing.T) {
	var d T

	for k, v := range cases {
		err := json.Unmarshal([]byte(k), &d)

		if err != nil {
			t.Error(err)
		}

		if time.Duration(d) != v {
			t.Errorf("got %s, expected %s", time.Duration(d), v)
		}
	}
}

func TestDurationMarshalJSON(t *testing.T) {
	for _, v := range cases {
		want := v.String()
		b, err := json.Marshal(want)

		if err != nil {
			t.Error(err)
		}

		var s string
		err = json.Unmarshal(b, &s)

		if s != want {
			t.Errorf("got %s, expected %s", s, want)
		}
	}
}
