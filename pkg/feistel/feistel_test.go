package feistel_test

import (
	"bytes"
	"evteev/cipher/pkg/feistel"
	"io"
	"strings"
	"testing"
)

func TestFeistelCipher(t *testing.T) {
	tests := []struct {
		input string
		key   string
	}{
		{"hello", "qwerty"},
		{"world", "hello"},
	}

	for _, test := range tests {

		key := feistel.GenerateKeyFromString(test.key)
		cipher := feistel.New(key)

		encrypted, err := cipher.Encrypt(strings.NewReader(test.input))
		if err != nil {
			t.Errorf("Encrypt(%q) throw error: %v", test.input, err)
		}

		buf := new(bytes.Buffer)
		io.Copy(buf, encrypted)

		decrypted, err := cipher.Decrypt(buf)
		if err != nil {
			t.Errorf("Decrypt(%q) throw error: %v", test.input, err)
		}

		dec, err := io.ReadAll(decrypted)
		if err != nil {
			t.Errorf("Ошибка при чтении: %v", err)
		}

		if string(dec) != test.input {
			t.Errorf(test.input + " не совпадает с " + string(dec))
		}

		t.Logf("%q -> %q", test.input, string(dec))
	}
}
