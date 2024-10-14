package json2ema

import (
	"encoding/json"
	"fmt"

	"quadrimus.com/ema/file"
	"quadrimus.com/ema/tool"
)

func init() {
	register()
}

func register() {
	if err := tool.RegisterWorker(WorkerInputFormat, WorkerOutputFormat, Worker{}); err != nil {
		panic(fmt.Errorf("json2ema.register: %w", err))
	}
}

const (
	WorkerInputFormat  = "json"
	WorkerOutputFormat = "ema"
)

type Worker struct{}

func (Worker) DoWork(w tool.Work) error {
	s := []any(nil)
	if err := json.Unmarshal(w.Input(), &s); err != nil {
		return fmt.Errorf("json2ema.Worker.DoWork: %w", err)
	}
	d, err := file.DocumentFromJSON(s)
	if err != nil {
		return fmt.Errorf("json2ema.Worker.DoWork: %w", err)
	}
	data, err := d.ToEMA()
	if err != nil {
		return fmt.Errorf("json2ema.Worker.DoWork: %w", err)
	}
	if err := w.Output(data); err != nil {
		return fmt.Errorf("json2ema.Worker.DoWork: %w", err)
	}
	return nil
}
