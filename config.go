package main

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeConfig() (*rest.Config, error) {
	// Check if we're in cluster
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// If not, fallback to local .kube/config
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
