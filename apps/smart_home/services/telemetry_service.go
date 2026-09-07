package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelemetryService struct {
	BaseURL    string
	HTTPClient *http.Client
}

type TelemetryIngest struct {
	RecordedAt time.Time `json:"recorded_at"`
	Unit       string    `json:"unit"`
	Value      float64   `json:"value"`
	DeviceID   string    `json:"device_id"`
}

func NewTelemetryService(baseURL string) *TelemetryService {
	return &TelemetryService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TelemetryService) Ingest(ctx context.Context, in TelemetryIngest) error {
	if s.BaseURL == "" || in.DeviceID == "" {
		return nil
	}

	body, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("encode reading: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.BaseURL+"/api/v1/telemetry", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("telemetry request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("telemetry ingest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telemetry ingest status: %d", resp.StatusCode)
	}

	return nil
}
