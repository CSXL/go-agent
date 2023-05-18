package resource

import (
	"errors"
	"sync"
)

// Semaphore is a simple implementation of a counting semaphore.
type Semaphore struct {
	count    int
	maxCount int
	mu       sync.Mutex
}

// NewSemaphore creates a new Semaphore with the given maximum count.
func NewSemaphore(maxCount int) *Semaphore {
	return &Semaphore{
		count:    0,
		maxCount: maxCount,
	}
}

// Allocate reserves a slot in the semaphore, blocking if the maximum count is reached.
func (s *Semaphore) Allocate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count >= s.maxCount {
		return errors.New("semaphore: maximum count reached")
	}

	s.count++
	return nil
}

// Deallocate releases a slot in the semaphore.
func (s *Semaphore) Deallocate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.count <= 0 {
		return errors.New("semaphore: no slots to release")
	}

	s.count--
	return nil
}
