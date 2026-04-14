package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"test_task/internal/domain"
	"test_task/internal/service"
)

type Handler struct {
	deviceService *service.DeviceService
	signService   *service.SignService
}

func NewHandler(deviceService *service.DeviceService, signService *service.SignService) *Handler {
	return &Handler{
		deviceService: deviceService,
		signService:   signService,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/devices", h.handleDevices)
	mux.HandleFunc("/devices/", h.handleDeviceByPath)
	return mux
}

func (h *Handler) handleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateDevice(w, r)
	case http.MethodGet:
		h.ListDevices(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleDeviceByPath(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/devices/")
	path = strings.Trim(path, "/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	parts := strings.Split(path, "/")
	switch {
	case len(parts) == 1 && r.Method == http.MethodGet:
		h.GetDevice(w, r, parts[0])
	case len(parts) == 2 && parts[1] == "sign" && r.Method == http.MethodPost:
		h.Sign(w, r, parts[0])
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	device, err := h.deviceService.CreateDevice(req.ID, domain.Algorithm(strings.ToUpper(req.Algorithm)), req.Label)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrDeviceAlreadyExists) {
			status = http.StatusConflict
		}
		writeError(w, status, err)
		return
	}

	writeJSON(w, http.StatusCreated, newDeviceResponse(device))
}

func (h *Handler) ListDevices(w http.ResponseWriter, _ *http.Request) {
	devices, err := h.deviceService.ListDevices()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	response := make([]DeviceResponse, 0, len(devices))
	for _, device := range devices {
		response = append(response, newDeviceResponse(device))
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *Handler) GetDevice(w http.ResponseWriter, _ *http.Request, deviceID string) {
	device, err := h.deviceService.GetDevice(deviceID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrDeviceNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}

	writeJSON(w, http.StatusOK, newDeviceResponse(device))
}

func (h *Handler) Sign(w http.ResponseWriter, r *http.Request, deviceID string) {
	var req SignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	signature, signedData, err := h.signService.Sign(deviceID, req.Data)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrDeviceNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}

	writeJSON(w, http.StatusOK, SignatureResponse{
		Signature:  signature,
		SignedData: signedData,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, ErrorResponse{Error: err.Error()})
}
