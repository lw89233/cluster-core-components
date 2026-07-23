package reconciler

import (
	"reflect"
	"testing"
)

func TestComputeActions(t *testing.T) {
	r := New()

	tests := []struct {
		name     string
		desired  State
		actual   State
		expected []Action
	}{
		{
			name: "Perfect synchronization",
			desired: State{Tasks: map[string]bool{"task-1": true, "task-2": true}},
			actual:  State{Tasks: map[string]bool{"task-1": true, "task-2": true}},
			expected: nil,
		},
		{
			name: "State drift - missing task",
			desired: State{Tasks: map[string]bool{"task-1": true, "task-2": true}},
			actual:  State{Tasks: map[string]bool{"task-1": true}},
			expected: []Action{
				{Type: ActionCreate, TaskID: "task-2"},
			},
		},
		{
			name: "State drift - orphaned task",
			desired: State{Tasks: map[string]bool{"task-1": true}},
			actual:  State{Tasks: map[string]bool{"task-1": true, "task-3": true}},
			expected: []Action{
				{Type: ActionDelete, TaskID: "task-3"},
			},
		},
		{
			name: "Complex drift - both missing and orphaned",
			desired: State{Tasks: map[string]bool{"task-new": true}},
			actual:  State{Tasks: map[string]bool{"task-old": true}},
			expected: []Action{
				{Type: ActionCreate, TaskID: "task-new"},
				{Type: ActionDelete, TaskID: "task-old"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.ComputeActions(tt.desired, tt.actual)

			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ComputeActions() = %v, expected %v", result, tt.expected)
			}
		})
	}
}