package service

import (
	"errors"
	"test_task/internal/crypto"
	"test_task/internal/domain"
)

type DefaultSignerFactory struct{}

func (f *DefaultSignerFactory) Create(d *domain.Device) (domain.Signer, error) {
	switch d.Algorithm {
	case domain.RSA:
		return crypto.NewRSASigner(d.PrivateKey)
	case domain.ECC:
		return crypto.NewECDSASigner(d.PrivateKey)
	default:
		return nil, errors.New("unsupported algorithm")
	}
}
