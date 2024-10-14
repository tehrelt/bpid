package steno

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"log/slog"
)

type bmpSteno struct {
}

func NewBMPSteno() Steno {
	return &bmpSteno{}
}

// Embed implements Steno.
func (b *bmpSteno) Embed(img image.Image, reader io.Reader) (*image.RGBA, error) {

	log := slog.With(slog.String("fn", "Embed"))
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	sizebuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizebuf, uint32(len(data)))
	data = append(sizebuf, data...)
	dsize := len(data)
	log.Debug("write data", slog.Any("data", data), slog.Int("size", dsize))

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	out := image.NewRGBA(image.Rect(0, 0, width, height))

	var di, bi int

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			if di < dsize {
				rNew := (r & 0xFFFC) | (uint32(data[di]) >> bi & 0x03)
				bi += 2
				if bi >= 8 {
					bi = 0
					log.Debug("written byte", slog.String("val", string(data[di])), slog.Int("di", di), slog.Int("dsize", dsize))
					di++
				}
				out.SetRGBA(x, y, color.RGBA{uint8(rNew), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
			} else {
				out.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
			}
		}
	}

	if di < dsize {
		return nil, fmt.Errorf("not such space")
	}

	return out, nil
}

// Extract implements Steno.
func (b *bmpSteno) Extract(in image.Image) (io.Reader, error) {
	log := slog.With(slog.String("fn", "Extract"))
	bounds := in.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	sizespace := 4
	sizebuf := make([]byte, sizespace)
	var dataidx, bitidx int

	for y := 0; y < height && dataidx < sizespace; y++ {
		for x := 0; x < width && dataidx < sizespace; x++ {
			r, _, _, _ := in.At(x, y).RGBA()
			sizebuf[dataidx] |= byte(r&0x03) << byte(bitidx)
			bitidx += 2
			if bitidx >= 8 {
				bitidx = 0
				log.Debug("read byte", slog.Int("x", x), slog.Int("y", y), slog.Int("dataidx", dataidx), slog.Int("bitidx", bitidx), slog.Any("val", sizebuf[dataidx]))
				dataidx++
			}
		}
	}

	datasize := binary.LittleEndian.Uint32(sizebuf) + uint32(sizespace)
	log.Debug("extracted size", slog.Int("size", int(datasize)))
	dataidx = 0
	bitidx = 0

	data := make([]byte, datasize)

	for y := 0; y < height && dataidx < int(datasize); y++ {
		for x := 0; x < width && dataidx < int(datasize); x++ {
			r, _, _, _ := in.At(x, y).RGBA()

			data[dataidx] |= byte(r&0x03) << bitidx
			bitidx += 2
			if bitidx >= 8 {
				bitidx = 0
				log.Debug("read byte", slog.Int("x", x), slog.Int("y", y), slog.Int("dataidx", dataidx), slog.Int("bitidx", bitidx), slog.Any("val", data[dataidx]))
				dataidx++
			}
		}
	}

	if dataidx < int(datasize) {
		return nil, fmt.Errorf("cannot extract data")
	}

	log.Debug("extracted data", slog.Any("data", data[sizespace:]))

	return bytes.NewReader(data[4:]), nil
}
