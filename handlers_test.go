package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleHealth)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandleCreateJob(t *testing.T) {
	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			payload: JobRequest{
				Name:    "test-job",
				Command: []string{"echo", "hello"},
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "invalid json",
			payload:        "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			payload:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty command",
			payload: JobRequest{
				Name:    "test-job",
				Command: []string{},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.payload != nil {
				if str, ok := tt.payload.(string); ok {
					body = []byte(str)
				} else {
					body, _ = json.Marshal(tt.payload)
				}
			}

			req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handleCreateJob)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}