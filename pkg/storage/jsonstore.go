package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var ErrOlderVersion = errors.New("snapshot version is older than current")

type Snapshot struct {
	Version     int64             `json:"version"`
	Assignments map[string]string `json:"assignments"`
}

type JSONStore struct {
	mu       sync.RWMutex
	filepath string
	version  int64
}

func NewJSONStore(path string) *JSONStore {
	return &JSONStore{
		filepath: path,
	}
}

func (s *JSONStore) Save(snapshot Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if snapshot.Version < s.version {
		return ErrOlderVersion
	}

	dir := filepath.Dir(s.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "state-tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()

	defer func() {
		tmpFile.Close()
		os.Remove(tmpName)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		return err
	}

	if err := tmpFile.Sync(); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, s.filepath); err != nil {
		return err
	}

	s.version = snapshot.Version
	return nil
}

func (s *JSONStore) Load() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var snap Snapshot

	data, err := os.ReadFile(s.filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return snap, nil
		}
		return snap, err
	}

	if err := json.Unmarshal(data, &snap); err != nil {
		return snap, err
	}

	s.version = snap.Version
	return snap, nil
}