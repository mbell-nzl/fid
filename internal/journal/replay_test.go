package journal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestReplayJournalRejectsNilConnection(t *testing.T) {
	t.Parallel()

	_, err := ReplayJournal(nil, filepath.Join(config.LedgerDir(t.TempDir()), "journal"), nil)
	if err == nil {
		t.Fatal("ReplayJournal() error = nil, want nil connection rejection")
	}
	if err.Error() != "replay journal: connection is nil" {
		t.Fatalf("ReplayJournal() error = %q, want %q", err.Error(), "replay journal: connection is nil")
	}
}

func TestReplayJournalClassifiesRetiredIssueEffectKeyAsStopped(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	journalPath := filepath.Join(config.LedgerDir(workspace), "journal", "2026-08-11.jsonl")
	writeReplayEvents(t, journalPath, Event{
		EventID:     "EVT-retired-effect-key",
		EventSchema: EventSchemaV1,
		EventType:   "issue.created",
		OccurredAt:  time.Date(2026, time.August, 11, 0, 0, 0, 0, time.UTC),
		Actor:       Actor{Role: "developer"},
		Command:     Command{Argv: []string{"fid", "issue", "create"}, CWD: ".", FIDVersion: "0.1.0"},
		Effects: []Effect{{
			Op:    "insert",
			Table: "issues",
			PK:    map[string]any{"id": "FID-retired-effect-key"},
			Values: map[string]any{
				"kind":        "task",
				"title":       "retired key",
				"workflow_id": "task-default",
				"status":      "todo",
				"in_inbox":    true,
			},
		}},
	})

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	_, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
	var corruptErr *CorruptJournalError
	if !errors.As(err, &corruptErr) {
		t.Fatalf("ReplayJournal() error = %T %v, want *CorruptJournalError", err, err)
	}
	if corruptErr.File != ".fid/journal/2026-08-11.jsonl" || corruptErr.Line != 1 {
		t.Fatalf("replay stop location = %s:%d, want .fid/journal/2026-08-11.jsonl:1", corruptErr.File, corruptErr.Line)
	}
}

func TestReplayJournalAppliesGoldenV1JournalInFileOrder(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	stats, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
	if err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}

	if stats.EventsReplayed != 2 {
		t.Fatalf("replayed events = %d, want %d", stats.EventsReplayed, 2)
	}
	if stats.FilesReplayed != 1 {
		t.Fatalf("replayed files = %d, want %d", stats.FilesReplayed, 1)
	}
	if stats.UpgradedEvents != 0 {
		t.Fatalf("upgraded events = %d, want %d", stats.UpgradedEvents, 0)
	}

	var labelCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = ?`, "developer").Scan(&labelCount); err != nil {
		t.Fatalf("query labels: %v", err)
	}
	if labelCount != 1 {
		t.Fatalf("projected label count = %d, want %d", labelCount, 1)
	}

	var title, workflowID, status string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&title); err != nil {
		t.Fatalf("query issues: %v", err)
	}
	if title != "Implement replay" {
		t.Fatalf("projected issue title = %q, want %q", title, "Implement replay")
	}
	if err := conn.QueryRow(`SELECT workflow_id, status FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&workflowID, &status); err != nil {
		t.Fatalf("query migrated issue workflow pair: %v", err)
	}
	if workflowID != "task-default" || status != "todo" {
		t.Fatalf("migrated issue workflow pair = (%q, %q), want (task-default, todo)", workflowID, status)
	}
	var category string
	if err := conn.QueryRow(`SELECT category FROM workflow_states WHERE workflow_id = ? AND state = ?`, workflowID, status).Scan(&category); err != nil {
		t.Fatalf("query migrated workflow state: %v", err)
	}
	if category != "unstarted" {
		t.Fatalf("migrated workflow category = %q, want unstarted", category)
	}
}

func TestReplayJournalAppliesPhase20IssueRevisionGolden(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "p20-issue-revision-effects.golden.jsonl", "2026-08-11.jsonl")
	journalPath := filepath.Join(config.LedgerDir(workspace), "journal", "2026-08-11.jsonl")

	payload, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", journalPath, err)
	}
	lines := strings.Split(strings.TrimSpace(string(payload)), "\n")
	wantTypes := []string{"issue.created", "issue.updated", "issue.transitioned", "issue.updated", "issue.transitioned", "issue.updated", "issue.ac_added", "label.added", "issue.label.added", "issue.label.removed"}
	wantPriorRevisions := []float64{0, 1, 2, 3, 4, 5, 6, 0, 7, 8}
	if len(lines) != len(wantTypes) {
		t.Fatalf("golden Event count = %d, want %d", len(lines), len(wantTypes))
	}
	for index, line := range lines {
		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("Unmarshal(golden Event %d) error = %v", index, err)
		}
		if event.EventType != wantTypes[index] {
			t.Fatalf("golden Event %d type = %q, want %q", index, event.EventType, wantTypes[index])
		}
		effect := event.Effects[0]
		if event.EventType == "issue.created" {
			if got := effect.Values["revision"]; got != float64(1) || len(effect.Prior) != 0 || len(effect.Preconditions) != 0 {
				t.Fatalf("issue.created revision Effect = %#v, want absolute revision 1 only", effect)
			}
			continue
		}
		if event.EventType == "label.added" {
			if effect.Table != "labels" || effect.Op != "insert" {
				t.Fatalf("label.added Effect = %#v, want labels insert", effect)
			}
			continue
		}
		prior := wantPriorRevisions[index]
		var revisionEffect *Effect
		for effectIndex := range event.Effects {
			candidate := &event.Effects[effectIndex]
			if candidate.Table == "issues" && candidate.Op == "update" {
				revisionEffect = candidate
			}
		}
		if revisionEffect == nil {
			t.Fatalf("%s has no Issue revision Effect: %#v", event.EventType, event.Effects)
		}
		if got := revisionEffect.Values["revision"]; got != prior+1 {
			t.Fatalf("%s revision value = %#v, want %v", event.EventType, got, prior+1)
		}
		if got := revisionEffect.Prior["revision"]; got != prior {
			t.Fatalf("%s revision prior = %#v, want %v", event.EventType, got, prior)
		}
		if got := revisionEffect.Preconditions["revision"]; got != prior {
			t.Fatalf("%s revision precondition = %#v, want %v", event.EventType, got, prior)
		}
	}

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	if _, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil); err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	var revision int
	if err := conn.QueryRow(`SELECT revision FROM issues WHERE id = 'FID-p20-revision'`).Scan(&revision); err != nil {
		t.Fatalf("query golden Issue revision: %v", err)
	}
	if revision != 9 {
		t.Fatalf("golden replay revision = %d, want 9", revision)
	}
	var criterionCount, issueLabelCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM acceptance_criteria WHERE issue_id = 'FID-p20-revision'`).Scan(&criterionCount); err != nil {
		t.Fatalf("count golden acceptance criteria: %v", err)
	}
	if err := conn.QueryRow(`SELECT COUNT(*) FROM issue_labels WHERE issue_id = 'FID-p20-revision'`).Scan(&issueLabelCount); err != nil {
		t.Fatalf("count golden Issue labels: %v", err)
	}
	if criterionCount != 1 || issueLabelCount != 0 {
		t.Fatalf("golden domain counts = (AC %d, labels %d), want (1, 0)", criterionCount, issueLabelCount)
	}
}

func TestRebuildWorkspaceRecreatesIdenticalStateAfterDeletingDatabase(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	conn := openReplayDB(t, dbPath)

	stats, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
	if err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if stats.EventsReplayed != 2 {
		t.Fatalf("initial replayed events = %d, want %d", stats.EventsReplayed, 2)
	}

	before := snapshotReplayState(t, conn)
	closeProjectDB(t, conn)

	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t))
	if err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}
	if rebuildStats.EventsReplayed != 2 {
		t.Fatalf("rebuilt events = %d, want %d", rebuildStats.EventsReplayed, 2)
	}

	rebuiltConn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}
	defer closeProjectDB(t, rebuiltConn)

	after := snapshotReplayState(t, rebuiltConn)
	if before != after {
		t.Fatalf("rebuilt state = %q, want %q", after, before)
	}
}

func TestReplayJournalReportsMalformedTrailingLine(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "corrupt-tail.jsonl", "2026-06-05.jsonl")

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	_, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
	if err == nil {
		t.Fatal("ReplayJournal() error = nil, want corruption error")
	}

	var corruptErr *CorruptJournalError
	if !errors.As(err, &corruptErr) {
		t.Fatalf("ReplayJournal() error = %T, want *CorruptJournalError", err)
	}
	if corruptErr.File != filepath.ToSlash(filepath.Join(".fid", "journal", "2026-06-05.jsonl")) {
		t.Fatalf("corrupt file = %q, want %q", corruptErr.File, filepath.ToSlash(filepath.Join(".fid", "journal", "2026-06-05.jsonl")))
	}
	if corruptErr.Line != 2 {
		t.Fatalf("corrupt line = %d, want %d", corruptErr.Line, 2)
	}
	if !strings.Contains(corruptErr.Error(), ".fid/journal/2026-06-05.jsonl:2") {
		t.Fatalf("CorruptJournalError.Error() = %q, want file:line context", corruptErr.Error())
	}
	if corruptErr.Unwrap() == nil {
		t.Fatal("CorruptJournalError.Unwrap() = nil, want underlying decode error")
	}
}

func TestReplayJournalStopsOnPreconditionConflict(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "conflicting-precondition.jsonl", "2026-06-05.jsonl")

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	_, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
	if err == nil {
		t.Fatal("ReplayJournal() error = nil, want precondition conflict")
	}

	var conflictErr *PreconditionConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("ReplayJournal() error = %T, want *PreconditionConflictError", err)
	}
	if conflictErr.EventID != "EVT-20260605-000002" {
		t.Fatalf("conflict event_id = %q, want %q", conflictErr.EventID, "EVT-20260605-000002")
	}
	if conflictErr.Table != "issues" {
		t.Fatalf("conflict table = %q, want %q", conflictErr.Table, "issues")
	}
	if conflictErr.Line != 2 {
		t.Fatalf("conflict line = %d, want %d", conflictErr.Line, 2)
	}
}

func TestReplayJournalRejectsStaleIssueRevisionBeforeApplyingEffect(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	journalPath := filepath.Join(config.LedgerDir(workspace), "journal", "2026-08-11.jsonl")
	issueID := "FID-stale-revision"
	writeReplayEvents(t, journalPath,
		Event{
			EventID:     "EVT-revision-created",
			EventSchema: EventSchemaV1,
			EventType:   "issue.created",
			OccurredAt:  time.Date(2026, time.August, 11, 0, 0, 0, 0, time.UTC),
			Actor:       Actor{Role: "developer"},
			Command:     Command{Argv: []string{"fid", "issue", "create"}, CWD: ".", FIDVersion: "0.1.0"},
			Effects: []Effect{{
				Op:    "insert",
				Table: "issues",
				PK:    map[string]any{"id": issueID},
				Values: map[string]any{
					"kind":        "task",
					"title":       "stale revision",
					"workflow_id": "task-default",
					"status":      "todo",
					"revision":    1,
				},
			}},
		},
		issueRevisionReplayEvent("EVT-revision-current", issueID, 1, 2),
		issueRevisionReplayEvent("EVT-revision-stale", issueID, 1, 2),
	)

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	_, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
	var conflictErr *PreconditionConflictError
	if !errors.As(err, &conflictErr) {
		t.Fatalf("ReplayJournal() error = %T %v, want *PreconditionConflictError", err, err)
	}
	if conflictErr.EventID != "EVT-revision-stale" {
		t.Fatalf("conflict event ID = %q, want EVT-revision-stale", conflictErr.EventID)
	}
	if got := conflictErr.Expected["revision"]; got != float64(1) {
		t.Fatalf("conflict expected revision = %#v, want 1", got)
	}
	if got := conflictErr.Actual["revision"]; got != int64(2) {
		t.Fatalf("conflict actual revision = %#v, want 2", got)
	}
	var revision int
	if err := conn.QueryRow(`SELECT revision FROM issues WHERE id = ?`, issueID).Scan(&revision); err != nil {
		t.Fatalf("query replayed revision: %v", err)
	}
	if revision != 2 {
		t.Fatalf("revision after stale replay = %d, want 2", revision)
	}
}

func issueRevisionReplayEvent(eventID, issueID string, prior, next int) Event {
	return Event{
		EventID:     eventID,
		EventSchema: EventSchemaV1,
		EventType:   "issue.updated",
		OccurredAt:  time.Date(2026, time.August, 11, 0, 1, 0, 0, time.UTC),
		Actor:       Actor{Role: "developer"},
		Command:     Command{Argv: []string{"fid", "issue", "update"}, CWD: ".", FIDVersion: "0.1.0"},
		Effects: []Effect{{
			Op:    "update",
			Table: "issues",
			PK:    map[string]any{"id": issueID},
			Values: map[string]any{
				"revision": next,
			},
			Prior: map[string]any{
				"revision": prior,
			},
			Preconditions: map[string]any{
				"revision": prior,
			},
		}},
	}
}

func TestEnsureProjectionFreshUsesProjectionMetaMatrix(t *testing.T) {
	t.Parallel()

	t.Run("clean match stays clean", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionFresh {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionFresh)
		}
		if result.ReplayStats.EventsReplayed != 0 {
			t.Fatalf("fresh replayed events = %d, want %d", result.ReplayStats.EventsReplayed, 0)
		}
	})

	t.Run("missing meta rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
		if _, err := conn.Exec(`DELETE FROM projection_meta`); err != nil {
			t.Fatalf("delete projection_meta: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
		if result.ReplayStats.EventsReplayed != 2 {
			t.Fatalf("rebuilt events = %d, want %d", result.ReplayStats.EventsReplayed, 2)
		}
	})

	t.Run("missing database rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
		if result.ReplayStats.EventsReplayed != 2 {
			t.Fatalf("rebuilt events = %d, want %d", result.ReplayStats.EventsReplayed, 2)
		}
	})

	t.Run("schema mismatch rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
		if _, err := conn.Exec(`UPDATE projection_meta SET projection_schema_version = ? WHERE key = 'default'`, 999); err != nil {
			t.Fatalf("update projection_meta schema version: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
	})

	t.Run("position mismatch rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
		if _, err := conn.Exec(`UPDATE projection_meta SET journal_offset = ? WHERE key = 'default'`, 99999); err != nil {
			t.Fatalf("update projection_meta journal_offset: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
	})

	t.Run("hash mismatch rebuilds", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
		if _, err := conn.Exec(`UPDATE projection_meta SET last_event_hash = ? WHERE key = 'default'`, "sha256:bad"); err != nil {
			t.Fatalf("update projection_meta last_event_hash: %v", err)
		}
		closeProjectDB(t, conn)

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionRebuild {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
		}
	})

	t.Run("journal advance catches up", func(t *testing.T) {
		t.Parallel()

		workspace := newReplayWorkspace(t)
		writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

		if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
			t.Fatalf("RebuildWorkspace() error = %v", err)
		}

		appendReplayFixtureLine(t, workspace, "replay-catchup-event.json", "2026-06-05.jsonl")

		result, err := EnsureProjectionFresh(workspace, testProjectMigrationsDir(t))
		if err != nil {
			t.Fatalf("EnsureProjectionFresh() error = %v", err)
		}

		if result.Decision != fiddb.ProjectionCatchUp {
			t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionCatchUp)
		}
		if result.ReplayStats.EventsReplayed != 1 {
			t.Fatalf("caught-up events = %d, want %d", result.ReplayStats.EventsReplayed, 1)
		}
	})
}

func TestEnsureProjectionFreshWithLockRebuildsMissingSQLiteProjection(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	journalPath := filepath.Join(config.LedgerDir(workspace), "journal", "2026-06-05.jsonl")
	before, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) before rebuild error = %v", journalPath, err)
	}

	var result EnsureProjectionResult
	err = WithLock(workspace, LockOwner{Command: "fid catch-up"}, func(lock *WriteLock) error {
		if _, statErr := os.Stat(filepath.Join(config.LedgerDir(workspace), lockDirectoryName)); statErr != nil {
			return fmt.Errorf("caller-held lock missing before catch-up: %w", statErr)
		}

		var ensureErr error
		result, ensureErr = EnsureProjectionFreshWithLock(lock, testProjectMigrationsDir(t))
		return ensureErr
	})
	if err != nil {
		t.Fatalf("EnsureProjectionFreshWithLock() error = %v", err)
	}

	if result.Decision != fiddb.ProjectionRebuild {
		t.Fatalf("freshness decision = %q, want %q", result.Decision, fiddb.ProjectionRebuild)
	}
	if result.ReplayStats.EventsReplayed != 2 {
		t.Fatalf("replayed events = %d, want %d", result.ReplayStats.EventsReplayed, 2)
	}

	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() rebuilt projection error = %v", err)
	}
	defer closeProjectDB(t, conn)

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&title); err != nil {
		t.Fatalf("query rebuilt issue: %v", err)
	}
	if title != "Implement replay" {
		t.Fatalf("rebuilt issue title = %q, want %q", title, "Implement replay")
	}

	after, err := os.ReadFile(journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) after rebuild error = %v", journalPath, err)
	}
	if !bytes.Equal(after, before) {
		t.Fatalf("EnsureProjectionFreshWithLock() changed journal bytes")
	}
}

func TestOrderedJournalFilesSkipsNonJournalEntries(t *testing.T) {
	t.Parallel()

	journalDir := filepath.Join(config.LedgerDir(t.TempDir()), "journal")
	mustMkdirAll(t, journalDir)

	files := []string{
		"2026-06-07.jsonl",
		"2026-06-06.jsonl",
		"2026-06-07.txt",
		"not-a-date.jsonl",
	}
	for _, name := range files {
		path := filepath.Join(journalDir, name)
		if err := os.WriteFile(path, []byte("{}\n"), 0o644); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", path, err)
		}
	}
	if err := os.Mkdir(filepath.Join(journalDir, "nested"), 0o755); err != nil {
		t.Fatalf("Mkdir(nested) error = %v", err)
	}

	got, err := orderedJournalFiles(journalDir)
	if err != nil {
		t.Fatalf("orderedJournalFiles() error = %v", err)
	}

	want := []string{"2026-06-06.jsonl", "2026-06-07.jsonl"}
	if !equalReplayStrings(got, want) {
		t.Fatalf("orderedJournalFiles() = %#v, want %#v", got, want)
	}
}

func TestReplayJournalHonorsStartOffset(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	journalPath := filepath.Join(config.LedgerDir(workspace), "journal", "2026-06-05.jsonl")
	first := replayRoleEvent("EVT-20260605-000001", "developer")
	second := replayRoleEvent("EVT-20260605-000002", "qa_lead")
	writeReplayEvents(t, journalPath, first, second)

	firstLine, err := marshalReplayEvent(first)
	if err != nil {
		t.Fatalf("marshalReplayEvent(first) error = %v", err)
	}

	conn := openReplayDB(t, filepath.Join(t.TempDir(), "fid.sqlite"))
	defer closeProjectDB(t, conn)

	stats, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), &ReplayOptions{
		StartFile:   "2026-06-05.jsonl",
		StartOffset: int64(len(firstLine) + 1),
	})
	if err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}

	if stats.EventsReplayed != 1 {
		t.Fatalf("replayed events = %d, want %d", stats.EventsReplayed, 1)
	}

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM labels`).Scan(&count); err != nil {
		t.Fatalf("count labels: %v", err)
	}
	if count != 1 {
		t.Fatalf("projected role count = %d, want %d", count, 1)
	}
}

func newReplayWorkspace(t *testing.T) string {
	t.Helper()

	workspace := t.TempDir()
	mustMkdirAll(t, filepath.Join(config.LedgerDir(workspace), "journal"))
	return workspace
}

func openReplayDB(t *testing.T, dbPath string) *sql.DB {
	t.Helper()

	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	return conn
}

func replayFixturePath(t *testing.T, name string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "testdata", "journal", name))
}

func writeReplayFixture(t *testing.T, workspaceRoot, fixtureName, journalFileName string) {
	t.Helper()

	payload, err := os.ReadFile(replayFixturePath(t, fixtureName))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixtureName, err)
	}

	journalPath := filepath.Join(config.LedgerDir(workspaceRoot), "journal", journalFileName)
	if err := os.WriteFile(journalPath, payload, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", journalPath, err)
	}
}

func appendReplayFixtureLine(t *testing.T, workspaceRoot, fixtureName, journalFileName string) {
	t.Helper()

	payload, err := os.ReadFile(replayFixturePath(t, fixtureName))
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixtureName, err)
	}

	journalPath := filepath.Join(config.LedgerDir(workspaceRoot), "journal", journalFileName)
	file, err := os.OpenFile(journalPath, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatalf("OpenFile(%q) error = %v", journalPath, err)
	}
	defer closeProjectDB(t, file)

	if _, err := file.Write(append(payload, '\n')); err != nil {
		t.Fatalf("append replay fixture line: %v", err)
	}
}

func snapshotReplayState(t *testing.T, conn *sql.DB) string {
	t.Helper()

	var roleName string
	if err := conn.QueryRow(`SELECT name FROM labels ORDER BY name LIMIT 1`).Scan(&roleName); err != nil {
		t.Fatalf("query snapshot label: %v", err)
	}

	var title string
	if err := conn.QueryRow(`SELECT title FROM issues WHERE id = ?`, "FID-a1b2c3").Scan(&title); err != nil {
		t.Fatalf("query snapshot issue: %v", err)
	}

	meta, found, err := fiddb.ReadProjectionMeta(conn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() found = false, want true")
	}

	return strings.Join([]string{
		roleName,
		title,
		meta.JournalFile,
		meta.LastEventID,
		meta.LastEventHash,
		meta.ReplayStatus,
	}, "|")
}

func replayRoleEvent(eventID, roleName string) Event {
	return Event{
		EventID:     eventID,
		EventSchema: EventSchemaV1,
		EventType:   "project.initialized",
		OccurredAt:  time.Date(2026, time.June, 5, 14, 23, 11, 0, time.UTC),
		Actor: Actor{
			Role: "human",
		},
		Command: Command{
			Argv:       []string{"fid", "init"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "labels",
				PK: map[string]any{
					"name": roleName,
				},
				Values: map[string]any{
					"description": roleName,
					"colour":      "#00ff00",
				},
			},
		},
	}
}

func writeReplayEvents(t *testing.T, path string, events ...Event) {
	t.Helper()

	var payload []byte
	for _, event := range events {
		line, err := marshalReplayEvent(event)
		if err != nil {
			t.Fatalf("marshalReplayEvent() error = %v", err)
		}
		payload = append(payload, line...)
		payload = append(payload, '\n')
	}

	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func marshalReplayEvent(event Event) ([]byte, error) {
	return MarshalCanonicalJSON(event)
}

func equalReplayStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}

	sortedGot := append([]string(nil), got...)
	sortedWant := append([]string(nil), want...)
	sort.Strings(sortedGot)
	sort.Strings(sortedWant)
	for index := range sortedGot {
		if sortedGot[index] != sortedWant[index] {
			return false
		}
	}

	return true
}
