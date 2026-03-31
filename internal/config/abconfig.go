package config

// ABConfig contains a collection of ABExperiment definitions.
type ABConfig struct {
	Experiments []ABExperiment
}

// ABExperiment defines a simple, config-driven AB testing experiment.
// Variants represent the variant names (e.g., ["A", "B"]).
// Weights, if provided, must align with Variants length and sum to distribution space.
type ABExperiment struct {
	Name     string
	Variants []string
	Weights  []int
	Sticky   bool
}
