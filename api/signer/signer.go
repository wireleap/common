// Copyright (c) 2022 Wireleap

/*
Package signer defines and implements a simple Signer opaque identity which
encapsulates the private key used to sign data and provides an accessor for its
public key. It is intentionally different from crypto.Signer and simpler.
*/
package signer

import "crypto/ed25519"

type Signer ed25519.PrivateKey

func (s Signer) Sign(data []byte) []byte { return ed25519.Sign(ed25519.PrivateKey(s), data) }

func (s Signer) Public() ed25519.PublicKey { return ed25519.PrivateKey(s).Public().(ed25519.PublicKey) }

func New(k ed25519.PrivateKey) Signer { return Signer(k) }
