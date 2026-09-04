package journal

import (
	"bufio"
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

	rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t))
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

	secondRebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t))
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

func TestScaleReplayFixtureUsesRoleFreeEventV2Facts(t *testing.T) {
	count := 0
	forEachScaleReplayEvent(100, func(event Event) {
		count++
		if !event.IsV2() {
			t.Fatalf("scale Event %q schema = %d, want Event-v2", event.EventID, event.EventSchema)
		}
		if event.Actor.Role != "" || len(event.Command.Argv) != 0 {
			t.Fatalf("scale Event %q retains role attribution", event.EventID)
		}
		for _, effect := range event.Effects {
			if effect.Table == "roles" || effect.Table == "work_claims" || effect.Table == "gate_runs" {
				t.Fatalf("scale Event %q writes removed table %q", event.EventID, effect.Table)
			}
		}
	})
	if count != 100 {
		t.Fatalf("scale Event count = %d, want 100", count)
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
	eventCount := writeScaleReplayWorkspace(tb, journalPath, targetEvents)

	return scaleReplaySeed{EventCount: eventCount}
}

func writeScaleReplayWorkspace(tb testing.TB, path string, targetEvents int) int {
	tb.Helper()

	file, err := os.Create(path)
	if err != nil {
		tb.Fatalf("create scale replay journal %q: %v", path, err)
	}
	writer := bufio.NewWriter(file)
	count := 0
	forEachScaleReplayEvent(targetEvents, func(event Event) {
		line, marshalErr := marshalReplayEvent(event)
		if marshalErr != nil {
			_ = file.Close()
			tb.Fatalf("marshal scale replay event: %v", marshalErr)
		}
		if _, writeErr := writer.Write(append(line, '\n')); writeErr != nil {
			_ = file.Close()
			tb.Fatalf("write scale replay event: %v", writeErr)
		}
		count++
	})
	if err := writer.Flush(); err != nil {
		_ = file.Close()
		tb.Fatalf("flush scale replay journal: %v", err)
	}
	if err := file.Close(); err != nil {
		tb.Fatalf("close scale replay journal: %v", err)
	}
	return count
}

func forEachScaleReplayEvent(targetEvents int, visit func(Event)) {
	labels := []string{"perf", "scale", "benchmark"}
	eventCount := 0
	issueSequences := make(map[int]int)
	appendEvent := func(event Event) bool {
		if eventCount >= targetEvents {
			return false
		}
		visit(event)
		eventCount++
		return eventCount < targetEvents
	}

	nextSequence := 1
	if !appendEvent(scaleBootstrapEvent(nextSequence)) {
		return
	}
	nextSequence++

	for _, label := range labels {
		if !appendEvent(scaleLabelSeedEvent(nextSequence, label)) {
			return
		}
		nextSequence++
	}

	for issueIndex := 1; eventCount < targetEvents; issueIndex++ {
		if !appendEvent(scaleIssueCreatedEvent(nextSequence, issueIndex, nextScaleIssueSequence(issueSequences, issueIndex))) {
			return
		}
		nextSequence++

		if issueIndex%3 == 0 {
			if !appendEvent(scaleAcceptanceCriterionEvent(nextSequence, issueIndex, nextScaleIssueSequence(issueSequences, issueIndex))) {
				return
			}
			nextSequence++
		}

		if issueIndex%4 == 0 {
			if !appendEvent(scaleIssueLabelEvent(nextSequence, issueIndex, labels[issueIndex%len(labels)], nextScaleIssueSequence(issueSequences, issueIndex))) {
				return
			}
			nextSequence++
		}

		if issueIndex%5 == 0 {
			if !appendEvent(scaleIssueTransitionEvent(nextSequence, issueIndex, nextScaleIssueSequence(issueSequences, issueIndex))) {
				return
			}
			nextSequence++
		}

		if issueIndex > 1 && issueIndex%7 == 0 {
			if !appendEvent(scaleIssueLinkEvent(nextSequence, issueIndex)) {
				return
			}
			nextSequence++
		}
	}
}

func nextScaleIssueSequence(sequences map[int]int, issueIndex int) int {
	sequences[issueIndex]++
	return sequences[issueIndex]
}

func scaleBootstrapEvent(sequence int) Event {
	return scaleEvent(sequence, "project.initialized", Aggregate{ID: "project-scale", Type: "project", Sequence: 1}, []Effect{})
}

func scaleLabelSeedEvent(sequence int, label string) Event {
	return scaleEvent(
		sequence,
		"label.added",
		Aggregate{ID: "LABEL-" + label, Type: "label", Sequence: 1},
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

func scaleIssueCreatedEvent(sequence, issueIndex, aggregateSequence int) Event {
	occurredAt := scaleEventTime(sequence)
	issueID := scaleIssueID(issueIndex)

	return scaleEvent(
		sequence,
		"issue.created",
		Aggregate{ID: issueID, Type: "issue", Sequence: aggregateSequence},
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
				"created_at":  occurredAt.Format(time.RFC3339),
				"updated_at":  occurredAt.Format(time.RFC3339),
			},
		}},
	)
}

func scaleAcceptanceCriterionEvent(sequence, issueIndex, aggregateSequence int) Event {
	return scaleEvent(
		sequence,
		"issue.ac-added",
		Aggregate{ID: scaleIssueID(issueIndex), Type: "issue", Sequence: aggregateSequence},
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

func scaleIssueLabelEvent(sequence, issueIndex int, label string, aggregateSequence int) Event {
	return scaleEvent(
		sequence,
		"issue.label-added",
		Aggregate{ID: scaleIssueID(issueIndex), Type: "issue", Sequence: aggregateSequence},
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

func scaleIssueTransitionEvent(sequence, issueIndex, aggregateSequence int) Event {
	occurredAt := scaleEventTime(sequence)

	return scaleEvent(
		sequence,
		"issue.transitioned",
		Aggregate{ID: scaleIssueID(issueIndex), Type: "issue", Sequence: aggregateSequence},
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
		"link.added",
		Aggregate{ID: fmt.Sprintf("LINK-scale-%05d", issueIndex), Type: "link", Sequence: 1},
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

func scaleEvent(sequence int, eventType string, aggregate Aggregate, effects []Effect) Event {
	return Event{
		EventID:     fmt.Sprintf("EVT-20260607-%06d", sequence),
		EventSchema: EventSchemaV2,
		EventType:   eventType,
		OccurredAt:  scaleEventTime(sequence),
		ActorV2: EventActor{
			ActorID:       "actor-benchmark",
			ActorType:     "service",
			Producer:      "fid",
			ProvenanceRef: "benchmark:rebuild",
		},
		Aggregate:   aggregate,
		CommandName: "benchmark.replay",
		Effects:     effects,
	}
}

func scaleEventTime(sequence int) time.Time {
	return time.Date(2026, time.June, 7, 16, 0, 0, 0, time.UTC).Add(time.Duration(sequence) * time.Second)
}

func scaleIssueID(issueIndex int) string {
	return fmt.Sprintf("FID-scale-%05d", issueIndex)
}

func snapshotScaleProjectionState(t *testing.T, conn *sql.DB) string {
	t.Helper()

	queries := []struct {
		name  string
		query string
	}{
		{
			name:  "labels",
			query: `SELECT name, COALESCE(description, 'NULL'), COALESCE(colour, 'NULL') FROM labels ORDER BY name`,
		},
		{
			name: "issues",
			query: `
				SELECT id, kind, title, workflow_id, status, COALESCE(priority, -1), COALESCE(created_at, 'NULL'), COALESCE(updated_at, 'NULL')
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
