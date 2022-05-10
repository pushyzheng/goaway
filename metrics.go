package main

import "github.com/prometheus/client_golang/prometheus"

var (
	entranceSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "entrance_duration_milliseconds",
			Help: "The summary metrics of entrance",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
	)

	routerCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "router_counter_vec",
			Help: "The counter metrics of router",
		},
		[]string{"path"},
	)

	reverseCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "reverse_counter_vec",
			Help: "The counter metrics of reverse proxy",
		},
		[]string{"appName"},
	)

	errorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_counter_vec",
			Help: "The counter metrics of error page",
		},
		[]string{"statusCode"},
	)
)

func init() {
	prometheus.MustRegister(entranceSummary)
	prometheus.MustRegister(routerCounter)
	prometheus.MustRegister(reverseCounter)
	prometheus.MustRegister(errorCounter)
}
