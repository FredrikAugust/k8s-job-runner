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
			name: "valid request with defaults",
			payload: JobRequest{
				Name:    "test-job",
				Command: []string{"echo", "hello"},
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "valid request with custom image",
			payload: JobRequest{
				Name:    "test-job",
				Command: []string{"echo", "hello"},
				Image:   "busybox:latest",
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "valid request with custom namespace",
			payload: JobRequest{
				Name:         "test-job",
				Command:      []string{"echo", "hello"},
				K8sNamespace: "custom-ns",
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name: "valid request with all fields",
			payload: JobRequest{
				Name:         "test-job",
				Command:      []string{"echo", "hello"},
				Image:        "python:3.9-slim",
				K8sNamespace: "test-namespace",
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
				CreateJobFunc: func(ctx context.Context, job *batchv1.Job, namespace string) (*batchv1.Job, error) {
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

func TestHandleCreateJobWithFields(t *testing.T) {
	tests := []struct {
		name              string
		payload           JobRequest
		expectedImage     string
		expectedNamespace string
	}{
		{
			name: "defaults are applied",
			payload: JobRequest{
				Name:    "test-job",
				Command: []string{"echo", "hello"},
			},
			expectedImage:     "alpine",
			expectedNamespace: "default",
		},
		{
			name: "custom image is used",
			payload: JobRequest{
				Name:    "test-job",
				Command: []string{"echo", "hello"},
				Image:   "python:3.9-slim",
			},
			expectedImage:     "python:3.9-slim",
			expectedNamespace: "default",
		},
		{
			name: "custom namespace is used",
			payload: JobRequest{
				Name:         "test-job",
				Command:      []string{"echo", "hello"},
				K8sNamespace: "custom-namespace",
			},
			expectedImage:     "alpine",
			expectedNamespace: "custom-namespace",
		},
		{
			name: "both custom fields are used",
			payload: JobRequest{
				Name:         "test-job",
				Command:      []string{"echo", "hello"},
				Image:        "busybox:latest",
				K8sNamespace: "test-ns",
			},
			expectedImage:     "busybox:latest",
			expectedNamespace: "test-ns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req, err := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			var capturedJob *batchv1.Job
			var capturedNamespace string
			mockJobService := &MockJobService{
				CreateJobFunc: func(ctx context.Context, job *batchv1.Job, namespace string) (*batchv1.Job, error) {
					capturedJob = job
					capturedNamespace = namespace
					return job, nil
				},
			}
			app := &App{jobService: mockJobService}

			rr := httptest.NewRecorder()
			app.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusAccepted {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusAccepted)
			}

			// Verify the job was created with correct image
			if capturedJob != nil && len(capturedJob.Spec.Template.Spec.Containers) > 0 {
				actualImage := capturedJob.Spec.Template.Spec.Containers[0].Image
				if actualImage != tt.expectedImage {
					t.Errorf("wrong image: got %v want %v", actualImage, tt.expectedImage)
				}
			}

			// Verify the correct namespace was used
			if capturedNamespace != tt.expectedNamespace {
				t.Errorf("wrong namespace: got %v want %v", capturedNamespace, tt.expectedNamespace)
			}
		})
	}
}
