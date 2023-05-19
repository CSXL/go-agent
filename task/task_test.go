package task_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/CSXL/go-agent/task"
)

func TestNewTask(t *testing.T) {
	t.Run("valid parameters", func(t *testing.T) {
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			return nil
		}, task.LowPriority)

		assert.NotNil(t, tsk)
		assert.Equal(t, "test_task", tsk.ID())
		assert.Equal(t, task.LowPriority, tsk.Priority())
	})

	t.Run("nil function", func(t *testing.T) {
		tsk := task.NewTask("test_task", nil, task.LowPriority)

		assert.NotNil(t, tsk)
		assert.Equal(t, "test_task", tsk.ID())
		assert.Equal(t, task.LowPriority, tsk.Priority())
	})
}

func TestTaskExecute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}, task.LowPriority)

		err := tsk.Execute(context.Background())

		assert.NoError(t, err)
	})

	t.Run("canceled context", func(t *testing.T) {
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		}, task.LowPriority)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := tsk.Execute(ctx)

		assert.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("nil function", func(t *testing.T) {
		tsk := task.NewTask("test_task", nil, task.LowPriority)

		err := tsk.Execute(context.Background())

		assert.Error(t, err)
		assert.Equal(t, "task function is nil", err.Error())
	})
}

func TestTaskCancel(t *testing.T) {
	t.Run("cancel before execute", func(t *testing.T) {
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		}, task.LowPriority)

		tsk.Cancel()
		err := tsk.Execute(context.Background())

		assert.NoError(t, err)
	})

	t.Run("cancel during execute", func(t *testing.T) {
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		}, task.LowPriority)

		go func() {
			time.Sleep(100 * time.Millisecond)
			tsk.Cancel()
		}()

		err := tsk.Execute(context.Background())

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestTaskWait(t *testing.T) {
	t.Run("wait for completion", func(t *testing.T) {
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			return nil
		}, task.LowPriority)

		go func() {
			_ = tsk.Execute(context.Background())
		}()

		err := tsk.Wait(context.Background())

		assert.NoError(t, err)
	})

	t.Run("context canceled", func(t *testing.T) {
		tsk := task.NewTask("test_task", func(ctx context.Context) error {
			time.Sleep(500 * time.Millisecond)
			return nil
		}, task.LowPriority)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		go func() {
			_ = tsk.Execute(context.Background())
		}()

		err := tsk.Wait(ctx)

		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestTaskDependencies(t *testing.T) {
	tsk := task.NewTask("test_task", func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}, task.LowPriority)

	otherTsk := task.NewTask("test_task", func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	}, task.LowPriority)
	t.Run("get/add/remove dependencies", func(t *testing.T) {
		tsk.AddDependency(otherTsk)

		assert.Equal(t, 1, len(tsk.Dependencies()))
		assert.Equal(t, otherTsk, tsk.Dependencies()[0])

		expectedDependencies := []*task.Task{otherTsk}
		assert.Equal(t, expectedDependencies, tsk.Dependencies())

		tsk.RemoveDependency(otherTsk)
		assert.Equal(t, 0, len(tsk.Dependencies()))
	})
}
