package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"path", "method"},
)

func init() {
	prometheus.MustRegister(requestCounter)
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
	http.Handle("/health", withMetrics(http.HandlerFunc(handleHealth)))
	http.Handle("/jobs", withMetrics(http.HandlerFunc(handleCreateJob)))

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on :8084")
	if err := http.ListenAndServe("0.0.0.0:8084", nil); err != nil {
		log.Fatal(err)
	}
}
