package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *App) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req JobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", 400)
		return
	}

	if len(req.Command) < 1 {
		http.Error(w, "command must include at least one part", 400)
		return
	}

	if req.Image == "" {
		req.Image = "alpine"
	}

	if req.K8sNamespace == "" {
		req.K8sNamespace = "default"
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name + "-" + fmt.Sprint(time.Now().Unix()),
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "runner",
							Image:   req.Image,
							Command: req.Command,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: resource.Quantity{
										Format: "1",
									},
									corev1.ResourceMemory: resource.Quantity{
										Format: "512Mi",
									},
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU: resource.Quantity{
										Format: "200m",
									},
									corev1.ResourceMemory: resource.Quantity{
										Format: "128Mi",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	createdJob, err := a.jobService.CreateJob(context.TODO(), job, req.K8sNamespace)
	if err != nil {
		http.Error(w, "failed to create job", 500)
		return
	}

	log.Println("Created job", createdJob.Name)

	w.WriteHeader(202)
	json.NewEncoder(w).Encode(map[string]string{
		"job_name": createdJob.Name,
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
