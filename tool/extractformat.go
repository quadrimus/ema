package tool

import (
	"path/filepath"
	"strings"
)

// ExtractFormat returns path from passed arg and derived format.
// Default format is derived from file extension.
// It can be overridden if arg contains # with following format.
func ExtractFormat(arg string) (path, format string) {
	if i := strings.LastIndexByte(arg, '#'); i >= 0 {
		return arg[:i], arg[i+1:]
	}
	ext := filepath.Ext(arg)
	if strings.HasPrefix(ext, ".") {
		return arg, ext[1:]
	}
	return arg, ext
}
