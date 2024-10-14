package file

import (
	"fmt"
)

// DocumentFromJSON creates new Document based on JSON array.
func DocumentFromJSON(array []any) (*Document, error) {
	content := make([]Content, len(array))
	for i, v := range array {
		switch x := v.(type) {
		case map[string]any:
			content[i] = Data(x)
		case string:
			content[i] = Text(x)
		default:
			return nil, fmt.Errorf("file.DocumentFromJSON: array value %d: invalid type %T", i, v)
		}
	}
	return &Document{
		Content: content,
	}, nil
}
