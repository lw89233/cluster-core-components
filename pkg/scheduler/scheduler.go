package scheduler

import (
	"sort"
	"sync"
)

type Task struct {
	ID       string
	GroupID  string
	Stateful bool
	ReqCPU   int
	ReqRAM   int
}

type Node struct {
	ID  string
	CPU int
	RAM int
}

type Scheduler struct {
	mu sync.RWMutex
}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Assign(tasks []Task, nodes []Node, previous map[string]string) map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	assignments := make(map[string]string)
	usedCPU := make(map[string]int)
	usedRAM := make(map[string]int)
	groupCounts := make(map[string]map[string]int)
	activeNodes := make(map[string]bool)
	nodeMap := make(map[string]Node)

	sortedNodeIDs := make([]string, 0, len(nodes))
	for _, n := range nodes {
		activeNodes[n.ID] = true
		nodeMap[n.ID] = n
		groupCounts[n.ID] = make(map[string]int)
		sortedNodeIDs = append(sortedNodeIDs, n.ID)
	}
	sort.Strings(sortedNodeIDs)

	var statelessTasks []Task

	for _, t := range tasks {
		if t.Stateful {
			assigned := false
			if prevNode, exists := previous[t.ID]; exists && activeNodes[prevNode] {
				assignments[t.ID] = prevNode
				usedCPU[prevNode] += t.ReqCPU
				usedRAM[prevNode] += t.ReqRAM
				groupCounts[prevNode][t.GroupID]++
				assigned = true
			}
			if !assigned && len(sortedNodeIDs) > 0 {
				assignments[t.ID] = sortedNodeIDs[0]
				usedCPU[sortedNodeIDs[0]] += t.ReqCPU
				usedRAM[sortedNodeIDs[0]] += t.ReqRAM
				groupCounts[sortedNodeIDs[0]][t.GroupID]++
			}
		} else {
			statelessTasks = append(statelessTasks, t)
		}
	}

	for _, t := range statelessTasks {
		if len(sortedNodeIDs) == 0 {
			break
		}

		var selectedNode string
		first := true
		bestScore := 0

		for _, nodeID := range sortedNodeIDs {
			freeCPU := nodeMap[nodeID].CPU - usedCPU[nodeID]
			freeRAM := nodeMap[nodeID].RAM - usedRAM[nodeID]

			if freeCPU < t.ReqCPU || freeRAM < t.ReqRAM {
				continue
			}

			score := freeCPU + freeRAM

			if groupCounts[nodeID][t.GroupID] > 0 {
				score -= 1000000
			}

			if first || score > bestScore {
				bestScore = score
				selectedNode = nodeID
				first = false
			}
		}

		if selectedNode != "" {
			assignments[t.ID] = selectedNode
			usedCPU[selectedNode] += t.ReqCPU
			usedRAM[selectedNode] += t.ReqRAM
			groupCounts[selectedNode][t.GroupID]++
		}
	}

	return assignments
}