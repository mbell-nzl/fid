package maintenance

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"
	"time"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

type structuralFixture struct {
	db          *sql.DB
	journalDir  string
	journalPath string
}

func newStructuralFixture(t *testing.T) structuralFixture {
	t.Helper()

	fixture := newMaintenanceFixture(t)
	seedJournalEvent(t, fixture.journalPath, "EVT-20260715-000001")
	seedImpossibleStructuralRows(t, fixture.db)
	return fixture
}

func newValidStructuralFixture(t *testing.T) structuralFixture {
	t.Helper()

	fixture := newMaintenanceFixture(t)
	seedJournalEvent(t, fixture.journalPath, "EVT-20260715-000001")
	seedValidStructuralRows(t, fixture.db)
	return fixture
}

func newMaintenanceFixture(t *testing.T) structuralFixture {
	t.Helper()

	root := t.TempDir()
	dbPath := filepath.Join(root, "fid.sqlite")
	if err := fiddb.MigrateUp(dbPath, maintenanceMigrationsDir(t)); err != nil {
		t.Fatalf("MigrateUp() error = %v", err)
	}
	db, err := fiddb.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close() error = %v", err)
		}
	})

	journalDir := filepath.Join(root, ".fid", "journal")
	if err := os.MkdirAll(journalDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", journalDir, err)
	}
	return structuralFixture{
		db:          db,
		journalDir:  journalDir,
		journalPath: filepath.Join(journalDir, "2026-07-15.jsonl"),
	}
}

func (fixture structuralFixture) journalBytes(t *testing.T) []byte {
	t.Helper()

	contents, err := os.ReadFile(fixture.journalPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", fixture.journalPath, err)
	}
	return contents
}

func (fixture structuralFixture) changeCount(t *testing.T) int {
	t.Helper()

	var count int
	if err := fixture.db.QueryRow(`SELECT total_changes()`).Scan(&count); err != nil {
		t.Fatalf("query total_changes: %v", err)
	}
	return count
}

func seedImpossibleStructuralRows(t *testing.T, db *sql.DB) {
	t.Helper()

	dropStructuralTriggers(t, db)
	setFixturePragmas(t, db, false, true)
	defer setFixturePragmas(t, db, true, false)

	seedIssue(t, db, "FID-parent-epic", "epic", nil, nil, "epic-default", "todo")
	seedIssue(t, db, "FID-parent-task", "task", nil, nil, "task-default", "todo")
	seedIssue(t, db, "FID-invalid-epic", "epic", stringPtr("FID-parent-task"), nil, "epic-default", "todo")
	seedIssue(t, db, "FID-invalid-task-parent", "task", stringPtr("FID-parent-task"), nil, "task-default", "todo")
	seedIssue(t, db, "FID-invalid-subtask-parent", "subtask", stringPtr("FID-parent-epic"), nil, "subtask-default", "todo")
	seedIssue(t, db, "FID-unattached-sequenced", "task", nil, intPtr(1), "task-default", "todo")
	seedIssue(t, db, "FID-no-metadata-epic", "epic", nil, nil, "epic-default", "todo")
	seedIssue(t, db, "FID-no-metadata-task", "task", stringPtr("FID-no-metadata-epic"), nil, "task-default", "todo")
	seedIssue(t, db, "FID-sequential-epic", "epic", nil, nil, "epic-default", "todo")
	seedEpicMetadata(t, db, "FID-sequential-epic", "sequential")
	seedIssue(t, db, "FID-sequential-task", "task", stringPtr("FID-sequential-epic"), nil, "task-default", "todo")
	seedIssue(t, db, "FID-prioritised-epic", "epic", nil, nil, "epic-default", "todo")
	seedEpicMetadata(t, db, "FID-prioritised-epic", "prioritised")
	seedIssue(t, db, "FID-prioritised-task", "task", stringPtr("FID-prioritised-epic"), intPtr(2), "task-default", "todo")
	seedIssue(t, db, "FID-workflow-state", "task", nil, nil, "task-default", "wrapped")
	seedIssue(t, db, "FID-workflow-unknown", "task", nil, nil, "unknown", "todo")

	seedAcceptanceCriterion(t, db, "AC-invalid", "FID-parent-epic")
	seedCompletion(t, db, "CMP-invalid", "FID-parent-epic")
	seedEpicMetadata(t, db, "FID-parent-task", "prioritised")
	seedDecision(t, db, "DEC-valid", "FID-parent-epic")
	seedDecision(t, db, "DEC-invalid", "FID-parent-task")
	seedEscalation(t, db, "ESC-valid")
	seedTriageEvent(t, db, "TRI-invalid", "FID-parent-task", "FID-parent-task")
	seedRequirement(t, db, "REQ-valid")

	seedLink(t, db, "issue", "FID-missing-link-from", "relates_to", "issue", "FID-parent-task")
	seedLink(t, db, "issue", "FID-parent-task", "relates_to", "requirement", "REQ-missing-link-to")
	seedLink(t, db, "issue", "FID-parent-task", "implements", "issue", "FID-parent-epic")
	seedLink(t, db, "decision", "DEC-valid", "relates_to", "issue", "FID-parent-task")
	seedActivity(t, db, "ACT-invalid-ref", "EVT-20260715-000001", "issue", "FID-missing-activity")
	seedActivity(t, db, "ACT-invalid-type", "EVT-20260715-000001", "unknown", "UNKNOWN-1")
	seedActivity(t, db, "ACT-missing-event", "EVT-20260715-999999", "decision", "DEC-valid")
	seedActivity(t, db, "ACT-missing-event-second", "EVT-20260715-999999", "escalation", "ESC-valid")

	seedForeignKeyViolation(t, db)
}

func seedValidStructuralRows(t *testing.T, db *sql.DB) {
	t.Helper()

	seedIssue(t, db, "FID-epic-sequential", "epic", nil, nil, "epic-default", "todo")
	seedEpicMetadata(t, db, "FID-epic-sequential", "sequential")
	seedIssue(t, db, "FID-task-sequential", "task", stringPtr("FID-epic-sequential"), intPtr(1), "task-default", "todo")
	seedIssue(t, db, "FID-epic-prioritised", "epic", nil, nil, "epic-default", "todo")
	seedEpicMetadata(t, db, "FID-epic-prioritised", "prioritised")
	seedIssue(t, db, "FID-task-prioritised", "task", stringPtr("FID-epic-prioritised"), nil, "task-default", "todo")
	seedIssue(t, db, "FID-task-unattached", "task", nil, nil, "task-default", "todo")
	seedIssue(t, db, "FID-subtask", "subtask", stringPtr("FID-task-sequential"), intPtr(1), "subtask-default", "todo")
	seedAcceptanceCriterion(t, db, "AC-task", "FID-task-sequential")
	seedAcceptanceCriterion(t, db, "AC-subtask", "FID-subtask")
	seedCompletion(t, db, "CMP-task", "FID-task-sequential")
	seedCompletion(t, db, "CMP-subtask", "FID-subtask")
	seedDecision(t, db, "DEC-valid", "FID-epic-sequential")
	seedEscalation(t, db, "ESC-valid")
	seedTriageEvent(t, db, "TRI-valid", "FID-epic-sequential", "FID-task-sequential")
	seedRequirement(t, db, "REQ-valid")
	seedLink(t, db, "issue", "FID-task-sequential", "implements", "requirement", "REQ-valid")
	seedLink(t, db, "requirement", "REQ-valid", "implemented_by", "issue", "FID-task-sequential")
	seedLink(t, db, "issue", "FID-task-prioritised", "relates_to", "requirement", "REQ-valid")
	seedActivity(t, db, "ACT-issue", "EVT-20260715-000001", "issue", "FID-task-sequential")
	seedActivity(t, db, "ACT-requirement", "EVT-20260715-000001", "requirement", "REQ-valid")
	seedActivity(t, db, "ACT-decision", "EVT-20260715-000001", "decision", "DEC-valid")
	seedActivity(t, db, "ACT-escalation", "EVT-20260715-000001", "escalation", "ESC-valid")
	seedAllWorkflowStates(t, db)
}

func dropStructuralTriggers(t *testing.T, db *sql.DB) {
	t.Helper()

	triggers := []string{
		"issues_parent_kind_guard_insert", "issues_parent_kind_guard_update",
		"acceptance_criteria_issue_kind_guard_insert", "task_completions_issue_kind_guard_insert",
		"epic_metadata_issue_kind_guard_insert", "triage_events_epic_kind_guard_insert",
		"decisions_epic_kind_guard_insert", "issues_kind_immutable",
		"links_endpoint_scope_guard_insert",
		"links_from_ref_integrity_guard_insert", "links_to_ref_integrity_guard_insert",
		"links_relation_compatibility_guard_insert", "activity_log_ref_integrity_guard_insert",
		"issues_sequence_scheduling_guard_insert",
		"issues_sequence_scheduling_guard_update",
	}
	for _, trigger := range triggers {
		if _, err := db.Exec("DROP TRIGGER " + trigger); err != nil {
			t.Fatalf("DROP TRIGGER %s error = %v", trigger, err)
		}
	}
}

func setFixturePragmas(t *testing.T, db *sql.DB, foreignKeys, ignoreChecks bool) {
	t.Helper()

	foreignKeyValue := "OFF"
	if foreignKeys {
		foreignKeyValue = "ON"
	}
	checkValue := "OFF"
	if ignoreChecks {
		checkValue = "ON"
	}
	if _, err := db.Exec("PRAGMA foreign_keys = " + foreignKeyValue); err != nil {
		t.Fatalf("set foreign_keys=%s: %v", foreignKeyValue, err)
	}
	if _, err := db.Exec("PRAGMA ignore_check_constraints = " + checkValue); err != nil {
		t.Fatalf("set ignore_check_constraints=%s: %v", checkValue, err)
	}
}

func seedIssue(t *testing.T, db *sql.DB, id, kind string, parentID *string, sequence *int, workflowID, status string) {
	t.Helper()

	if _, err := db.Exec(`
		INSERT INTO issues (id, kind, parent_id, sequence, title, workflow_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, kind, parentID, sequence, id, workflowID, status, "2026-07-15T00:00:00Z", "2026-07-15T00:00:00Z"); err != nil {
		t.Fatalf("insert issue %q: %v", id, err)
	}
}

func seedEpicMetadata(t *testing.T, db *sql.DB, issueID, scheduling string) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO epic_metadata (issue_id, scheduling) VALUES (?, ?)`, issueID, scheduling); err != nil {
		t.Fatalf("insert epic_metadata %q: %v", issueID, err)
	}
}

func seedAcceptanceCriterion(t *testing.T, db *sql.DB, id, issueID string) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO acceptance_criteria (id, issue_id, ordinal, text) VALUES (?, ?, 1, 'criterion')`, id, issueID); err != nil {
		t.Fatalf("insert acceptance criterion %q: %v", id, err)
	}
}

func seedCompletion(t *testing.T, db *sql.DB, id, issueID string) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO task_completions (id, issue_id, attempt_number) VALUES (?, ?, 1)`, id, issueID); err != nil {
		t.Fatalf("insert completion %q: %v", id, err)
	}
}

func seedDecision(t *testing.T, db *sql.DB, id, epicID string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO decisions (id, epic_id, author_role, summary, body_md_path, status)
		VALUES (?, ?, 'tl', 'decision', 'decision.md', 'active')
	`, id, epicID); err != nil {
		t.Fatalf("insert decision %q: %v", id, err)
	}
}

func seedEscalation(t *testing.T, db *sql.DB, id string) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO escalations (id, status) VALUES (?, 'open')`, id); err != nil {
		t.Fatalf("insert escalation %q: %v", id, err)
	}
}

func seedTriageEvent(t *testing.T, db *sql.DB, id, epicID, issueID string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO triage_events (id, epic_id, issue_id, triager_role, source, verdict)
		VALUES (?, ?, ?, 'tl_lite', 'manual', 'resume')
	`, id, epicID, issueID); err != nil {
		t.Fatalf("insert triage event %q: %v", id, err)
	}
}

func seedRequirement(t *testing.T, db *sql.DB, id string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO requirements (id, type, summary, source_doc)
		VALUES (?, 'functional', 'requirement', '01_PRD.md')
	`, id); err != nil {
		t.Fatalf("insert requirement %q: %v", id, err)
	}
}

func seedLink(t *testing.T, db *sql.DB, fromType, fromID, relation, toType, toID string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO links (from_type, from_id, relation, to_type, to_id)
		VALUES (?, ?, ?, ?, ?)
	`, fromType, fromID, relation, toType, toID); err != nil {
		t.Fatalf("insert link %s:%s: %v", fromType, fromID, err)
	}
}

func seedActivity(t *testing.T, db *sql.DB, id, journalEventID, refType, refID string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO activity_log (id, journal_event_id, ref_type, ref_id, event)
		VALUES (?, ?, ?, ?, 'test')
	`, id, journalEventID, refType, refID); err != nil {
		t.Fatalf("insert activity %q: %v", id, err)
	}
}

func seedForeignKeyViolation(t *testing.T, db *sql.DB) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO issue_labels (issue_id, label) VALUES ('FID-missing-fk', 'LABEL-missing-fk')`); err != nil {
		t.Fatalf("insert foreign key violation: %v", err)
	}
}

func seedAllWorkflowStates(t *testing.T, db *sql.DB) {
	t.Helper()

	states := map[string][]string{
		"task-default":    {"todo", "in_progress", "dev_complete", "qa_review", "qa_block", "failed_gate", "qa_pass", "done", "cancelled"},
		"epic-default":    {"todo", "in_progress", "wrapped", "closed", "cancelled"},
		"subtask-default": {"todo", "in_progress", "done", "cancelled"},
	}
	for workflowID, workflowStates := range states {
		kind := workflowID[:len(workflowID)-len("-default")]
		for index, state := range workflowStates {
			id := workflowID + "-" + state
			if kind == "subtask" {
				seedIssue(t, db, id, kind, stringPtr("FID-task-sequential"), intPtr(index+2), workflowID, state)
			} else {
				seedIssue(t, db, id, kind, nil, nil, workflowID, state)
			}
		}
	}
}

func seedJournalEvent(t *testing.T, journalPath, eventID string) {
	t.Helper()

	event := journal.Event{
		EventID:     eventID,
		EventSchema: journal.EventSchemaV1,
		EventType:   "test.seeded",
		OccurredAt:  time.Date(2026, time.July, 15, 0, 0, 0, 0, time.UTC),
		Actor:       journal.Actor{Role: "developer"},
		Command:     journal.Command{Argv: []string{"fid", "test"}, CWD: ".", FIDVersion: "0.1.0"},
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal journal event: %v", err)
	}
	if err := os.WriteFile(journalPath, append(encoded, '\n'), 0o644); err != nil {
		t.Fatalf("write journal event: %v", err)
	}
}

func findingChecks(findings []Finding) []string {
	checks := make([]string, len(findings))
	for index, finding := range findings {
		checks[index] = finding.Check
	}
	return checks
}

func assertFindingsOrdered(t *testing.T, findings []Finding) {
	t.Helper()

	sorted := append([]Finding(nil), findings...)
	sort.Slice(sorted, func(left, right int) bool {
		if sorted[left].Check != sorted[right].Check {
			return sorted[left].Check < sorted[right].Check
		}
		if sorted[left].EntityID != sorted[right].EntityID {
			return sorted[left].EntityID < sorted[right].EntityID
		}
		return sorted[left].Message < sorted[right].Message
	})
	if !reflect.DeepEqual(findings, sorted) {
		t.Fatalf("findings are not ordered: %#v", findings)
	}
}

func maintenanceMigrationsDir(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller(0) failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "migrations"))
}

func stringPtr(value string) *string { return &value }

func intPtr(value int) *int { return &value }

func newWarningFixture(t *testing.T, now time.Time) structuralFixture {
	t.Helper()

	fixture := newMaintenanceFixture(t)
	seedJournalEvent(t, fixture.journalPath, "EVT-20260715-000001")
	seedIssueAt(t, fixture.db, "FID-stale-claim", "task", "todo", now)
	seedIssueAt(t, fixture.db, "FID-live-claim", "task", "todo", now)
	seedIssueAt(t, fixture.db, "FID-start-live", "task", "in_progress", now.Add(-time.Hour))
	seedIssueAt(t, fixture.db, "FID-start-before", "task", "in_progress", now.Add(-299*time.Second))
	seedIssueAt(t, fixture.db, "FID-start-at", "task", "in_progress", now)
	seedIssueAt(t, fixture.db, "FID-start-after", "task", "in_progress", now)
	seedIssueAt(t, fixture.db, "FID-start-fallback", "task", "in_progress", now.Add(-300*time.Second))
	seedIssueAt(t, fixture.db, "FID-start-non-status-update", "task", "in_progress", now.Add(-time.Hour))
	seedIssueAt(t, fixture.db, "FID-start-zero", "task", "in_progress", now)
	seedIssueAt(t, fixture.db, "FID-completion", "task", "todo", now)

	seedClaimAt(t, fixture.db, "FID-stale-claim", now.Add(-300*time.Second), 300)
	seedClaimAt(t, fixture.db, "FID-live-claim", now.Add(-time.Second), 300)
	seedClaimAt(t, fixture.db, "FID-start-live", now.Add(-time.Second), 300)
	seedActivityAt(t, fixture.db, "ACT-start-at", "FID-start-at", "status_changed", now.Add(-300*time.Second))
	seedActivityAt(t, fixture.db, "ACT-start-after", "FID-start-after", "status_changed", now.Add(-301*time.Second))
	seedActivityAt(t, fixture.db, "ACT-start-update", "FID-start-non-status-update", "status_changed", now.Add(-300*time.Second))
	seedActivityAt(t, fixture.db, "ACT-non-status-update", "FID-start-non-status-update", "issue_updated", now.Add(-time.Second))

	seedCompletionWithStatus(t, fixture.db, "TC-superseded", "FID-completion", 1, "superseded")
	seedCompletionWithStatus(t, fixture.db, "TC-current", "FID-completion", 2, "dev_complete")
	seedCompletionWithStatus(t, fixture.db, "TC-missing-run", "FID-completion", 3, "dev_complete")
	seedCompletionWithStatus(t, fixture.db, "TC-issue-only-run", "FID-completion", 4, "dev_complete")
	seedCompletionWithStatus(t, fixture.db, "TC-linked-run", "FID-completion", 5, "dev_complete")
	seedAcceptanceCriterion(t, fixture.db, "AC-warning", "FID-completion")
	seedVerification(t, fixture.db, "AV-stale", "AC-warning", stringPtr("TC-superseded"))
	seedVerification(t, fixture.db, "AV-current", "AC-warning", stringPtr("TC-current"))
	seedVerification(t, fixture.db, "AV-nil", "AC-warning", nil)
	seedAgentRun(t, fixture.db, "RUN-superseded", "FID-completion", stringPtr("TC-superseded"))
	seedAgentRun(t, fixture.db, "RUN-current", "FID-completion", stringPtr("TC-current"))
	seedAgentRun(t, fixture.db, "RUN-issue-only", "FID-completion", nil)
	seedAgentRun(t, fixture.db, "RUN-linked", "FID-completion", stringPtr("TC-linked-run"))
	return fixture
}

func newInfoFixture(t *testing.T, now time.Time) structuralFixture {
	t.Helper()

	fixture := newMaintenanceFixture(t)
	seedJournalEvent(t, fixture.journalPath, "EVT-20260715-000001")
	seedIssueAt(t, fixture.db, "FID-info-live", "task", "todo", now)
	seedIssueAt(t, fixture.db, "FID-info-stale", "task", "todo", now)
	seedIssueAt(t, fixture.db, "FID-info-activity", "task", "todo", now)
	seedClaimAt(t, fixture.db, "FID-info-live", now.Add(-time.Second), 2)
	seedClaimAt(t, fixture.db, "FID-info-stale", now.Add(-time.Second), 1)
	seedActivityAt(t, fixture.db, "ACT-info-now", "FID-info-activity", "issue_updated", now)
	seedActivityAt(t, fixture.db, "ACT-info-boundary", "FID-info-activity", "issue_updated", now.Add(-7*24*time.Hour))
	seedActivityAt(t, fixture.db, "ACT-info-old", "FID-info-activity", "issue_updated", now.Add(-7*24*time.Hour-time.Second))
	seedActivityAt(t, fixture.db, "ACT-info-future", "FID-info-activity", "issue_updated", now.Add(time.Second))
	return fixture
}

func seedIssueAt(t *testing.T, db *sql.DB, id, kind, status string, createdAt time.Time) {
	t.Helper()
	workflowID := kind + "-default"
	if _, err := db.Exec(`
		INSERT INTO issues (id, kind, title, workflow_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, kind, id, workflowID, status, createdAt.Format(time.RFC3339), createdAt.Format(time.RFC3339)); err != nil {
		t.Fatalf("insert dated issue %q: %v", id, err)
	}
}

func seedClaimAt(t *testing.T, db *sql.DB, issueID string, claimedAt time.Time, ttlSeconds int) {
	t.Helper()
	envelopeID := "ENV-" + issueID
	if _, err := db.Exec(`
		INSERT INTO work_envelopes (id, current_revision, status, created_at)
		VALUES (?, 1, 'active', ?)
	`, envelopeID, claimedAt.Format(time.RFC3339)); err != nil {
		t.Fatalf("insert claim envelope %q: %v", envelopeID, err)
	}
	if _, err := db.Exec(`
		INSERT INTO envelope_revisions (
			envelope_id, revision, outcome, non_goals_json, risk_result_json,
			repository_scope_json, status, created_actor_id, approved_actor_id,
			created_at, activated_at
		) VALUES (?, 1, 'maintenance claim fixture', '[]', '{}', '{}', 'active', 'agent:developer', 'agent:developer', ?, ?)
	`, envelopeID, claimedAt.Format(time.RFC3339), claimedAt.Format(time.RFC3339)); err != nil {
		t.Fatalf("insert claim envelope revision %q: %v", envelopeID, err)
	}
	if _, err := db.Exec(`
		INSERT INTO claims (
			id, issue_id, envelope_id, envelope_revision, actor_id, worktree_id,
			scope_json, mode, version, status, acquired_at, expires_at
		) VALUES (?, ?, ?, 1, 'agent:developer', 'maintenance', '{}', 'mutation', 1, 'active', ?, ?)
	`, "CLM-"+issueID, issueID, envelopeID, claimedAt.Format(time.RFC3339), claimedAt.Add(time.Duration(ttlSeconds)*time.Second).Format(time.RFC3339)); err != nil {
		t.Fatalf("insert dated claim %q: %v", issueID, err)
	}
}

func seedActivityAt(t *testing.T, db *sql.DB, id, issueID, event string, occurredAt time.Time) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO activity_log (id, journal_event_id, ref_type, ref_id, event, occurred_at)
		VALUES (?, 'EVT-20260715-000001', 'issue', ?, ?, ?)
	`, id, issueID, event, occurredAt.Format(time.RFC3339)); err != nil {
		t.Fatalf("insert activity %q: %v", id, err)
	}
}

func seedCompletionWithStatus(t *testing.T, db *sql.DB, id, issueID string, attempt int, status string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO task_completions (id, issue_id, attempt_number, status)
		VALUES (?, ?, ?, ?)
	`, id, issueID, attempt, status); err != nil {
		t.Fatalf("insert completion %q: %v", id, err)
	}
}

func seedVerification(t *testing.T, db *sql.DB, id, acID string, completionID *string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO ac_verifications (id, ac_id, verifier_role, status, completion_id, verified_at)
		VALUES (?, ?, 'qa_reviewer', 'verified', ?, '2026-07-15T12:00:00Z')
	`, id, acID, completionID); err != nil {
		t.Fatalf("insert verification %q: %v", id, err)
	}
}

func seedAgentRun(t *testing.T, db *sql.DB, id, issueID string, completionID *string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO agent_runs (id, issue_id, completion_id, role, model)
		VALUES (?, ?, ?, 'developer', 'gpt-5.4-mini')
	`, id, issueID, completionID); err != nil {
		t.Fatalf("insert agent run %q: %v", id, err)
	}
}

func findingEntityIDs(findings []Finding) []string {
	ids := make([]string, len(findings))
	for index, finding := range findings {
		ids[index] = finding.EntityID
	}
	return ids
}
