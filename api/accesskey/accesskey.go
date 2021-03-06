// Copyright (c) 2022 Wireleap

package accesskey

import (
	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/pof"
	"github.com/wireleap/common/api/texturl"

	"github.com/blang/semver"
)

const VERSION_STRING = "0.1.0"

var VERSION = semver.MustParse(VERSION_STRING)

type T struct {
	Version  *semver.Version `json:"version"`
	Contract *Contract       `json:"contract,omitempty"`
	Pofs     []*pof.T        `json:"pofs,omitempty"`
}

type Contract struct {
	Endpoint  *texturl.URL `json:"endpoint,omitempty"`
	PublicKey jsonb.PK     `json:"public_key,omitempty"`
}
