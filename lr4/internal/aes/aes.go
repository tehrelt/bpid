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
	nonceSize int
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
		nonceSize: aes.BlockSize,
	}, nil
}

func (c *Cipher) Encrypt(in io.Reader) (out io.Reader, err error) {
	nonce := make([]byte, c.nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(c.block, nonce)
	out = &cipher.StreamReader{S: stream, R: in}

	output := make([]byte, c.nonceSize)
	copy(output, nonce)
	return io.MultiReader(bytes.NewReader(output), out), nil
}

func (c *Cipher) Decrypt(in io.Reader) (out io.Reader, err error) {
	nonce := make([]byte, c.nonceSize)
	if _, err := io.ReadFull(in, nonce); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBDecrypter(c.block, nonce)
	out = &cipher.StreamReader{S: stream, R: in}
	return out, nil
}
