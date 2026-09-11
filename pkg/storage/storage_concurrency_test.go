package storage

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestStorage_Concurrency(t *testing.T) {
	testFile := "test_concurrent_state.json"
	defer os.Remove(testFile)

	store := NewJSONStore(testFile)

	initialSnap := Snapshot{
		Version:     1,
		Assignments: map[string]string{"init": "node-0"},
	}
	_ = store.Save(initialSnap)

	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 50

	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go func() {
			defer wg.Done()
			_, _ = store.Load()
		}()
	}

	wg.Add(numWriters)
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			snap := Snapshot{
				Version:     int64(id + 2),
				Assignments: map[string]string{fmt.Sprintf("task-%d", id): "node-1"},
			}
			_ = store.Save(snap)
		}(i)
	}

	wg.Wait()

	finalSnap, err := store.Load()
	if err != nil {
		t.Errorf("Failed to load final state: %v", err)
	}
	if finalSnap.Version < 1 {
		t.Errorf("Invalid final version: %d", finalSnap.Version)
	}
}