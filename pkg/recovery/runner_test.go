package recovery

import (
	"sync"
	"testing"
	"time"
)

func TestSafeGo_RecoversFromPanic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	SafeGo(func() {
		defer wg.Done()
		panic("simulated fatal memory error")
	})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout: SafeGo did not recover properly, WaitGroup never finished")
	}
}

func TestSafeGo_NormalExecution(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	executed := false
	SafeGo(func() {
		defer wg.Done()
		executed = true
	})

	wg.Wait()

	if !executed {
		t.Errorf("Expected function to execute normally")
	}
}