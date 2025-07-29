package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
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

func (a *App) ServeHTTP(rr http.ResponseWriter, req *http.Request) {
	handler := http.HandlerFunc(a.handleCreateJob)
	handler.ServeHTTP(rr, req)
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

			mockJobService := &MockJobService{
				CreateJobFunc: func(ctx context.Context, job *batchv1.Job) (*batchv1.Job, error) {
					return job, nil
				},
			}
			app := &App{jobService: mockJobService}

			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}