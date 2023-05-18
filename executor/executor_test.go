package executor_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/CSXL/go-agent/executor"
	"github.com/CSXL/go-agent/task"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutor(t *testing.T) {
	ex := executor.NewExecutor(4)
	assert.NotNil(t, ex)
}

func TestExecutorStartStop(t *testing.T) {
	ex := executor.NewExecutor(4)

	ex.Start()
	ex.Stop()

	// Ensure that the executor can be started and stopped without any issues.
}

func TestExecutorSubmit(t *testing.T) {
	ex := executor.NewExecutor(4)
	ex.Start()
	defer ex.Stop()

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			wg.Done()
			return nil
		}, task.LowPriority)

		ex.Submit(tsk)
	}

	wg.Wait()

	// Ensure that the submitted tasks are executed by the executor.
}
