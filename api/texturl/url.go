// Copyright (c) 2021 Wireleap

package texturl

import (
	"net/url"
)

// https://github.com/golang/go/issues/25705

type URL struct{ url.URL }

func (u *URL) UnmarshalText(b []byte) error { return u.UnmarshalBinary(b) }
func (u URL) MarshalText() ([]byte, error)  { return u.MarshalBinary() }

func URLMustParse(rawurl string) *URL {
	u, err := url.Parse(rawurl)

	if err != nil {
		panic(err)
	}

	return &URL{*u}
}
