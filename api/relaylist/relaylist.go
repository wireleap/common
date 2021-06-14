// Copyright (c) 2021 Wireleap Ltd.

package relaylist

import "github.com/wireleap/common/api/relayentry"

type T map[string]*relayentry.T

// All returns all relays in t.
func (t T) All() (rs []*relayentry.T) {
	for _, v := range t {
		rs = append(rs, v)
	}
	return
}
