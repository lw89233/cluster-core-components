package main

import (
	"fmt"

	"github.com/lw89233/cluster-core-components/pkg/reconciler"
	"github.com/lw89233/cluster-core-components/pkg/scheduler"
)

func main() {
	nodes := []string{"node-1", "node-2", "node-3"}

	tasks := []scheduler.Task{
		{ID: "web-1", Stateful: false},
		{ID: "web-2", Stateful: false},
		{ID: "db-1", Stateful: true},
	}

	previousAssignments := map[string]string{
		"db-1": "node-2",
	}

	sched := scheduler.New()
	assignments := sched.Assign(tasks, nodes, previousAssignments)

	desiredState := reconciler.State{
		Tasks: make(map[string]bool),
	}
	for taskID := range assignments {
		desiredState.Tasks[taskID] = true
	}

	actualState := reconciler.State{
		Tasks: map[string]bool{
			"web-1":  true,
			"db-old": true,
		},
	}

	rec := reconciler.New()
	actions := rec.ComputeActions(desiredState, actualState)

	for taskID, node := range assignments {
		fmt.Printf("ASSIGN\t%s\t%s\n", taskID, node)
	}

	for _, action := range actions {
		fmt.Printf("%s\t%s\n", action.Type, action.TaskID)
	}
}