package api

import "test_task/internal/domain"

type CreateDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type SignRequest struct {
	Data string `json:"data"`
}

type DeviceResponse struct {
	ID               string           `json:"id"`
	Algorithm        domain.Algorithm `json:"algorithm"`
	Label            string           `json:"label,omitempty"`
	SignatureCounter uint64           `json:"signature_counter"`
	LastSignature    string           `json:"last_signature,omitempty"`
}

type SignatureResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func newDeviceResponse(device *domain.Device) DeviceResponse {
	return DeviceResponse{
		ID:               device.ID,
		Algorithm:        device.Algorithm,
		Label:            device.Label,
		SignatureCounter: device.SignatureCounter,
		LastSignature:    device.LastSignature,
	}
}
