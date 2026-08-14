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

	nodes := []scheduler.Node{
		{ID: "node-1", CPU: 4, RAM: 8},
		{ID: "node-3", CPU: 8, RAM: 16},
	}

	tasks := []scheduler.Task{
		{ID: "db-1", GroupID: "db", Stateful: true, ReqCPU: 2, ReqRAM: 4},
		{ID: "web-1", GroupID: "web", Stateful: false, ReqCPU: 1, ReqRAM: 2},
		{ID: "web-2", GroupID: "web", Stateful: false, ReqCPU: 1, ReqRAM: 2},
		{ID: "web-3", GroupID: "web", Stateful: false, ReqCPU: 1, ReqRAM: 2},
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

	for taskID, nodeID := range newAssignments {
		fmt.Printf("ASSIGN\t%s\t%s\n", taskID, nodeID)
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