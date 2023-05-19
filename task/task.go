package task

import (
	"context"
	"errors"
)

// Priority represents the priority level of a task.
type Priority int

const (
	// LowPriority represents a low-priority task.
	LowPriority Priority = iota

	// MediumPriority represents a medium-priority task.
	MediumPriority

	// HighPriority represents a high-priority task.
	HighPriority
)

// Task represents a unit of work that can be executed concurrently.
type Task struct {
	id           string
	fn           func(context.Context) error
	priority     Priority
	dependencies []*Task
	cancel       context.CancelFunc
	done         chan struct{}
}

// NewTask creates a new task with the given ID, function, and priority.
func NewTask(id string, fn func(context.Context) error, priority Priority) *Task {
	return &Task{
		id:       id,
		fn:       fn,
		priority: priority,
		done:     make(chan struct{}),
	}
}

// AddDependency adds a dependency to the task.
func (t *Task) AddDependency(dependency *Task) {
	t.dependencies = append(t.dependencies, dependency)
}

// Dependencies returns the task's dependencies.
func (t *Task) Dependencies() []*Task {
	return t.dependencies
}

// RemoveDependency removes a dependency from the task.
func (t *Task) RemoveDependency(dependency *Task) {
	for i, dep := range t.dependencies {
		if dep == dependency {
			t.dependencies = append(t.dependencies[:i], t.dependencies[i+1:]...)
			break
		}
	}
}

// ID returns the task's unique identifier.
func (t *Task) ID() string {
	return t.id
}

// Priority returns the task's priority level.
func (t *Task) Priority() Priority {
	return t.priority
}

// SetPriority sets the task's priority level.
func (t *Task) SetPriority(priority Priority) {
	t.priority = priority
}

// Execute runs the task and returns any error encountered during execution.
func (t *Task) Execute(ctx context.Context) error {
	ctx, t.cancel = context.WithCancel(ctx)
	defer close(t.done)
	if t.fn == nil {
		return errors.New("task function is nil")
	}

	// Wrap the task function to handle context cancellation.
	errChan := make(chan error, 1)
	go func() {
		errChan <- t.fn(ctx)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Cancel stops the task's execution if it's currently running.
func (t *Task) Cancel() {
	if t.cancel != nil {
		t.cancel()
	}
}

// Wait blocks until the task has completed execution or the context is canceled.
func (t *Task) Wait(ctx context.Context) error {
	select {
	case <-t.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
