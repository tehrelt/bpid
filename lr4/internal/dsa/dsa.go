package dsa

import (
	"bytes"
	dsa "crypto/ed25519"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

var (
	ErrSignatureVerificationFailed = errors.New("signature failed")
	ErrSignatureMismatch           = errors.New("signatures mismatch")
)

type Signer struct {
	key dsa.PrivateKey
}

type Verifier struct {
	key dsa.PublicKey
}

func NewSigner(key dsa.PrivateKey) *Signer {
	return &Signer{
		key: key,
	}
}

func NewVerifier(key dsa.PublicKey) *Verifier {
	return &Verifier{
		key: key,
	}
}

func (s *Signer) Sign(in io.Reader) (io.Reader, error) {
	hasher := sha1.New()
	if _, err := io.Copy(hasher, in); err != nil {
		return nil, err
	}

	hashed := hasher.Sum(nil)

	signature := dsa.Sign(s.key, hashed)

	return bytes.NewReader(signature), nil
}

func (v *Verifier) Verify(in io.Reader, sign io.Reader) error {

	hasher := sha1.New()
	if _, err := io.Copy(hasher, in); err != nil {
		return err
	}
	hashed := hasher.Sum(nil)

	signature, err := io.ReadAll(sign)
	if err != nil {
		return fmt.Errorf("failed to read a signature: %w", err)
	}

	// Verify the signature
	if !dsa.Verify(v.key, hashed, signature) {
		return ErrSignatureVerificationFailed
	}

	return nil
}
