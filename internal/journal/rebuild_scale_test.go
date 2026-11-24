package journal

import (
	"bufio"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestRebuildWorkspaceIsDeterministicAtSingleProjectScale(t *testing.T) {
	t.Parallel()

	const targetEvents = 1000

	workspace := newReplayWorkspace(t)
	seed := seedScaleReplayWorkspace(t, workspace, targetEvents)

	dbPath := filepath.Join(config.LedgerDir(workspace), "fid.sqlite")
	conn := openReplayDB(t, dbPath)

	stats, err := ReplayJournal(conn, filepath.Join(config.LedgerDir(workspace), "journal"), nil)
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

	rebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsFS(t))
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

	secondRebuildStats, err := RebuildWorkspace(workspace, testProjectMigrationsFS(t))
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
	t.Parallel()
	count := 0
	forEachScaleReplayEvent(t, 100, func(event Event) {
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

	journalPath := filepath.Join(config.LedgerDir(workspace), "journal", "2026-06-07.jsonl")
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
	forEachScaleReplayEvent(tb, targetEvents, func(event Event) {
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

func forEachScaleReplayEvent(tb testing.TB, targetEvents int, visit func(Event)) {
	tb.Helper()
	templates := scaleWriterTemplates(tb)
	revisions := map[string]int{}
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
		if !appendEvent(scaleTemplateEvent(tb, templates, revisions, nextSequence, "label.added", 0, 1, label)) {
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
			if !appendEvent(scaleTemplateEvent(tb, templates, revisions, nextSequence, "issue.ac-added", issueIndex, nextScaleIssueSequence(issueSequences, issueIndex), "")) {
				return
			}
			nextSequence++
		}

		if issueIndex%4 == 0 {
			if !appendEvent(scaleTemplateEvent(tb, templates, revisions, nextSequence, "issue.label.added", issueIndex, nextScaleIssueSequence(issueSequences, issueIndex), labels[issueIndex%len(labels)])) {
				return
			}
			nextSequence++
		}

		if issueIndex%5 == 0 {
			if !appendEvent(scaleTemplateEvent(tb, templates, revisions, nextSequence, "issue.transitioned", issueIndex, nextScaleIssueSequence(issueSequences, issueIndex), "")) {
				return
			}
			nextSequence++
		}

		if issueIndex > 1 && issueIndex%7 == 0 {
			if !appendEvent(scaleTemplateEvent(tb, templates, revisions, nextSequence, "link.added", issueIndex, 1, "")) {
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

func scaleEvent(sequence int, eventType string, aggregate Aggregate, effects []Effect) Event {
	event := Event{
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
	// Admission requires each family's whole Effect grammar, as its writer
	// produces it: the Issue's own update, its activity and the retry result.
	at := scaleEventTime(sequence).Format(time.RFC3339)
	if eventType == "issue.ac-added" || eventType == "issue.label-added" {
		event.Effects = append(event.Effects, Effect{Op: "update", Table: "issues", PK: map[string]any{"id": aggregate.ID}, Values: map[string]any{"updated_at": at}})
	}
	if strings.HasPrefix(eventType, "issue.") || strings.HasPrefix(eventType, "link.") {
		issue := aggregate.ID
		if eventType == "link.added" {
			issue, _ = effects[0].PK["from_id"].(string)
		}
		event.Effects = append(event.Effects, Effect{Op: "activity", Table: "activity_log", PK: map[string]any{"id": fmt.Sprintf("ACT-scale-%06d", sequence)}, Values: map[string]any{
			"journal_event_id": event.EventID, "ref_type": "issue", "ref_id": issue, "event": eventType, "details": nil, "actor": "service:actor-benchmark", "occurred_at": at}})
	}
	if eventType == "label.added" {
		request := IdempotencyRequest{OperationScope: "benchmark.replay:project-scale", IdempotencyKey: fmt.Sprintf("scale-%06d", sequence), RequestDigest: "sha256:" + strings.Repeat("0", 64)}
		event.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, RequestDigest: request.RequestDigest, Scope: request.OperationScope}
		event.Effects = append(event.Effects, idempotencyResultEffect(request, event, aggregate.ID, `{"replayed":false}`))
	}
	return event
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
			query: `SELECT id, COALESCE(triage_status, 'NULL') FROM issues_enriched ORDER BY id`,
		},
		{
			name:  "projection_meta",
			query: `SELECT key, COALESCE(journal_file, 'NULL'), journal_offset, COALESCE(last_event_id, 'NULL'), COALESCE(last_event_hash, 'NULL'), event_schema, projection_schema_version, replay_status, updated_at FROM projection_meta ORDER BY key`,
		},
	}

	parts := make([]string, 0, len(queries))
	for _, spec := range queries {
		parts = append(parts, dumpProjectionQuery(t, conn, spec.name, spec.query))
	}

	return strings.Join(parts, "\n")
}

//go:embed testdata/row-writer-templates.json
var scaleWriterJSON []byte

type scaleTemplates struct {
	First, Second string
	Templates     map[string]Event
}

func scaleWriterTemplates(tb testing.TB) scaleTemplates {
	tb.Helper()
	var templates scaleTemplates
	if err := json.Unmarshal(scaleWriterJSON, &templates); err != nil {
		tb.Fatal(err)
	}
	return templates
}

// The templates were captured from issue/label/link commands in a guarded
// temporary repository. Scale changes identifiers, clocks, label/criterion text and observed
// revisions; the owning commands supply each Effect's shape and order.
func scaleTemplateEvent(tb testing.TB, templates scaleTemplates, revisions map[string]int, sequence int, family string, issueIndex, aggregateSequence int, label string) Event {
	tb.Helper()
	raw, err := MarshalCanonicalJSON(templates.Templates[family])
	if err != nil {
		tb.Fatal(err)
	}
	text := strings.NewReplacer(templates.First, scaleIssueID(issueIndex), templates.Second, scaleIssueID(issueIndex-1)).Replace(string(raw))
	if label != "" {
		text = strings.ReplaceAll(text, "perf", label)
	}
	var e Event
	if err := e.UnmarshalJSON([]byte(text)); err != nil {
		tb.Fatal(err)
	}
	e.EventID = fmt.Sprintf("EVT-20260607-%06d", sequence)
	e.OccurredAt = scaleEventTime(sequence)
	e.Aggregate.Sequence = aggregateSequence
	e.Idempotency = nil
	at := e.OccurredAt.Format(time.RFC3339)
	effects := make([]Effect, 0, len(e.Effects))
	for _, fx := range e.Effects {
		if fx.Table == "idempotency_results" {
			continue
		}
		if fx.Values == nil {
			fx.Values = map[string]any{}
		}
		for _, column := range []string{"created_at", "updated_at", "occurred_at", "decided_at"} {
			if _, has := fx.Values[column]; has {
				fx.Values[column] = at
			}
		}
		switch fx.Table {
		case "issues":
			if fx.Op == "update" {
				if _, changes := fx.Values["revision"]; changes {
					id := fx.PK["id"].(string)
					prior := revisions[id]
					if prior == 0 {
						prior = 1
					}
					fx.Values["revision"] = prior + 1
					fx.Prior["revision"] = prior
					fx.Preconditions["revision"] = prior
					revisions[id] = prior + 1
				}
			}
		case "labels":
			fx.Values["colour"] = fmt.Sprintf("#%06d", sequence)
			fx.Values["description"] = "scale seed " + label
		case "acceptance_criteria":
			fx.Values["text"] = fmt.Sprintf("Acceptance criterion for scale task %04d", issueIndex)
			fx.PK["id"] = fmt.Sprintf("AC-scale-%05d", issueIndex)
		case "person_decisions":
			fx.PK["id"] = fmt.Sprintf("PDN-%06x", sequence)
		case "activity_log":
			fx.PK["id"] = fmt.Sprintf("ACT-scale-%06d", sequence)
			fx.Values["journal_event_id"] = e.EventID
			if details, ok := fx.Values["details"].(map[string]any); ok {
				if _, has := details["ac_id"]; has {
					details["ac_id"] = fmt.Sprintf("AC-scale-%05d", issueIndex)
				}
			}
		}
		effects = append(effects, fx)
	}
	e.Effects = effects
	if family == "label.added" {
		request := IdempotencyRequest{OperationScope: e.CommandName + ":project-scale", IdempotencyKey: fmt.Sprintf("scale-%06d", sequence), RequestDigest: "sha256:" + strings.Repeat("0", 64)}
		e.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest}
		e.Effects = append(e.Effects, idempotencyResultEffect(request, e, e.Aggregate.ID, `{"replayed":false}`))
	}
	return e
}
