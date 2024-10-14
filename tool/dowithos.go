package tool

import (
	"context"
	"fmt"
	"os"
)

// DoWithOS is wrapper for Do with NewContext(context.Background()) and os.Args.
// It also handles returned error and exit code.
func DoWithOS() {
	ctx := NewContext(context.Background())
	defer ctx.Cancel()
	if err := do(ctx, os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(ExitCode(err))
	}
}
