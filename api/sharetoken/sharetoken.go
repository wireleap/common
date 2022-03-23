// Copyright (c) 2022 Wireleap

package sharetoken

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/api/nonce"
	"github.com/wireleap/common/api/servicekey"
)

// T is the struct holding the sharetoken data.
type T struct {
	Version     int64    `json:"version"`
	PublicKey   jsonb.PK `json:"public_key"`
	Timestamp   int64    `json:"timestamp"`
	RelayPubkey jsonb.PK `json:"relay_pubkey"`
	ShareKey    string   `json:"share_key"` // TODO
	Signature   jsonb.B  `json:"signature"`
	Nonce       string   `json:"nonce"`

	Contract *servicekey.Contract `json:"contract,omitempty"`
}

// New creates a new sharetoken.
func New(sk *servicekey.T, pub ed25519.PublicKey) (*T, error) {
	nonce, err := nonce.New(32)

	if err != nil {
		return nil, fmt.Errorf("could not generate nonce: %w", err)
	}

	st := &T{
		PublicKey:   sk.PublicKey,
		Contract:    sk.Contract,
		Timestamp:   time.Now().Unix(),
		RelayPubkey: jsonb.PK(pub),
		Nonce:       nonce,
	}

	st.Signature = ed25519.Sign(
		ed25519.PrivateKey(sk.PrivateKey),
		[]byte(st.Digest()),
	)

	return st, nil
}

func (t *T) Digest() string {
	e := base64.RawURLEncoding.EncodeToString
	i := func(i64 int64) string { return strconv.FormatInt(i64, 10) }

	return strings.Join([]string{
		i(t.Version),
		e(t.PublicKey),
		i(t.Timestamp),
		e(t.RelayPubkey),
		t.ShareKey,
		t.Nonce,
		e(t.Contract.PublicKey),
		e(t.Contract.Signature),
		i(t.Contract.SettlementOpen),
		i(t.Contract.SettlementClose),
	}, ":")
}

func (t *T) Verify() error {
	if !ed25519.Verify(ed25519.PublicKey(t.PublicKey), []byte(t.Digest()), t.Signature) {
		return fmt.Errorf("sharetoken client signature is invalid")
	}

	if t.Contract.PublicKey == nil {
		return fmt.Errorf("contract public key is null")
	} else {
		err := t.Contract.Verify()

		if err != nil {
			return fmt.Errorf(
				"error while verifying sharetoken contract data: %w",
				err,
			)
		}
	}

	return nil
}

func (t *T) IsExpiredAt(utime int64) bool {
	return t.Contract.SettlementOpen <= utime
}

func (t *T) IsSettlingAt(utime int64) bool {
	return t.Contract.SettlementOpen <= utime && t.Contract.SettlementClose > utime
}
