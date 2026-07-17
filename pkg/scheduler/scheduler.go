package scheduler

import (
	"errors"
	"sort"
)

var ErrNoNodesAvailable = errors.New("no nodes available to schedule the task")

type Node struct {
	ID          string
	ActiveTasks int
	MaxCapacity int
}

type Task struct {
	ID string
}

type Scheduler interface {
	Schedule(task Task, nodes []Node) (string, error)
}

type LeastContainers struct{}

func (lc *LeastContainers) Schedule(task Task, nodes []Node) (string, error) {
	if len(nodes) == 0 {
		return "", ErrNoNodesAvailable
	}

	var availableNodes []Node
	for _, n := range nodes {
		if n.ActiveTasks < n.MaxCapacity {
			availableNodes = append(availableNodes, n)
		}
	}

	if len(availableNodes) == 0 {
		return "", ErrNoNodesAvailable
	}

	sort.Slice(availableNodes, func(i, j int) bool {
		if availableNodes[i].ActiveTasks == availableNodes[j].ActiveTasks {
			return availableNodes[i].ID < availableNodes[j].ID
		}
		return availableNodes[i].ActiveTasks < availableNodes[j].ActiveTasks
	})

	return availableNodes[0].ID, nil
}