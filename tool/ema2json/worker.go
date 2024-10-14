package ema2json

import (
	"bytes"
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
		panic(fmt.Errorf("ema2json.register: %w", err))
	}
}

const (
	WorkerInputFormat  = "ema"
	WorkerOutputFormat = "json"
)

type Worker struct{}

func (Worker) DoWork(w tool.Work) error {
	d, err := file.ParseDocument("", bytes.NewReader(w.Input()))
	if err != nil {
		return fmt.Errorf("ema2json.Worker.DoWork: %w", err)
	}
	data, err := json.Marshal(d.ToJSON())
	if err != nil {
		return fmt.Errorf("ema2json.Worker.DoWork: %w", err)
	}
	if err := w.Output(data); err != nil {
		return fmt.Errorf("ema2json.Worker.DoWork: %w", err)
	}
	return nil
}
