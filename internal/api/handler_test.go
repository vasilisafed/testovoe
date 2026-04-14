package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"test_task/internal/repository"
	"test_task/internal/service"
)

func newTestHandler() http.Handler {
	repo := repository.NewInMemoryDeviceRepository()
	deviceService := service.NewDeviceService(repo)
	signService := service.NewSignService(repo, &service.DefaultSignerFactory{})
	return NewHandler(deviceService, signService).Routes()
}

func TestDeviceLifecycleEndpoints(t *testing.T) {
	handler := newTestHandler()

	createBody := bytes.NewBufferString(`{"id":"device-1","algorithm":"RSA","label":"Primary"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/devices", createBody)
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	var created DeviceResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID != "device-1" || created.Algorithm != "RSA" {
		t.Fatalf("unexpected create response: %+v", created)
	}

	getRec := httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/devices/device-1", nil)
	handler.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRec.Code)
	}

	listRec := httptest.NewRecorder()
	listReq := httptest.NewRequest(http.MethodGet, "/devices", nil)
	handler.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, listRec.Code)
	}

	signRec := httptest.NewRecorder()
	signReq := httptest.NewRequest(http.MethodPost, "/devices/device-1/sign", bytes.NewBufferString(`{"data":"hello"}`))
	handler.ServeHTTP(signRec, signReq)
	if signRec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, signRec.Code, signRec.Body.String())
	}

	var signResp SignatureResponse
	if err := json.NewDecoder(signRec.Body).Decode(&signResp); err != nil {
		t.Fatalf("decode sign response: %v", err)
	}
	if signResp.Signature == "" || signResp.SignedData == "" {
		t.Fatalf("unexpected sign response: %+v", signResp)
	}
}
