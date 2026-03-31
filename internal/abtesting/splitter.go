package abtesting

import (
	cfg "github.com/openclaw/agent-proxy/internal/config"
)

// VariantForUser is a thin wrapper to PickVariant for external callers.
func VariantForUser(userID string, exp cfg.ABExperiment) string {
	return PickVariant(userID, exp)
}
