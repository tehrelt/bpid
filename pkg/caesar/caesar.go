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
		// Шифруем только русские буквы (а-я, А-Я)
		if r >= 'а' && r <= 'я' {
			r = ((r-'а'+rune(shift))%32+32)%32 + 'а' // модуль 33 для обертывания
		} else if r >= 'А' && r <= 'Я' {
			r = ((r-'А'+rune(shift))%32+32)%32 + 'А' // модуль 33 для обертывания
		}
		result.WriteRune(r)
	}

	return strings.NewReader(result.String()), nil
}
