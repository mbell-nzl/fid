package journal

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
)

func TestProjectEventPopulatesPhaseTwoRepresentativeTablesViaJournalEffects(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}

	conn, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeProjectDB(t, conn)

	event := Event{
		EventID:     "EVT-20260607-000101",
		EventSchema: EventSchemaV1,
		EventType:   "phase2.schema.seeded",
		OccurredAt:  time.Date(2026, time.June, 7, 9, 30, 0, 0, time.UTC),
		Actor: Actor{
			Role: "developer",
		},
		Command: Command{
			Argv:       []string{"fid", "maintenance", "replay"},
			CWD:        ".",
			FIDVersion: "0.1.0",
		},
		Effects: []Effect{
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-epic01",
				},
				Values: map[string]any{
					"kind":        "epic",
					"title":       "Phase 2 epic",
					"workflow_id": "epic-default",
					"status":      "todo",
					"created_at":  "2026-06-07T09:00:00Z",
					"updated_at":  "2026-06-07T09:00:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "epic_metadata",
				PK: map[string]any{
					"issue_id": "FID-epic01",
				},
				Values: map[string]any{
					"scheduling":    "sequential",
					"brief_md_path": "project-docs/epics/phase-2.md",
					"target_start":  "2026-06-07",
					"target_end":    "2026-06-21",
				},
			},
			{
				Op:    "insert",
				Table: "issues",
				PK: map[string]any{
					"id": "FID-task01",
				},
				Values: map[string]any{
					"kind":        "task",
					"parent_id":   "FID-epic01",
					"sequence":    1,
					"title":       "Implement migration 0002",
					"workflow_id": "task-default",
					"status":      "todo",
					"created_at":  "2026-06-07T09:05:00Z",
					"updated_at":  "2026-06-07T09:05:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "requirements",
				PK: map[string]any{
					"id": "PRD-OPS-003",
				},
				Values: map[string]any{
					"type":             "PRD",
					"summary":          "Rebuild over the widened schema",
					"source_doc":       "01_PRD.md",
					"source_section":   "OPS",
					"last_synced_hash": "sha256:req",
					"status":           "active",
				},
			},
			{
				Op:    "insert",
				Table: "labels",
				PK: map[string]any{
					"name": "schema",
				},
				Values: map[string]any{
					"colour":      "#005f73",
					"description": "Schema work",
				},
			},
			{
				Op:    "insert",
				Table: "acceptance_criteria",
				PK: map[string]any{
					"id": "AC-000001",
				},
				Values: map[string]any{
					"issue_id": "FID-task01",
					"ordinal":  1,
					"text":     "Migration 0002 creates the remaining v0.1 tables",
				},
			},
			{
				Op:    "insert",
				Table: "task_completions",
				PK: map[string]any{
					"id": "TC-000001",
				},
				Values: map[string]any{
					"issue_id":           "FID-task01",
					"attempt_number":     1,
					"developer_role":     "developer",
					"model":              "gpt-5",
					"status":             "dev_complete",
					"gate_results":       `["go test ./..."]`,
					"human_intervention": false,
					"completed_at":       "2026-06-07T09:20:00Z",
				},
			},
			{
				Op:    "insert",
				Table: "concerns",
				PK: map[string]any{
					"id": "CON-000001",
				},
				Values: map[string]any{
					"issue_id":                "FID-task01",
					"raised_in_completion_id": "TC-000001",
					"raised_by_role":          "developer",
					"summary":                 "Double-check down migration order",
					"body_md_path":            "project-docs/concerns/CON-000001.md",
					"status":                  "open",
					"created_at":              "2026-06-07T09:21:00Z",
					"updated_at":              "2026-06-07T09:21:00Z",
				},
			},
		},
	}

	eventHash, err := HashEvent(event)
	if err != nil {
		t.Fatalf("HashEvent() error = %v", err)
	}

	if err := ProjectEvent(conn, ProjectionInput{
		Event:     event,
		EventHash: eventHash,
		File:      "2026-06-07.jsonl",
		Line:      1,
		Offset:    1024,
	}); err != nil {
		t.Fatalf("ProjectEvent() error = %v", err)
	}

	assertRepresentativeRequirement(t, conn)
	assertRepresentativeLabel(t, conn)
	assertRepresentativeEpicMetadata(t, conn)
	assertRepresentativeAcceptanceCriterion(t, conn)
	assertRepresentativeTaskCompletion(t, conn)
	assertRepresentativeConcern(t, conn)
}

func TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion(t *testing.T) {
	t.Parallel()

	workspace := newReplayWorkspace(t)
	writeReplayFixture(t, workspace, "replay-v1.golden.jsonl", "2026-06-05.jsonl")

	if _, err := RebuildWorkspace(workspace, testProjectMigrationsDir(t)); err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}

	conn := openReplayDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	if _, err := conn.Exec(`UPDATE projection_meta SET projection_schema_version = ? WHERE key = 'default'`, 1); err != nil {
		t.Fatalf("downgrade projection_meta schema version: %v", err)
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

	rebuiltConn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() rebuilt error = %v", err)
	}
	defer closeProjectDB(t, rebuiltConn)

	meta, found, err := fiddb.ReadProjectionMeta(rebuiltConn)
	if err != nil {
		t.Fatalf("ReadProjectionMeta() error = %v", err)
	}
	if !found {
		t.Fatal("ReadProjectionMeta() found = false, want true")
	}
	if meta.ProjectionSchemaVersion != fiddb.ProjectionSchemaVersion {
		t.Fatalf("projection_schema_version = %d, want %d", meta.ProjectionSchemaVersion, fiddb.ProjectionSchemaVersion)
	}
}

func assertRepresentativeRequirement(t *testing.T, conn *sql.DB) {
	t.Helper()

	var (
		id             string
		reqType        string
		summary        string
		sourceDoc      string
		sourceSection  string
		lastSyncedHash string
		status         string
	)
	if err := conn.QueryRow(`
		SELECT id, type, summary, source_doc, source_section, last_synced_hash, status
		FROM requirements
		WHERE id = ?
	`, "PRD-OPS-003").Scan(&id, &reqType, &summary, &sourceDoc, &sourceSection, &lastSyncedHash, &status); err != nil {
		t.Fatalf("query requirements: %v", err)
	}
	if id != "PRD-OPS-003" || reqType != "PRD" || summary != "Rebuild over the widened schema" || sourceDoc != "01_PRD.md" || sourceSection != "OPS" || lastSyncedHash != "sha256:req" || status != "active" {
		t.Fatalf("requirements row = %#v", []string{id, reqType, summary, sourceDoc, sourceSection, lastSyncedHash, status})
	}
}

func assertRepresentativeLabel(t *testing.T, conn *sql.DB) {
	t.Helper()

	var name, colour, description string
	if err := conn.QueryRow(`
		SELECT name, colour, description
		FROM labels
		WHERE name = ?
	`, "schema").Scan(&name, &colour, &description); err != nil {
		t.Fatalf("query labels: %v", err)
	}
	if name != "schema" || colour != "#005f73" || description != "Schema work" {
		t.Fatalf("labels row = %#v", []string{name, colour, description})
	}
}

func assertRepresentativeEpicMetadata(t *testing.T, conn *sql.DB) {
	t.Helper()

	var issueID, scheduling, briefPath, startDate, endDate string
	if err := conn.QueryRow(`
		SELECT issue_id, scheduling, brief_md_path, date(target_start), date(target_end)
		FROM epic_metadata
		WHERE issue_id = ?
	`, "FID-epic01").Scan(&issueID, &scheduling, &briefPath, &startDate, &endDate); err != nil {
		t.Fatalf("query epic_metadata: %v", err)
	}
	if issueID != "FID-epic01" || scheduling != "sequential" || briefPath != "project-docs/epics/phase-2.md" || startDate != "2026-06-07" || endDate != "2026-06-21" {
		t.Fatalf("epic_metadata row = %#v", []string{issueID, scheduling, briefPath, startDate, endDate})
	}
}

func assertRepresentativeAcceptanceCriterion(t *testing.T, conn *sql.DB) {
	t.Helper()

	var id, issueID, text string
	var ordinal int
	if err := conn.QueryRow(`
		SELECT id, issue_id, ordinal, text
		FROM acceptance_criteria
		WHERE id = ?
	`, "AC-000001").Scan(&id, &issueID, &ordinal, &text); err != nil {
		t.Fatalf("query acceptance_criteria: %v", err)
	}
	if id != "AC-000001" || issueID != "FID-task01" || ordinal != 1 || text != "Migration 0002 creates the remaining v0.1 tables" {
		t.Fatalf("acceptance_criteria row = %#v", []any{id, issueID, ordinal, text})
	}
}

func assertRepresentativeTaskCompletion(t *testing.T, conn *sql.DB) {
	t.Helper()

	var id, issueID, developerRole, model, status, gateResults string
	var attemptNumber int
	var humanIntervention bool
	if err := conn.QueryRow(`
		SELECT id, issue_id, attempt_number, developer_role, model, status, gate_results, human_intervention
		FROM task_completions
		WHERE id = ?
	`, "TC-000001").Scan(&id, &issueID, &attemptNumber, &developerRole, &model, &status, &gateResults, &humanIntervention); err != nil {
		t.Fatalf("query task_completions: %v", err)
	}
	if id != "TC-000001" || issueID != "FID-task01" || attemptNumber != 1 || developerRole != "developer" || model != "gpt-5" || status != "dev_complete" || gateResults != `["go test ./..."]` || humanIntervention {
		t.Fatalf("task_completions row = %#v", []any{id, issueID, attemptNumber, developerRole, model, status, gateResults, humanIntervention})
	}
}

func assertRepresentativeConcern(t *testing.T, conn *sql.DB) {
	t.Helper()

	var id, issueID, completionID, raisedByRole, summary, status string
	if err := conn.QueryRow(`
		SELECT id, issue_id, raised_in_completion_id, raised_by_role, summary, status
		FROM concerns
		WHERE id = ?
	`, "CON-000001").Scan(&id, &issueID, &completionID, &raisedByRole, &summary, &status); err != nil {
		t.Fatalf("query concerns: %v", err)
	}
	if id != "CON-000001" || issueID != "FID-task01" || completionID != "TC-000001" || raisedByRole != "developer" || summary != "Double-check down migration order" || status != "open" {
		t.Fatalf("concerns row = %#v", []string{id, issueID, completionID, raisedByRole, summary, status})
	}
}
