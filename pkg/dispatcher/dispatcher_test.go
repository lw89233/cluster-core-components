package dispatcher

import (
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	d := New()

	sub1 := d.Subscribe()
	sub2 := d.Subscribe()

	d.Broadcast("test-message")

	select {
	case msg := <-sub1:
		if msg != "test-message" {
			t.Errorf("Expected 'test-message', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout waiting for sub1")
	}

	select {
	case msg := <-sub2:
		if msg != "test-message" {
			t.Errorf("Expected 'test-message', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout waiting for sub2")
	}

	d.Unsubscribe(sub1)

	d.Broadcast("second-message")

	select {
	case <-sub1:
	default:
	}

	select {
	case msg := <-sub2:
		if msg != "second-message" {
			t.Errorf("Expected 'second-message', got %v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout waiting for sub2")
	}
}