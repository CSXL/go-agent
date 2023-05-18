// File: go-agent/priority_queue/priority_queue_test.go
package priority_queue_test

import (
	"testing"

	"github.com/CSXL/go-agent/priority_queue"
	"github.com/CSXL/go-agent/task"
	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	pq := priority_queue.NewPriorityQueue()

	tsk1 := task.NewTask("task1", nil, task.LowPriority)
	tsk2 := task.NewTask("task2", nil, task.MediumPriority)
	tsk3 := task.NewTask("task3", nil, task.HighPriority)

	pq.Push(tsk1)
	pq.Push(tsk2)
	pq.Push(tsk3)

	assert.Equal(t, 3, pq.Len())

	peekedTask := pq.Peek()
	assert.NotNil(t, peekedTask)
	assert.Equal(t, "task3", peekedTask.ID())

	poppedTask := pq.Pop().(*task.Task)
	assert.NotNil(t, poppedTask)
	assert.Equal(t, "task3", poppedTask.ID())

	pq.Remove(tsk1)
	assert.Equal(t, 1, pq.Len())

	pq.UpdatePriority(tsk2, task.HighPriority)
	updatedTask := pq.Pop().(*task.Task)
	assert.NotNil(t, updatedTask)
	assert.Equal(t, "task2", updatedTask.ID())
	assert.Equal(t, task.HighPriority, updatedTask.Priority())
}
