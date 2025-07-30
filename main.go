package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
)

type App struct {
	jobService JobService
}

var requestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"path", "method"},
)

func init() {
	log.Println("Initializing metrics")

	prometheus.MustRegister(
		requestCounter,
		runningGauge,
	)
}

func withMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		log.Println(method, path)

		requestCounter.WithLabelValues(path, method).Inc()
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting server")

	log.Println("Getting kubernetes configuration")
	config, err := getKubeConfig()
	if err != nil {
		log.Println("could not get kubernetes configuration", err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("could not get kubernetes connection", err.Error())
		return
	}

	log.Println("Creating informer")
	createInformer(clientset)

	app := &App{
		jobService: NewK8sClient(clientset, "default"),
	}

	http.Handle("/health", withMetrics(http.HandlerFunc(handleHealth)))
	http.Handle("/jobs", withMetrics(http.HandlerFunc(app.handleCreateJob)))

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Ready to serve on :8084")
	if err := http.ListenAndServe("0.0.0.0:8084", nil); err != nil {
		log.Fatal(err)
	}
}
