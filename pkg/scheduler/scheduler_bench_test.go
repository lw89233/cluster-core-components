package scheduler

import (
	"fmt"
	"testing"
)

func BenchmarkAssign(b *testing.B) {
	s := New()

	nodes := make([]Node, 100)
	for i := 0; i < 100; i++ {
		nodes[i] = Node{
			ID:  fmt.Sprintf("node-%d", i),
			CPU: 64,
			RAM: 256,
		}
	}

	tasks := make([]Task, 5000)
	for i := 0; i < 5000; i++ {
		tasks[i] = Task{
			ID:       fmt.Sprintf("task-%d", i),
			GroupID:  fmt.Sprintf("app-%d", i%50),
			Stateful: i%10 == 0,
			ReqCPU:   1,
			ReqRAM:   2,
		}
	}

	previous := make(map[string]string)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Assign(tasks, nodes, previous)
	}
}