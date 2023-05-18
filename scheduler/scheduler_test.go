package scheduler_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/CSXL/go-agent/executor"
	"github.com/CSXL/go-agent/resource"
	"github.com/CSXL/go-agent/scheduler"
	"github.com/CSXL/go-agent/task"
)

func TestScheduler(t *testing.T) {
	t.Run("submit", func(t *testing.T) {
		ex := executor.NewExecutor(4)
		rm := resource.NewManager()
		sch := scheduler.NewScheduler(ex, rm)
		sch.Start()
		defer sch.Stop()

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			tsk := task.NewTask("test_task", func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				wg.Done()
				return nil
			}, task.LowPriority)
			sch.Submit(tsk)
		}
		wg.Wait()
		// Ensure that the submitted tasks are executed by the scheduler.
	})

	t.Run("softstop", func(t *testing.T) {
		ex := executor.NewExecutor(4)
		rm := resource.NewManager()
		sch := scheduler.NewScheduler(ex, rm)
		sch.Start()

		for i := 0; i < 10; i++ {
			tsk := task.NewTask("test_task", func(ctx context.Context) error {
				time.Sleep(100 * time.Millisecond)
				return nil
			}, task.LowPriority)
			sch.Submit(tsk)
		}
		sch.SoftStop()
	})
}
