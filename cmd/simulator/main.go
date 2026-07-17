package main

import (
	"fmt"
	"log"

	"github.com/lw89233/cluster-core-components/pkg/scheduler"
	"github.com/lw89233/cluster-core-components/pkg/storage"
)

type ClusterState struct {
	Nodes []scheduler.Node `json:"nodes"`
}

func main() {
	store := storage.NewJSONStore()
	stateFile := "state.json"

	state := ClusterState{
		Nodes: []scheduler.Node{
			{ID: "node-worker-01", ActiveTasks: 2, MaxCapacity: 10},
			{ID: "node-worker-02", ActiveTasks: 0, MaxCapacity: 10},
			{ID: "node-worker-03", ActiveTasks: 1, MaxCapacity: 10},
		},
	}

	fmt.Println("Initial state defined. Saving to disk...")
	if err := store.AtomicWrite(stateFile, state); err != nil {
		log.Fatalf("Failed to write initial state: %v", err)
	}

	var loadedState ClusterState
	if err := store.Read(stateFile, &loadedState); err != nil {
		log.Fatalf("Failed to read state: %v", err)
	}

	lcScheduler := &scheduler.LeastContainers{}
	newTask := scheduler.Task{ID: "web-server-container"}

	fmt.Println("Attempting to schedule new task...")
	selectedNodeID, err := lcScheduler.Schedule(newTask, loadedState.Nodes)
	if err != nil {
		log.Fatalf("Scheduling failed: %v", err)
	}

	fmt.Printf("Task scheduled successfully on: %s\n", selectedNodeID)

	for i := range loadedState.Nodes {
		if loadedState.Nodes[i].ID == selectedNodeID {
			loadedState.Nodes[i].ActiveTasks++
			break
		}
	}

	fmt.Println("Updating cluster state on disk...")
	if err := store.AtomicWrite(stateFile, loadedState); err != nil {
		log.Fatalf("Failed to update state: %v", err)
	}

	fmt.Println("Simulation completed successfully.")
}