// Copyright (c) 2021 Wireleap

package dirinfo

import (
	"sync"

	"github.com/blang/semver"
	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/texturl"
)

type T struct {
	PublicKey  jsonb.PK     `json:"public_key"`
	Version    string       `json:"version"`
	Endpoint   *texturl.URL `json:"endpoint"`
	Info       *texturl.URL `json:"info,omitempty"`
	Enrollment Enrollment   `json:"enrollment"`
	// NOTE: update_channels is deprecated
	Channels        ChannelMap      `json:"update_channels,omitempty"`
	UpgradeChannels UpgradeChannels `json:"upgrade_channels,omitempty"`
}

type ChannelMap map[string]semver.Version

type UpgradeChannels struct {
	Relay  ChannelMap `json:"relay,omitempty"`
	Client ChannelMap `json:"client,omitempty"`
}

type Enrollment struct {
	Fronting RoleInfo `json:"fronting"`
	Entropic RoleInfo `json:"entropic"`
	Backing  RoleInfo `json:"backing"`
}

func (t *Enrollment) Role(role string) (ri *RoleInfo) {
	switch role {
	case "fronting":
		ri = &t.Fronting
	case "entropic":
		ri = &t.Entropic
	case "backing":
		ri = &t.Backing
	}

	return
}

func (t *Enrollment) Restrict(keys map[string][]string) {
	t.Fronting.Restricted = false
	t.Entropic.Restricted = false
	t.Backing.Restricted = false

	for role, ks := range keys {
		if role != "" && len(ks) > 0 {
			switch role {
			case "fronting":
				t.Fronting.Restricted = true
			case "entropic":
				t.Entropic.Restricted = true
			case "backing":
				t.Backing.Restricted = true
			}
		}
	}
}

type RoleInfo struct {
	sync.Mutex `json:"-"`

	Count      int  `json:"count"`
	Restricted bool `json:"restricted"`
}

func (ri *RoleInfo) Incr() {
	ri.Lock()
	ri.Count++
	ri.Unlock()
}

func (ri *RoleInfo) Decr() {
	ri.Lock()
	ri.Count--
	ri.Unlock()
}
