package scheduler

import "testing"

func TestLeastContainers_Schedule(t *testing.T) {
	scheduler := &LeastContainers{}
	task := Task{ID: "task-1"}

	nodes := []Node{
		{ID: "node-1", ActiveTasks: 5, MaxCapacity: 10},
		{ID: "node-2", ActiveTasks: 2, MaxCapacity: 10},
		{ID: "node-3", ActiveTasks: 2, MaxCapacity: 10},
	}

	selectedNode, err := scheduler.Schedule(task, nodes)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if selectedNode != "node-2" {
		t.Errorf("expected node-2, got %s", selectedNode)
	}
}