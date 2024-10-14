package caesar_test

import (
	"tehrelt/bpid/cipher/pkg/caesar"
	"io"
	"strings"
	"testing"
)

func TestCaesarEncrypt(t *testing.T) {
	tests := []struct {
		input          string
		shift          int
		expectedOutput string
	}{
		{"привет", 1, "рсйгжу"},
		{"привет", 2, "сткдзф"},
	}

	for _, test := range tests {
		encrypted, err := caesar.Encrypt(strings.NewReader(test.input), test.shift)
		if err != nil {
			t.Errorf("Ошибка при чтении из Encrypt: %v", err)
		}

		result, err := io.ReadAll(encrypted)
		if err != nil {
			t.Errorf("Ошибка при чтении: %v", err)
		}

		if string(result) != test.expectedOutput {
			t.Errorf("Encrypt(%q): Ожидалось %q, но получилось %q", test.input, test.expectedOutput, string(result))
		}
	}
}

func TestCaesarCipher(t *testing.T) {
	tests := []struct {
		input          string
		shift          int
		expectedOutput string
	}{
		{"привет", 1, "рсйгжу"},
		{"привет", 2, "сткдзф"},
		{"используется", 1, "йтрпмэифжута"},
	}

	for _, test := range tests {
		encrypted, err := caesar.Encrypt(strings.NewReader(test.input), test.shift)
		if err != nil {
			t.Errorf("Ошибка при чтении из Encrypt: %v", err)
		}

		enc, err := io.ReadAll(encrypted)
		if err != nil {
			t.Errorf("Ошибка при чтении: %v", err)
		}

		if string(enc) != test.expectedOutput {
			t.Errorf("Encrypt(%q): Ожидалось %q, но получилось %q", test.input, test.expectedOutput, string(enc))
		}

		decrypted, err := caesar.Decrypt(strings.NewReader(string(enc)), test.shift)
		if err != nil {
			t.Errorf("Ошибка при чтении из Encrypt: %v", err)
		}

		dec, err := io.ReadAll(decrypted)
		if err != nil {
			t.Errorf("Ошибка при чтении: %v", err)
		}

		if string(dec) != test.input {
			t.Errorf("Decrypt(%q): Ожидалось %q, но получилось %q", test.input, test.expectedOutput, string(dec))
		}
	}
}
