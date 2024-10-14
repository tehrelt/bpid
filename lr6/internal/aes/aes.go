package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	ErrInvalidKey = errors.New("invalid key size. key must be 16, 24, 32 byte len")
)

type Cipher struct {
	key       []byte
	block     cipher.Block
	blockSize int
}

// New initializes a new Cipher with the provided key.
func New(key []byte) (*Cipher, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, ErrInvalidKey
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &Cipher{
		key:       key,
		block:     block,
		blockSize: aes.BlockSize,
	}, nil
}

func (c *Cipher) Encrypt(in io.Reader) (io.Reader, error) {
	input, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	input = pad(input, c.blockSize)
	output := make([]byte, len(input)+c.blockSize)

	iv := output[:c.blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	encrypter := cipher.NewCBCEncrypter(c.block, iv)
	encrypter.CryptBlocks(output[c.blockSize:], input)

	return bytes.NewReader(output), nil
}

func (c *Cipher) Decrypt(in io.Reader) (io.Reader, error) {
	input, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	iv := input[:c.blockSize]
	block := cipher.NewCBCDecrypter(c.block, iv)
	input = input[c.blockSize:]
	block.CryptBlocks(input, input)
	output := unpad(input)

	return bytes.NewReader(output), nil
}

func pad(in []byte, blockSize int) []byte {
	padding := blockSize - len(in)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(in, padText...)
}

func unpad(in []byte) []byte {
	length := len(in)
	unpadding := int(in[length-1])
	return in[:length-unpadding]
}
