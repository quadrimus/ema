package json

import (
	"encoding/json"
	"fmt"
	"strings"

	"quadrimus.com/ema/file/internal"
)

func Parse(reader *internal.Reader) (map[string]any, error) {
	o, err := parseObject(reader)
	if err != nil {
		return nil, fmt.Errorf("json.Parse: %w", err)
	}
	return o, nil
}

func parseObject(reader *internal.Reader) (map[string]any, error) {
	for o, start, lastKey := map[string]any{}, true, (*string)(nil); ; {
		if err := reader.SkipSpaces(); err != nil {
			return nil, err
		}
		r, err := reader.ReadRune()
		if err != nil {
			return nil, err
		}
		switch {
		case r == '}':
			return o, nil
		case r == ':' || r == ',':
			// ignore
		case r == '{':
			if lastKey == nil {
				return nil, reader.FileErrorf("expected object key instead of object")
			}
			p, err := parseObject(reader)
			if err != nil {
				return nil, err
			}
			o[*lastKey] = p
			lastKey = nil
		case r == '[':
			if lastKey == nil {
				return nil, reader.FileErrorf("expected object key instead of array")
			}
			p, err := parseArray(reader, nil)
			if err != nil {
				return nil, err
			}
			o[*lastKey] = p
			lastKey = nil
		case r == '"':
			p, err := parseString(reader)
			if err != nil {
				return nil, err
			}
			if lastKey == nil {
				lastKey = &p
				start = false
			} else {
				o[*lastKey] = p
				lastKey = nil
			}
		case isNumberStart(r):
			if lastKey == nil {
				return nil, reader.FileErrorf("expected object key instead of number")
			}
			p, err := parseNumber(reader)
			if err != nil {
				return nil, err
			}
			o[*lastKey] = p
			lastKey = nil
		case isLetter(r):
			if start {
				p, err := parseCommand(reader)
				if err != nil {
					return nil, err
				}
				o["use"] = p
				return o, nil
			}
			if lastKey == nil {
				return nil, reader.FileErrorf("expected object key instead of keyword")
			}
			p, err := parseKeyword(reader)
			if err != nil {
				return nil, err
			}
			o[*lastKey] = p
			lastKey = nil
		}
	}
}

func parseArray(reader *internal.Reader, a []any) ([]any, error) {
	for comma := false; ; {
		if err := reader.SkipSpaces(); err != nil {
			return nil, err
		}
		r, err := reader.ReadRune()
		if err != nil {
			return nil, err
		}
		if comma {
			switch r {
			case ']':
				return a, nil
			case ',':
				comma = false
				continue
			default:
				return nil, reader.FileErrorf("expected array end or comma")
			}
		}
		switch {
		case r == ']':
			return a, nil
		case r == ',':
			return nil, reader.FileErrorf("unexpected comma")
		case r == '{':
			p, err := parseObject(reader)
			if err != nil {
				return nil, err
			}
			a = append(a, p)
			comma = true
		case r == '[':
			p, err := parseArray(reader, nil)
			if err != nil {
				return nil, err
			}
			a = append(a, p)
			comma = true
		case r == '"':
			p, err := parseString(reader)
			if err != nil {
				return nil, err
			}
			a = append(a, p)
			comma = true
		case isNumberStart(r):
			p, err := parseNumber(reader)
			if err != nil {
				return nil, err
			}
			a = append(a, p)
			comma = true
		case isLetter(r):
			p, err := parseKeyword(reader)
			if err != nil {
				return nil, err
			}
			a = append(a, p)
			comma = true
		}
	}
}

func hexDigit(r rune) int {
	if r >= '0' && r <= '9' {
		return int(r - '0')
	}
	if r >= 'A' && r <= 'F' {
		return int(r - ('A' - 10))
	}
	if r >= 'a' && r <= 'f' {
		return int(r - ('a' - 10))
	}
	return -1
}

func parseString(reader *internal.Reader) (string, error) {
	for sb, escape, u := (strings.Builder{}), 0, uint32(0); ; {
		r, err := reader.ReadRune()
		if err != nil {
			return "", err
		}
		switch escape {
		case 1: // means \ found
			switch r {
			case 'u':
				escape = 2
				u = 0
				continue
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			case '/':
				sb.WriteByte('/')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			default:
				return "", reader.FileErrorf("invalid string escape %c", r)
			}
			escape = 0
			continue
		case 2, 3, 4, 5: // means \u found
			v := hexDigit(r)
			if r < 0 {
				return "", reader.FileErrorf("invalid string escape %c", r)
			}
			u |= uint32(v << (4 * (5 - escape)))
			if escape == 5 {
				sb.WriteRune(rune(u))
				escape = 0
			} else {
				escape++
			}
			continue
		}
		switch r {
		case '"':
			return sb.String(), nil
		case '\\':
			escape = 1
		default:
			sb.WriteRune(r)
		}
	}
}

func isNumberStart(r rune) bool {
	return r == '-' || (r >= '0' && r <= '9')
}

type number int

const (
	numberStart = number(iota)
	numberMinus
	numberZero
	numberDigit
	numberDot
	numberFraction
	numberE
	numberExponentSign
	numberExponent
)

func parseNumber(reader *internal.Reader) (json.Number, error) {
	n := numberStart
	reader.UnreadRune()
	for sb := (strings.Builder{}); ; {
		r, err := reader.ReadRune()
		if err != nil {
			return "", err
		}
		switch n {
		case numberStart:
			switch {
			case r == '-':
				n = numberMinus
			case r == '0':
				n = numberZero
			case r >= '1' && r <= '9':
				n = numberDigit
			default:
				return "", reader.FileErrorf("expected minus or number digit instead of %c", r)
			}
		case numberMinus:
			switch {
			case r == '0':
				n = numberZero
			case r >= '1' && r <= '9':
				n = numberDigit
			default:
				return "", reader.FileErrorf("expected number digit instead of %c", r)
			}
		case numberZero:
			switch {
			case r == '.':
				n = numberDot
			case r == 'E' || r == 'e':
				n = numberE
			default:
				reader.UnreadRune()
				return json.Number(sb.String()), nil
			}
		case numberDigit:
			switch {
			case r >= '0' && r <= '9':
				n = numberDigit
			case r == '.':
				n = numberDot
			case r == 'E' || r == 'e':
				n = numberE
			default:
				reader.UnreadRune()
				return json.Number(sb.String()), nil
			}
		case numberDot:
			switch {
			case r >= '0' && r <= '9':
				n = numberDigit
			default:
				return "", reader.FileErrorf("expected number fraction digit instead of %c", r)
			}
		case numberFraction:
			switch {
			case r >= '0' && r <= '9':
				n = numberDigit
			case r == 'E' || r == 'e':
				n = numberE
			default:
				reader.UnreadRune()
				return json.Number(sb.String()), nil
			}
		case numberE:
			switch {
			case r >= '0' && r <= '9':
				n = numberExponent
			case r == '+' || r == '-':
				n = numberExponentSign
			default:
				return "", reader.FileErrorf("expected number exponent digit instead of %c", r)
			}
		case numberExponentSign:
			switch {
			case r >= '0' && r <= '9':
				n = numberExponent
			default:
				return "", reader.FileErrorf("expected number exponent digit instead of %c", r)
			}
		case numberExponent:
			switch {
			case r >= '0' && r <= '9':
				n = numberDigit
			default:
				reader.UnreadRune()
				return json.Number(sb.String()), nil
			}
		}
		sb.WriteRune(r)
	}
}

func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func parseCommand(reader *internal.Reader) ([]any, error) {
	reader.UnreadRune()
	sb := strings.Builder{}
	for {
		r, err := reader.ReadRune()
		if err != nil {
			return nil, err
		}
		if isLetter(r) || (r >= '0' && r <= '9') || r == '-' {
			sb.WriteRune(r)
			continue
		}
		reader.UnreadRune()
		break
	}
	a := []any{sb.String()}
	if err := reader.SkipSpaces(); err != nil {
		return nil, err
	}
	r, err := reader.ReadRune()
	if err != nil {
		return nil, err
	}
	switch r {
	case '{':
		p := any(nil)
		p, err = parseObject(reader)
		if err != nil {
			return nil, err
		}
		a = append(a, p)
		if err := reader.SkipSpaces(); err != nil {
			return nil, err
		}
		r, err = reader.ReadRune()
		if err != nil {
			return nil, err
		}
	case '[':
		a, err = parseArray(reader, a)
		if err != nil {
			return nil, err
		}
		if err := reader.SkipSpaces(); err != nil {
			return nil, err
		}
		r, err = reader.ReadRune()
		if err != nil {
			return nil, err
		}
	}
	if r == '}' {
		return a, nil
	}
	return nil, reader.FileErrorf("unexpected character %c", r)
}

type keyword struct {
	name  string
	value any
}

var (
	keywordFalse = keyword{
		name:  "false",
		value: false,
	}

	keywordNull = keyword{
		name:  "null",
		value: nil,
	}

	keywordTrue = keyword{
		name:  "true",
		value: true,
	}
)

func keywordByNameStart(first rune) keyword {
	switch first {
	case 'f':
		return keywordFalse
	case 'n':
		return keywordNull
	case 't':
		return keywordTrue
	default:
		return keyword{}
	}
}

func parseKeyword(reader *internal.Reader) (any, error) {
	reader.UnreadRune()
	r, err := reader.ReadRune()
	if err != nil {
		return nil, err
	}
	k := keywordByNameStart(r)
	if k.name == "" {
		return nil, fmt.Errorf("unexpected character %c", r)
	}
	for i := 1; i < len(k.name); i++ {
		r, err = reader.ReadRune()
		if err != nil {
			return nil, err
		}
		if rune(k.name[i]) != r {
			return nil, fmt.Errorf("expected %s", k)
		}
	}
	return k.value, nil
}
