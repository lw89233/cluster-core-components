package scheduler

import (
	"sort"
	"sync"
)

type Task struct {
	ID       string
	Stateful bool
}

type Scheduler struct {
	mu sync.RWMutex
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Assign(tasks []Task, nodes []string, previous map[string]string) map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	assignments := make(map[string]string)
	nodeLoads := make(map[string]int)
	activeNodes := make(map[string]bool)

	for _, n := range nodes {
		activeNodes[n] = true
		nodeLoads[n] = 0
	}

	sortedNodes := make([]string, len(nodes))
	copy(sortedNodes, nodes)
	sort.Strings(sortedNodes)

	var statelessTasks []Task

	for _, t := range tasks {
		if t.Stateful {
			assigned := false
			if prevNode, exists := previous[t.ID]; exists && activeNodes[prevNode] {
				assignments[t.ID] = prevNode
				nodeLoads[prevNode]++
				assigned = true
			}
			if !assigned && len(sortedNodes) > 0 {
				assignments[t.ID] = sortedNodes[0]
				nodeLoads[sortedNodes[0]]++
			}
		} else {
			statelessTasks = append(statelessTasks, t)
		}
	}

	for _, t := range statelessTasks {
		if len(sortedNodes) == 0 {
			break
		}

		var selectedNode string
		minLoad := -1

		for _, n := range sortedNodes {
			load := nodeLoads[n]
			if minLoad == -1 || load < minLoad {
				minLoad = load
				selectedNode = n
			}
		}

		assignments[t.ID] = selectedNode
		nodeLoads[selectedNode]++
	}

	return assignments
}