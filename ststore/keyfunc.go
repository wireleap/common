// Copyright (c) 2021 Wireleap

package ststore

import "github.com/wireleap/common/api/sharetoken"

// KeyFunc is the type of key derivation functions for this sharetoken store.
// On taking a sharetoken as argument it should produce keys of 1st, 2nd and
// 3rd order.
type KeyFunc func(*sharetoken.T) (string, string, string)

// ContractKeyFunc is the standard key derivation function for service
// contracts. The store is organized as such:
// servicekey public key -> relay public key -> sharetoken signature + .json
var ContractKeyFunc = func(st *sharetoken.T) (k1, k2, k3 string) {
	return st.PublicKey.String(), st.RelayPubkey.String(), st.Signature.String()
}

// RelayKeyFunc is the standard key derivation function for relays.  The store
// is organized as such:
// servicekey public key -> service contract public key -> sharetoken signature + .json
var RelayKeyFunc = func(st *sharetoken.T) (k1, k2, k3 string) {
	return st.Contract.PublicKey.String(), st.PublicKey.String(), st.Signature.String()
}
