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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req JobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", 400)
		return
	}

	if len(req.Command) < 1 {
		http.Error(w, "command must include at least one part", 400)
		return
	}

	config, err := getKubeConfig()
	if err != nil {
		log.Println("could not get kubernetes connection", err.Error())
		http.Error(w, "could not get valid kubernetes configuration", 500)
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		http.Error(w, "failed to initialize connection to kubernetes cluster", 500)
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

	jobClient := clientset.BatchV1().Jobs("default")
	createdJob, err := jobClient.Create(context.TODO(), job, metav1.CreateOptions{})
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
