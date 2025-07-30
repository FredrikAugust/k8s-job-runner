package main

import (
	"context"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var runningGauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "running_jobs",
})

func createInformer(clientset *kubernetes.Clientset) {
	informer := informers.NewSharedInformerFactory(clientset, time.Second*30).Batch().V1().Jobs().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			job := obj.(*batchv1.Job)

			if job.Labels["app.kubernetes.io/managed-by"] != "k8s-job-runner" {
				return
			}

			if job.Status.CompletionTime != nil {
				return
			}

			runningGauge.Inc()
			log.Println("Job added", job.Name)
		},
		UpdateFunc: func(oldObj, newObj any) {
			newJob := newObj.(*batchv1.Job)
			oldJob := oldObj.(*batchv1.Job)

			if newJob == nil || oldJob == nil {
				log.Println("i thought i saw something fishy")
			}

			if newJob.Labels["app.kubernetes.io/managed-by"] != "k8s-job-runner" {
				return
			}

			if oldJob.Status.CompletionTime == nil && newJob.Status.CompletionTime != nil {
				runningGauge.Dec()
				log.Println("Job completed", newJob.Name)
				return
			}
		},
	})

	go informer.RunWithContext(context.TODO())

	log.Println("Informer started")
}
