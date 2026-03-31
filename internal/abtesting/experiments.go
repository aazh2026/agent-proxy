package abtesting

import (
	cfg "github.com/openclaw/agent-proxy/internal/config"
	"hash/fnv"
)

// PickVariant deterministically assigns a user to a variant for a given experiment.
// - userID: unique user identifier used for stickiness
// - exp: experiment configuration (name, variants, optional weights)
// Returns the selected variant name or an empty string if no variants configured.
func PickVariant(userID string, exp cfg.ABExperiment) string {
	if len(exp.Variants) == 0 {
		return ""
	}

	// Determine if weighting is specified and valid
	totalVariants := len(exp.Variants)
	hasWeights := len(exp.Weights) == totalVariants && totalVariants > 0

	// Stable hash based on userID and experiment name for determinism
	h := fnv.New32a()
	h.Write([]byte(userID))
	h.Write([]byte("|"))
	h.Write([]byte(exp.Name))
	hashVal := int(h.Sum32())

	if hasWeights {
		// Weighted distribution across variants
		bucket := hashVal % 0x7fffffff // avoid negatives, keep large space
		// compute total weight
		total := 0
		for _, w := range exp.Weights {
			if w > 0 {
				total += w
			}
		}
		if total <= 0 {
			total = totalVariants
		}
		bucket = bucket % total
		acc := 0
		for i, w := range exp.Weights {
			acc += w
			if bucket < acc {
				return exp.Variants[i]
			}
		}
		// Fallback to last variant
		return exp.Variants[len(exp.Variants)-1]
	}

	// Uniform distribution when no weights provided
	idx := hashVal % totalVariants
	return exp.Variants[idx]
}
