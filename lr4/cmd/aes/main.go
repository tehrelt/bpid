package main

import (
	"evteev/bpid/lr4/internal/aes"
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
	inputPath  string
	outputPath string
	mode       string
	key        string
)

func init() {
	flag.StringVar(&inputPath, "in", DEFAULT_INPUT, "input file")
	flag.StringVar(&outputPath, "out", DEFAULT_OUTPUT, "output file")
	flag.StringVar(&mode, "mode", DEFAULT_MODE, "mode <(encrypt|e)|(decrypt|d)>")
	flag.StringVar(&key, "key", DEFAULT_KEY, "")
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

	if key == DEFAULT_KEY {
		log.Fatal("key is required")
	}

	if mode != "e" && mode != "encrypt" && mode != "decrypt" && mode != "d" {
		log.Fatal("mode must be <(encryption|e)|(decryption|d)>")
	}

	cipher, err := aes.New([]byte(key))
	if err != nil {
		log.Fatalf("failed to create cipher: %v", err)
	}

	in, err := os.Open(inputPath)
	if err != nil {
		panic("cannot open file")
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
