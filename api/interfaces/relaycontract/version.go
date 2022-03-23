// Copyright (c) 2022 Wireleap

package relaycontract

import (
	"github.com/blang/semver"
	"github.com/wireleap/common/api/interfaces"
)

// The version of this interface is declared here.

const VERSION_STRING = "0.1.0"

var VERSION = semver.MustParse(VERSION_STRING)

var T = interfaces.T{
	Consumer: interfaces.Relay,
	Provider: interfaces.Contract,
	Version:  VERSION,
}
