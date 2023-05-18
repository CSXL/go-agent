package scheduler

import (
	"sync"

	"github.com/CSXL/go-agent/executor"
	"github.com/CSXL/go-agent/resource"
	"github.com/CSXL/go-agent/task"
)

// Scheduler is responsible for managing and scheduling tasks for execution.
type Scheduler struct {
	mu          sync.Mutex
	executor    *executor.Executor
	resourceMgr *resource.Manager
	priorityMap map[task.Priority][]*task.Task
}

// NewScheduler creates a new Scheduler with the given executor and resource manager.
func NewScheduler(executor *executor.Executor, resourceMgr *resource.Manager) *Scheduler {
	return &Scheduler{
		executor:    executor,
		resourceMgr: resourceMgr,
		priorityMap: make(map[task.Priority][]*task.Task),
	}
}

// Start initializes the scheduler and starts the executor.
func (s *Scheduler) Start() {
	s.executor.Start()
}

// Stop gracefully shuts down the scheduler and its executor.
func (s *Scheduler) Stop() {
	s.executor.Stop()
}

// SoftStop gracefully shuts down the scheduler and its executor after all tasks have completed.
func (s *Scheduler) SoftStop() {
	// Wait for the priorty map to be empty.
	for {
		s.mu.Lock()
		empty := len(s.priorityMap[task.HighPriority]) == 0 &&
			len(s.priorityMap[task.MediumPriority]) == 0 &&
			len(s.priorityMap[task.LowPriority]) == 0
		s.mu.Unlock()

		if empty {
			break
		}
	}
	s.Stop()
}

// Submit adds a task to the scheduler's task queue.
func (s *Scheduler) Submit(t *task.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.priorityMap[t.Priority()] = append(s.priorityMap[t.Priority()], t)
	s.scheduleTasks()
}

// scheduleTasks schedules tasks from the priority map for execution.
func (s *Scheduler) scheduleTasks() {
	for _, priority := range []task.Priority{task.HighPriority, task.MediumPriority, task.LowPriority} {
		taskList := s.priorityMap[priority]
		for len(taskList) > 0 {
			t := taskList[0]
			taskList = taskList[1:]
			s.priorityMap[priority] = taskList

			s.executor.Submit(t)
		}
	}
}
