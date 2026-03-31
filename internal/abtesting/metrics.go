package abtesting

import "time"

// Metrics scaffolding for AB testing MVP.
type Metrics struct {
	LatencyMs  int
	Cost       float64
	Success    bool
	Variant    string
	Experiment string
	TS         time.Time
}

// NewMetrics creates a basic metrics instance for a given experiment/variant.
func NewMetrics(latencyMs int, cost float64, success bool, experiment, variant string) Metrics {
	return Metrics{
		LatencyMs:  latencyMs,
		Cost:       cost,
		Success:    success,
		Experiment: experiment,
		Variant:    variant,
		TS:         time.Now(),
	}
}
