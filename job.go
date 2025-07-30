package main

type JobRequest struct {
	Name         string   `json:"name"`
	Command      []string `json:"command"`
	Image        string   `json:"image"`
	K8sNamespace string   `json:"k8s_namespace"`
}
