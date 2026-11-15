package main

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestJournalAssertionsAcrossUTCMidnight(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		appendFirst bool
	}{
		{name: "append after snapshot"},
		{name: "copied fixture already spans dates", appendFirst: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			workspace := clonePristineWorkspace(t)
			backdateGenesis(t, workspace, time.Now().UTC().AddDate(0, 0, -1))
			if tc.appendFirst {
				runFIDCommand(t, workspace, "label", "add", "earlier", "--by", "developer")
			}

			before := splitJournalLines(mustReadJournal(t, workspace))
			args := []string{"--json", "label", "add", "midnight", "--by", "developer", "--idempotency-key", "midnight-add"}
			result := executeCommandForTest(workspace, args...)
			if result.err != nil || result.stderr != "" {
				t.Fatalf("label add = %#v", result)
			}
			after := mustReadJournal(t, workspace)
			lines := splitJournalLines(after)
			if len(lines) != len(before)+1 {
				t.Fatalf("journal events = %d, want %d", len(lines), len(before)+1)
			}
			for i, line := range before {
				if lines[i] != line {
					t.Fatalf("journal event %d changed", i+1)
				}
			}
			last := mustReadLastJournalEvent(t, workspace)
			if last.EventType != "label.added" {
				t.Fatalf("last event = %s, want label.added", last.EventType)
			}
			stdout, _ := runFIDCommand(t, workspace, "journal", "inspect", last.EventID)
			var inspected journalInspectOutput
			if err := json.Unmarshal([]byte(stdout), &inspected); err != nil {
				t.Fatal(err)
			}
			if inspected.EventID != last.EventID || inspected.EventType != last.EventType {
				t.Fatalf("inspected event = %+v, want %s", inspected, last.EventID)
			}
			runFIDCommand(t, workspace, args...)
			runFIDCommand(t, workspace, "rebuild")
			if got := mustReadJournal(t, workspace); !bytes.Equal(got, after) {
				t.Fatal("retry or rebuild changed journal bytes")
			}
		})
	}
}
