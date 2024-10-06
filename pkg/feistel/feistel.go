package feistel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type FeisterCipher struct {
	keys []uint32
}

func New(keys []uint32) *FeisterCipher {
	return &FeisterCipher{keys: keys}
}

// F - простая функция преобразования, которая используется в сети Фейстеля
func F(right uint32, key uint32) uint32 {
	// Простая функция F: здесь можно использовать любую другую функцию
	return right ^ key // Пример: XOR с ключом
}

// process - основной алгоритм Фейстеля
func (c *FeisterCipher) process(left uint32, right uint32, decrypt bool) (uint32, uint32) {
	if decrypt {
		// Если дешифруем, то перебираем ключи в обратном порядке
		for i := len(c.keys) - 1; i >= 0; i-- {
			temp := left
			left = right ^ F(left, c.keys[i])
			right = temp
		}
	} else {
		// Если шифруем, то перебираем ключи в прямом порядке
		for _, key := range c.keys {
			temp := right
			right = left ^ F(right, key)
			left = temp
		}
	}
	return left, right
}

// Encrypt шифрует данные с использованием сети Фейстеля
func (c *FeisterCipher) Encrypt(input io.Reader) (output io.Reader, err error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// Разделяем данные на 64 бита (8 байт) для шифрования
	if len(data) < 8 {
		return nil, fmt.Errorf("недостаточно данных для шифрования")
	}

	left := binary.BigEndian.Uint32(data[:4])   // Левая половина
	right := binary.BigEndian.Uint32(data[4:8]) // Правая половина

	// Шифрование
	left, right = c.process(left, right, false)

	// Собираем зашифрованные данные обратно
	var result bytes.Buffer
	binary.Write(&result, binary.BigEndian, left)
	binary.Write(&result, binary.BigEndian, right)

	return bytes.NewReader(result.Bytes()), nil
}

// Decrypt дешифрует данные с использованием сети Фейстеля
func (c *FeisterCipher) Decrypt(input io.Reader) (output io.Reader, err error) {
	data, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	// Разделяем данные на 64 бита (8 байт) для дешифрования
	if len(data) < 8 {
		return nil, fmt.Errorf("недостаточно данных для дешифрования")
	}

	left := binary.BigEndian.Uint32(data[:4])   // Левая половина
	right := binary.BigEndian.Uint32(data[4:8]) // Правая половина

	// Дешифрование
	left, right = c.process(left, right, true)

	// Собираем расшифрованные данные обратно
	var result bytes.Buffer
	binary.Write(&result, binary.BigEndian, left)
	binary.Write(&result, binary.BigEndian, right)

	return bytes.NewReader(result.Bytes()), nil
}
