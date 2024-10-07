package feistel

import (
	"bytes"
	"crypto/sha256"
	"io"
)

const (
	BLOCK_SIZE  = 16
	PARTS_COUNT = 4
)

type FeistelCipher struct {
	keys   [][]byte
	rounds int
}

func New(keys [][]byte) *FeistelCipher {
	return &FeistelCipher{keys: keys, rounds: len(keys)}
}

func (c *FeistelCipher) f(in []byte, key []byte) []byte {
	out := bytes.Clone(in)

	for i := 0; i < len(in); i++ {
		out[i] = in[i] ^ key[i%BLOCK_SIZE]
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

func (c *FeistelCipher) split(in []byte, count int) [][]byte {
	size := len(in) / count

	parts := make([][]byte, count)
	for i := 0; i < count; i++ {
		parts[i] = in[i*size : (i+1)*size]
	}

	return parts
}

func (c *FeistelCipher) encrypt(in []byte) []byte {

	p := c.split(in, PARTS_COUNT)
	size := len(p[0])

	for round := range c.rounds {
		key := c.keys[round]

		fout := c.f(p[0], key)
		x2 := make([]byte, size)
		x3 := make([]byte, size)
		x4 := make([]byte, size)

		for i, b := range fout {
			x2[i] = p[1][i] ^ b
			x3[i] = p[2][i] ^ b
			x4[i] = p[3][i] ^ b
		}

		p = [][]byte{x2, x3, x4, p[0]}
	}

	p = [][]byte{p[3], p[0], p[1], p[2]}

	out := make([]byte, 0, len(in))

	for i := range PARTS_COUNT {
		out = append(out, p[i]...)
	}

	return out
}

func (c *FeistelCipher) decrypt(in []byte) []byte {
	p := c.split(in, PARTS_COUNT)
	size := len(p[0])

	p = [][]byte{p[1], p[2], p[3], p[0]}

	for round := range c.rounds {
		p = [][]byte{p[3], p[0], p[1], p[2]}

		key := c.keys[c.rounds-round-1]

		fout := c.f(p[0], key)
		x2 := make([]byte, size)
		x3 := make([]byte, size)
		x4 := make([]byte, size)

		for i, b := range fout {
			x2[i] = p[1][i] ^ b
			x3[i] = p[2][i] ^ b
			x4[i] = p[3][i] ^ b
		}

		p = [][]byte{p[0], x2, x3, x4}
	}

	out := make([]byte, 0, len(in))

	for i := range PARTS_COUNT {
		out = append(out, p[i]...)
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

func GenerateKeysFromString(in []string) [][]byte {
	keys := make([][]byte, 0, len(in))

	for _, key := range in {
		hash := sha256.Sum256([]byte(key))
		keys = append(keys, hash[:16], hash[16:32])
	}

	return keys
}
