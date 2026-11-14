package main

import (
	"encoding/json"
	"os"
	"testing"
)

func writeActivationInput(t *testing.T, path string, value map[string]any) {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal activation input: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("write activation input: %v", err)
	}
}
