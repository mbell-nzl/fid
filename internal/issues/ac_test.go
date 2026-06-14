package issues

import (
	"database/sql"
	"encoding/json"
	"errors"
	"path/filepath"
	"strings"
	"testing"
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
