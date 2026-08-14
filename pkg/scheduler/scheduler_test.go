package scheduler

import (
	"reflect"
	"testing"
)

func TestAssign(t *testing.T) {
	s := New()

	tests := []struct {
		name     string
		tasks    []Task
		nodes    []Node
		previous map[string]string
		expected map[string]string
	}{
		{
			name: "Resource limits filtering",
			tasks: []Task{
				{ID: "heavy", GroupID: "g1", Stateful: false, ReqCPU: 8, ReqRAM: 16},
				{ID: "light", GroupID: "g2", Stateful: false, ReqCPU: 1, ReqRAM: 1},
			},
			nodes: []Node{
				{ID: "node-1", CPU: 2, RAM: 4},
				{ID: "node-2", CPU: 16, RAM: 32},
			},
			previous: nil,
			expected: map[string]string{
				"heavy": "node-2",
				"light": "node-2",
			},
		},
		{
			name: "Anti-affinity spreading",
			tasks: []Task{
				{ID: "t-1", GroupID: "web", Stateful: false, ReqCPU: 1, ReqRAM: 1},
				{ID: "t-2", GroupID: "web", Stateful: false, ReqCPU: 1, ReqRAM: 1},
			},
			nodes: []Node{
				{ID: "node-1", CPU: 4, RAM: 8},
				{ID: "node-2", CPU: 4, RAM: 8},
			},
			previous: nil,
			expected: map[string]string{
				"t-1": "node-1",
				"t-2": "node-2",
			},
		},
		{
			name: "Stateful task pinned and affects scoring",
			tasks: []Task{
				{ID: "db-1", GroupID: "db", Stateful: true, ReqCPU: 2, ReqRAM: 4},
				{ID: "web-1", GroupID: "web", Stateful: false, ReqCPU: 2, ReqRAM: 4},
			},
			nodes: []Node{
				{ID: "node-1", CPU: 4, RAM: 8},
				{ID: "node-2", CPU: 4, RAM: 8},
			},
			previous: map[string]string{
				"db-1": "node-2",
			},
			expected: map[string]string{
				"db-1":  "node-2",
				"web-1": "node-1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.Assign(tt.tasks, tt.nodes, tt.previous)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Assign() = %v, expected %v", result, tt.expected)
			}
		})
	}
}