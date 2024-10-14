package file

import (
	"io"
)

// Content represents part of document.
type Content interface {
	// WriterTo allows to write EMA form of content.
	io.WriterTo

	// ContentType always returns "data" or "text".
	// Value "data" represents type Data, value "text" represents type Text.
	ContentType() string

	content() any
}
