package tool

import (
	"context"
	"errors"
	"fmt"
	"runtime"

	"quadrimus.com/ema"
)

// Do tool with specific arguments.
func Do(ctx context.Context, args []string) error {
	if err := do(ctx, args); err != nil {
		return fmt.Errorf("tool.Do: %w", err)
	}
	return nil
}

func do(ctx context.Context, args []string) error {
	switch len(args) {
	case 3:
		// continue
	case 2:
		if args[1] == "version" {
			fmt.Print("ema " + ema.Version() + " " + runtime.GOOS + " " + runtime.GOARCH)
			return nil
		}
		fallthrough
	default:
		fmt.Print(Help())
		return NewExitError(ExitCodeInvalidArgs, errors.New("required input and output argument"))
	}

	w, err := NewWorkFromOS(ctx, args[1], args[2])
	if err != nil {
		fmt.Print(Help())
		return err
	}

	if worker := RegisteredWorker(w.Formats()); worker != nil {
		return worker.DoWork(w)
	}
	inputFormat, outputFormat := w.Formats()
	return NewExitError(ExitCodeInvalidFormats, fmt.Errorf("no worker for %q to %q", inputFormat, outputFormat))
}
