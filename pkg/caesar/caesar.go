package caesar

import (
	"io"
	"strings"
)

type CaesarReader struct {
	r     io.Reader
	shift int
}

func (c *CaesarReader) Reader() io.Reader {
	return c.r
}

func (cr *CaesarReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)
	if err != nil && err != io.EOF {
		return n, err
	}

	// Используем strings.Builder для сбора зашифрованного текста
	var result strings.Builder

	for i := 0; i < n; {
		// Получаем следующую руну
		r, size := decodeRune(p[i:n])
		if r == 0 {
			break // Если руна равна 0, выходим из цикла
		}
		// Шифруем руну и добавляем в результат
		result.WriteRune(shiftRune(r, cr.shift))
		i += size
	}

	// Копируем результат в p
	return strings.NewReader(result.String()).Read(p)
}

func decodeRune(b []byte) (r rune, size int) {
	if len(b) == 0 {
		return 0, 0
	}

	r, size = rune(b[0]), 1
	if r >= 0x80 {
		if (b[0] & 0xF0) == 0xF0 {
			size = 4
		} else if (b[0] & 0xE0) == 0xE0 {
			size = 3
		} else if (b[0] & 0xC0) == 0xC0 {
			size = 2
		} else {
			size = 1
		}

		if size > len(b) {
			return 0, 0
		}

		r = 0
		for i := 0; i < size; i++ {
			r = (r << 6) | rune(b[i]&0x3F)
		}
		r |= rune(b[0]&(1<<(8-size)-1)) << (6 * (size - 1))
	}
	return r, size
}

func shiftRune(r rune, shift int) rune {
	const (
		startLower = 'а'
		endLower   = 'я'
		startUpper = 'А'
		endUpper   = 'Я'
		count      = 32
	)

	if r >= startLower && r <= endLower {
		return startLower + (r-startLower+rune(shift)+count)%count
	} else if r >= startUpper && r <= endUpper {
		return startUpper + (r-startUpper+rune(shift)+count)%count
	} else if r == 'ё' {
		return 'ё'
	} else if r == 'Ё' {
		return 'Ё'
	}
	// Если символ не буква, возвращаем как есть
	return r
}

func NewCaesarReader(r io.Reader, shift int) io.Reader {
	return &CaesarReader{r, shift}
}
