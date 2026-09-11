package scheduler

import (
	"fmt"
	"sync"
	"testing"
)

func TestScheduler_Concurrency(t *testing.T) {
	s := New()
	var wg sync.WaitGroup

	nodes := make([]Node, 50)
	for i := 0; i < 50; i++ {
		nodes[i] = Node{
			ID:  fmt.Sprintf("node-%d", i),
			CPU: 32,
			RAM: 128,
		}
	}

	numRoutines := 100
	wg.Add(numRoutines)

	for i := 0; i < numRoutines; i++ {
		go func(routineID int) {
			defer wg.Done()

			tasks := make([]Task, 100)
			for j := 0; j < 100; j++ {
				tasks[j] = Task{
					ID:       fmt.Sprintf("task-%d-%d", routineID, j),
					GroupID:  "concurrent-app",
					Stateful: false,
					ReqCPU:   1,
					ReqRAM:   2,
				}
			}

			previous := make(map[string]string)
			_ = s.Assign(tasks, nodes, previous)
		}(i)
	}

	wg.Wait()
}