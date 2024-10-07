package feistel

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
)

const (
	BLOCK_SIZE  = 16
	PARTS_COUNT = 4
)

type FeistelCipher struct {
	key    []byte
	rounds int
}

func New(key []byte, rounds ...int) *FeistelCipher {
	if len(key) != 16 {
		panic("key must be 16 bytes long")
	}

	if rounds == nil {
		rounds = []int{8}
	}

	return &FeistelCipher{key: key, rounds: rounds[0]}
}

func (c *FeistelCipher) f(in []byte) []byte {
	out := bytes.Clone(in)

	for i := 0; i < len(in); i++ {
		out[i] = in[i] ^ c.key[i%BLOCK_SIZE]
	}

	return out
}

func (c *FeistelCipher) Encrypt(in io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	data = c.pad(data)
	encoded := c.encrypt(data)

	return bytes.NewReader(encoded), nil
}

func (c *FeistelCipher) Decrypt(in io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	decoded := c.decrypt(data)
	decoded = c.unpad(decoded)

	return bytes.NewReader(decoded), nil
}

func (c *FeistelCipher) split(in []byte, count int) ([][]byte, int) {
	size := len(in) / count

	parts := make([][]byte, count)
	for i := 0; i < count; i++ {
		parts[i] = in[i*size : (i+1)*size]
	}

	return parts, size
}

func (c *FeistelCipher) encrypt(in []byte) []byte {

	parts, partSize := c.split(in, PARTS_COUNT)

	for i := 0; i < c.rounds; i++ {
		fout := c.f(parts[0])

		x2 := make([]byte, partSize)
		x3 := make([]byte, partSize)
		x4 := make([]byte, partSize)

		for i, b := range fout {
			x2[i] = parts[1][i] ^ b
			x3[i] = parts[2][i] ^ b
			x4[i] = parts[3][i] ^ b
		}

		parts = [][]byte{x2, x3, x4, parts[0]}
	}

	out := make([]byte, 0, len(in))

	for i := range PARTS_COUNT {
		out = append(out, parts[i]...)
	}

	return out
}

func (c *FeistelCipher) decrypt(in []byte) []byte {
	parts, partSize := c.split(in, PARTS_COUNT)

	for i := 0; i < c.rounds; i++ {
		parts = [][]byte{parts[3], parts[0], parts[1], parts[2]}
		fout := c.f(parts[0])
		x2 := make([]byte, partSize)
		x3 := make([]byte, partSize)
		x4 := make([]byte, partSize)

		for i, b := range fout {
			x2[i] = parts[1][i] ^ b
			x3[i] = parts[2][i] ^ b
			x4[i] = parts[3][i] ^ b
		}

		parts[1] = x2
		parts[2] = x3
		parts[3] = x4
	}

	out := make([]byte, 0, len(in))

	for i := range PARTS_COUNT {
		out = append(out, parts[i]...)
	}

	if len(in) != len(out) {
		panic(fmt.Sprintf("watafak len(in)=(%d) != len(out)=(%d)", len(in), len(out)))
	}

	return out
}

func (c *FeistelCipher) pad(in []byte) []byte {
	pad_len := BLOCK_SIZE - len(in)%BLOCK_SIZE
	return append(in, bytes.Repeat([]byte{byte(pad_len)}, pad_len)...)
}

func (c *FeistelCipher) unpad(in []byte) []byte {
	pad_len := int(in[len(in)-1])

	if pad_len <= len(in) {
		return in[:len(in)-pad_len]
	}

	return in
}

func GenerateKeyFromString(in string) []byte {
	hash := sha256.Sum256([]byte(in))
	return hash[:16]
}
