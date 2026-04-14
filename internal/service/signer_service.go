package service

import (
	"encoding/base64"
	"fmt"
	"sync"
	"test_task/internal/domain"
	"test_task/internal/repository"
)

type SignerFactory interface {
	Create(device *domain.Device) (domain.Signer, error)
}

type SignService struct {
	repo    repository.DeviceRepository
	factory SignerFactory
	mu      sync.Mutex // критично!
}

func NewSignService(repo repository.DeviceRepository, factory SignerFactory) *SignService {
	return &SignService{
		repo:    repo,
		factory: factory,
	}
}

func (s *SignService) Sign(deviceID string, data string) (string, string, error) {
	if deviceID == "" {
		return "", "", fmt.Errorf("device id is required")
	}
	if data == "" {
		return "", "", fmt.Errorf("data is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	device, err := s.repo.GetByID(deviceID)
	if err != nil {
		return "", "", ErrDeviceNotFound
	}

	var last string
	if device.SignatureCounter == 0 {
		last = base64.StdEncoding.EncodeToString([]byte(device.ID))
	} else {
		last = device.LastSignature
	}

	signedData := fmt.Sprintf("%d\n%s\n%s",
		device.SignatureCounter,
		data,
		last,
	)

	signer, err := s.factory.Create(device)
	if err != nil {
		return "", "", err
	}

	sigBytes, err := signer.Sign([]byte(signedData))
	if err != nil {
		return "", "", err
	}

	signature := base64.StdEncoding.EncodeToString(sigBytes)

	device.LastSignature = signature
	device.SignatureCounter++

	if err := s.repo.Save(device); err != nil {
		return "", "", err
	}

	return signature, signedData, nil
}
