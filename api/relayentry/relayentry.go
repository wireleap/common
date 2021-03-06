// Copyright (c) 2022 Wireleap

package relayentry

import (
	"fmt"
	"strings"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/texturl"

	"github.com/blang/semver"
)

type T struct {
	Role     string       `json:"role,omitempty"`
	Addr     *texturl.URL `json:"address,omitempty"`
	Pubkey   jsonb.PK     `json:"pubkey,omitempty"`
	Versions Versions     `json:"versions,omitempty"`
	Key      string       `json:"key,omitempty"`
	// Update channel for the relay, used in determining the update version to
	// Upgrade channel for the relay, used in determining the upgrade version to
	// push from the directory. Can be empty, in which case it's taken to be
	// default. (DEPRECATED)
	Channel string `json:"update_channel,omitempty"`
	// Upgrade channel for the relay, used in determining the upgrade version to
	// push from the directory. Can be empty, in which case it's taken to be
	// "default".
	UpgradeChannel string `json:"upgrade_channel,omitempty"`
}

type Versions struct {
	Software      *semver.Version `json:"software,omitempty"`
	ClientRelay   *semver.Version `json:"client-relay,omitempty"`
	RelayRelay    *semver.Version `json:"relay-relay,omitempty"`
	RelayDir      *semver.Version `json:"relay-dir,omitempty"`
	RelayContract *semver.Version `json:"relay-contract,omitempty"`
}

func (r *T) String() string {
	if r == nil {
		return "<nil relay entry>"
	} else {
		user := "-" // no user

		if strings.ContainsRune(r.Key, ':') {
			userpass := strings.SplitN(r.Key, ":", 2)
			user = userpass[0]
		}

		return fmt.Sprintf(
			"%s|%s|%s|%s",
			r.Role,
			r.Addr,
			r.Pubkey,
			user,
		)
	}
}

func (r *T) Validate() error {
	if r == nil {
		return fmt.Errorf("relay entry is null")
	}

	switch r.Role {
	case "fronting", "entropic", "backing":
		// OK
	default:
		return fmt.Errorf("invalid relay role: %s", r.Role)
	}

	if r.Addr == nil {
		return fmt.Errorf("relay entry address is null")
	}

	switch r.Addr.Scheme {
	case "wireleap", "https":
		// OK
	default:
		return fmt.Errorf("invalid relay URL scheme: %s", r.Addr.Scheme)
	}

	// NOTE: update_channel is deprecated
	if r.Channel == "" {
		r.Channel = "default"
	}

	return nil
}
