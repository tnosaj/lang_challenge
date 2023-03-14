package domain

import "github.com/prometheus/client_golang/prometheus"

type Order struct {
	ID     string `json:"ID" validate:"required"`
	Status string `json:"Status" validate:"required"`
}

type OrderMetrics struct {
	RedisLatency *prometheus.HistogramVec
	RedisErrors  *prometheus.CounterVec
}
