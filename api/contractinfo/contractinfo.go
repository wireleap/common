// Copyright (c) 2022 Wireleap

package contractinfo

import (
	"math/big"

	"github.com/blang/semver"
	"github.com/wireleap/common/api/duration"
	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/texturl"
)

// T describes the fields of the wireleap-contract config which are publicly
// accessible via the /info endpoint.
type T struct {
	// Pubkey is the public key used to verify the issued PoF for validity.
	Pubkey jsonb.PK `json:"pubkey"`
	// Version is the version of this wireleap-contract.
	Version semver.Version `json:"version"`
	// Endpoint is the publicly accessible URL of this wireleap-contract instance.
	Endpoint *texturl.URL `json:"endpoint,omitempty"`
	// Pof is the section describing proof of funding mechanisms.
	Pofs []*Pof `json:"proof_of_funding,omitempty"`
	// Servicekey is the section describing servicekey parameters.
	Servicekey Servicekey `json:"servicekey,omitempty"`
	// Settlement is the section describing settlement parameters.
	Settlement Settlement `json:"settlement,omitempty"`
	// Payout is the section describing the configured payout method.
	Payout Payout `json:"payout,omitempty"`
	// Directory is the section describing this contract's directory.
	Directory Directory `json:"directory,omitempty"`
	// Metadata is the section describing this contract's metadata such as
	// logo, ToS etc. All fields are optional.
	Metadata Metadata `json:"metadata,omitempty"`
}

// Pof is the section describing proof of funding mechanisms.
type Pof struct {
	// Endpoint is the URL of the PoF purchasing location.
	Endpoint *texturl.URL `json:"endpoint,omitempty"`
	// Type is one of a set of predefined PoF types.
	// Currently this is either "stripe" or "dummy".
	Type string `json:"type,omitempty"`
	// Pubkey is the public key used to verify the issued PoF for validity.
	Pubkey jsonb.PK `json:"pubkey"`
}

// Servicekey is the section describing servicekey parameters.
type Servicekey struct {
	// Currency is the backing currency of the issued servicekey's value.
	Currency string `json:"currency,omitempty"`
	// Value is the monetary value of the issued servicekey.
	Value *big.Rat `json:"value"`
	// Duration is the length of time a servicekey remains valid for the
	// purpose of issuing sharetokens. After it elapses, submission of share
	// tokens by relays to a service contract to increase the relay's balance
	// becomes possible.
	Duration duration.T `json:"duration,omitempty"`
}

// Settlement is the section describing settlement parameters.
type Settlement struct {
	// FeePercent is the optional value in percent of the fee withheld by the
	// service contract operator for their services.
	FeePercent *big.Rat `json:"fee_percent,omitempty"`
	// SubmissionWindow is the length of time within which relays can submit
	// sharetokens after a servicekey's duration ends.
	SubmissionWindow duration.T `json:"submission_window"`
}

// Payout is the section describing the configured payout method.
type Payout struct {
	// Endpoint is the URL of the payment system (usually wireleap-auth for now).
	Endpoint *texturl.URL `json:"endpoint"`
	// Type is one of a set of predefined payout types. Currently, "stripe" is
	// possible.
	Type string `json:"type"`
	// CheckPeriod is the period at which contract checks for withdrawal status
	// updates.
	CheckPeriod duration.T `json:"check_period,omitempty"`
	// MinWithdrawal is the minimum withdrawal amount.
	MinWithdrawal int64 `json:"min_withdrawal,omitempty"`
	// MaxWithdrawal is the maximum withdrawal amount.
	MaxWithdrawal int64 `json:"max_withdrawal,omitempty"`
	// Info is the URL which provides more information about the supported
	// payout mechanism for the relay operator, if available.
	Info *texturl.URL `json:"info,omitempty"`
}

// Directory is the section describing this contract's directory.
type Directory struct {
	// Endpoint is the URL of the wireleap-dir instance keeping a directory of
	// this service contract's relays.
	Endpoint *texturl.URL `json:"endpoint"`
	// PublicKey is the public key of the directory for verification.
	PublicKey jsonb.PK `json:"public_key"`
}

// Metadata is the section describing this contract's metadata such as
// logo, ToS etc. All fields are optional.
type Metadata struct {
	// Operator is the name of the entity operating this service contract.
	Operator string `json:"operator,omitempty"`
	// OperatorURL is the URL of the entity operating this service contract.
	OperatorURL *texturl.URL `json:"operator_url,omitempty"`
	// Name is this service contract's name.
	Name string `json:"name,omitempty"`
	// ToS is the URL of the Terms of Service of this service contract.
	ToS *texturl.URL `json:"terms_of_service,omitempty"`
	// PrivPolicy is the URL of the privacy policy of this service contract.
	PrivPolicy *texturl.URL `json:"privacy_policy,omitempty"`
}
