package models

import "time"

// BaseLog BaseLog
type BaseLog struct {
	Log  string    `json:"log"`
	Time time.Time `json:"time"`
}

// Kubernetes Kubernetes
type Kubernetes struct {
	PodName        string `json:"pod_name"`
	ContainerImage string `json:"container_image"`
	ContainerName  string `json:"container_name"`
	Labels         Labels `json:"labels"`
}

// Labels Labels
type Labels struct {
	Task         string `json:"tekton.dev/task"`
	PipelineTask string `json:"tekton.dev/pipelineTask"`
}

// LogVO LogVO
type LogVO struct {
	Run       string `json:"run"`
	Step      string `json:"step"`
	Timestamp int64  `json:"timestamp"`
	Log       string `json:"log"`
	PodName   string `json:"podName"`
}
