package chaos

import (
	"testing"
	"time"
)

func TestInjector_Delay(t *testing.T) {
	inj := NewInjector(Config{MaxDelay: 50 * time.Millisecond})

	start := time.Now()
	err := inj.Execute(func() error { return nil })
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if duration > 100*time.Millisecond {
		t.Errorf("Execution took too long: %v", duration)
	}
}

func TestInjector_Error(t *testing.T) {
	inj := NewInjector(Config{ErrorProbability: 1.0})

	err := inj.Execute(func() error { return nil })

	if err != ErrSimulated {
		t.Errorf("Expected ErrSimulated, got %v", err)
	}
}

func TestInjector_Success(t *testing.T) {
	inj := NewInjector(Config{ErrorProbability: 0.0, MaxDelay: 0})

	executed := false
	err := inj.Execute(func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !executed {
		t.Errorf("Expected function to be executed")
	}
}