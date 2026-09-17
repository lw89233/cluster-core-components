package reconciler

import "testing"

func TestValidateState_Valid(t *testing.T) {
	validJSON := []byte(`{"version": 1, "assignments": {"task-1": "node-1"}}`)

	_, err := ValidateState(validJSON)
	if err != nil {
		t.Errorf("Expected valid state, got error: %v", err)
	}
}

func TestValidateState_MalformedJSON(t *testing.T) {
	badJSON := []byte(`{"version": 1, "assignments": {"task-1": "node-1"`)

	_, err := ValidateState(badJSON)
	if err != ErrMalformedData {
		t.Errorf("Expected ErrMalformedData, got %v", err)
	}
}

func TestValidateState_MissingFields(t *testing.T) {
	missingJSON := []byte(`{"version": 0, "assignments": null}`)

	_, err := ValidateState(missingJSON)
	if err != ErrMissingFields {
		t.Errorf("Expected ErrMissingFields, got %v", err)
	}
}