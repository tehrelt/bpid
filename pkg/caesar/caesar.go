package caesar

import (
	"bytes"
	"io"
	"strings"
)

func Encrypt(input io.Reader, shift int) (output io.Reader, err error) {
	return process(input, shift)
}

func Decrypt(input io.Reader, shift int) (output io.Reader, err error) {
	return process(input, -shift)
}

func process(input io.Reader, shift int) (output io.Reader, err error) {
	// Читаем все данные из входного Reader
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer

	for _, r := range string(data) {
		if r >= 'а' && r <= 'я' {
			r = ((r-'а'+rune(shift))%32+32)%32 + 'а'
		} else if r >= 'А' && r <= 'Я' {
			r = ((r-'А'+rune(shift))%32+32)%32 + 'А'
		} else if r >= 'a' && r <= 'z' {
			r = ((r-'a'+rune(shift))%26+26)%26 + 'a'
		} else if r >= 'A' && r <= 'Z' {
			r = ((r-'A'+rune(shift))%26+26)%26 + 'A'
		}
		result.WriteRune(r)
	}

	return strings.NewReader(result.String()), nil
}
