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
		nodes    []string
		previous map[string]string
		expected map[string]string
	}{
		{
			name: "Stateless tasks least containers",
			tasks: []Task{
				{ID: "task-1", Stateful: false},
				{ID: "task-2", Stateful: false},
				{ID: "task-3", Stateful: false},
			},
			nodes:    []string{"node-2", "node-1"},
			previous: nil,
			expected: map[string]string{
				"task-1": "node-1",
				"task-2": "node-2",
				"task-3": "node-1",
			},
		},
		{
			name: "Stateful task pinned",
			tasks: []Task{
				{ID: "db-1", Stateful: true},
			},
			nodes: []string{"node-1", "node-2", "node-3"},
			previous: map[string]string{
				"db-1": "node-2",
			},
			expected: map[string]string{
				"db-1": "node-2",
			},
		},
		{
			name: "Stateful task repinned when node dies",
			tasks: []Task{
				{ID: "db-1", Stateful: true},
			},
			nodes: []string{"node-1", "node-3"},
			previous: map[string]string{
				"db-1": "node-2",
			},
			expected: map[string]string{
				"db-1": "node-1",
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