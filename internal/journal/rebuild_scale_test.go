package journal

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRebuildWorkspaceIsDeterministicAtSingleProjectScale(t *testing.T) {
	t.Parallel()

	const targetEvents = 1000

	workspace := newReplayWorkspace(t)
	seed := seedScaleReplayWorkspace(t, workspace, targetEvents)

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	conn := openReplayDB(t, dbPath)

	stats, err := ReplayJournal(conn, filepath.Join(workspace, ".fid", "journal"), nil)
	if err != nil {
		t.Fatalf("ReplayJournal() error = %v", err)
	}
	if stats.EventsReplayed != seed.EventCount {
		t.Fatalf("initial replayed events = %d, want %d", stats.EventsReplayed, seed.EventCount)
	}

	before := snapshotScaleProjectionState(t, conn)
	closeProjectDB(t, conn)

	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer")
	if err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}
	if rebuildStats.EventsReplayed != seed.EventCount {
		t.Fatalf("rebuilt events = %d, want %d", rebuildStats.EventsReplayed, seed.EventCount)
	}

	rebuiltConn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}

	afterDeleteRebuild := snapshotScaleProjectionState(t, rebuiltConn)
	closeProjectDB(t, rebuiltConn)

	if before != afterDeleteRebuild {
		t.Fatalf("rebuilt state after delete = %q, want %q", afterDeleteRebuild, before)
	}

	secondRebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t), "developer")
	if err != nil {
		t.Fatalf("RebuildWorkspace() second pass error = %v", err)
	}
	if secondRebuildStats.EventsReplayed != seed.EventCount {
		t.Fatalf("second rebuilt events = %d, want %d", secondRebuildStats.EventsReplayed, seed.EventCount)
	}

	secondConn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() second rebuilt error = %v", err)
	}
	defer closeProjectDB(t, secondConn)

	afterSecondRebuild := snapshotScaleProjectionState(t, secondConn)
	if afterDeleteRebuild != afterSecondRebuild {
		t.Fatalf("rebuilt state after second rebuild = %q, want %q", afterSecondRebuild, afterDeleteRebuild)
	}
}

type scaleReplaySeed struct {
	EventCount int
}

func seedScaleReplayWorkspace(tb testing.TB, workspace string, targetEvents int) scaleReplaySeed {
	tb.Helper()

	if targetEvents <= 0 {
		tb.Fatalf("target events = %d, want positive", targetEvents)
	}

	journalPath := filepath.Join(workspace, ".fid", "journal", "2026-06-07.jsonl")
	events := buildScaleReplayEvents(targetEvents)
	writeReplayEventsTB(tb, journalPath, events...)

	return scaleReplaySeed{EventCount: len(events)}
}

func buildScaleReplayEvents(targetEvents int) []Event {
	events := make([]Event, 0, targetEvents)
	labels := []string{"perf", "scale", "benchmark"}

	appendEvent := func(event Event) {
		if len(events) < targetEvents {
			events = append(events, event)
		}
	}

	nextSequence := 1
	appendEvent(scaleRoleSeedEvent(nextSequence))
	nextSequence++

	for _, label := range labels {
		appendEvent(scaleLabelSeedEvent(nextSequence, label))
		nextSequence++
	}

	for issueIndex := 1; len(events) < targetEvents; issueIndex++ {
		appendEvent(scaleIssueCreatedEvent(nextSequence, issueIndex))
		nextSequence++
		if len(events) >= targetEvents {
			break
		}

		if issueIndex%3 == 0 {
			appendEvent(scaleAcceptanceCriterionEvent(nextSequence, issueIndex))
			nextSequence++
		}
		if len(events) >= targetEvents {
			break
		}

		if issueIndex%4 == 0 {
			appendEvent(scaleIssueLabelEvent(nextSequence, issueIndex, labels[issueIndex%len(labels)]))
			nextSequence++
		}
		if len(events) >= targetEvents {
			break
		}

		if issueIndex%5 == 0 {
			appendEvent(scaleIssueTransitionEvent(nextSequence, issueIndex))
			nextSequence++
		}
		if len(events) >= targetEvents {
			break
		}

		if issueIndex > 1 && issueIndex%7 == 0 {
			appendEvent(scaleIssueLinkEvent(nextSequence, issueIndex))
			nextSequence++
		}
	}

	return events
}

func scaleRoleSeedEvent(sequence int) Event {
	return scaleEvent(
		sequence,
		"project.initialized",
		[]string{"fid", "init"},
		[]Effect{{
			Op:    "insert",
			Table: "roles",
			PK:    map[string]any{"name": "developer"},
			Values: map[string]any{
				"description": "developer",
				"active":      true,
				"is_human":    false,
			},
		}},
	)
}

func scaleLabelSeedEvent(sequence int, label string) Event {
	return scaleEvent(
		sequence,
		"label.created",
		[]string{"fid", "label", "add", label, "--by", "developer"},
		[]Effect{{
			Op:    "insert",
			Table: "labels",
			PK:    map[string]any{"name": label},
			Values: map[string]any{
				"colour":      fmt.Sprintf("#%06d", sequence),
				"description": "scale seed " + label,
			},
		}},
	)
}

func scaleIssueCreatedEvent(sequence, issueIndex int) Event {
	occurredAt := scaleEventTime(sequence)
	issueID := scaleIssueID(issueIndex)

	return scaleEvent(
		sequence,
		"issue.created",
		[]string{"fid", "issue", "create", "--kind", "task", "--title", "Scale task", "--by", "developer"},
		[]Effect{{
			Op:    "insert",
			Table: "issues",
			PK:    map[string]any{"id": issueID},
			Values: map[string]any{
				"kind":        "task",
				"title":       fmt.Sprintf("Scale task %04d", issueIndex),
				"workflow_id": "task-default",
				"status":      "todo",
				"priority":    issueIndex % 5,
				"owner_role":  "developer",
				"created_at":  occurredAt.Format(time.RFC3339),
				"updated_at":  occurredAt.Format(time.RFC3339),
			},
		}},
	)
}

func scaleAcceptanceCriterionEvent(sequence, issueIndex int) Event {
	return scaleEvent(
		sequence,
		"issue.ac.added",
		[]string{"fid", "issue", "ac", "add", scaleIssueID(issueIndex), "--by", "developer"},
		[]Effect{{
			Op:    "insert",
			Table: "acceptance_criteria",
			PK:    map[string]any{"id": fmt.Sprintf("AC-scale-%05d", issueIndex)},
			Values: map[string]any{
				"issue_id": scaleIssueID(issueIndex),
				"ordinal":  1,
				"text":     fmt.Sprintf("Acceptance criterion for scale task %04d", issueIndex),
			},
		}},
	)
}

func scaleIssueLabelEvent(sequence, issueIndex int, label string) Event {
	return scaleEvent(
		sequence,
		"issue.label.added",
		[]string{"fid", "issue", "label", "add", scaleIssueID(issueIndex), label, "--by", "developer"},
		[]Effect{{
			Op:    "insert",
			Table: "issue_labels",
			PK: map[string]any{
				"issue_id": scaleIssueID(issueIndex),
				"label":    label,
			},
			Values: map[string]any{},
		}},
	)
}

func scaleIssueTransitionEvent(sequence, issueIndex int) Event {
	occurredAt := scaleEventTime(sequence)

	return scaleEvent(
		sequence,
		"issue.transitioned",
		[]string{"fid", "issue", "transition", scaleIssueID(issueIndex), "--to", "in_progress", "--by", "developer"},
		[]Effect{{
			Op:    "update",
			Table: "issues",
			PK:    map[string]any{"id": scaleIssueID(issueIndex)},
			Values: map[string]any{
				"status":     "in_progress",
				"updated_at": occurredAt.Format(time.RFC3339),
			},
			Prior: map[string]any{
				"status": "todo",
			},
			Preconditions: map[string]any{
				"status":      "todo",
				"workflow_id": "task-default",
			},
		}},
	)
}

func scaleIssueLinkEvent(sequence, issueIndex int) Event {
	occurredAt := scaleEventTime(sequence)

	return scaleEvent(
		sequence,
		"issue.linked",
		[]string{"fid", "link", "add", scaleIssueID(issueIndex), "--relation", "blocks", "--target", scaleIssueID(issueIndex - 1), "--by", "developer"},
		[]Effect{{
			Op:    "insert",
			Table: "links",
			PK: map[string]any{
				"from_type": "issue",
				"from_id":   scaleIssueID(issueIndex),
				"relation":  "blocks",
				"to_type":   "issue",
				"to_id":     scaleIssueID(issueIndex - 1),
			},
			Values: map[string]any{
				"created_at": occurredAt.Format(time.RFC3339),
			},
		}},
	)
}

func scaleEvent(sequence int, eventType string, argv []string, effects []Effect) Event {
	return Event{
		EventID:     fmt.Sprintf("EVT-20260607-%06d", sequence),
		EventSchema: EventSchemaV1,
		EventType:   eventType,
		OccurredAt:  scaleEventTime(sequence),
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       argv,
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: effects,
	}
}

func scaleEventTime(sequence int) time.Time {
	return time.Date(2026, time.June, 7, 16, 0, 0, 0, time.UTC).Add(time.Duration(sequence) * time.Second)
}

func scaleIssueID(issueIndex int) string {
	return fmt.Sprintf("FID-scale-%05d", issueIndex)
}

func writeReplayEventsTB(tb testing.TB, path string, events ...Event) {
	tb.Helper()

	var payload []byte
	for _, event := range events {
		line, err := marshalReplayEvent(event)
		if err != nil {
			tb.Fatalf("marshalReplayEvent() error = %v", err)
		}
		payload = append(payload, line...)
		payload = append(payload, '\n')
	}

	if err := os.WriteFile(path, payload, 0o644); err != nil {
		tb.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func snapshotScaleProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()

	queries := []struct {
		name  string
		query string
	}{
		{
			name:  "roles",
			query: `SELECT name, COALESCE(description, 'NULL'), active, is_human FROM roles ORDER BY name`,
		},
		{
			name: "issues",
			query: `
				SELECT id, kind, title, workflow_id, status, COALESCE(priority, -1), COALESCE(owner_role, 'NULL'), COALESCE(created_at, 'NULL'), COALESCE(updated_at, 'NULL')
				FROM issues
				ORDER BY id
			`,
		},
		{
			name:  "acceptance_criteria",
			query: `SELECT id, issue_id, ordinal, text FROM acceptance_criteria ORDER BY id`,
		},
		{
			name:  "labels",
			query: `SELECT name, COALESCE(colour, 'NULL'), COALESCE(description, 'NULL') FROM labels ORDER BY name`,
		},
		{
			name:  "issue_labels",
			query: `SELECT issue_id, label FROM issue_labels ORDER BY issue_id, label`,
		},
		{
			name:  "links",
			query: `SELECT from_type, from_id, relation, to_type, to_id, COALESCE(created_at, 'NULL') FROM links ORDER BY from_type, from_id, relation, to_type, to_id`,
		},
		{
			name:  "issues_enriched",
			query: `SELECT id, concerns_raised, open_concerns_count, COALESCE(triage_status, 'NULL') FROM issues_enriched ORDER BY id`,
		},
		{
			name:  "projection_meta",
			query: `SELECT key, COALESCE(journal_file, 'NULL'), journal_offset, COALESCE(last_event_id, 'NULL'), COALESCE(last_event_hash, 'NULL'), event_schema, projection_schema_version, replay_status, updated_at FROM projection_meta ORDER BY key`,
		},
	}

	parts := make([]string, 0, len(queries))
	for _, spec := range queries {
		parts = append(parts, dumpPhaseTwoQuery(t, conn, spec.name, spec.query))
	}

	return strings.Join(parts, "\n")
}
