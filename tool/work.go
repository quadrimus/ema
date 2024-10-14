package tool

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Work represents one specific run with arguments.
type Work interface {
	Context() context.Context
	Formats() (input, output string)
	Input() []byte
	Output(data []byte) error
}

// NewWorkFromOS returns new Work based on passed parameters.
// Values input and output can be path to file or - (representing stdin and stdout).
// Format is derived from file extension or overridden if # is present.
// See also ExtractFormat.
func NewWorkFromOS(ctx context.Context, input, output string) (w Work, err error) {
	inputPath, inputFormat := ExtractFormat(input)
	outputPath, outputFormat := ExtractFormat(output)

	inputData, err := readInput(inputPath)
	if err != nil {
		return nil, fmt.Errorf("tool.NewWorkFromOS: %w", err)
	}

	return &work{
		ctx:          ctx,
		inputFormat:  inputFormat,
		outputFormat: outputFormat,
		inputData:    inputData,
		outputPath:   outputPath,
	}, nil
}

type work struct {
	ctx          context.Context
	inputFormat  string
	outputFormat string
	inputData    []byte
	outputPath   string
}

func (w *work) Context() context.Context {
	return w.ctx
}

func (w *work) Formats() (input, output string) {
	return w.inputFormat, w.outputFormat
}

func (w *work) Input() []byte {
	return w.inputData
}

func (w *work) Output(data []byte) error {
	if w.outputPath == "" || w.outputPath == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(w.outputPath, data, 0644)
}

func readInput(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

var _ Work = (*work)(nil)
