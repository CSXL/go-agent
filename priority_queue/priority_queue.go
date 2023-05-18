package priority_queue

import (
	"container/heap"
	"sync"

	"github.com/CSXL/go-agent/task"
)

// PriorityQueue represents a task priority queue.
type PriorityQueue struct {
	items []*task.Task
	mu    sync.Mutex
}

// NewPriorityQueue creates a new PriorityQueue.
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		items: make([]*task.Task, 0),
	}
	heap.Init(pq)
	return pq
}

// acquireLock acquires the lock on the priority queue.
func (pq *PriorityQueue) acquireLock() {
	pq.mu.Lock()
}

// releaseLock releases the lock on the priority queue.
func (pq *PriorityQueue) releaseLock() {
	pq.mu.Unlock()
}

// Len returns the number of items in the priority queue.
func (pq *PriorityQueue) Len() int {
	pq.acquireLock()
	defer pq.releaseLock()
	return len(pq.items)
}

// Less compares the priority of two tasks in the priority queue.
func (pq *PriorityQueue) Less(i, j int) bool {
	pq.acquireLock()
	defer pq.releaseLock()
	return pq.items[i].Priority() > pq.items[j].Priority()
}

// Swap swaps two tasks in the priority queue.
func (pq *PriorityQueue) Swap(i, j int) {
	pq.acquireLock()
	defer pq.releaseLock()
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	pq.acquireLock()
	defer pq.releaseLock()
	pq.items = append(pq.items, x.(*task.Task))
}

// Pop removes the highest-priority item from the priority queue.
func (pq *PriorityQueue) Pop() interface{} {
	pq.acquireLock()
	defer pq.releaseLock()
	old := pq.items
	n := len(old)
	item := old[n-1]
	pq.items = old[0 : n-1]
	return item
}

// Peek returns the highest-priority item from the priority queue without removing it.
func (pq *PriorityQueue) Peek() *task.Task {
	pq.acquireLock()
	defer pq.releaseLock()
	if len(pq.items) == 0 {
		return nil
	}
	n := len(pq.items) - 1
	item := pq.items[n]
	return item
}

// Remove removes a specific task from the priority queue.
func (pq *PriorityQueue) Remove(t *task.Task) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	for i, item := range pq.items {
		if item.ID() == t.ID() {
			pq.mu.Unlock() // Unlock the mutex before calling heap.Remove
			heap.Remove(pq, i)
			pq.mu.Lock() // Lock the mutex again after heap.Remove is done
			break
		}
	}
}

// UpdatePriority updates the priority of a specific task in the priority queue.
func (pq *PriorityQueue) UpdatePriority(t *task.Task, newPriority task.Priority) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	for i, item := range pq.items {
		if item.ID() == t.ID() {
			item.SetPriority(newPriority)
			pq.mu.Unlock() // Unlock the mutex before calling heap.Fix
			heap.Fix(pq, i)
			pq.mu.Lock() // Lock the mutex again after heap.Fix is done
			break
		}
	}
}
