package liveness

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestTracker(t *testing.T) {
	tr := New()
	now := time.Now()

	tr.mu.Lock()
	tr.nodes["node-1"] = &NodeInfo{ID: "node-1", LastSeen: now.Add(-5 * time.Second)}
	tr.nodes["node-2"] = &NodeInfo{ID: "node-2", LastSeen: now.Add(-20 * time.Second)}
	tr.nodes["node-3"] = &NodeInfo{ID: "node-3", LastSeen: now.Add(-40 * time.Second)}
	tr.mu.Unlock()

	if status := tr.GetStatus("node-1", now); status != StatusActive {
		t.Errorf("Expected ACTIVE, got %v", status)
	}
	if status := tr.GetStatus("node-2", now); status != StatusStale {
		t.Errorf("Expected STALE, got %v", status)
	}
	if status := tr.GetStatus("node-3", now); status != StatusDead {
		t.Errorf("Expected DEAD, got %v", status)
	}
	if status := tr.GetStatus("node-unknown", now); status != StatusDead {
		t.Errorf("Expected DEAD, got %v", status)
	}

	active := tr.GetActiveNodes(now)
	sort.Strings(active)
	expected := []string{"node-1"}

	if !reflect.DeepEqual(active, expected) {
		t.Errorf("Expected %v, got %v", expected, active)
	}

	tr.Update("node-2")
	if status := tr.GetStatus("node-2", time.Now()); status != StatusActive {
		t.Errorf("Expected ACTIVE after update, got %v", status)
	}
}