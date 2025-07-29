package main

type JobRequest struct {
	Name    string   `json:"name"`
	Command []string `json:"command"`
}
