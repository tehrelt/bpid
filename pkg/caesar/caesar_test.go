package caesar_test

import (
	"bytes"
	"evteev/caesar/pkg/caesar"
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
		{"привет", 3, "тулйззг"},
		{"Привет", 5, "Фулзййг"},
		{"Ёжик", 1, "Ёзйл"},
		{"абвгд", -1, "яабвг"},
		{"Привет мир!", 7, "Фсрнйлч тлх!"},
	}

	for _, test := range tests {
		reader := caesar.NewCaesarReader(strings.NewReader(test.input), test.shift)
		result, err := io.ReadAll(reader)
		if err != nil {
			t.Errorf("Ошибка при чтении из CaesarReader: %v", err)
		}
		if string(result) != test.expectedOutput {
			t.Errorf("Encrypt(%q): Ожидалось %q, но получилось %q", test.input, test.expectedOutput, string(result))
		}
	}
}

func TestCaesarDecrypt(t *testing.T) {
	tests := []struct {
		encrypted      string
		shift          int
		expectedOutput string
	}{
		{"тулйззг", -3, "привет"},
		{"Фулзййг", -5, "Привет"},
		{"Ёзйл", -1, "Ёжик"},
		{"яабвг", 1, "абвгд"},
		{"Фсрнйлч тлх!", -7, "Привет мир!"},
	}

	for _, test := range tests {
		reader := caesar.NewCaesarReader(strings.NewReader(test.encrypted), test.shift)
		result, err := io.ReadAll(reader)
		if err != nil {
			t.Errorf("Ошибка при чтении из CaesarReader: %v", err)
		}
		if string(result) != test.expectedOutput {
			t.Errorf("Decrypt(%q): Ожидалось %q, но получилось %q", test.encrypted, test.expectedOutput, string(result))
		}
	}
}

func TestCaesarCipher(t *testing.T) {
	tests := []struct {
		input    string
		shift    int
		expected string
	}{
		{"привет", 3, "тулйззг"},
		{"Привет", 5, "Фулзййг"},
		{"Ёжик", 1, "Ёзйл"},
		{"абвгд", -1, "яабвг"},
		{"Привет мир!", 7, "Фсрнйлч тлх!"},
	}

	for _, test := range tests {
		reader := strings.NewReader(test.input)
		encryptor := caesar.NewCaesarReader(reader, test.shift)

		encrypted, err := io.ReadAll(encryptor)
		if err != nil {
			t.Errorf("Ошибка при чтении из CaesarReader: %v", err)
		}

		if string(encrypted) != test.expected {
			t.Errorf("Encrypt(%q) Ожидалось %q, но получилось %q", test.input, test.expected, string(encrypted))
		}

		decryptor := caesar.NewCaesarReader(bytes.NewBuffer(encrypted), -test.shift)
		decrypted, err := io.ReadAll(decryptor)
		if err != nil {
			t.Errorf("Ошибка при чтении из CaesarReader: %v", err)
		}

		if string(decrypted) != test.input {
			t.Errorf("Decrypt(%q) Ожидалось %q, но получилось %q", encrypted, test.expected, string(decrypted))
		}
	}
}
