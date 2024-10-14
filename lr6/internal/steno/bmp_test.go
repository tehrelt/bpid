package steno_test

import (
	"bytes"
	"evteev/bpd/lr6/internal/steno"
	"image"
	"image/color"
	"image/draw"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func generateImage(width, height int, imgColor color.Color) image.Image {
	// Create a new blank RGBA image with the specified dimensions.
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with the given color.
	draw.Draw(img, img.Bounds(), &image.Uniform{imgColor}, image.Point{}, draw.Src)

	return img
}

func compareReaders(r1, r2 io.Reader) (bool, error) {
	const chunkSize = 1024 // Define the chunk size to read data in parts.

	buf1 := make([]byte, chunkSize)
	buf2 := make([]byte, chunkSize)

	for {
		// Read chunks from both readers.
		n1, err1 := r1.Read(buf1)
		n2, err2 := r2.Read(buf2)

		// If the number of bytes read is different, the readers are not equal.
		if n1 != n2 {
			return false, nil
		}

		// If both reached EOF, the readers are equal.
		if err1 == io.EOF && err2 == io.EOF {
			return true, nil
		}

		// If one reached EOF but the other didn't, they are not equal.
		if (err1 == io.EOF || err2 == io.EOF) && err1 != err2 {
			return false, nil
		}

		// If there are other read errors, return the error.
		if err1 != nil && err1 != io.EOF {
			return false, err1
		}
		if err2 != nil && err2 != io.EOF {
			return false, err2
		}

		// Compare the bytes read from both readers.
		if !bytes.Equal(buf1[:n1], buf2[:n2]) {
			return false, nil
		}
	}
}

func TestBmpSteno(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"t", "1"},
		{"t", "maxim bogomolov"},
		{"t", "fdwoqjfowpqjnfwqopfjnqwpofhwqopfopqwjof"},
		{"t", "lol"},
	}

	img := generateImage(200, 200, color.Black)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			steno := steno.NewBMPSteno()
			embed, _ := steno.Embed(img, strings.NewReader(c.in))
			extracted, _ := steno.Extract(embed)

			actual, _ := io.ReadAll(extracted)

			equals := strings.Compare(string(c.in), string(actual)) == 0
			if !equals {
				t.Logf("in:  %s", c.in)
				t.Logf("out: %s", actual)
				t.Errorf("reader dont equals")
			}
		})
	}
}
