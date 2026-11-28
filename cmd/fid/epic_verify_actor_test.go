package main

import (
	"testing"
	"time"
)

func TestEpicValidationVerifyWorkActors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		actor, code string
	}{
		{"agent:builder", ""},
		{"human:mark", ""},
		{"service:ci", "authority-denied"},
		{"agent:", "invalid-argument"},
	} {
		t.Run(tc.actor, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			before := snapshotWorkspacePersistence(t, f.workspace)
			r := f.run("--json", "epic", "validation", "verify", epic, "--actor", tc.actor, "--idempotency-key", "verify")
			if tc.code == "" {
				result, _ := lifecycleResult(t, r, "epic.validation.verify")
				if result["epic_id"] != epic {
					t.Fatalf("verified another epic: %v", result)
				}
				return
			}
			if r.err == nil {
				t.Fatal("unauthorized actor verified an epic")
			}
			assertJSONFailureCode(t, r.stderr, "epic.validation.verify", tc.code)
			if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
				t.Fatal("denied verification changed the ledger")
			}
		})
	}
}
