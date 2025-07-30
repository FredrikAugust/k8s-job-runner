package main

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
)

type MockJobService struct {
	CreateJobFunc func(ctx context.Context, job *batchv1.Job, namespace string) (*batchv1.Job, error)
}

func (m *MockJobService) CreateJob(ctx context.Context, job *batchv1.Job, namespace string) (*batchv1.Job, error) {
	if m.CreateJobFunc != nil {
		return m.CreateJobFunc(ctx, job, namespace)
	}
	return job, nil
}