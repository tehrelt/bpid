package main

import (
	"bytes"
	"crypto/md5"
	"evteev/bpd/lr6/internal/aes"
	"evteev/bpd/lr6/internal/steno"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	"golang.org/x/image/bmp"
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
	imagePath  string
	mode       string
	key        string

	debug bool
)

func init() {
	flag.StringVar(&inputPath, "in", DEFAULT_INPUT, "input file")
	flag.StringVar(&outputPath, "out", DEFAULT_OUTPUT, "output file")

	flag.StringVar(&imagePath, "img", DEFAULT_INPUT, "picture .bmp file")

	flag.StringVar(&mode, "mode", DEFAULT_MODE, "mode <(encode|e)|(decode|d)>")
	flag.StringVar(&key, "key", DEFAULT_KEY, "")

	flag.BoolVar(&debug, "debug", false, "debug mode")
}

func main() {
	flag.Parse()

	if debug {
		logfile, err := os.OpenFile(fmt.Sprintf("%s.log", time.Now().Local().Format("2006_01_02T15_04_05")), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Errorf("cannot create logfile: %w", err))
		}
		defer logfile.Close()
		slog.SetDefault(slog.New(slog.NewJSONHandler(logfile, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}

	if mode == DEFAULT_MODE {
		log.Fatal("mode is required")
	}

	if key == DEFAULT_KEY {
		log.Fatal("key is required")
	}

	if mode != "e" && mode != "encode" && mode != "decode" && mode != "d" {
		log.Fatal("mode must be <(encryption|e)|(decryption|d)>")
	}

	cipher, err := aes.New([]byte(key))
	if err != nil {
		panic(err)
	}

	if mode == "encode" || mode == "e" {
		Encode(cipher)
	} else if mode == "decode" || mode == "d" {
		Decode(cipher)
	}
}

func outname(in, out, mbout string) string {
	if out == "" {
		return mbout
	}

	return out
}

func Encode(cipher *aes.Cipher) {

	if inputPath == DEFAULT_INPUT {
		log.Fatal("input file is required")
	}

	if imagePath == DEFAULT_INPUT {
		log.Fatal("input image file is required")
	}

	in, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	img, err := readBmp(imagePath)
	if err != nil {
		panic(err)
	}

	out, err := encode(in, img, cipher)
	if err != nil {
		panic(err)
	}

	if err := writeBmp(outname(imagePath, outputPath, fmt.Sprintf("%s.embedded", imagePath)), out); err != nil {
		panic(fmt.Errorf("failed to write output file: %w", err))
	}
}

func Decode(cipher *aes.Cipher) {

	if imagePath == DEFAULT_INPUT {
		log.Fatal("input image file is required")
	}

	img, err := readBmp(imagePath)
	if err != nil {
		panic(err)
	}

	decoded, hashr, err := decode(img, cipher)
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(decoded)
	if err != nil {
		panic(err)
	}

	hash, err := io.ReadAll(hashr)
	if err != nil {
		panic(err)
	}

	ok, err := verifysum(content, hash)
	if err != nil {
		panic(err)
	}

	if !ok {
		fmt.Printf("хеш-сумма не совпала")
		return
	}

	fmt.Printf("хеш сумма совпадает")

	outfile, err := os.OpenFile(outname(inputPath, outputPath, fmt.Sprintf("%s.extracted", inputPath)), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	slog.Debug("writing result", slog.String("content", string(content)))
	outfile.Write(content)
}

func _md5(in []byte) []byte {
	hash := md5.Sum(in)
	return hash[:]
}

func verifysum(data, expected []byte) (bool, error) {
	actual := _md5(data)
	eq := bytes.Equal(actual, expected)
	slog.Debug("md5 hashsum verified", slog.Any("expected", expected), slog.Any("actual", actual), slog.Any("equal", eq))
	return eq, nil
}

func encode(input io.Reader, img image.Image, cipher *aes.Cipher) (image.Image, error) {
	in, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	hash := _md5(in)
	saturated := append(in, hash...)

	encrypted, err := cipher.Encrypt(bytes.NewReader(saturated))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt file data: %w", err)
	}

	steno := steno.NewBMPSteno()

	embeded, err := steno.Embed(img, encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to embed image: %w", err)
	}

	return embeded, nil
}

func decode(img image.Image, cipher *aes.Cipher) (io.Reader, io.Reader, error) {
	steno := steno.NewBMPSteno()

	encrypted, err := steno.Extract(img)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract: %w", err)
	}

	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	content, err := io.ReadAll(decrypted)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read decrypted: %w", err)
	}

	data := content[:len(content)-md5.Size]
	hash := content[len(content)-md5.Size:]

	return bytes.NewReader(data), bytes.NewReader(hash), nil
}

func readBmp(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return bmp.Decode(f)
}

func writeBmp(filename string, data image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return bmp.Encode(f, data)
}
