package steno

import (
	"image"
	"io"
)

type Steno interface {
	Embed(in image.Image, data io.Reader) (*image.RGBA, error)
	Extract(in image.Image) (io.Reader, error)
}
