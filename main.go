package main

import (
	"evteev/caesar/pkg/caesar"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	var path string
	var shift int

	fmt.Print("enter the filepath: ")
	fmt.Scan(&path)

	var file *os.File
	var err error

	file, err = os.Open(path)
	if err != nil {
		slog.Error("cannot open file", slog.String("err", err.Error()), slog.Any("filename", path))
		return
	}
	defer file.Close()
	slog.Debug("file opened", slog.String("filename", path))

	fmt.Print("enter the shift: ")
	fmt.Scan(&shift)

	encrypted, err := caesar.Encrypt(file, shift)
	if err != nil {
		slog.Error("cannot encrypt content", slog.String("err", err.Error()), slog.Any("file", file), slog.Int("shift", shift))
		panic(err)
	}

	content, err := io.ReadAll(encrypted)
	if err != nil {
		slog.Error("cannot read encrypted content", slog.String("err", err.Error()), slog.Any("reader", encrypted))
		panic(err)
	}
	if err := os.WriteFile(fmt.Sprintf("%s/enc_%s", filepath.Dir(path), filepath.Base(path)), content, 0666); err != nil {
		slog.Error("cannot write file", slog.String("err", err.Error()))
		panic(err)
	}

	slog.Info("Успешно зашифровано")
}
