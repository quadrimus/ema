package file

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDocument(t *testing.T) {
	source := `{EMA{"x": -1, "y": 2.03e+1, "color": [3, 4, 5], "enabled": true, "text": "x\ny\""}}This is -{{{bold[1]}EMA{x} file format!`
	d, err := ParseDocument("", strings.NewReader(source))
	if assert.NoError(t, err) {
		assert.Equal(t, &Document{
			Content: []Content{
				Data{
					"use": []any{
						"EMA",
						map[string]any{
							"x": json.Number("-1"),
							"y": json.Number("2.03e+1"),
							"color": []any{
								json.Number("3"),
								json.Number("4"),
								json.Number("5"),
							},
							"enabled": true,
							"text":    "x\ny\"",
						},
					},
				},
				Text("This is -{"),
				Data{
					"use": []any{"bold", json.Number("1")},
				},
				Text("EMA"),
				Data{
					"use": []any{"x"},
				},
				Text(" file format!"),
			},
		}, d)
	}
}
