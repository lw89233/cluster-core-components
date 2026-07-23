package reconciler

import (
	"sort"
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionDelete ActionType = "DELETE"
)

type Action struct {
	Type   ActionType
	TaskID string
}

type State struct {
	Tasks map[string]bool
}

type Reconciler struct{}

func New() *Reconciler {
	return &Reconciler{}
}

func (r *Reconciler) ComputeActions(desired, actual State) []Action {
	var actions []Action

	for taskID := range desired.Tasks {
		if _, exists := actual.Tasks[taskID]; !exists {
			actions = append(actions, Action{Type: ActionCreate, TaskID: taskID})
		}
	}

	for taskID := range actual.Tasks {
		if _, exists := desired.Tasks[taskID]; !exists {
			actions = append(actions, Action{Type: ActionDelete, TaskID: taskID})
		}
	}

	sort.Slice(actions, func(i, j int) bool {
		if actions[i].Type == actions[j].Type {
			return actions[i].TaskID < actions[j].TaskID
		}
		return actions[i].Type < actions[j].Type
	})

	return actions
}