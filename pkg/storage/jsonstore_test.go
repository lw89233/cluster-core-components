package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestJSONStore_SaveAndLoad(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "store-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "state.json")
	store := NewJSONStore(storePath)

	snap1 := Snapshot{
		Version: 1,
		Assignments: map[string]string{
			"task-1": "node-1",
		},
	}

	if err := store.Save(snap1); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !reflect.DeepEqual(loaded, snap1) {
		t.Errorf("Loaded %v, expected %v", loaded, snap1)
	}

	snap2 := Snapshot{
		Version: 0,
		Assignments: map[string]string{
			"task-1": "node-2",
		},
	}

	err = store.Save(snap2)
	if err != ErrOlderVersion {
		t.Errorf("Expected ErrOlderVersion, got %v", err)
	}

	snap3 := Snapshot{
		Version: 2,
		Assignments: map[string]string{
			"task-1": "node-2",
		},
	}

	if err := store.Save(snap3); err != nil {
		t.Fatalf("Save newer version failed: %v", err)
	}

	loaded3, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !reflect.DeepEqual(loaded3, snap3) {
		t.Errorf("Loaded %v, expected %v", loaded3, snap3)
	}
}