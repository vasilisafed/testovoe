package repository

import "test_task/internal/domain"

type DeviceRepository interface {
	Save(device *domain.Device) error
	GetByID(id string) (*domain.Device, error)
	List() ([]*domain.Device, error)
}
