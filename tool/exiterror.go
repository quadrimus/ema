package tool

import (
	"fmt"
)

// ExitError represents error with specified exit code.
type ExitError struct {
	Code int
	Err  error
}

// NewExitError creates new ExitError.
func NewExitError(code int, err error) *ExitError {
	return &ExitError{
		Code: code,
		Err:  err,
	}
}

// ErrString returns string representation of Err.
// If e or Err is nil, it returns value "unknown error".
func (e *ExitError) ErrString() string {
	if e != nil && e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// Error returns string representation of error.
func (e *ExitError) Error() string {
	return fmt.Sprintf("exit code %d: %s", e.ExitCode(), e.ErrString())
}

// ExitCode returns specified exit code.
func (e *ExitError) ExitCode() int {
	if e != nil {
		return e.Code
	}
	return 1
}

// Unwrap returns wrapped error if any.
func (e *ExitError) Unwrap() error {
	if e != nil {
		return e.Err
	}
	return nil
}

var _ error = (*ExitError)(nil)
