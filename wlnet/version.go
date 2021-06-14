// Copyright (c) 2021 Wireleap

package wlnet

import (
	"fmt"

	"github.com/blang/semver"
)

var PROTO_VERSION_STRING = "<unset>"

// PROTO_VERSION is the current protocol version string according to
// semver v2.0.0: https://semver.org/spec/v2.0.0.html
var PROTO_VERSION semver.Version = semver.MustParse(PROTO_VERSION_STRING)

// versionCheck checks if the given version is parsable and compatible with
// the current protocol version (PROTO_VERSION). The current implementation
// checks if the major versions are different as that is the definition of
// backwards incompatibility in semver v2.0.0.
func VersionCheck(v2 *semver.Version) error {
	v1 := &PROTO_VERSION

	if v1.Minor != v2.Minor {
		return fmt.Errorf("expecting version 0.%d.x, got %s", v1.Minor, v2.String())
	}

	return nil
}
