package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

// Both subjects reach done through their owning commands, with receipts.
func TestFinishedWorkCannotBeCancelled(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "epic"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			var f *statusFixture
			var subject string
			if kind == "epic" {
				f, subject = newEpicCancellationFixture(t)
				f.clock = time.Now
				lifecycleResult(t, f.run("--json", "done", subject, "--actor", "human:mark", "--reason", "All work finished"), "task.done")
			} else {
				var hold string
				f, subject, hold = doneBuildFixture(t, "standard", true)
				lifecycleResult(t, f.run(append([]string{"--json"}, doneHandoffArgs(hold)...)...), "task.handoff")
				lifecycleResult(t, recordReviewTask(t, f, reviewInput(), "review", "human:courier", subject), "review.record")
				lifecycleResult(t, runDone(t, f, doneRunner(t, doneFacts(t, f, subject), nil), "--json", "done", subject, "--actor", "agent:finisher", "--idempotency-key", "finish"), "task.done")
			}
			assertFinishedCancellationAdmission(t, f, subject, kind)
			before := snapshotWorkspacePersistence(t, f.workspace)
			got := f.run("--json", "issue", "transition", subject, "--to", "cancelled", "--actor", "human:mark", "--reason", "Undo merged work")
			assertJSONFailureCode(t, got.stderr, "issue.transition", "transition-not-permitted")
			err := func() error {
				_, err := issues.Transition(issues.TransitionOptions{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, ID: subject, ToState: "cancelled", Actor: "human:mark", Reason: "Undo merged work"})
				return err
			}()
			var refused *issues.CommandTransitionError
			if !errors.As(err, &refused) || refused.Rule != "finished-work-stays-done" {
				t.Fatalf("cancel done %s = %#v, want typed finished-work-stays-done", kind, got)
			}
			for _, text := range []string{"finished work stays done", "new task", "## Rollback"} {
				if !strings.Contains(refused.Error(), text) {
					t.Errorf("refusal %q must include %q", refused.Error(), text)
				}
			}
			if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
				t.Fatal("refused cancellation changed persistence")
			}
		})
	}
}

func assertFinishedCancellationAdmission(t *testing.T, f *statusFixture, subject, kind string) {
	t.Helper()
	ledger := config.LedgerDir(f.workspace)
	dir := filepath.Join(ledger, "journal")
	conn, err := fiddb.Open(filepath.Join(ledger, "fid.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	var revision int
	if err := conn.QueryRow(`SELECT revision FROM issues WHERE id=?`, subject).Scan(&revision); err != nil {
		t.Fatal(err)
	}
	seq, err := journal.NextAggregateSequence(dir, journal.Aggregate{Type: "issue", ID: subject})
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	stamp := now.Format(time.RFC3339)
	e := journal.Event{EventSchema: 2, EventType: "issue.transitioned", OccurredAt: now, ActorV2: journal.EventActor{ActorType: "human", ActorID: "mark", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{Type: "issue", ID: subject, Sequence: seq}, CommandName: "issue.transition", Effects: []journal.Effect{
		{Op: "update", Table: "issues", PK: map[string]any{"id": subject}, Values: map[string]any{"status": "cancelled", "updated_at": stamp}, Preconditions: map[string]any{"status": "done"}},
		{Op: "update", Table: "issues", PK: map[string]any{"id": subject}, Values: map[string]any{"revision": revision + 1}, Prior: map[string]any{"revision": revision}, Preconditions: map[string]any{"revision": revision}},
		{Op: "insert", Table: "person_decisions", PK: map[string]any{"id": "PDN-cancel"}, Values: map[string]any{"kind": "cancel", "actor": "human:mark", "subject_type": kind, "subject_id": subject, "reason": "Undo merged work", "binding_id": subject, "decided_at": stamp}},
	}}
	if err := journal.WithLock(f.workspace, journal.LockOwner{Command: "malformed test"}, func(lock *journal.WriteLock) error { var err error; e, err = journal.PrepareEvent(lock, e); return err }); err != nil {
		t.Fatal(err)
	}
	e.Effects = append(e.Effects, journal.Effect{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "ACT-cancel"}, Values: map[string]any{"journal_event_id": e.EventID, "ref_type": "issue", "ref_id": subject, "event": "status-changed", "details": nil, "actor": "mark", "occurred_at": stamp}})
	files, err := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	if err != nil || len(files) == 0 {
		t.Fatalf("journal files: %v", err)
	}
	file := files[len(files)-1]
	prefix, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	line, err := journal.MarshalCanonicalJSON(e)
	if err != nil {
		t.Fatal(err)
	}
	hash, err := journal.HashEvent(e)
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"commit", "dry-run", "project", "replay", "catch-up", "freshen", "rebuild"} {
		t.Run(mode, func(t *testing.T) {
			before := snapshotWorkspacePersistence(t, f.workspace)
			var err error
			switch mode {
			case "commit":
				err = journal.WithLock(f.workspace, journal.LockOwner{Command: "malformed test"}, func(lock *journal.WriteLock) error { _, err := journal.CommitEvent(lock, conn, e); return err })
			case "dry-run":
				err = journal.DryRunProjectEvent(conn, e)
			case "project":
				err = journal.ProjectEvent(conn, journal.ProjectionInput{Event: e, EventHash: hash, File: filepath.Base(file), Line: bytes.Count(prefix, []byte{'\n'}) + 1, Offset: int64(len(prefix) + len(line) + 1)})
			default:
				if err := os.WriteFile(file, append(append(bytes.Clone(prefix), line...), '\n'), 0o644); err != nil {
					t.Fatal(err)
				}
				switch mode {
				case "replay":
					path := filepath.Join(t.TempDir(), "replay.sqlite")
					if err := fiddb.MigrateUp(path, fidmigrations.FS); err != nil {
						t.Fatal(err)
					}
					replay, openErr := fiddb.Open(path)
					if openErr != nil {
						t.Fatal(openErr)
					}
					_, err = journal.ReplayJournal(replay, dir, nil)
					var state string
					if scanErr := replay.QueryRow(`SELECT status FROM issues WHERE id=?`, subject).Scan(&state); scanErr != nil || state != "done" {
						t.Fatalf("refused replay state=%q: %v", state, scanErr)
					}
					_ = replay.Close()
				case "catch-up":
					_, err = journal.ReplayJournal(conn, dir, &journal.ReplayOptions{StartFile: filepath.Base(file), StartOffset: int64(len(prefix))})
				case "freshen":
					_, err = journal.EnsureProjectionFresh(f.workspace, fidmigrations.FS)
				case "rebuild":
					_, err = journal.RebuildWorkspace(f.workspace, fidmigrations.FS)
				}
				stored, readErr := os.ReadFile(file)
				if readErr != nil || !bytes.Equal(stored, append(append(bytes.Clone(prefix), line...), '\n')) {
					t.Fatalf("refusal changed journal: %v", readErr)
				}
				if err := os.WriteFile(file, prefix, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			var refused *journal.StateMoveError
			if !errors.As(err, &refused) || !strings.Contains(refused.Error(), "finished work stays done") {
				t.Fatalf("%s = %v (%T), want typed finished cancellation refusal", mode, err, err)
			}
			if !before.equal(snapshotWorkspacePersistence(t, f.workspace)) {
				t.Fatal("refusal changed persistence")
			}
		})
	}
	// A later legitimate write still succeeds against the original ledger.
	lifecycleResult(t, f.run("--json", "issue", "update", subject, "--title", fmt.Sprint("Finished ", kind), "--actor", "human:mark", "--idempotency-key", "retitle"), "issue.update")
}
