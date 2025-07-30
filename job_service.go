package main

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type JobService interface {
	CreateJob(ctx context.Context, job *batchv1.Job, namespace string) (*batchv1.Job, error)
}

type K8sClient struct {
	client    *kubernetes.Clientset
	namespace string
}

func NewK8sClient(client *kubernetes.Clientset, namespace string) *K8sClient {
	return &K8sClient{
		client: client, namespace: namespace,
	}
}

func (k *K8sClient) CreateJob(ctx context.Context, job *batchv1.Job, namespace string) (*batchv1.Job, error) {
	return k.client.BatchV1().Jobs(namespace).Create(ctx, job, metav1.CreateOptions{})
}
