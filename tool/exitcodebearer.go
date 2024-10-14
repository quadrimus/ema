package tool

// ExitCodeBearer represents value with specified exit code.
// See also ExitCode or ExitError.
type ExitCodeBearer interface {
	ExitCode() int
}
