package model

import "time"

type Report struct {
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Environment string    `json:"environment"`
	Checks      []Check   `json:"checks"`
}

type Check struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	ResponseTime string `json:"response_time"`
	Error        string `json:"error,omitempty"`
}
