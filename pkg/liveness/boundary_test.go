package liveness

import (
	"testing"
	"time"
)

func TestTracker_BoundaryDelaysAndClockSkew(t *testing.T) {
	tracker := New()

	tracker.Update("node-edge")

	activeNow := tracker.GetActiveNodes(time.Now())
	if len(activeNow) != 1 {
		t.Errorf("Expected exactly 1 node to be active immediately, got %d", len(activeNow))
	}

	futureTime := time.Now().Add(1 * time.Hour)
	activeFuture := tracker.GetActiveNodes(futureTime)

	if len(activeFuture) != 0 {
		t.Errorf("Expected node to be strictly marked as dead after extreme 1h delay, got %v", activeFuture)
	}

	pastTime := time.Now().Add(-1 * time.Hour)
	activePast := tracker.GetActiveNodes(pastTime)

	if len(activePast) != 1 {
		t.Errorf("Clock skew (past): expected node updated just now to still be considered active, got %d", len(activePast))
	}
}