package dsa

import (
	dsa "crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func ExtractPrivateKeyFromFile(path string) (dsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key: block.Type = %s", block.Type)
	}

	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	pkey, ok := privKey.(dsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not a DSA private key")
	}

	return pkey, nil
}

func ExtractPublicKeyFromFile(path string) (dsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub, ok := pubKey.(dsa.PublicKey)
	if !ok {
		return nil, errors.New("not an DSA public key")
	}

	return pub, nil
}
