package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"quadrimus.com/ema/file/internal"
	"quadrimus.com/ema/file/internal/json"
	"quadrimus.com/ema/file/internal/text"
)

// Document consists of data or text parts.
type Document struct {
	Content []Content
}

// ParseDocument read document from passed source.
func ParseDocument(file string, source io.Reader) (*Document, error) {
	d, err := parseDocument(file, source)
	if err != nil {
		return nil, fmt.Errorf("file.ParseDocument: %w", err)
	}
	return d, nil
}

func parseDocument(file string, source io.Reader) (*Document, error) {
	d := &Document{}
	for reader := internal.NewReader(file, source); ; {
		r, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if r == '{' {
			r, err = reader.ReadRune()
			if err != nil {
				return nil, err
			}
			if r == '{' {
				reader.UnreadRune()
				reader.UnreadRune()
				t, err := text.Parse(reader)
				if err != nil {
					return nil, err
				}
				d.Content = append(d.Content, Text(t))
				continue
			}
			reader.UnreadRune()
			o, err := json.Parse(reader)
			if err != nil {
				return nil, err
			}
			d.Content = append(d.Content, Data(o))
		} else {
			reader.UnreadRune()
			t, err := text.Parse(reader)
			if err != nil {
				return nil, err
			}
			d.Content = append(d.Content, Text(t))
		}
	}
	return d, nil
}

func (d *Document) ToEMA() ([]byte, error) {
	b := &bytes.Buffer{}
	for _, v := range d.Content {
		if _, err := v.WriteTo(b); err != nil {
			return nil, fmt.Errorf("file.Document.ToEMA: %w", err)
		}
	}
	return b.Bytes(), nil
}

func (d *Document) ToJSON() []any {
	list := make([]any, len(d.Content))
	for i, v := range d.Content {
		list[i] = v.content()
	}
	return list
}
