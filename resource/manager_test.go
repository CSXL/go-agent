package resource_test

import (
	"sync"
	"testing"

	"github.com/CSXL/go-agent/resource"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	t.Run("register and unregister", func(t *testing.T) {
		m := resource.NewManager()

		sem := resource.NewSemaphore(1)
		err := m.Register("sem", sem)
		assert.NoError(t, err)

		err = m.Unregister("sem")
		assert.NoError(t, err)
	})

	t.Run("allocate and deallocate", func(t *testing.T) {
		m := resource.NewManager()

		sem := resource.NewSemaphore(1)
		err := m.Register("sem", sem)
		assert.NoError(t, err)

		err = m.Allocate("sem")
		assert.NoError(t, err)

		err = m.Deallocate("sem")
		assert.NoError(t, err)
	})

	t.Run("concurrent access", func(t *testing.T) {
		m := resource.NewManager()

		sem := resource.NewSemaphore(5)
		err := m.Register("sem", sem)
		assert.NoError(t, err)

		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := m.Allocate("sem")
				assert.NoError(t, err)

				err = m.Deallocate("sem")
				assert.NoError(t, err)
			}()
		}

		wg.Wait()
	})
}
