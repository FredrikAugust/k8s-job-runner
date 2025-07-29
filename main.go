package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
)

func init() {
	prometheus.MustRegister(requestCounter)
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		requestCounter.Inc()
		fmt.Fprintln(w, "ok")
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Listening on :8084")
	if err := http.ListenAndServe("0.0.0.0:8084", nil); err != nil {
		log.Fatal(err)
	}
}
