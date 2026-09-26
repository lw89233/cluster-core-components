package graceful

import (
	"testing"
	"time"
)

func TestShutdownManager(t *testing.T) {
	m := NewManager()

	workerFinished := false

	m.AddWorker()
	go func() {
		defer m.Done()
		<-m.Context().Done()
		time.Sleep(50 * time.Millisecond)
		workerFinished = true
	}()

	m.Stop()
	m.Wait()

	if !workerFinished {
		t.Errorf("Manager exited before worker could finish its task")
	}
}