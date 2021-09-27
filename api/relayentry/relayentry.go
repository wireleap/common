// Copyright (c) 2021 Wireleap

package relayentry

import (
	"fmt"
	"strings"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/texturl"

	"github.com/blang/semver"
)

type T struct {
	Role     string       `json:"role"`
	Addr     *texturl.URL `json:"address"`
	Pubkey   jsonb.PK     `json:"pubkey"`
	Versions Versions     `json:"versions,omitempty"`
	Key      string       `json:"key,omitempty"`
	// Update channel for the relay, used in determining the update version to
	// push from the directory. Can be empty, in which case it's taken to be
	// default.
	Channel string `json:"update_channel,omitempty"`
}

type Versions struct {
	Software      *semver.Version `json:"software"`
	ClientRelay   *semver.Version `json:"client-relay"`
	RelayRelay    *semver.Version `json:"relay-relay"`
	RelayDir      *semver.Version `json:"relay-dir"`
	RelayContract *semver.Version `json:"relay-contract"`
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

	switch r.Addr.Scheme {
	case "wireleap", "https":
		// OK
	default:
		return fmt.Errorf("invalid relay URL scheme: %s", r.Addr.Scheme)
	}

	if r.Channel == "" {
		r.Channel = "default"
	}

	return nil
}
