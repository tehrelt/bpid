package main

import (
	"evteev/cipher/pkg/feistel"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <(encrypt|e)|(decrypt|d)> <path> <key> [out]\n", os.Args[0])
		return
	}

	mode := os.Args[1]
	path := os.Args[2]
	key := os.Args[3]

	var file *os.File
	var err error

	file, err = os.Open(path)
	if err != nil {
		slog.Error("cannot open file", slog.String("err", err.Error()), slog.Any("filename", path))
		return
	}
	defer file.Close()
	slog.Debug("file opened", slog.String("filename", path))

	keys := feistel.GenerateKeyFromString(key)

	cipher := feistel.New(keys)

	if strings.Compare(mode, "encrypt") == 0 || strings.Compare(mode, "e") == 0 {

		encoded, err := cipher.Encrypt(file)
		if err != nil {
			slog.Error("cannot encrypt file", slog.String("err", err.Error()))
			return
		}

		content, err := io.ReadAll(encoded)
		if err != nil {
			slog.Error("cannot read encrypted content", slog.String("err", err.Error()))
			return
		}

		outname := func() string {
			if len(os.Args) > 4 {
				return os.Args[4]
			}
			return fmt.Sprintf("%s/enc_feistel_%s", filepath.Dir(path), filepath.Base(path))
		}()

		if err := os.WriteFile(outname, content, 0666); err != nil {
			slog.Error("cannot write file", slog.String("err", err.Error()))
			return
		}

	} else if strings.Compare(mode, "decrypt") == 0 || strings.Compare(mode, "d") == 0 {

		decoded, err := cipher.Decrypt(file)
		if err != nil {
			slog.Error("cannot encrypt file", slog.String("err", err.Error()))
			return
		}

		content, err := io.ReadAll(decoded)
		if err != nil {
			slog.Error("cannot read encrypted content", slog.String("err", err.Error()))
			return
		}

		outname := func() string {
			if len(os.Args) > 4 {
				return os.Args[4]
			}
			return fmt.Sprintf("%s/dec_feistel_%s", filepath.Dir(path), filepath.Base(path))
		}()

		if err := os.WriteFile(outname, content, 0666); err != nil {
			slog.Error("cannot write file", slog.String("err", err.Error()))
			return
		}
	} else {
		log.Fatal(mode + " is not a valid mode\n")
	}
}
