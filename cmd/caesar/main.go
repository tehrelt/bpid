package main

import (
	"evteev/cipher/pkg/caesar"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <path> <shift>", os.Args[0])
		return
	}

	log := slog.With(slog.String("alg", "caesar"))

	path := os.Args[1]
	shiftRaw := os.Args[2]

	shift, err := strconv.Atoi(shiftRaw)
	if err != nil {
		log.Error("cannot parse shift", slog.String("err", err.Error()), slog.Any("shift", shiftRaw))
		return
	}

	file, err := os.Open(path)
	if err != nil {
		slog.Error("cannot open file", slog.String("err", err.Error()), slog.Any("filename", path))
		return
	}
	defer file.Close()
	slog.Debug("file opened", slog.String("filename", path))

	encrypted, err := caesar.Encrypt(file, shift)
	if err != nil {
		log.Error("cannot encrypt content", slog.String("err", err.Error()), slog.Any("file", file), slog.Int("shift", shift))
		return
	}

	content, err := io.ReadAll(encrypted)
	if err != nil {
		log.Error("cannot read encrypted content", slog.String("err", err.Error()), slog.Any("reader", encrypted))
		return
	}
	if err := os.WriteFile(fmt.Sprintf("%s/caesar_%s", filepath.Dir(path), filepath.Base(path)), content, 0666); err != nil {
		log.Error("cannot write file", slog.String("err", err.Error()))
		return
	}

	log.Info("Успешно зашифровано")
}
