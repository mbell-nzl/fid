package issues

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestAddAcceptanceCriterionPersistsAndReadsBackInOrdinalOrder(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Implement AC add",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	linesBefore := countJournalLines(journalPathForWorkspace(t, workspace))

	second, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "Second acceptance criterion",
		Ordinal:       2,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion(second) error = %v", err)
	}
	if !strings.HasPrefix(second.ID, "AC-") {
		t.Fatalf("AddAcceptanceCriterion(second) ID = %q, want AC- prefix", second.ID)
	}

	first, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "First acceptance criterion",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion(first) error = %v", err)
	}
	if !strings.HasPrefix(first.ID, "AC-") {
		t.Fatalf("AddAcceptanceCriterion(first) ID = %q, want AC- prefix", first.ID)
	}

	if got := countJournalLines(journalPathForWorkspace(t, workspace)); got != linesBefore+2 {
		t.Fatalf("journal line count = %d, want %d", got, linesBefore+2)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got.AcceptanceCriteria) != 2 {
		t.Fatalf("Get() acceptance_criteria len = %d, want 2", len(got.AcceptanceCriteria))
	}
	if got.AcceptanceCriteria[0].ID != first.ID || got.AcceptanceCriteria[0].Ordinal != 1 || got.AcceptanceCriteria[0].Text != "First acceptance criterion" {
		t.Fatalf("first read-back = %#v, want first ordinal/text/id", got.AcceptanceCriteria[0])
	}
	if got.AcceptanceCriteria[1].ID != second.ID || got.AcceptanceCriteria[1].Ordinal != 2 || got.AcceptanceCriteria[1].Text != "Second acceptance criterion" {
		t.Fatalf("second read-back = %#v, want second ordinal/text/id", got.AcceptanceCriteria[1])
	}

	stored := readIssueStoredEvents(t, journalPathForWorkspace(t, workspace))
	if len(stored) < 2 {
		t.Fatalf("stored events len = %d, want at least 2", len(stored))
	}
	last := stored[len(stored)-1]
	if last.EventType != "issue.ac_added" {
		t.Fatalf("last event type = %q, want issue.ac_added", last.EventType)
	}
	if len(last.Effects) != 1 {
		t.Fatalf("last event effects len = %d, want 1", len(last.Effects))
	}
	if last.Effects[0].Table != "acceptance_criteria" || last.Effects[0].Op != "insert" {
		t.Fatalf("last event effect = %#v, want acceptance_criteria insert", last.Effects[0])
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var issueID string
	var ordinal int
	var text string
	if err := conn.QueryRow(`SELECT issue_id, ordinal, text FROM acceptance_criteria WHERE id = ?`, first.ID).Scan(&issueID, &ordinal, &text); err != nil {
		t.Fatalf("query acceptance_criteria: %v", err)
	}
	if issueID != taskCreated.ID || ordinal != 1 || text != "First acceptance criterion" {
		t.Fatalf("stored row = (%q, %d, %q), want (%q, 1, %q)", issueID, ordinal, text, taskCreated.ID, "First acceptance criterion")
	}

	var activityRows int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityRows); err != nil {
		t.Fatalf("count activity_log: %v", err)
	}
	if activityRows != 0 {
		t.Fatalf("activity_log rows = %d, want 0", activityRows)
	}
}

func TestAddAcceptanceCriterionRejectsBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(t *testing.T, workspace string) string
		options func(workspace, issueID string) AddAcceptanceCriterionOptions
		wantErr string
		verify  func(t *testing.T, workspace, issueID string)
	}{
		{
			name: "epic kind rejected before append",
			prepare: func(t *testing.T, workspace string) string {
				t.Helper()
				created, err := Create(CreateOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Kind:          "epic",
					Title:         "Epic without ACs",
					ActorRole:     "tl",
					Scheduling:    "prioritised",
				})
				if err != nil {
					t.Fatalf("Create(epic) error = %v", err)
				}
				return created.ID
			},
			options: func(workspace, issueID string) AddAcceptanceCriterionOptions {
				return AddAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					Text:          "Rejected on epic",
					Ordinal:       1,
					ActorRole:     "developer",
				}
			},
			wantErr: `add acceptance criterion: issue "`,
			verify: func(t *testing.T, workspace, issueID string) {
				t.Helper()
				conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
				defer closeIssueTestDB(t, conn)
				assertAcceptanceCriterionCount(t, conn, issueID, 0)
			},
		},
		{
			name: "unknown role rejected before append",
			prepare: func(t *testing.T, workspace string) string {
				t.Helper()
				created, err := Create(CreateOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Kind:          "task",
					Title:         "Task for role validation",
					ActorRole:     "developer",
				})
				if err != nil {
					t.Fatalf("Create(task) error = %v", err)
				}
				return created.ID
			},
			options: func(workspace, issueID string) AddAcceptanceCriterionOptions {
				return AddAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					Text:          "Unknown role rejection",
					Ordinal:       1,
					ActorRole:     "intern",
				}
			},
			wantErr: `role "intern" not in project profile`,
			verify: func(t *testing.T, workspace, issueID string) {
				t.Helper()
				conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
				defer closeIssueTestDB(t, conn)
				assertAcceptanceCriterionCount(t, conn, issueID, 0)
			},
		},
		{
			name: "invalid ordinal rejected before append",
			prepare: func(t *testing.T, workspace string) string {
				t.Helper()
				created, err := Create(CreateOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Kind:          "task",
					Title:         "Task for ordinal validation",
					ActorRole:     "developer",
				})
				if err != nil {
					t.Fatalf("Create(task) error = %v", err)
				}
				return created.ID
			},
			options: func(workspace, issueID string) AddAcceptanceCriterionOptions {
				return AddAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					Text:          "Ordinal rejection",
					Ordinal:       0,
					ActorRole:     "developer",
				}
			},
			wantErr: "add acceptance criterion: ordinal must be a positive integer",
			verify: func(t *testing.T, workspace, issueID string) {
				t.Helper()
				conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
				defer closeIssueTestDB(t, conn)
				assertAcceptanceCriterionCount(t, conn, issueID, 0)
			},
		},
		{
			name: "not found rejected before append",
			prepare: func(t *testing.T, workspace string) string {
				t.Helper()
				return "FID-missing"
			},
			options: func(workspace, issueID string) AddAcceptanceCriterionOptions {
				return AddAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					Text:          "Not found rejection",
					Ordinal:       1,
					ActorRole:     "developer",
				}
			},
			wantErr: `issue "FID-missing" not found`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newIssueWorkspace(t)
			issueID := tt.prepare(t, workspace)
			journalPath := journalPathForWorkspace(t, workspace)
			linesBefore := countJournalLines(journalPath)

			_, err := AddAcceptanceCriterion(tt.options(workspace, issueID))
			if err == nil {
				t.Fatalf("AddAcceptanceCriterion() error = nil, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("AddAcceptanceCriterion() error = %q, want substring %q", err.Error(), tt.wantErr)
			}
			if got := countJournalLines(journalPath); got != linesBefore {
				t.Fatalf("journal line count after rejection = %d, want unchanged %d", got, linesBefore)
			}
			if tt.verify != nil {
				tt.verify(t, workspace, issueID)
			}
		})
	}
}

func TestAddAcceptanceCriterionInputValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		options AddAcceptanceCriterionOptions
		wantErr string
	}{
		{
			name: "requires workspace root",
			options: AddAcceptanceCriterionOptions{
				MigrationsDir: testIssueMigrationsDir(t),
				IssueID:       "FID-123",
				Text:          "Missing workspace",
				Ordinal:       1,
				ActorRole:     "developer",
			},
			wantErr: "add acceptance criterion: workspace root is required",
		},
		{
			name: "requires migrations dir",
			options: AddAcceptanceCriterionOptions{
				WorkspaceRoot: t.TempDir(),
				IssueID:       "FID-123",
				Text:          "Missing migrations",
				Ordinal:       1,
				ActorRole:     "developer",
			},
			wantErr: "add acceptance criterion: migrations dir is required",
		},
		{
			name: "requires issue id",
			options: AddAcceptanceCriterionOptions{
				WorkspaceRoot: t.TempDir(),
				MigrationsDir: testIssueMigrationsDir(t),
				Text:          "Missing issue id",
				Ordinal:       1,
				ActorRole:     "developer",
			},
			wantErr: "add acceptance criterion: issue id is required",
		},
		{
			name: "requires text",
			options: AddAcceptanceCriterionOptions{
				WorkspaceRoot: t.TempDir(),
				MigrationsDir: testIssueMigrationsDir(t),
				IssueID:       "FID-123",
				Ordinal:       1,
				ActorRole:     "developer",
			},
			wantErr: "add acceptance criterion: text is required",
		},
		{
			name: "requires positive ordinal",
			options: AddAcceptanceCriterionOptions{
				WorkspaceRoot: t.TempDir(),
				MigrationsDir: testIssueMigrationsDir(t),
				IssueID:       "FID-123",
				Text:          "Bad ordinal",
				ActorRole:     "developer",
			},
			wantErr: "add acceptance criterion: ordinal must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := AddAcceptanceCriterion(tt.options)
			if err == nil {
				t.Fatalf("AddAcceptanceCriterion() error = nil, want %q", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("AddAcceptanceCriterion() error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestAddAcceptanceCriterionTieBreaksByID(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	if _, err := conn.Exec(`
		INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, "FID-ordering", "task", "Ordering task", "todo", "task-default", "developer", "2026-06-14T00:00:00Z", "2026-06-14T00:00:00Z"); err != nil {
		t.Fatalf("insert task: %v", err)
	}
	if _, err := conn.Exec(`
		INSERT INTO acceptance_criteria (id, issue_id, ordinal, text)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		"AC-z-last", "FID-ordering", 1, "Later ID",
		"AC-a-first", "FID-ordering", 1, "Earlier ID",
	); err != nil {
		t.Fatalf("insert acceptance criteria: %v", err)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            "FID-ordering",
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got.AcceptanceCriteria) != 2 {
		t.Fatalf("Get() acceptance_criteria len = %d, want 2", len(got.AcceptanceCriteria))
	}
	if got.AcceptanceCriteria[0].ID != "AC-a-first" || got.AcceptanceCriteria[1].ID != "AC-z-last" {
		t.Fatalf("Get() acceptance_criteria ids = [%s %s], want [AC-a-first AC-z-last]", got.AcceptanceCriteria[0].ID, got.AcceptanceCriteria[1].ID)
	}
}

func TestAddAcceptanceCriterionEventRoundTripStaysCanonical(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Round-trip AC add",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	result, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "Canonical event",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion() error = %v", err)
	}

	stored := readIssueStoredEvents(t, journalPathForWorkspace(t, workspace))
	last := stored[len(stored)-1]
	payload, err := json.Marshal(last)
	if err != nil {
		t.Fatalf("Marshal(last event) error = %v", err)
	}

	var roundTripped struct {
		EventType string `json:"event_type"`
		Effects   []struct {
			Op    string         `json:"op"`
			Table string         `json:"table"`
			PK    map[string]any `json:"pk"`
		} `json:"effects"`
	}
	if err := json.Unmarshal(payload, &roundTripped); err != nil {
		t.Fatalf("Unmarshal(round-trip) error = %v", err)
	}
	if roundTripped.EventType != "issue.ac_added" {
		t.Fatalf("round-tripped event type = %q, want issue.ac_added", roundTripped.EventType)
	}
	if len(roundTripped.Effects) != 1 || roundTripped.Effects[0].PK["id"] != result.ID {
		t.Fatalf("round-tripped event effects = %#v, want one effect keyed by %q", roundTripped.Effects, result.ID)
	}
}

func TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Implement AC verify",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	criterion, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "Acceptance criterion can be verified",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion() error = %v", err)
	}

	linesBefore := countJournalLines(journalPathForWorkspace(t, workspace))

	result, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   criterion.ID,
		ActorRole:     "developer",
		Status:        "pass",
		EvidenceType:  acStringPointer("test"),
		EvidenceRef:   acStringPointer("internal/issues/ac_test.go::TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent"),
		Notes:         acStringPointer("Verified by targeted engine test"),
	})
	if err != nil {
		t.Fatalf("VerifyAcceptanceCriterion() error = %v", err)
	}
	if !strings.HasPrefix(result.ID, "ACV-") {
		t.Fatalf("VerifyAcceptanceCriterion() ID = %q, want ACV- prefix", result.ID)
	}
	if result.IssueID != taskCreated.ID {
		t.Fatalf("VerifyAcceptanceCriterion() IssueID = %q, want %q", result.IssueID, taskCreated.ID)
	}
	if result.CriterionID != criterion.ID {
		t.Fatalf("VerifyAcceptanceCriterion() CriterionID = %q, want %q", result.CriterionID, criterion.ID)
	}
	if result.Status != "pass" {
		t.Fatalf("VerifyAcceptanceCriterion() Status = %q, want pass", result.Status)
	}
	if result.EvidenceType == nil || *result.EvidenceType != "test" {
		t.Fatalf("VerifyAcceptanceCriterion() EvidenceType = %#v, want test", result.EvidenceType)
	}
	if result.EvidenceRef == nil || *result.EvidenceRef != "internal/issues/ac_test.go::TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent" {
		t.Fatalf("VerifyAcceptanceCriterion() EvidenceRef = %#v, want targeted test ref", result.EvidenceRef)
	}
	if result.Notes == nil || *result.Notes != "Verified by targeted engine test" {
		t.Fatalf("VerifyAcceptanceCriterion() Notes = %#v, want targeted note", result.Notes)
	}
	if strings.TrimSpace(result.VerifiedAt) == "" {
		t.Fatal("VerifyAcceptanceCriterion() VerifiedAt = empty, want timestamp")
	}

	if got := countJournalLines(journalPathForWorkspace(t, workspace)); got != linesBefore+1 {
		t.Fatalf("journal line count = %d, want %d", got, linesBefore+1)
	}

	stored := readIssueStoredEvents(t, journalPathForWorkspace(t, workspace))
	if len(stored) < 3 {
		t.Fatalf("stored events len = %d, want at least 3", len(stored))
	}
	last := stored[len(stored)-1]
	if last.EventType != "issue.ac_verified" {
		t.Fatalf("last event type = %q, want issue.ac_verified", last.EventType)
	}
	if len(last.Effects) != 1 {
		t.Fatalf("last event effects len = %d, want 1", len(last.Effects))
	}
	if last.Effects[0].Table != "ac_verifications" || last.Effects[0].Op != "insert" {
		t.Fatalf("last event effect = %#v, want ac_verifications insert", last.Effects[0])
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var (
		acID         string
		verifierRole string
		status       string
		evidenceType sql.NullString
		evidenceRef  sql.NullString
		notes        sql.NullString
		completionID sql.NullString
		verifiedAt   string
	)
	if err := conn.QueryRow(`
		SELECT ac_id, verifier_role, status, evidence_type, evidence_ref, notes, completion_id, verified_at
		FROM ac_verifications
		WHERE id = ?
	`, result.ID).Scan(&acID, &verifierRole, &status, &evidenceType, &evidenceRef, &notes, &completionID, &verifiedAt); err != nil {
		t.Fatalf("query ac_verifications: %v", err)
	}
	if acID != criterion.ID || verifierRole != "developer" || status != "pass" {
		t.Fatalf("stored verification = (%q, %q, %q), want (%q, developer, pass)", acID, verifierRole, status, criterion.ID)
	}
	if !evidenceType.Valid || evidenceType.String != "test" {
		t.Fatalf("stored evidence_type = %#v, want test", evidenceType)
	}
	if !evidenceRef.Valid || evidenceRef.String != "internal/issues/ac_test.go::TestVerifyAcceptanceCriterionPersistsAllFieldsAndAppendsEvent" {
		t.Fatalf("stored evidence_ref = %#v, want targeted test ref", evidenceRef)
	}
	if !notes.Valid || notes.String != "Verified by targeted engine test" {
		t.Fatalf("stored notes = %#v, want targeted note", notes)
	}
	if completionID.Valid {
		t.Fatalf("stored completion_id = %#v, want NULL", completionID)
	}
	if verifiedAt != result.VerifiedAt {
		t.Fatalf("stored verified_at = %q, want %q", verifiedAt, result.VerifiedAt)
	}

	var activityRows int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM activity_log`).Scan(&activityRows); err != nil {
		t.Fatalf("count activity_log: %v", err)
	}
	if activityRows != 0 {
		t.Fatalf("activity_log rows = %d, want 0", activityRows)
	}
}

func TestLatestVerificationSelectionPrefersAppendOrderWithinTimestampTie(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Deterministic AC latest ordering",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	criterion, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "Latest-per-role follows append order",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion() error = %v", err)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	if _, err := conn.Exec(`
		INSERT INTO ac_verifications (id, ac_id, verifier_role, status, verified_at)
		VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?)
	`,
		"ACV-z-fail",
		criterion.ID,
		"developer",
		"fail",
		"2026-06-14T00:00:00Z",
		"ACV-a-pass",
		criterion.ID,
		"developer",
		"pass",
		"2026-06-14T00:00:00Z",
	); err != nil {
		t.Fatalf("insert ac_verifications: %v", err)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	developerLatest, ok := got.AcceptanceCriteria[0].LatestVerifications["developer"].(map[string]any)
	if !ok {
		t.Fatalf("developer latest verification = %#v, want map", got.AcceptanceCriteria[0].LatestVerifications["developer"])
	}
	if developerLatest["status"] != "pass" {
		t.Fatalf("developer latest verification = %#v, want pass from later append", developerLatest)
	}
}

func TestVerifyAcceptanceCriterionPreservesHistoryAndLatestPerRole(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Append-only AC verification history",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	criterion, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "History is append only",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion() error = %v", err)
	}

	developerFail, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   criterion.ID,
		ActorRole:     "developer",
		Status:        "fail",
		Notes:         acStringPointer("Initial developer failure"),
	})
	if err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(developer fail) error = %v", err)
	}

	afterFail, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(after fail) error = %v", err)
	}
	developerLatest, ok := afterFail.AcceptanceCriteria[0].LatestVerifications["developer"].(map[string]any)
	if !ok {
		t.Fatalf("developer latest after fail = %#v, want map", afterFail.AcceptanceCriteria[0].LatestVerifications["developer"])
	}
	if developerLatest["status"] != "fail" {
		t.Fatalf("developer latest after fail = %#v, want fail", developerLatest)
	}

	qaLeadFail, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   criterion.ID,
		ActorRole:     "qa_lead",
		Status:        "fail",
		Notes:         acStringPointer("Independent QA lead failure"),
	})
	if err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(qa_lead fail) error = %v", err)
	}

	developerPass, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   criterion.ID,
		ActorRole:     "developer",
		Status:        "pass",
		Notes:         acStringPointer("Developer fixed the failure"),
	})
	if err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(developer pass) error = %v", err)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(after pass) error = %v", err)
	}
	if len(got.AcceptanceCriteria) != 1 {
		t.Fatalf("Get() acceptance_criteria len = %d, want 1", len(got.AcceptanceCriteria))
	}

	developerLatest, ok = got.AcceptanceCriteria[0].LatestVerifications["developer"].(map[string]any)
	if !ok {
		t.Fatalf("developer latest after pass = %#v, want map", got.AcceptanceCriteria[0].LatestVerifications["developer"])
	}
	if developerLatest["status"] != "pass" || developerLatest["verified_at"] != developerPass.VerifiedAt {
		t.Fatalf("developer latest after pass = %#v, want latest pass at %q", developerLatest, developerPass.VerifiedAt)
	}

	qaLeadLatest, ok := got.AcceptanceCriteria[0].LatestVerifications["qa_lead"].(map[string]any)
	if !ok {
		t.Fatalf("qa_lead latest after pass = %#v, want map", got.AcceptanceCriteria[0].LatestVerifications["qa_lead"])
	}
	if qaLeadLatest["status"] != "fail" || qaLeadLatest["verified_at"] != qaLeadFail.VerifiedAt {
		t.Fatalf("qa_lead latest after pass = %#v, want independent fail at %q", qaLeadLatest, qaLeadFail.VerifiedAt)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	assertAcceptanceVerificationCount(t, conn, criterion.ID, 3)

	var developerRows int
	if err := conn.QueryRow(`
		SELECT COUNT(*)
		FROM ac_verifications
		WHERE ac_id = ? AND verifier_role = ?
	`, criterion.ID, "developer").Scan(&developerRows); err != nil {
		t.Fatalf("count developer ac_verifications: %v", err)
	}
	if developerRows != 2 {
		t.Fatalf("developer ac_verifications rows = %d, want 2", developerRows)
	}

	var developerStatuses []string
	rows, err := conn.Query(`
		SELECT status
		FROM ac_verifications
		WHERE ac_id = ? AND verifier_role = ?
		ORDER BY rowid
	`, criterion.ID, "developer")
	if err != nil {
		t.Fatalf("query developer verification history: %v", err)
	}
	defer closeIssueRows(rows)
	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err != nil {
			t.Fatalf("scan developer verification history: %v", err)
		}
		developerStatuses = append(developerStatuses, status)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate developer verification history: %v", err)
	}
	if want := []string{"fail", "pass"}; !equalStrings(developerStatuses, want) {
		t.Fatalf("developer verification history = %#v, want %#v", developerStatuses, want)
	}

	if developerFail.ID == developerPass.ID {
		t.Fatalf("developer verification ids = %q and %q, want distinct rows", developerFail.ID, developerPass.ID)
	}
}

func TestVerifyAcceptanceCriterionLatestPerRoleSurvivesRebuild(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Rebuild AC latest-per-role parity",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	criterion, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "Rebuild keeps latest verification",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion() error = %v", err)
	}

	if _, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   criterion.ID,
		ActorRole:     "developer",
		Status:        "fail",
		Notes:         acStringPointer("Initial failure before immediate pass"),
	}); err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(developer fail) error = %v", err)
	}
	developerPass, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   criterion.ID,
		ActorRole:     "developer",
		Status:        "pass",
		Notes:         acStringPointer("Immediate pass before rebuild"),
	})
	if err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(developer pass) error = %v", err)
	}
	qaLeadFail, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   criterion.ID,
		ActorRole:     "qa_lead",
		Status:        "fail",
		Notes:         acStringPointer("Independent QA lead failure"),
	})
	if err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(qa_lead fail) error = %v", err)
	}

	before, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(before rebuild) error = %v", err)
	}

	dbPath := filepath.Join(workspace, ".fid", "fid.sqlite")
	if err := os.Remove(dbPath); err != nil {
		t.Fatalf("Remove(%q) error = %v", dbPath, err)
	}

	if _, err := journal.RebuildWorkspace(workspace, testIssueMigrationsDir(t), "developer"); err != nil {
		t.Fatalf("RebuildWorkspace() error = %v", err)
	}

	after, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get(after rebuild) error = %v", err)
	}

	beforeDeveloper := before.AcceptanceCriteria[0].LatestVerifications["developer"].(map[string]any)
	afterDeveloper := after.AcceptanceCriteria[0].LatestVerifications["developer"].(map[string]any)
	if afterDeveloper["status"] != "pass" || afterDeveloper["verified_at"] != developerPass.VerifiedAt {
		t.Fatalf("developer latest after rebuild = %#v, want pass at %q", afterDeveloper, developerPass.VerifiedAt)
	}
	if beforeDeveloper["status"] != afterDeveloper["status"] || beforeDeveloper["verified_at"] != afterDeveloper["verified_at"] {
		t.Fatalf("developer latest before/after rebuild = %#v / %#v, want identical", beforeDeveloper, afterDeveloper)
	}

	beforeQALead := before.AcceptanceCriteria[0].LatestVerifications["qa_lead"].(map[string]any)
	afterQALead := after.AcceptanceCriteria[0].LatestVerifications["qa_lead"].(map[string]any)
	if afterQALead["status"] != "fail" || afterQALead["verified_at"] != qaLeadFail.VerifiedAt {
		t.Fatalf("qa_lead latest after rebuild = %#v, want fail at %q", afterQALead, qaLeadFail.VerifiedAt)
	}
	if beforeQALead["status"] != afterQALead["status"] || beforeQALead["verified_at"] != afterQALead["verified_at"] {
		t.Fatalf("qa_lead latest before/after rebuild = %#v / %#v, want identical", beforeQALead, afterQALead)
	}
}

func TestVerifyAcceptanceCriterionRejectsBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(t *testing.T, workspace string) (issueID string, criterionID string)
		options func(workspace, issueID, criterionID string) VerifyAcceptanceCriterionOptions
		wantErr string
		verify  func(t *testing.T, workspace, issueID, criterionID string)
	}{
		{
			name: "invalid status rejected before append",
			prepare: func(t *testing.T, workspace string) (string, string) {
				return seedAcceptanceCriterionForVerification(t, workspace)
			},
			options: func(workspace, issueID, criterionID string) VerifyAcceptanceCriterionOptions {
				return VerifyAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					CriterionID:   criterionID,
					ActorRole:     "developer",
					Status:        "bogus",
				}
			},
			wantErr: `verify acceptance criterion: status "bogus" must be one of pass, fail, blocked, not_applicable`,
		},
		{
			name: "invalid evidence type rejected before append",
			prepare: func(t *testing.T, workspace string) (string, string) {
				return seedAcceptanceCriterionForVerification(t, workspace)
			},
			options: func(workspace, issueID, criterionID string) VerifyAcceptanceCriterionOptions {
				return VerifyAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					CriterionID:   criterionID,
					ActorRole:     "developer",
					Status:        "pass",
					EvidenceType:  acStringPointer("pdf"),
				}
			},
			wantErr: `verify acceptance criterion: evidence type "pdf" must be one of test, source, screenshot, manual, command, commit`,
		},
		{
			name: "missing acceptance criterion rejected before append",
			prepare: func(t *testing.T, workspace string) (string, string) {
				issueID, _ := seedAcceptanceCriterionForVerification(t, workspace)
				return issueID, "AC-missing"
			},
			options: func(workspace, issueID, criterionID string) VerifyAcceptanceCriterionOptions {
				return VerifyAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					CriterionID:   criterionID,
					ActorRole:     "developer",
					Status:        "pass",
				}
			},
			wantErr: `verify acceptance criterion: acceptance criterion "AC-missing" not found`,
		},
		{
			name: "issue and acceptance criterion mismatch rejected before append",
			prepare: func(t *testing.T, workspace string) (string, string) {
				_, criterionID := seedAcceptanceCriterionForVerification(t, workspace)
				otherCreated, err := Create(CreateOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					Kind:          "task",
					Title:         "Other issue",
					ActorRole:     "developer",
				})
				if err != nil {
					t.Fatalf("Create(other task) error = %v", err)
				}
				return otherCreated.ID, criterionID
			},
			options: func(workspace, issueID, criterionID string) VerifyAcceptanceCriterionOptions {
				return VerifyAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					CriterionID:   criterionID,
					ActorRole:     "developer",
					Status:        "pass",
				}
			},
			wantErr: `verify acceptance criterion: acceptance criterion "`,
			verify: func(t *testing.T, workspace, issueID, criterionID string) {
				t.Helper()
				conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
				defer closeIssueTestDB(t, conn)
				assertAcceptanceVerificationCount(t, conn, criterionID, 0)
			},
		},
		{
			name: "unknown role rejected before append",
			prepare: func(t *testing.T, workspace string) (string, string) {
				return seedAcceptanceCriterionForVerification(t, workspace)
			},
			options: func(workspace, issueID, criterionID string) VerifyAcceptanceCriterionOptions {
				return VerifyAcceptanceCriterionOptions{
					WorkspaceRoot: workspace,
					MigrationsDir: testIssueMigrationsDir(t),
					IssueID:       issueID,
					CriterionID:   criterionID,
					ActorRole:     "intern",
					Status:        "pass",
				}
			},
			wantErr: `role "intern" not in project profile`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newIssueWorkspace(t)
			issueID, criterionID := tt.prepare(t, workspace)
			journalPath := journalPathForWorkspace(t, workspace)
			linesBefore := countJournalLines(journalPath)

			_, err := VerifyAcceptanceCriterion(tt.options(workspace, issueID, criterionID))
			if err == nil {
				t.Fatalf("VerifyAcceptanceCriterion() error = nil, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("VerifyAcceptanceCriterion() error = %q, want substring %q", err.Error(), tt.wantErr)
			}
			if got := countJournalLines(journalPath); got != linesBefore {
				t.Fatalf("journal line count after rejection = %d, want unchanged %d", got, linesBefore)
			}

			conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
			defer closeIssueTestDB(t, conn)
			assertAcceptanceVerificationCount(t, conn, criterionID, 0)

			if tt.verify != nil {
				tt.verify(t, workspace, issueID, criterionID)
			}
		})
	}
}

func acStringPointer(value string) *string {
	return &value
}

func seedAcceptanceCriterionForVerification(t *testing.T, workspace string) (string, string) {
	t.Helper()

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Verification task",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	criterion, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "Criterion for verification",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion() error = %v", err)
	}

	return taskCreated.ID, criterion.ID
}

func assertAcceptanceCriterionCount(t *testing.T, conn *sql.DB, issueID string, want int) {
	t.Helper()

	var got int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM acceptance_criteria WHERE issue_id = ?`, issueID).Scan(&got); err != nil {
		t.Fatalf("count acceptance_criteria(%q): %v", issueID, err)
	}
	if got != want {
		t.Fatalf("acceptance_criteria count for %q = %d, want %d", issueID, got, want)
	}
}

func assertAcceptanceVerificationCount(t *testing.T, conn *sql.DB, criterionID string, want int) {
	t.Helper()

	var got int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM ac_verifications WHERE ac_id = ?`, criterionID).Scan(&got); err != nil {
		t.Fatalf("count ac_verifications(%q): %v", criterionID, err)
	}
	if got != want {
		t.Fatalf("ac_verifications count for %q = %d, want %d", criterionID, got, want)
	}
}

func TestAddAcceptanceCriterionNotFoundErrorShape(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	_, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       "FID-missing",
		Text:          "Missing issue",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("AddAcceptanceCriterion() error = nil, want not found")
	}

	var notFoundErr *NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("AddAcceptanceCriterion() error = %T, want *NotFoundError", err)
	}
	if notFoundErr.IssueID != "FID-missing" {
		t.Fatalf("NotFoundError.IssueID = %q, want FID-missing", notFoundErr.IssueID)
	}
}

func TestListAcceptanceCriteriaReturnsOrdinalLatestPerRoleAndMatchesIssueGet(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "List acceptance criteria",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	second, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "Second criterion",
		Ordinal:       2,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion(second) error = %v", err)
	}

	first, err := AddAcceptanceCriterion(AddAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		Text:          "First criterion",
		Ordinal:       1,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddAcceptanceCriterion(first) error = %v", err)
	}

	if _, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   first.ID,
		ActorRole:     "developer",
		Status:        "fail",
	}); err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(fail) error = %v", err)
	}
	if _, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   first.ID,
		ActorRole:     "developer",
		Status:        "pass",
		EvidenceType:  acStringPointer("test"),
		EvidenceRef:   acStringPointer("internal/issues/ac_test.go::TestListAcceptanceCriteriaReturnsOrdinalLatestPerRoleAndMatchesIssueGet"),
	}); err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(pass) error = %v", err)
	}
	if _, err := VerifyAcceptanceCriterion(VerifyAcceptanceCriterionOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
		CriterionID:   first.ID,
		ActorRole:     "qa_lead",
		Status:        "blocked",
	}); err != nil {
		t.Fatalf("VerifyAcceptanceCriterion(qa_lead blocked) error = %v", err)
	}

	got, err := ListAcceptanceCriteria(ListAcceptanceCriteriaOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("ListAcceptanceCriteria() error = %v", err)
	}

	if len(got.Criteria) != 2 {
		t.Fatalf("ListAcceptanceCriteria() len = %d, want 2", len(got.Criteria))
	}
	if got.Criteria[0].ID != first.ID || got.Criteria[0].Ordinal != 1 || got.Criteria[0].Text != "First criterion" {
		t.Fatalf("first list row = %#v, want first ordinal/text/id", got.Criteria[0])
	}
	if got.Criteria[1].ID != second.ID || got.Criteria[1].Ordinal != 2 || got.Criteria[1].Text != "Second criterion" {
		t.Fatalf("second list row = %#v, want second ordinal/text/id", got.Criteria[1])
	}

	developerLatest, ok := got.Criteria[0].LatestVerifications["developer"].(map[string]any)
	if !ok {
		t.Fatalf("developer latest verification = %#v, want map", got.Criteria[0].LatestVerifications["developer"])
	}
	if developerLatest["status"] != "pass" {
		t.Fatalf("developer latest verification status = %#v, want pass", developerLatest["status"])
	}
	qaLatest, ok := got.Criteria[0].LatestVerifications["qa_lead"].(map[string]any)
	if !ok {
		t.Fatalf("qa_lead latest verification = %#v, want map", got.Criteria[0].LatestVerifications["qa_lead"])
	}
	if qaLatest["status"] != "blocked" {
		t.Fatalf("qa_lead latest verification status = %#v, want blocked", qaLatest["status"])
	}

	issue, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !reflect.DeepEqual(got.Criteria, issue.AcceptanceCriteria) {
		t.Fatalf("ListAcceptanceCriteria() criteria = %#v, want issue acceptance %#v", got.Criteria, issue.AcceptanceCriteria)
	}
}

func TestListAcceptanceCriteriaEmptyAndNotFound(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	taskCreated, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Empty acceptance criteria",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create(task) error = %v", err)
	}

	got, err := ListAcceptanceCriteria(ListAcceptanceCriteriaOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       taskCreated.ID,
	})
	if err != nil {
		t.Fatalf("ListAcceptanceCriteria(empty) error = %v", err)
	}
	if len(got.Criteria) != 0 {
		t.Fatalf("ListAcceptanceCriteria(empty) len = %d, want 0", len(got.Criteria))
	}

	_, err = ListAcceptanceCriteria(ListAcceptanceCriteriaOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       "FID-missing",
	})
	if err == nil {
		t.Fatal("ListAcceptanceCriteria(missing) error = nil, want not found")
	}

	var notFoundErr *NotFoundError
	if !errors.As(err, &notFoundErr) {
		t.Fatalf("ListAcceptanceCriteria(missing) error = %T, want *NotFoundError", err)
	}
	if notFoundErr.IssueID != "FID-missing" {
		t.Fatalf("NotFoundError.IssueID = %q, want FID-missing", notFoundErr.IssueID)
	}
}
