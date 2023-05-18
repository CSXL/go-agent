package resource

import (
	"errors"
	"sync"
)

// Manager is responsible for managing shared resources among concurrent tasks.
type Manager struct {
	resources map[string]Resource
	mu        sync.Mutex
}

// NewManager creates a new resource Manager.
func NewManager() *Manager {
	return &Manager{
		resources: make(map[string]Resource),
	}
}

// Register registers a shared resource with the Manager.
func (m *Manager) Register(name string, res Resource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.resources[name]; exists {
		return errors.New("resource already registered")
	}

	m.resources[name] = res
	return nil
}

// Unregister removes a shared resource from the Manager.
func (m *Manager) Unregister(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.resources[name]; !exists {
		return errors.New("resource not found")
	}

	delete(m.resources, name)
	return nil
}

// Allocate reserves the specified resource for use by a task.
func (m *Manager) Allocate(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	res, exists := m.resources[name]
	if !exists {
		return errors.New("resource not found")
	}

	return res.Allocate()
}

// Deallocate releases the specified resource previously reserved by a task.
func (m *Manager) Deallocate(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	res, exists := m.resources[name]
	if !exists {
		return errors.New("resource not found")
	}

	return res.Deallocate()
}
