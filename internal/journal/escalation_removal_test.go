package journal

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLatestProjectionRejectsRetiredEscalationEffects(t *testing.T) {
	for _, test := range []struct {
		name   string
		effect Effect
	}{
		{"table", Effect{Op: "insert", Table: "escalations", PK: map[string]any{"id": "ESC-retired"}, Values: map[string]any{"status": "open"}}},
		{"triage column", Effect{Op: "update", Table: "triage_events", PK: map[string]any{"id": "TRI-retired"}, Values: map[string]any{"esc_id": "ESC-retired"}}},
		{"concern column", Effect{Op: "update", Table: "concerns", PK: map[string]any{"id": "CON-retired"}, Values: map[string]any{"escalated_to_esc_id": "ESC-retired"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace := newReplayWorkspace(t)
			conn := openReplayDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
			defer closeProjectDB(t, conn)
			event := phaseTwoFullSchemaEvent()
			event.Effects = []Effect{{Op: "insert", Table: "labels", PK: map[string]any{"name": "must-rollback"}}, test.effect}
			hash, err := HashEvent(event)
			if err != nil {
				t.Fatal(err)
			}
			if err := ProjectEvent(conn, ProjectionInput{Event: event, EventHash: hash}); err == nil || !strings.Contains(err.Error(), "no such") {
				t.Fatalf("latest ProjectEvent admitted retired Effect or wrong rejection: %v", err)
			}
			var count int
			if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = 'must-rollback'`).Scan(&count); err != nil || count != 0 {
				t.Fatalf("rejected Event leaked preceding Effect: count%d err%v", count, err)
			}
			journalDir := filepath.Join(workspace, ".fid", "journal")
			writeReplayEvents(t, filepath.Join(journalDir, "2026-06-07.jsonl"), event)
			if _, err := ReplayJournal(conn, journalDir, nil); err == nil || !strings.Contains(err.Error(), "no such") {
				t.Fatalf("latest replay admitted retired Effect or wrong rejection: %v", err)
			}
			if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = 'must-rollback'`).Scan(&count); err != nil || count != 0 {
				t.Fatalf("rejected replay leaked preceding Effect: count%d err%v", count, err)
			}
			control := event
			control.Effects = []Effect{{Op: "insert", Table: "labels", PK: map[string]any{"name": "healthy-control"}}}
			controlHash, err := HashEvent(control)
			if err != nil {
				t.Fatal(err)
			}
			if err := ProjectEvent(conn, ProjectionInput{Event: control, EventHash: controlHash}); err != nil {
				t.Fatalf("healthy retained Effect rejected: %v", err)
			}
		})
	}
}
