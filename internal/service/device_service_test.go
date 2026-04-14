package service

import (
	"sync"
	"testing"

	"test_task/internal/domain"
	"test_task/internal/repository"
)

func newTestServices(t *testing.T) (*DeviceService, *SignService) {
	t.Helper()

	repo := repository.NewInMemoryDeviceRepository()
	deviceService := NewDeviceService(repo)
	signService := NewSignService(repo, &DefaultSignerFactory{})
	return deviceService, signService
}

func TestCreateDeviceAndGetDevice(t *testing.T) {
	deviceService, _ := newTestServices(t)

	device, err := deviceService.CreateDevice("device-1", domain.RSA, "Primary")
	if err != nil {
		t.Fatalf("CreateDevice() error = %v", err)
	}

	if device.SignatureCounter != 0 {
		t.Fatalf("expected signature counter to start at 0, got %d", device.SignatureCounter)
	}

	got, err := deviceService.GetDevice("device-1")
	if err != nil {
		t.Fatalf("GetDevice() error = %v", err)
	}

	if got.ID != "device-1" || got.Label != "Primary" || got.Algorithm != domain.RSA {
		t.Fatalf("unexpected device returned: %+v", got)
	}
}

func TestCreateDeviceRejectsDuplicateID(t *testing.T) {
	deviceService, _ := newTestServices(t)

	if _, err := deviceService.CreateDevice("device-1", domain.RSA, "One"); err != nil {
		t.Fatalf("first CreateDevice() error = %v", err)
	}

	if _, err := deviceService.CreateDevice("device-1", domain.ECC, "Two"); err != ErrDeviceAlreadyExists {
		t.Fatalf("expected ErrDeviceAlreadyExists, got %v", err)
	}
}

func TestSignUsesBaseCaseAndIncrementsCounter(t *testing.T) {
	deviceService, signService := newTestServices(t)

	if _, err := deviceService.CreateDevice("device-1", domain.RSA, "Signer"); err != nil {
		t.Fatalf("CreateDevice() error = %v", err)
	}

	signature, signedData, err := signService.Sign("device-1", "payload")
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	expectedSignedData := "0\npayload\nZGV2aWNlLTE="
	if signedData != expectedSignedData {
		t.Fatalf("expected signed data %q, got %q", expectedSignedData, signedData)
	}
	if signature == "" {
		t.Fatal("expected non-empty signature")
	}

	device, err := deviceService.GetDevice("device-1")
	if err != nil {
		t.Fatalf("GetDevice() error = %v", err)
	}

	if device.SignatureCounter != 1 {
		t.Fatalf("expected signature counter to be 1, got %d", device.SignatureCounter)
	}
	if device.LastSignature != signature {
		t.Fatalf("expected last signature to be updated")
	}
}

func TestSignIsSafeForConcurrentClients(t *testing.T) {
	deviceService, signService := newTestServices(t)

	if _, err := deviceService.CreateDevice("device-1", domain.ECC, "Concurrent"); err != nil {
		t.Fatalf("CreateDevice() error = %v", err)
	}

	const workers = 25
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, _, err := signService.Sign("device-1", "payload"); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("concurrent Sign() error = %v", err)
		}
	}

	device, err := deviceService.GetDevice("device-1")
	if err != nil {
		t.Fatalf("GetDevice() error = %v", err)
	}

	if device.SignatureCounter != workers {
		t.Fatalf("expected signature counter %d, got %d", workers, device.SignatureCounter)
	}
}
