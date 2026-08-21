package liveness

import (
	"sync"
	"time"
)

type Status string

const (
	StatusActive Status = "ACTIVE"
	StatusStale  Status = "STALE"
	StatusDead   Status = "DEAD"

	ActiveThreshold = 15 * time.Second
	DeadThreshold   = 30 * time.Second
)

type NodeInfo struct {
	ID       string
	LastSeen time.Time
}

type Tracker struct {
	mu    sync.RWMutex
	nodes map[string]*NodeInfo
}

func New() *Tracker {
	return &Tracker{
		nodes: make(map[string]*NodeInfo),
	}
}

func (t *Tracker) Update(id string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.nodes[id] = &NodeInfo{
		ID:       id,
		LastSeen: time.Now(),
	}
}

func (t *Tracker) GetStatus(id string, now time.Time) Status {
	t.mu.RLock()
	defer t.mu.RUnlock()

	info, exists := t.nodes[id]
	if !exists {
		return StatusDead
	}

	diff := now.Sub(info.LastSeen)
	if diff <= ActiveThreshold {
		return StatusActive
	}
	if diff <= DeadThreshold {
		return StatusStale
	}

	return StatusDead
}

func (t *Tracker) GetActiveNodes(now time.Time) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var active []string
	for id, info := range t.nodes {
		if now.Sub(info.LastSeen) <= ActiveThreshold {
			active = append(active, id)
		}
	}
	return active
}