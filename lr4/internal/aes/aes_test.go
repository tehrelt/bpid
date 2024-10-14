package aes_test

import (
	"evteev/bpid/lr4/internal/aes"
	"io"
	"strings"
	"testing"
)

func TestKeySizeError(t *testing.T) {
	cases := []struct {
		name string
		key  string
	}{
		{"len<16", "reltrlet"},
		{"16<len<24", "reltreltreltreltrelt"},
		{"24<len<32", "reltreltreltreltreltreltreltreltreltreltrelt"},
		{"32<len", "reltreltreltreltreltreltreltreltreltreltreltreltrelt"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := aes.New([]byte(c.key)); err != aes.ErrInvalidKey {
				t.Error("expected error for key size, got nil")
			}
		})
	}
}

func TestCipher(t *testing.T) {
	cases := []struct {
		name string
		key  string
		in   string
	}{
		{"aes128", "examplekey123456", "lorem ipsum"},
		{"aes192", "a24bytelongkeyforaes192m", "hello world"},
		{"aes256", "a32bytelongkeyforaesencryption12", "maxim bogomolov"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			cipher, err := aes.New([]byte(c.key))
			if err != nil {
				t.Errorf("failed to create cipher: %v", err)
			}

			encrypted, err := cipher.Encrypt(strings.NewReader(c.in))
			if err != nil {
				t.Errorf("failed to encrypt: %v", err)
			}

			decrypted, err := cipher.Decrypt(encrypted)
			if err != nil {
				t.Errorf("failed to decrypt: %v", err)
			}

			out, err := io.ReadAll(decrypted)
			if err != nil {
				t.Errorf("failed to read decrypted: %v", err)
			}

			if string(out) != c.in {
				t.Errorf("decrypted value does not match input: got %s, want %s", decrypted, c.in)
			}
		})
	}
}
