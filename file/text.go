package file

import (
	"fmt"
	"io"
)

// Text represents part of document.
// It can be stored as plain text.
type Text string

// ContentType always returns "text".
func (Text) ContentType() string {
	return "text"
}

// WriteTo write EMA form of Text to passed io.Writer.
func (t Text) WriteTo(w io.Writer) (n int64, err error) {
	i, err := w.Write([]byte(t))
	if err != nil {
		return int64(i), fmt.Errorf("file.Text.WriteTo: %w", err)
	}
	return int64(i), nil
}

func (t Text) content() any {
	return string(t)
}

var _ Content = Text("")
