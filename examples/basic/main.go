package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CSXL/go-agent"
	"github.com/CSXL/go-agent/resource"
	"github.com/CSXL/go-agent/task"
)

func main() {
	workerCount := 4
	a := agent.NewAgent(workerCount)
	a.Start()

	// Register a shared semaphore resource.
	sem := resource.NewSemaphore(2)      // Allow 2 concurrent tasks to use the resource.
	a.RegisterResource("semaphore", sem) // nolint: errcheck

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
	time.Sleep(time.Second * 2)
}
