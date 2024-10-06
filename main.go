package main

import (
	"evteev/caesar/pkg/caesar"
	"fmt"
	"io"
	"log/slog"
	"os"
)

func main() {
	var filename string
	var shift int

	fmt.Print("enter the filename: ")
	fmt.Scan(&filename)
	fmt.Print("enter the shift: ")
	fmt.Scan(&shift)

	var file *os.File
	var err error

	file, err = os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	slog.Debug("file opened", slog.String("filename", filename))

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
	if err := os.WriteFile(fmt.Sprintf("enc_%s", filename), content, 0666); err != nil {
		slog.Error("cannot write file", slog.String("err", err.Error()))
		panic(err)
	}
}
