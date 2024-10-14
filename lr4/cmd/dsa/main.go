package main

import (
	"errors"
	"evteev/bpid/lr4/internal/dsa"
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
	signaturePath  string
	mode           string
	privateKeyPath string
	publicKeyPath  string
)

func init() {
	flag.StringVar(&inputPath, "in", DEFAULT_INPUT, "path to input file")
	flag.StringVar(&signaturePath, "sign", DEFAULT_OUTPUT, "path to signature")
	flag.StringVar(&mode, "mode", DEFAULT_MODE, "mode <(sign|s)|(verify|v)>")
	flag.StringVar(&privateKeyPath, "priv", DEFAULT_KEY, "path to private key. need for sign")
	flag.StringVar(&publicKeyPath, "pub", DEFAULT_KEY, "path to public key. need for verify")
}

func main() {

	flag.Parse()

	if mode == DEFAULT_MODE {
		log.Fatal("mode is required")
	}

	if mode == "s" || mode == "sign" {
		if inputPath == DEFAULT_INPUT {
			log.Fatal("input file is required")
		}

		if signaturePath == DEFAULT_OUTPUT {
			signaturePath = inputPath + ".sig"
		}

		if privateKeyPath == DEFAULT_KEY {
			log.Fatal("private key is required")
		}

		priv, err := dsa.ExtractPrivateKeyFromFile(privateKeyPath)
		if err != nil {
			panic(fmt.Errorf("failed to read private key: %w", err))
		}

		signer := dsa.NewSigner(priv)
		if err := sign(inputPath, signaturePath, signer); err != nil {
			panic(fmt.Errorf("failed to sign: %w", err))
		}

	} else if mode == "v" || mode == "verify" {
		if inputPath == DEFAULT_INPUT {
			log.Fatal("input file is required")
		}

		if signaturePath == DEFAULT_OUTPUT {
			log.Fatal("output file is required")
		}

		if publicKeyPath == DEFAULT_KEY {
			log.Fatal("public key is required")
		}

		pub, err := dsa.ExtractPublicKeyFromFile(publicKeyPath)
		if err != nil {
			panic(fmt.Sprintf("failed to read public key: %v", err))
		}

		verifier := dsa.NewVerifier(pub)
		if err := verify(inputPath, signaturePath, verifier); err != nil {
			panic(fmt.Errorf("failed to verify: %w", err))
		}

	} else {
		panic(fmt.Errorf("unknown mode: %s", mode))
	}
}

func verify(inputPath, signaturePath string, verifier *dsa.Verifier) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer in.Close()

	sign, err := os.Open(signaturePath)
	if err != nil {
		return fmt.Errorf("cannot open signature file: %w", err)
	}
	defer sign.Close()

	if err := verifier.Verify(in, sign); err != nil {
		if errors.Is(err, dsa.ErrSignatureMismatch) {
			fmt.Printf("Подпись не совпадает")
			return nil
		} else if errors.Is(err, dsa.ErrSignatureVerificationFailed) {
			fmt.Printf("Подпись не прошла проверку")
			return nil
		} else {
			return err
		}
	}

	fmt.Printf("Подпись верна")

	return nil
}

func sign(inputPath, signaturePath string, signer *dsa.Signer) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open input file: %w", err)
	}
	defer in.Close()

	signature, err := signer.Sign(in)
	if err != nil {
		return fmt.Errorf("cannot sign: %w", err)
	}

	outfile, err := os.OpenFile(signaturePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot open output file: %w", err)
	}

	if _, err := io.Copy(outfile, signature); err != nil {
		return fmt.Errorf("cannot copy: %w", err)
	}

	return nil
}
