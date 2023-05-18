package executor

import (
	"context"
	"sync"

	"github.com/CSXL/go-agent/task"
)

// Executor is responsible for managing and executing tasks concurrently.
type Executor struct {
	taskQueue   chan *task.Task
	workerCount int
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewExecutor creates a new Executor with the given number of workers.
func NewExecutor(workerCount int) *Executor {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &Executor{
		taskQueue:   make(chan *task.Task),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start initializes the executor and starts the worker goroutines.
func (e *Executor) Start() {
	e.wg.Add(e.workerCount)
	for i := 0; i < e.workerCount; i++ {
		go e.worker()
	}
}

// Stop gracefully shuts down the executor and waits for all workers to finish.
func (e *Executor) Stop() {
	close(e.taskQueue)
	e.wg.Wait()
}

// SoftStop gracefully shuts down the executor and waits for all workers to finish executing their current tasks.
func (e *Executor) SoftStop() {
	// Wait for the task queue to be empty.
	for {
		empty := len(e.taskQueue) == 0
		if empty {
			break
		}
	}
	e.cancel()
	e.Stop()
}

// Submit adds a task to the executor's task queue for execution.
func (e *Executor) Submit(t *task.Task) {
	e.taskQueue <- t
}

// worker represents a background goroutine that executes tasks.
func (e *Executor) worker() {
	defer e.wg.Done()
	for t := range e.taskQueue {
		_ = t.Execute(e.ctx)
	}
}
