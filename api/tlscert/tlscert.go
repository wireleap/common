// Copyright (c) 2021 Wireleap

// Package tlscert provides a few helper functions for working with TLS
// certificates.
package tlscert

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"
)

// Generate generates a PEM TLS certificate and private key based on the ed25519 private key privkey. The certifate is stored at certPath and the key is stored at keyPath.
func Generate(certPath, keyPath string, privkey ed25519.PrivateKey) error {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Wireleap"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(100 * (365 * (24 * time.Hour))),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"wireleap.com"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, privkey.Public(), privkey)

	if err != nil {
		return err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privkey)

	if err != nil {
		return err
	}

	certOut, err := os.Create(certPath)

	if err != nil {
		return err
	}

	log.Println("Writing ed25519 certificate to", certPath)

	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	if err != nil {
		return err
	}

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}

	log.Println("Writing ed25519 certificate key to", keyPath)

	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	if err != nil {
		return err
	}

	return nil
}
