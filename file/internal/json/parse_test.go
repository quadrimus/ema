package json

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"quadrimus.com/ema/file/internal"
)

func assertParseArray(t *testing.T, expected []any, source string) {
	t.Helper()
	reader := internal.NewReader("", strings.NewReader(source))
	a, err := parseArray(reader, nil)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, a)
	}
}

func assertParseArrayError(t *testing.T, expected, source string) {
	t.Helper()
	reader := internal.NewReader("", strings.NewReader(source))
	a, err := parseArray(reader, nil)
	if assert.EqualError(t, err, expected) {
		assert.Nil(t, a)
	}
}

func TestParseArray(t *testing.T) {
	assertParseArray(t, nil, `]`)
	assertParseArrayError(t, "0:1: unexpected comma", `,]`)
	assertParseArray(t, []any{json.Number("1")}, `1]`)
	assertParseArray(t, []any{json.Number("1")}, `1,]`) // accept trailing comma
	assertParseArrayError(t, "0:3: unexpected comma", `1,,]`)
	assertParseArray(t, []any{json.Number("1"), json.Number("2")}, `1,2,]`)
	assertParseArray(t, []any{json.Number("1"), json.Number("2")}, `1,2]`)
}

func TestHexDigit(t *testing.T) {
	assert.Equal(t, -1, hexDigit(0))
	assert.Equal(t, -1, hexDigit(' '))
	for r := 'G'; r <= 'Z'; r++ {
		assert.Equal(t, -1, hexDigit(r))
		assert.Equal(t, -1, hexDigit(r+('a'-'A')))
	}
	assert.Equal(t, 0, hexDigit('0'))
	assert.Equal(t, 1, hexDigit('1'))
	assert.Equal(t, 2, hexDigit('2'))
	assert.Equal(t, 3, hexDigit('3'))
	assert.Equal(t, 4, hexDigit('4'))
	assert.Equal(t, 5, hexDigit('5'))
	assert.Equal(t, 6, hexDigit('6'))
	assert.Equal(t, 7, hexDigit('7'))
	assert.Equal(t, 8, hexDigit('8'))
	assert.Equal(t, 9, hexDigit('9'))
	assert.Equal(t, 10, hexDigit('A'))
	assert.Equal(t, 11, hexDigit('B'))
	assert.Equal(t, 12, hexDigit('C'))
	assert.Equal(t, 13, hexDigit('D'))
	assert.Equal(t, 14, hexDigit('E'))
	assert.Equal(t, 15, hexDigit('F'))
	assert.Equal(t, 10, hexDigit('a'))
	assert.Equal(t, 11, hexDigit('b'))
	assert.Equal(t, 12, hexDigit('c'))
	assert.Equal(t, 13, hexDigit('d'))
	assert.Equal(t, 14, hexDigit('e'))
	assert.Equal(t, 15, hexDigit('f'))
}
