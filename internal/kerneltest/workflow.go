// Package kerneltest supplies the obsolete task graph solely to kernel tests.
// It is removed with those tests in RESET 1.5c. Product commands never use it.
package kerneltest

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

// Bind gives an existing fixture task an isolated legacy workflow. Its first
// working state is called building so tests exercise the old execution guard
// without trying to use the fixed task's start edge.
func Bind(root, migrations, id string) error {
	return bind(root, migrations, id, "kernel-task", "todo", true)
}

// BindTo attaches a declared legacy test workflow without a public task assignment.
func BindTo(root, migrations, id, workflowID, status string) error {
	return bind(root, migrations, id, workflowID, status, false)
}

func bind(root, migrations, id, workflowID, status string, defineDefault bool) error {
	return journal.WithLock(root, journal.LockOwner{Command: "kernel test fixture"}, func(lock *journal.WriteLock) error {
		if _, err := journal.EnsureProjectionFreshWithLock(lock, migrations); err != nil {
			return err
		}
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer func() { _ = conn.Close() }()
		var n int
		if err := conn.QueryRow(`SELECT count(*) FROM workflows WHERE id='kernel-task'`).Scan(&n); err != nil {
			return err
		}
		if n == 0 && defineDefault {
			if err := define(lock, conn); err != nil {
				return err
			}
		}
		sequence, err := journal.NextAggregateSequence(filepath.Join(config.LedgerDir(root), "journal"), journal.Aggregate{ID: id, Type: "issue"})
		if err != nil {
			return err
		}
		event := journal.Event{EventSchema: journal.EventSchemaV2, EventType: "issue.updated", CommandName: "issue.update", ActorV2: actor(), OccurredAt: time.Now().UTC(), Aggregate: journal.Aggregate{ID: id, Type: "issue", Sequence: sequence}, Effects: []journal.Effect{
			{Op: "update", Table: "issues", PK: map[string]any{"id": id}, Values: map[string]any{"workflow_id": workflowID, "status": status}, Preconditions: map[string]any{"workflow_id": "task-default", "status": "todo"}},
			{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "kernel-fixture-" + id}, Values: map[string]any{"ref_type": "issue", "ref_id": id, "event": "issue-updated", "actor": "fixture", "occurred_at": time.Now().UTC().Format(time.RFC3339)}},
		}}
		prepared, err := journal.PrepareEvent(lock, event)
		if err != nil {
			return err
		}
		event.Effects[1].Values["journal_event_id"] = prepared.EventID
		_, err = journal.CommitEvent(lock, conn, event)
		return err
	})
}

func actor() journal.EventActor {
	return journal.EventActor{ActorType: "service", ActorID: "fixture", Producer: "test", ProvenanceRef: "kernel-only"}
}

func define(lock *journal.WriteLock, conn *sql.DB) error {
	effects := []journal.Effect{{Op: "insert", Table: "workflows", PK: map[string]any{"id": "kernel-task"}, Values: map[string]any{"name": "Kernel test fixture", "applies_to": "task", "active": true}}}
	for i, state := range []string{"todo", "building", "dev_complete", "qa_review", "qa_block", "failed_gate", "qa_pass", "done", "cancelled"} {
		category := "started"
		switch state {
		case "todo":
			category = "unstarted"
		case "qa_pass", "done":
			category = "completed"
		case "cancelled":
			category = "cancelled"
		}
		effects = append(effects, journal.Effect{Op: "insert", Table: "workflow_states", PK: map[string]any{"workflow_id": "kernel-task", "state": state}, Values: map[string]any{"category": category, "ordinal": i + 1}})
	}
	for _, edge := range [][2]string{{"todo", "building"}, {"building", "dev_complete"}, {"building", "failed_gate"}, {"failed_gate", "building"}, {"dev_complete", "qa_review"}, {"qa_review", "qa_pass"}, {"qa_review", "qa_block"}, {"qa_block", "building"}, {"qa_pass", "done"}} {
		effects = append(effects, transition(edge[0], edge[1]))
	}
	for _, state := range []string{"todo", "building", "dev_complete", "qa_review", "qa_block", "failed_gate", "qa_pass", "done"} {
		effects = append(effects, transition(state, "cancelled"))
	}
	_, err := journal.CommitEvent(lock, conn, journal.Event{EventSchema: journal.EventSchemaV2, EventType: "workflow.defined", CommandName: "workflow.define", ActorV2: actor(), OccurredAt: time.Now().UTC(), Aggregate: journal.Aggregate{ID: "kernel-task", Type: "workflow", Sequence: 1}, Effects: effects})
	if err != nil {
		return fmt.Errorf("define kernel test fixture: %w", err)
	}
	return nil
}

func transition(from, to string) journal.Effect {
	return journal.Effect{Op: "insert", Table: "workflow_transitions", PK: map[string]any{"workflow_id": "kernel-task", "from_state": from, "to_state": to}, Values: map[string]any{}}
}
