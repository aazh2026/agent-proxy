package abtesting

import (
	cfg "github.com/openclaw/agent-proxy/internal/config"
	"testing"
)

func TestPickVariant_DeterministicSticky(t *testing.T) {
	exp := cfg.ABExperiment{
		Name:     "test-sticky",
		Variants: []string{"A", "B"},
		Sticky:   true,
	}

	v1 := PickVariant("user-123", exp)
	v2 := PickVariant("user-123", exp)
	if v1 != v2 {
		t.Fatalf("expected sticky variant to be same across calls, got %s and %s", v1, v2)
	}
}

func TestPickVariant_Bounds(t *testing.T) {
	exp := cfg.ABExperiment{
		Name:     "test-bound",
		Variants: []string{"A", "B"},
	}
	// Basic call should return either A or B
	v := PickVariant("user-xyz", exp)
	if v != "A" && v != "B" {
		t.Fatalf("unexpected variant: %s", v)
	}
}
