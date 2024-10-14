package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Data represents part of document.
// It can be stored as JSON object.
type Data map[string]any

// Command returns command name, args and true if Data is possible to represent as command.
func (d Data) Command() (name string, args any, ok bool) {
	if len(d) == 1 {
		for k, v := range d {
			if k != "use" {
				break
			}
			if s, ok := v.([]any); ok && len(s) > 0 {
				name, ok = s[0].(string)
				if !ok || !IsCommandName(name) {
					break
				}
				s = s[1:]
				switch len(s) {
				case 0:
					return name, nil, true
				case 1:
					if o, ok := s[0].(map[string]any); ok {
						return name, o, true
					}
					return name, s, true
				default:
					return name, s, true
				}
			}
		}
	}
	return "", nil, false
}

// ContentType always returns "data".
func (Data) ContentType() string {
	return "data"
}

// WriteTo write EMA form of Data to passed io.Writer.
func (d Data) WriteTo(w io.Writer) (n int64, err error) {
	data, err := d.marshal()
	if err != nil {
		return 0, fmt.Errorf("file.Data.WriteTo: %w", err)
	}
	i, err := w.Write(data)
	if err != nil {
		return int64(i), fmt.Errorf("file.Data.WriteTo: %w", err)
	}
	return int64(i), nil
}

func (d Data) content() any {
	return map[string]any(d)
}

func (d Data) marshal() ([]byte, error) {
	if name, args, ok := d.Command(); ok {
		b := &bytes.Buffer{}
		b.WriteByte('{')
		b.WriteString(name)
		if args == nil {
			b.WriteByte('}')
			return b.Bytes(), nil
		}
		if err := json.NewEncoder(b).Encode(args); err != nil {
			return nil, err
		}
		data := b.Bytes()
		data[len(data)-1] = '}'
		return data, nil
	}
	return json.Marshal(d)
}

var _ Content = (Data)(nil)
