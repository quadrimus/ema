package internal

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Rune struct {
	r        rune
	position int
	line     int
}

type Reader struct {
	file   string
	reader *bufio.Reader
	buffer []Rune
	start  int
}

func NewReader(file string, r io.Reader) *Reader {
	if file != "" {
		file += ":"
	}
	return &Reader{
		file:   file,
		reader: bufio.NewReader(r),
		buffer: make([]Rune, 2),
		start:  0,
	}
}

func (r *Reader) FileErrorf(format string, args ...any) error {
	s := fmt.Sprintf(format, args...)
	if len(r.buffer) != 0 {
		return fmt.Errorf("%s%d:%d: %s", r.file, r.buffer[0].line, r.buffer[0].position, s)
	}
	return fmt.Errorf("%s%s", r.file, s) // TODO space
}

func (r *Reader) SkipSpaces() error {
	for {
		x, err := r.readRune()
		if err != nil {
			return fmt.Errorf("internal.Reader.SkipSpaces: %w", err)
		}
		if !unicode.IsSpace(x) {
			r.start++
			return nil
		}
	}
}

func (r *Reader) ReadRune() (rune, error) {
	x, err := r.readRune()
	if err != nil {
		return x, fmt.Errorf("internal.Reader.ReadRune: %w", err)
	}
	return x, err
}

func (r *Reader) UnreadRune() {
	r.start++
}

func (r *Reader) readRune() (rune, error) {
	if r.start > 0 {
		r.start--
		return r.buffer[r.start].r, nil
	}
	copy(r.buffer[1:], r.buffer)
	v := &r.buffer[0]
	err := error(nil)
	v.r, _, err = r.reader.ReadRune()
	v.position = r.buffer[1].position
	v.line = r.buffer[1].line
	if v.r == '\n' {
		v.position = 0
		v.line++
	} else {
		v.position++
	}
	return v.r, err
}
