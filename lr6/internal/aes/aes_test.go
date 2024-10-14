package aes_test

import (
	"evteev/bpd/lr6/internal/aes"
	"io"
	"strings"
	"testing"
)

func TestKeySizeError(t *testing.T) {
	cases := []struct {
		name string
		key  string
	}{
		{"len<16", "reltrlet"},
		{"16<len<24", "reltreltreltreltrelt"},
		{"24<len<32", "reltreltreltreltreltreltreltreltreltreltrelt"},
		{"32<len", "reltreltreltreltreltreltreltreltreltreltreltreltrelt"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := aes.New([]byte(c.key)); err != aes.ErrInvalidKey {
				t.Error("expected error for key size, got nil")
			}
		})
	}
}

func TestCipher(t *testing.T) {
	cases := []struct {
		name string
		key  string
		in   string
	}{
		{"aes128", "examplekey123456", "lorem ipsum"},
		{"aes128", "examplekey123456", "ofjwqofpgjwqfo[pwqjfo[wqjfopwqjfowqphjfwqoiphfqwpoihfqwpofhqwopfhqwopf]]"},
		{"aes128", "examplekey123456", "12345678gkpewg goewjpg poewjgepwp[g mnewpgo jnewgpo ewngpnewpg jnjewpng ewping ewipng ewing iewng iewng ioweng ing iowe]"},
		{"aes128", "examplekey123456", "ogeqjgojeqgopenjgope ngeqpogn oeqpng poeqng poeqng poegeqpo"},
		{"aes128", "examplekey123456", "ofnmewqopfnqwo pfnwqfoiphn 2-09 fh192f h1-2 f21 f2=1 f1212="},
		{"aes128", "examplekey123456", "loaldakfjwqpfowqhjfoqwfhjwqfoqwopfjhwq"},
		{"aes192", "a24bytelongkeyforaes192m", "hello world"},
		{"aes256", "a32bytelongkeyforaesencryption12", `module evteev/bpd/lr6

go 1.22.1

require golang.org/x/image v0.21.0
`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			cipher, err := aes.New([]byte(c.key))
			if err != nil {
				t.Errorf("failed to create cipher: %v", err)
			}

			encrypted, err := cipher.Encrypt(strings.NewReader(c.in))
			if err != nil {
				t.Errorf("failed to encrypt: %v", err)
			}

			decrypted, err := cipher.Decrypt(encrypted)
			if err != nil {
				t.Errorf("failed to decrypt: %v", err)
			}

			out, err := io.ReadAll(decrypted)
			if err != nil {
				t.Errorf("failed to read decrypted: %v", err)
			}

			if string(out) != c.in {
				t.Errorf("decrypted value does not match input: got %s, want %s", decrypted, c.in)
			}

			t.Logf("in:  %s", c.in)
			t.Logf("out: %s", string(out))
		})
	}
}
