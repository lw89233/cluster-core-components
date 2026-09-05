package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lw89233/cluster-core-components/pkg/chaos"
	"github.com/lw89233/cluster-core-components/pkg/dispatcher"
	"github.com/lw89233/cluster-core-components/pkg/liveness"
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
	tracker.Update("node-1")
	tracker.Update("node-2")

	activeNodes := tracker.GetActiveNodes(time.Now())
	fmt.Printf("ACTIVE_NODES\t%v\n", activeNodes)

	storePath := "state.json"
	store := storage.NewJSONStore(storePath)

	snap, err := store.Load()
	if err != nil {
		snap.Version = 1
	}

	newSnap := storage.Snapshot{
		Version:     snap.Version + 1,
		Assignments: map[string]string{"web-1": "node-1", "db-1": "node-2"},
	}

	diskChaos := chaos.NewInjector(chaos.Config{
		ErrorProbability: 0.7,
		MaxDelay:         50 * time.Millisecond,
	})

	maxRetries := 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Printf("SAVE_ATTEMPT\t%d\n", attempt)

		err := diskChaos.Execute(func() error {
			return store.Save(newSnap)
		})

		if err == nil {
			fmt.Println("SAVE_SUCCESS\tState persisted successfully!")
			disp.Broadcast(fmt.Sprintf("STATE_SAVED_V%d", newSnap.Version))
			break
		}

		fmt.Printf("SAVE_FAILED\tError: %v. Retrying in background...\n", err)
		time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)

		if attempt == maxRetries {
			log.Fatalf("FATAL_ERROR\tCould not save state after %d attempts", maxRetries)
		}
	}

	time.Sleep(200 * time.Millisecond)
	disp.Unsubscribe(sub)
}