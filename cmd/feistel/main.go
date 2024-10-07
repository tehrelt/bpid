package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"tehrelt/bpid/cipher/pkg/feistel"
)

type arrayFlag []string

func (i *arrayFlag) String() string {
	buf := new(strings.Builder)

	for _, v := range *i {
		fmt.Fprintf(buf, "%s ", v)
	}

	return buf.String()
}
func (i *arrayFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	mode string
	in   string
	keys arrayFlag
	out  string
)

func init() {
	flag.StringVar(&in, "in", "", "input file")
	flag.StringVar(&out, "out", "", "out file")
	flag.StringVar(&mode, "mode", "", "encrypt or decrypt")
	flag.Var(&keys, "k", "list of keys. usage: -k key1 -k key2")
}

func main() {

	flag.Parse()

	if in == "" {
		log.Fatal("input file is not set. try -h for help\n")
	}

	if mode == "" {
		log.Fatal("mode is not set. try -h for help\n")
	}

	if len(keys) == 0 {
		log.Fatal("keys are not set. try -h for help\n")
	}

	var file *os.File
	var err error

	file, err = os.Open(in)
	if err != nil {
		slog.Error("cannot open file", slog.String("err", err.Error()), slog.Any("filename", in))
		return
	}
	defer file.Close()

	keys := feistel.GenerateKeysFromString(keys)

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
			if out != "" {
				return out
			}

			return fmt.Sprintf("%s/enc_feistel_%s", filepath.Dir(in), filepath.Base(in))
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
			if out != "" {
				return out
			}
			return fmt.Sprintf("%s/dec_feistel_%s", filepath.Dir(in), filepath.Base(in))
		}()

		if err := os.WriteFile(outname, content, 0666); err != nil {
			slog.Error("cannot write file", slog.String("err", err.Error()))
			return
		}
	} else {
		log.Fatal(mode + " is not a valid mode\n")
	}
}
