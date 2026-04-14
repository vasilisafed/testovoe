package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"test_task/internal/domain"
	"test_task/internal/repository"
)

var ErrDeviceNotFound = errors.New("device not found")
var ErrDeviceAlreadyExists = errors.New("device already exists")

type DeviceService struct {
	repo repository.DeviceRepository
}

func NewDeviceService(repo repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) CreateDevice(id string, algo domain.Algorithm, label string) (*domain.Device, error) {
	if id == "" {
		return nil, fmt.Errorf("device id is required")
	}
	if _, err := s.repo.GetByID(id); err == nil {
		return nil, ErrDeviceAlreadyExists
	}

	device := &domain.Device{
		ID:        id,
		Label:     label,
		Algorithm: algo,
	}

	switch algo {
	case domain.RSA:
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		device.PrivateKey = privateKey
		device.PublicKey = &privateKey.PublicKey
	case domain.ECC:
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, err
		}
		device.PrivateKey = privateKey
		device.PublicKey = &privateKey.PublicKey
	default:
		return nil, fmt.Errorf("unsupported algorithm %q", algo)
	}

	if err := s.repo.Save(device); err != nil {
		return nil, err
	}

	return device, nil
}

func (s *DeviceService) GetDevice(id string) (*domain.Device, error) {
	device, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}

func (s *DeviceService) ListDevices() ([]*domain.Device, error) {
	return s.repo.List()
}
