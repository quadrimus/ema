package tool

import (
	"errors"
)

const (
	ExitCodeUnknownError   = 1
	ExitCodeInvalidArgs    = 2
	ExitCodeInvalidFormats = 3
)

// ExitCode returns exit code from passed error.
// If passed error or its wrapped error doesn't implement ExitCode, it returns ExitCodeUnknownError.
// See also ExitError.
func ExitCode(from error) int {
	if b := ExitCodeBearer(nil); errors.As(from, &b) {
		return b.ExitCode()
	}
	return ExitCodeUnknownError
}
