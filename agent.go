package agent

import (
	"github.com/CSXL/go-agent/executor"
	"github.com/CSXL/go-agent/resource"
	"github.com/CSXL/go-agent/scheduler"
	"github.com/CSXL/go-agent/task"
)

// Agent is a high-level interface for managing concurrent tasks and shared resources.
type Agent struct {
	executor    *executor.Executor
	resourceMgr *resource.Manager
	scheduler   *scheduler.Scheduler
}

// NewAgent creates a new Agent with the given number of workers.
func NewAgent(workerCount int) *Agent {
	exec := executor.NewExecutor(workerCount)
	resMgr := resource.NewManager()
	sched := scheduler.NewScheduler(exec, resMgr)

	return &Agent{
		executor:    exec,
		resourceMgr: resMgr,
		scheduler:   sched,
	}
}

// Start initializes the agent and starts its scheduler and executor.
func (a *Agent) Start() {
	a.scheduler.Start()
}

// Stop gracefully shuts down the agent and its scheduler and executor.
func (a *Agent) Stop() {
	a.scheduler.Stop()
}

// SoftStop gracefully shuts down the agent and its scheduler and executor after all tasks have completed.
func (a *Agent) SoftStop() {
	a.scheduler.SoftStop()
}

// RegisterResource registers a shared resource with the agent's resource manager.
func (a *Agent) RegisterResource(name string, res resource.Resource) error {
	return a.resourceMgr.Register(name, res)
}

// UnregisterResource removes a shared resource from the agent's resource manager.
func (a *Agent) UnregisterResource(name string) error {
	return a.resourceMgr.Unregister(name)
}

// AllocateResource reserves the specified resource for use by a task.
func (a *Agent) AllocateResource(name string) error {
	return a.resourceMgr.Allocate(name)
}

// DeallocateResource releases the specified resource previously reserved by a task.
func (a *Agent) DeallocateResource(name string) error {
	return a.resourceMgr.Deallocate(name)
}

// SubmitTask submits a task to the agent's scheduler for execution.
func (a *Agent) SubmitTask(t *task.Task) {
	a.scheduler.Submit(t)
}
