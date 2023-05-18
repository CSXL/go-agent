package main

import (
	"context"
	"fmt"
	"time"

	"github.com/CSXL/go-agent/executor"
	"github.com/CSXL/go-agent/resource"
	"github.com/CSXL/go-agent/scheduler"
	"github.com/CSXL/go-agent/task"
)

func main() {
	// Create a resource manager and register a semaphore resource.
	resourceMgr := resource.NewManager()
	resourceMgr.Register("semaphore", resource.NewSemaphore(2)) // nolint: errcheck

	// Create an executor with 4 worker goroutines.
	exec := executor.NewExecutor(4)

	// Create a task scheduler with the executor and resource manager.
	sch := scheduler.NewScheduler(exec, resourceMgr)
	sch.Start()

	// Define a sample task function.
	taskFunc := func(ctx context.Context) error {
		// Allocate the semaphore resource before executing the task.
		if err := resourceMgr.Allocate("semaphore"); err != nil {
			return err
		}

		// Defer deallocation of the semaphore resource after executing the task.
		defer resourceMgr.Deallocate("semaphore") // nolint: errcheck

		// Simulate task execution.
		fmt.Println("Task started")
		time.Sleep(1 * time.Second)
		fmt.Println("Task completed")

		return nil
	}

	// Create and submit tasks to the scheduler.
	for i := 0; i < 5; i++ {
		t := task.NewTask(fmt.Sprintf("task-%d", i), taskFunc, task.MediumPriority)
		sch.Submit(t)
	}

	sch.SoftStop()
}
