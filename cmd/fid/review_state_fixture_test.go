package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Hand-off arrives in 1.5a7. Seed its state through the journal so the hold's
// completion/replay/export contract is still tested on the fixed task graph.
func seedReviewState(t *testing.T, root, task string) {
	t.Helper()
	err := journal.WithLock(root, journal.LockOwner{Command: "review state fixture"}, func(lock *journal.WriteLock) error {
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		var from string
		if err := conn.QueryRow(`SELECT status FROM issues WHERE id=?`, task).Scan(&from); err != nil {
			return err
		}
		seq, err := journal.NextAggregateSequence(filepath.Join(config.LedgerDir(root), "journal"), journal.Aggregate{ID: task, Type: "issue"})
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "issue.transitioned", CommandName: "issue.transition", ActorV2: journal.EventActor{ActorType: "agent", ActorID: "fixture", Producer: "test", ProvenanceRef: "review-fixture"}, OccurredAt: now, Aggregate: journal.Aggregate{ID: task, Type: "issue", Sequence: seq}, Effects: []journal.Effect{
			{Op: "update", Table: "issues", PK: map[string]any{"id": task}, Values: map[string]any{"status": "in_review"}, Preconditions: map[string]any{"status": from}},
			{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "review-fixture-" + task}, Values: map[string]any{"ref_type": "issue", "ref_id": task, "event": "status-changed", "actor": "fixture", "occurred_at": now.Format(time.RFC3339)}},
		}}
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event.Effects[1].Values["journal_event_id"] = prepared.EventID
		_, err = journal.CommitEvent(lock, conn, event)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
}
