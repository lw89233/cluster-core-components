package reconciler

import (
	"encoding/json"
	"errors"
)

var ErrMalformedData = errors.New("malformed state data")
var ErrMissingFields = errors.New("missing required fields in state")

type StateSnapshot struct {
	Version     int64             `json:"version"`
	Assignments map[string]string `json:"assignments"`
}

func ValidateState(rawJSON []byte) (*StateSnapshot, error) {
	var snap StateSnapshot

	if err := json.Unmarshal(rawJSON, &snap); err != nil {
		return nil, ErrMalformedData
	}

	if snap.Version < 1 {
		return nil, ErrMissingFields
	}

	if snap.Assignments == nil {
		return nil, ErrMissingFields
	}

	return &snap, nil
}