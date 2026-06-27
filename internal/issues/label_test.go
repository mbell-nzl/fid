package issues

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/journal"
)

func TestAddLabelAddsOneRowWithoutReplacingOthersAndProjectsActivity(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "api", "core", "legacy")

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Domain issue label add",
		ActorRole:     "developer",
		Labels:        []string{"core", "legacy"},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	result, err := AddLabel(AddLabelOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		Label:         "api",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("AddLabel() error = %v", err)
	}
	if result.IssueID != created.ID {
		t.Fatalf("AddLabelResult.IssueID = %q, want %q", result.IssueID, created.ID)
	}
	if result.Label != "api" {
		t.Fatalf("AddLabelResult.Label = %q, want %q", result.Label, "api")
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if want := []string{"api", "core", "legacy"}; !equalStrings(got.Labels, want) {
		t.Fatalf("labels = %#v, want %#v", got.Labels, want)
	}
	if got := countJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after AddLabel = %d, want %d", got, linesBefore+1)
	}

	event := readIssueStoredEvents(t, journalPath)
	last := event[len(event)-1]
	if last.EventType != "issue.label.added" {
		t.Fatalf("event type = %q, want %q", last.EventType, "issue.label.added")
	}

	var (
		insertEffect   *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "issue_labels" && effect.Op == "insert":
			insertEffect = effect
		case effect.Table == "activity_log" && effect.Op == "activity":
			activityEffect = effect
		}
	}
	if insertEffect == nil {
		t.Fatal("missing issue_labels insert effect")
	}
	if activityEffect == nil {
		t.Fatal("missing activity effect")
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	assertIssueLabelActivity(t, conn, last.EventID, "issue_label_added", "developer", `"label":"api"`, `"direction":"added"`)
}

func TestAddLabelRejectsMissingOrExistingLabelBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		label   string
		wantErr string
	}{
		{
			name:    "missing label",
			label:   "missing",
			wantErr: `issue label add: label "missing" does not exist`,
		},
		{
			name:    "already present",
			label:   "core",
			wantErr: `issue label add: issue "ISSUE-ID" already has label "core"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newIssueWorkspace(t)
			seedIssueLabels(t, workspace, "core")

			created, err := Create(CreateOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				Kind:          "task",
				Title:         "Reject issue label add",
				ActorRole:     "developer",
				Labels:        []string{"core"},
			})
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			journalPath := journalPathForWorkspace(t, workspace)
			linesBefore := countJournalLines(journalPath)

			_, err = AddLabel(AddLabelOptions{
				WorkspaceRoot: workspace,
				MigrationsDir: testIssueMigrationsDir(t),
				IssueID:       created.ID,
				Label:         tt.label,
				ActorRole:     "developer",
			})
			if err == nil {
				t.Fatal("AddLabel() error = nil, want rejection")
			}

			wantErr := strings.ReplaceAll(tt.wantErr, "ISSUE-ID", created.ID)
			if err.Error() != wantErr {
				t.Fatalf("AddLabel() error = %q, want %q", err.Error(), wantErr)
			}
			if got := countJournalLines(journalPath); got != linesBefore {
				t.Fatalf("journal lines after rejected AddLabel = %d, want %d", got, linesBefore)
			}
		})
	}
}

func TestRemoveLabelRemovesOneRowPreservesLabelTableAndProjectsActivity(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "api", "core", "legacy")

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Domain issue label remove",
		ActorRole:     "developer",
		Labels:        []string{"api", "core", "legacy"},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	result, err := RemoveLabel(RemoveLabelOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		Label:         "legacy",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("RemoveLabel() error = %v", err)
	}
	if result.IssueID != created.ID || result.Label != "legacy" {
		t.Fatalf("RemoveLabelResult = %#v, want issue %q label %q", result, created.ID, "legacy")
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if want := []string{"api", "core"}; !equalStrings(got.Labels, want) {
		t.Fatalf("labels = %#v, want %#v", got.Labels, want)
	}
	if got := countJournalLines(journalPath); got != linesBefore+1 {
		t.Fatalf("journal lines after RemoveLabel = %d, want %d", got, linesBefore+1)
	}

	events := readIssueStoredEvents(t, journalPath)
	last := events[len(events)-1]
	if last.EventType != "issue.label.removed" {
		t.Fatalf("event type = %q, want %q", last.EventType, "issue.label.removed")
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var labelCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = ?`, "legacy").Scan(&labelCount); err != nil {
		t.Fatalf("count labels row: %v", err)
	}
	if labelCount != 1 {
		t.Fatalf("labels row count for legacy = %d, want 1", labelCount)
	}

	assertIssueLabelActivity(t, conn, last.EventID, "issue_label_removed", "developer", `"label":"legacy"`, `"direction":"removed"`)
}

func TestRemoveLabelRejectsAbsentLabelBeforeAppend(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	seedIssueLabels(t, workspace, "core")

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Reject issue label remove",
		ActorRole:     "developer",
		Labels:        []string{"core"},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	journalPath := journalPathForWorkspace(t, workspace)
	linesBefore := countJournalLines(journalPath)

	_, err = RemoveLabel(RemoveLabelOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		Label:         "legacy",
		ActorRole:     "developer",
	})
	if err == nil {
		t.Fatal("RemoveLabel() error = nil, want rejection")
	}
	if err.Error() != `issue label remove: issue "`+created.ID+`" does not have label "legacy"` {
		t.Fatalf("RemoveLabel() error = %q, want absent-label rejection", err.Error())
	}
	if got := countJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected RemoveLabel = %d, want %d", got, linesBefore)
	}
}

func assertIssueLabelActivity(t *testing.T, conn *sql.DB, eventID, wantEvent, wantActor string, wantDetails ...string) {
	t.Helper()

	var (
		activityCount int
		eventName     string
		actor         string
		detailsRaw    string
	)
	if err := conn.QueryRow(`
		SELECT COUNT(*), event, COALESCE(actor, ''), COALESCE(details, '')
		FROM activity_log
		WHERE journal_event_id = ?
		GROUP BY event, actor, details
	`, eventID).Scan(&activityCount, &eventName, &actor, &detailsRaw); err != nil {
		t.Fatalf("query activity_log row: %v", err)
	}
	if activityCount != 1 {
		t.Fatalf("activity_log grouped row count = %d, want 1", activityCount)
	}
	if eventName != wantEvent {
		t.Fatalf("activity_log event = %q, want %q", eventName, wantEvent)
	}
	if actor != wantActor {
		t.Fatalf("activity_log actor = %q, want %q", actor, wantActor)
	}
	for _, want := range wantDetails {
		if !strings.Contains(detailsRaw, want) {
			t.Fatalf("activity_log details = %q, want substring %q", detailsRaw, want)
		}
	}
}
