package text

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"quadrimus.com/ema/file/internal"
)

func Parse(reader *internal.Reader) (string, error) {
	for sb := (strings.Builder{}); ; {
		r, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return sb.String(), nil
			}
			return "", fmt.Errorf("text.Parse: %w", err)
		}
		if r == '{' {
			r, err = reader.ReadRune()
			if err != nil {
				return "", fmt.Errorf("text.Parse: %w", err)
			}
			if r != '{' {
				reader.UnreadRune()
				reader.UnreadRune()
				return sb.String(), nil
			}
		}
		sb.WriteRune(r)
	}
}
