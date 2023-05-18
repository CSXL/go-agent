package agent_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/CSXL/go-agent"
	"github.com/CSXL/go-agent/resource"
	"github.com/CSXL/go-agent/task"
)

func TestAgent(t *testing.T) {
	workerCount := 4
	a := agent.NewAgent(workerCount)
	a.Start()

	// Register a shared semaphore resource.
	sem := resource.NewSemaphore(2) // Allow 2 concurrent tasks to use the resource.
	err := a.RegisterResource("semaphore", sem)
	if err != nil {
		t.Fatal("failed to register resource:", err)
	}

	// Create and submit tasks.
	taskCount := 8
	// Waitgroup to wait for all tasks to be allocated.
	for i := 0; i < taskCount; i++ {
		fn := func(ctx context.Context) error {
			fmt.Println("Task", i, "started")
			err := a.AllocateResource("semaphore")
			if err != nil {
				return err
			}
			defer a.DeallocateResource("semaphore") // nolint: errcheck

			// Simulate some work.
			time.Sleep(time.Millisecond * 100)
			fmt.Println("Task", i, "finished")
			return nil
		}

		task := task.NewTask(fmt.Sprintf("task-%d", i), fn, task.MediumPriority)
		a.SubmitTask(task)
	}
	a.SoftStop()
}

func TestAgentResourceError(t *testing.T) {
	workerCount := 4
	a := agent.NewAgent(workerCount)
	a.Start()

	// Create and submit a task that tries to allocate a non-existent resource.
	fn := func(ctx context.Context) error {
		err := a.AllocateResource("non-existent")
		if err == nil {
			return errors.New("expected error when allocating non-existent resource")
		}
		return nil
	}

	task := task.NewTask("task", fn, task.MediumPriority)
	a.SubmitTask(task)
	a.SoftStop()
}

func TestAgentTaskCancellation(t *testing.T) {
	workerCount := 4
	a := agent.NewAgent(workerCount)
	a.Start()

	// Create and submit a long-running task.
	fn := func(ctx context.Context) error {
		fmt.Println("Task started")
		select {
		case <-ctx.Done():
			fmt.Println("Task cancelled")
			return ctx.Err()
		case <-time.After(time.Second * 10):
			fmt.Println("Task finished")
			return nil
		}
	}

	task := task.NewTask("task", fn, task.MediumPriority)
	a.SubmitTask(task)

	// Cancel the task after 1 seconds.
	time.Sleep(time.Millisecond * 500)
	task.Cancel()
	a.SoftStop()
}
