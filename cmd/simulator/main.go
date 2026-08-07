package main

import (
	"fmt"
	"log"

	"github.com/lw89233/cluster-core-components/pkg/reconciler"
	"github.com/lw89233/cluster-core-components/pkg/scheduler"
	"github.com/lw89233/cluster-core-components/pkg/storage"
)

func main() {
	storePath := "state.json"
	store := storage.NewJSONStore(storePath)

	snap, err := store.Load()
	if err != nil {
		log.Fatalf("Load err: %v", err)
	}

	previousAssignments := snap.Assignments
	if previousAssignments == nil {
		previousAssignments = map[string]string{
			"db-1":  "node-2",
			"web-1": "node-1",
		}
		snap.Version = 1
	}

	nodes := []string{"node-1", "node-3"}

	tasks := []scheduler.Task{
		{ID: "db-1", Stateful: true},
		{ID: "web-1", Stateful: false},
		{ID: "web-2", Stateful: false},
	}

	sched := scheduler.New()
	newAssignments := sched.Assign(tasks, nodes, previousAssignments)

	desiredState := reconciler.State{
		Tasks: make(map[string]bool),
	}
	for taskID := range newAssignments {
		desiredState.Tasks[taskID] = true
	}

	actualState := reconciler.State{
		Tasks: make(map[string]bool),
	}
	for taskID := range previousAssignments {
		actualState.Tasks[taskID] = true
	}

	rec := reconciler.New()
	actions := rec.ComputeActions(desiredState, actualState)

	for taskID, node := range newAssignments {
		fmt.Printf("ASSIGN\t%s\t%s\n", taskID, node)
	}

	for _, action := range actions {
		fmt.Printf("%s\t%s\n", action.Type, action.TaskID)
	}

	newSnap := storage.Snapshot{
		Version:     snap.Version + 1,
		Assignments: newAssignments,
	}

	if err := store.Save(newSnap); err != nil {
		log.Fatalf("Save err: %v", err)
	}

	fmt.Printf("STATE_SAVED\tV:%d\n", newSnap.Version)
}