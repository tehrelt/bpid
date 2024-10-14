package main

import (
	"evteev/bpid/lr4/internal/rsa"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	DEFAULT_INPUT  = ""
	DEFAULT_OUTPUT = ""
	DEFAULT_MODE   = ""
	DEFAULT_KEY    = ""
)

var (
	inputPath      string
	outputPath     string
	mode           string
	privateKeyPath string
	publicKeyPath  string
)

func init() {
	flag.StringVar(&inputPath, "in", DEFAULT_INPUT, "input file")
	flag.StringVar(&outputPath, "out", DEFAULT_OUTPUT, "output file")
	flag.StringVar(&mode, "mode", DEFAULT_MODE, "mode <(encrypt|e)|(decrypt|d)>")
	flag.StringVar(&privateKeyPath, "priv", DEFAULT_KEY, "path to private key")
	flag.StringVar(&publicKeyPath, "pub", DEFAULT_KEY, "path to public key")
}

func main() {

	flag.Parse()

	if inputPath == DEFAULT_INPUT {
		log.Fatal("input file is required")
	}

	if outputPath == DEFAULT_OUTPUT {
		log.Fatal("output file is required")
	}

	if mode == DEFAULT_MODE {
		log.Fatal("mode is required")
	}

	if privateKeyPath == DEFAULT_KEY {
		log.Fatal("private key is required")
	}

	if publicKeyPath == DEFAULT_KEY {
		log.Fatal("public key is required")
	}

	if mode != "e" && mode != "encrypt" && mode != "decrypt" && mode != "d" {
		log.Fatal("mode must be <(encryption|e)|(decryption|d)>")
	}

	priv, err := rsa.ExtractPrivateKeyFromFile(privateKeyPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read private key: %v", err))
	}

	pub, err := rsa.ExtractPublicKeyFromFile(publicKeyPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read public key: %v", err))
	}

	cipher := rsa.New(priv, pub)

	in, err := os.Open(inputPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open input file: %v", err))
	}
	defer in.Close()

	var out io.Reader

	if mode == "encrypt" || mode == "e" {
		out, err = cipher.Encrypt(in)
		if err != nil {
			panic(fmt.Sprintf("cannot encrypt file: %v", err.Error()))
		}
	} else if mode == "decrypt" || mode == "d" {
		out, err = cipher.Decrypt(in)
		if err != nil {
			panic(fmt.Sprintf("cannot decrypt file: %s", err.Error()))
		}
	} else {
		panic("invalid mode")
	}

	outfile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic("cannot open output file")
	}
	defer outfile.Close()

	if _, err := io.Copy(outfile, out); err != nil {
		panic("cannot write to output file")
	}
}
