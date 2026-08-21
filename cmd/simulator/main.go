package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lw89233/cluster-core-components/pkg/dispatcher"
	"github.com/lw89233/cluster-core-components/pkg/liveness"
	"github.com/lw89233/cluster-core-components/pkg/reconciler"
	"github.com/lw89233/cluster-core-components/pkg/scheduler"
	"github.com/lw89233/cluster-core-components/pkg/storage"
)

func main() {
	disp := dispatcher.New()
	sub := disp.Subscribe()

	go func() {
		for msg := range sub {
			fmt.Printf("EVENT_RECEIVED\t%s\n", msg)
		}
	}()

	tracker := liveness.New()
	now := time.Now()

	tracker.Update("node-1")
	tracker.Update("node-2")

	future := now.Add(20 * time.Second)
	activeNodes := tracker.GetActiveNodes(future)

	fmt.Printf("ACTIVE_NODES\t%v\n", activeNodes)

	storePath := "state.json"
	store := storage.NewJSONStore(storePath)

	snap, err := store.Load()
	if err != nil {
		log.Fatalf("Load err: %v", err)
	}

	previousAssignments := snap.Assignments
	if previousAssignments == nil {
		previousAssignments = make(map[string]string)
		snap.Version = 1
	}

	var schedNodes []scheduler.Node
	for _, id := range activeNodes {
		schedNodes = append(schedNodes, scheduler.Node{ID: id, CPU: 4, RAM: 8})
	}

	tasks := []scheduler.Task{
		{ID: "web-1", GroupID: "web", Stateful: false, ReqCPU: 1, ReqRAM: 2},
		{ID: "db-1", GroupID: "db", Stateful: true, ReqCPU: 2, ReqRAM: 4},
	}

	sched := scheduler.New()
	newAssignments := sched.Assign(tasks, schedNodes, previousAssignments)

	desiredState := reconciler.State{Tasks: make(map[string]bool)}
	for taskID := range newAssignments {
		desiredState.Tasks[taskID] = true
	}

	actualState := reconciler.State{Tasks: make(map[string]bool)}
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

	disp.Broadcast(fmt.Sprintf("STATE_SAVED_V%d", newSnap.Version))

	time.Sleep(50 * time.Millisecond)
	disp.Unsubscribe(sub)
}