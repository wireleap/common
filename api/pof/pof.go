// Copyright (c) 2021 Wireleap

// Package pof provides a mechanism to generate proofs of funding.
package pof

import (
	"strconv"
	"strings"
	"time"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/nonce"
	"github.com/wireleap/common/api/signer"
)

type T struct {
	Type       string  `json:"type,omitempty"`
	Expiration int64   `json:"expiration,omitempty"`
	Nonce      string  `json:"nonce,omitempty"`
	Signature  jsonb.B `json:"signature,omitempty"`
}

func (t *T) Digest() string {
	return strings.Join([]string{
		t.Type,
		strconv.FormatInt(t.Expiration, 10),
		t.Nonce,
	}, ":")
}

func New(s signer.Signer, poftype string, duration int64) (*T, error) {
	exp := time.Now().Unix() + duration
	nonce, err := nonce.New(18)

	if err != nil {
		return nil, err
	}

	pof := &T{
		Type:       poftype,
		Expiration: exp,
		Nonce:      nonce,
	}

	pof.Signature = s.Sign([]byte(pof.Digest()))
	return pof, nil
}

type SKActivationRequest struct {
	Pubkey jsonb.PK `json:"pubkey,omitempty"`
	Pof    *T       `json:"pof,omitempty"`
}

func (t *T) IsExpiredAt(utime int64) bool {
	return t.Expiration <= utime
}
