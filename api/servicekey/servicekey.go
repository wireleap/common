// Copyright (c) 2022 Wireleap

package servicekey

import (
	"crypto/ed25519"
	"fmt"
	"strconv"
	"strings"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/signer"
)

// T is the struct holding the servicekey data.
type T struct {
	PrivateKey jsonb.SK  `json:"private_key,omitempty"`
	PublicKey  jsonb.PK  `json:"public_key,omitempty"`
	Contract   *Contract `json:"contract,omitempty"`
}

// New creates a new servicekey from pk.
func New(pk ed25519.PrivateKey) *T {
	return &T{
		PrivateKey: jsonb.SK(pk),
		PublicKey:  jsonb.PK(pk.Public().(ed25519.PublicKey)),
		Contract:   &Contract{},
	}
}

func (t *T) IsExpiredAt(utime int64) bool {
	return t.Contract.SettlementOpen <= utime
}

// Contract encloses the data of the contract that activated this servicekey.
type Contract struct {
	SettlementOpen  int64    `json:"settlement_open,omitempty"`
	SettlementClose int64    `json:"settlement_close,omitempty"`
	PublicKey       jsonb.PK `json:"public_key,omitempty"`
	Signature       jsonb.B  `json:"signature,omitempty"`
}

func (c *Contract) Digest() string {
	return strings.Join([]string{
		c.PublicKey.String(),
		strconv.FormatInt(c.SettlementOpen, 10),
		strconv.FormatInt(c.SettlementClose, 10),
	}, ":")
}

func (c *Contract) Sign(s signer.Signer) {
	c.PublicKey = jsonb.PK(s.Public())
	c.Signature = s.Sign([]byte(c.Digest()))
}

func (c *Contract) Verify() error {
	if !ed25519.Verify(ed25519.PublicKey(c.PublicKey), []byte(c.Digest()), c.Signature) {
		return fmt.Errorf(
			"contract data `%s` are not signed by signature `%s`",
			c.Digest(),
			c.Signature.String(),
		)
	}

	return nil
}

//  Source is the type of servicekey sources.
type Source func() (*T, error)
