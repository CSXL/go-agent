package resource

// Resource represents a shared resource that can be used by concurrent tasks.
type Resource interface {
	// Allocate reserves the resource for use by a task.
	Allocate() error

	// Deallocate releases the resource previously reserved by a task.
	Deallocate() error
}
