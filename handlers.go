package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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
							Image:   "alpine",
							Command: req.Command,
						},
					},
				},
			},
		},
	}

	createdJob, err := a.jobService.CreateJob(context.TODO(), job)
	if err != nil {
		http.Error(w, "failed to create job", 500)
		return
	}

	w.WriteHeader(202)
	json.NewEncoder(w).Encode(map[string]string{
		"job_name": createdJob.Name,
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
