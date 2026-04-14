package repository

import (
	"fmt"
	"sort"
	"sync"
	"test_task/internal/domain"
)

type InMemoryDeviceRepository struct {
	mu      sync.RWMutex
	devices map[string]*domain.Device
}

func NewInMemoryDeviceRepository() *InMemoryDeviceRepository {
	devices := make(map[string]*domain.Device)
	return &InMemoryDeviceRepository{sync.RWMutex{}, devices}
}

func (r *InMemoryDeviceRepository) Save(d *domain.Device) error {
	if d == nil {
		return fmt.Errorf("device is nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices[d.ID] = d
	return nil
}

func (r *InMemoryDeviceRepository) GetByID(id string) (*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	device, ok := r.devices[id]
	if !ok {
		return nil, fmt.Errorf("device %q not found", id)
	}

	return device, nil
}

func (r *InMemoryDeviceRepository) List() ([]*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	devices := make([]*domain.Device, 0, len(r.devices))
	for _, device := range r.devices {
		devices = append(devices, device)
	}
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].ID < devices[j].ID
	})

	return devices, nil
}
