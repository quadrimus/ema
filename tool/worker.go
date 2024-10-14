package tool

import (
	"fmt"
	"sync"
)

var (
	registeredWorkersMutex = sync.Mutex{}
	registeredWorkers      = map[[2]string]Worker{}
)

// Worker allows to do Work.
type Worker interface {
	DoWork(Work) error
}

func RegisterWorker(inputFormat, outputFormat string, worker Worker) error {
	formats := [2]string{inputFormat, outputFormat}

	registeredWorkersMutex.Lock()
	defer registeredWorkersMutex.Unlock()

	if _, ok := registeredWorkers[formats]; ok {
		return fmt.Errorf("tool.RegisterWorker: worker for %q to %q already registered", inputFormat, outputFormat)
	}

	registeredWorkers[formats] = worker
	return nil
}

func RegisteredWorker(inputFormat, outputFormat string) Worker {
	formats := [2]string{inputFormat, outputFormat}

	registeredWorkersMutex.Lock()
	defer registeredWorkersMutex.Unlock()

	worker, _ := registeredWorkers[formats]
	return worker
}
