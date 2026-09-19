package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DeviceService struct {
	BaseURL    string
	HTTPClient *http.Client
}

type Device struct {
	TypeCode     string `json:"type_code"`
	SerialNumber string `json:"serial_number"`
	Address      string `json:"address"`
	Status       string `json:"status"`
	ExternalID   string `json:"external_id"`
	ID           string `json:"id"`
	TypeID       string `json:"type_id"`
	HouseID      string `json:"house_id"`
}

type DeviceUpsert struct {
	SerialNumber string `json:"serial_number"`
	Address      string `json:"address"`
	TypeCode     string `json:"type_code"`
	Status       string `json:"status"`
	ExternalID   string `json:"external_id"`
}

func NewDeviceService(baseURL string) *DeviceService {
	return &DeviceService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *DeviceService) Upsert(ctx context.Context, in DeviceUpsert) (*Device, error) {
	if s.BaseURL == "" {
		return nil, nil
	}

	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("encode device: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.BaseURL+"/api/v1/devices", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("device request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device upsert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("device upsert status: %d", resp.StatusCode)
	}

	var device Device
	if err = json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, fmt.Errorf("decode device: %w", err)
	}

	return &device, nil
}
