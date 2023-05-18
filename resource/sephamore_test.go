package resource_test

import (
	"sync"
	"testing"

	"github.com/CSXL/go-agent/resource"
	"github.com/stretchr/testify/assert"
)

func TestSemaphore(t *testing.T) {
	t.Run("allocate and deallocate", func(t *testing.T) {
		sem := resource.NewSemaphore(2)

		err := sem.Allocate()
		assert.NoError(t, err)

		err = sem.Deallocate()
		assert.NoError(t, err)
	})

	t.Run("max count reached", func(t *testing.T) {
		sem := resource.NewSemaphore(1)

		err := sem.Allocate()
		assert.NoError(t, err)

		err = sem.Allocate()
		assert.Error(t, err)
	})

	t.Run("no slots to release", func(t *testing.T) {
		sem := resource.NewSemaphore(1)

		err := sem.Deallocate()
		assert.Error(t, err)
	})

	t.Run("concurrent access", func(t *testing.T) {
		sem := resource.NewSemaphore(5)
		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := sem.Allocate()
				assert.NoError(t, err)

				err = sem.Deallocate()
				assert.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}
