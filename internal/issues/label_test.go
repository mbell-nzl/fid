package issues

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		Labels: []string{"core", "legacy"},
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
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
	if result.IssueRevision != 2 || result.EventID == "" || result.OccurredAt == "" {
		t.Fatalf("AddLabelResult = %#v, want revision 2 and Event identity", result)
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
	if last.ActorV2.ActorID != "developer" || last.ActorV2.ActorType != "agent" || len(last.Effects) != 4 {
		t.Fatalf("Issue-label add Event actor/effects = (%#v, %d), want agent developer and 4 effects", last.ActorV2, len(last.Effects))
	}

	var (
		insertEffect   *journal.Effect
		revisionEffect *journal.Effect
		activityEffect *journal.Effect
	)
	for index := range last.Effects {
		effect := &last.Effects[index]
		switch {
		case effect.Table == "issue_labels" && effect.Op == "insert":
			insertEffect = effect
		case effect.Table == "issues" && effect.Op == "update":
			revisionEffect = effect
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
	assertIssueLabelRevisionEffect(t, revisionEffect, created.ID, 1, 2)

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	assertIssueLabelActivity(t, conn, last.EventID, "issue-label-added", "developer", `"label":"api"`, `"direction":"added"`)
}

func TestAddLabelRejectsMissingOrExistingLabelBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		label   string
		wantErr any
	}{
		{
			name:    "missing label",
			label:   "missing",
			wantErr: &LabelNotFoundError{Label: "missing"},
		},
		{
			name:    "already present",
			label:   "core",
			wantErr: &DuplicateIssueLabelError{Label: "core"},
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
				Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
				Labels: []string{"core"},
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
				Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
			})
			if err == nil {
				t.Fatal("AddLabel() error = nil, want rejection")
			}

			switch want := tt.wantErr.(type) {
			case *LabelNotFoundError:
				var got *LabelNotFoundError
				if !errors.As(err, &got) || got.Label != want.Label {
					t.Fatalf("AddLabel() error = %#v, want LabelNotFoundError(%q)", err, want.Label)
				}
			case *DuplicateIssueLabelError:
				var got *DuplicateIssueLabelError
				if !errors.As(err, &got) || got.IssueID != created.ID || got.Label != want.Label {
					t.Fatalf("AddLabel() error = %#v, want DuplicateIssueLabelError(%q, %q)", err, created.ID, want.Label)
				}
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		Labels: []string{"api", "core", "legacy"},
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err != nil {
		t.Fatalf("RemoveLabel() error = %v", err)
	}
	if result.IssueID != created.ID || result.Label != "legacy" {
		t.Fatalf("RemoveLabelResult = %#v, want issue %q label %q", result, created.ID, "legacy")
	}
	if result.IssueRevision != 2 || result.EventID == "" || result.OccurredAt == "" {
		t.Fatalf("RemoveLabelResult = %#v, want revision 2 and Event identity", result)
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
	if last.ActorV2.ActorID != "developer" || last.ActorV2.ActorType != "agent" || len(last.Effects) != 4 {
		t.Fatalf("Issue-label remove Event actor/effects = (%#v, %d), want agent developer and 4 effects", last.ActorV2, len(last.Effects))
	}
	var revisionEffect *journal.Effect
	for index := range last.Effects {
		effect := &last.Effects[index]
		if effect.Table == "issues" && effect.Op == "update" {
			revisionEffect = effect
		}
	}
	assertIssueLabelRevisionEffect(t, revisionEffect, created.ID, 1, 2)

	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	var labelCount int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM labels WHERE name = ?`, "legacy").Scan(&labelCount); err != nil {
		t.Fatalf("count labels row: %v", err)
	}
	if labelCount != 1 {
		t.Fatalf("labels row count for legacy = %d, want 1", labelCount)
	}

	assertIssueLabelActivity(t, conn, last.EventID, "issue-label-removed", "developer", `"label":"legacy"`, `"direction":"removed"`)
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		Labels: []string{"core"},
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
		Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
	})
	if err == nil {
		t.Fatal("RemoveLabel() error = nil, want rejection")
	}
	var notFound *LabelNotFoundError
	if !errors.As(err, &notFound) || notFound.Label != "legacy" {
		t.Fatalf("RemoveLabel() error = %#v, want LabelNotFoundError(%q)", err, "legacy")
	}
	if got := countJournalLines(journalPath); got != linesBefore {
		t.Fatalf("journal lines after rejected RemoveLabel = %d, want %d", got, linesBefore)
	}
}

func TestIssueLabelMutationsApplyClosedTokenBoundariesBeforeAppend(t *testing.T) {
	t.Parallel()

	t.Run("exact maximum succeeds for add and remove", func(t *testing.T) {
		workspace := newIssueWorkspace(t)
		label := strings.Repeat("l", 255)
		seedIssueLabels(t, workspace, label)
		created, err := Create(CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Kind: "task", Title: "Issue label boundary", Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if _, err := AddLabel(AddLabelOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), IssueID: created.ID, Label: label, Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()}); err != nil {
			t.Fatalf("AddLabel(exact maximum) error = %v", err)
		}
		if _, err := RemoveLabel(RemoveLabelOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), IssueID: created.ID, Label: label, Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()}); err != nil {
			t.Fatalf("RemoveLabel(exact maximum) error = %v", err)
		}
	})

	tests := []struct {
		name    string
		label   string
		wantErr string
	}{
		{name: "empty", label: " \t", wantErr: "label is required"},
		{name: "over maximum", label: strings.Repeat("l", 256), wantErr: "issue label exceeds 255 bytes"},
		{name: "control character", label: "bad\x7flabel", wantErr: "issue label contains ASCII control character"},
		{name: "invalid UTF-8", label: string([]byte{0xff}), wantErr: "issue label must be valid UTF-8"},
	}
	for _, operation := range []struct {
		name string
		run  func(string, string, string) error
	}{
		{name: "add", run: func(workspace, issueID, label string) error {
			_, err := AddLabel(AddLabelOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), IssueID: issueID, Label: label, Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
			return err
		}},
		{name: "remove", run: func(workspace, issueID, label string) error {
			_, err := RemoveLabel(RemoveLabelOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), IssueID: issueID, Label: label, Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
			return err
		}},
	} {
		t.Run(operation.name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					workspace := newIssueWorkspace(t)
					created, err := Create(CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Kind: "task", Title: "Issue label rejection", Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
					if err != nil {
						t.Fatalf("Create() error = %v", err)
					}
					journalPath := journalPathForWorkspace(t, workspace)
					before := countJournalLines(journalPath)
					err = operation.run(workspace, created.ID, test.label)
					wantErr := test.wantErr
					if test.name == "empty" {
						wantErr = "issue label " + operation.name + ": " + wantErr
					}
					if err == nil || err.Error() != wantErr {
						t.Fatalf("%s error = %v, want %q", operation.name, err, wantErr)
					}
					if got := countJournalLines(journalPath); got != before {
						t.Fatalf("journal lines after rejection = %d, want %d", got, before)
					}
				})
			}
		})
	}
}

func TestIssueLabelAddRejectsNestedLimitWithoutAppend(t *testing.T) {
	t.Parallel()
	workspace := newIssueWorkspace(t)
	created, err := Create(CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), Kind: "task", Title: "Issue label limit", Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Begin() error = %v", err)
	}
	for index := 0; index <= 100; index++ {
		label := fmt.Sprintf("limit-%03d", index)
		if _, err := tx.Exec(`INSERT INTO labels (name) VALUES (?)`, label); err != nil {
			t.Fatalf("insert label %d: %v", index, err)
		}
		if index < 100 {
			if _, err := tx.Exec(`INSERT INTO issue_labels (issue_id, label) VALUES (?, ?)`, created.ID, label); err != nil {
				t.Fatalf("insert issue label %d: %v", index, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
	closeIssueTestDB(t, conn)

	journalPath := journalPathForWorkspace(t, workspace)
	before := countJournalLines(journalPath)
	_, err = AddLabel(AddLabelOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), IssueID: created.ID, Label: "limit-100", Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()})
	var limit *ResultLimitExceededError
	if !errors.As(err, &limit) || limit.Collection != "issue_labels" || limit.Limit != 100 || limit.Observed != 101 {
		t.Fatalf("AddLabel() error = %#v, want issue_labels 101/100 limit", err)
	}
	if got := countJournalLines(journalPath); got != before {
		t.Fatalf("journal lines after limit rejection = %d, want %d", got, before)
	}
	if _, err := RemoveLabel(RemoveLabelOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), IssueID: created.ID, Label: "limit-000", Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()}); err != nil {
		t.Fatalf("RemoveLabel(at limit) error = %v", err)
	}
	if _, err := AddLabel(AddLabelOptions{WorkspaceRoot: workspace, MigrationsDir: testIssueMigrationsDir(t), IssueID: created.ID, Label: "limit-100", Actor: testActorForRole("developer"), IdempotencyKey: testIdempotencyKey()}); err != nil {
		t.Fatalf("AddLabel(after removal) error = %v", err)
	}
}

func TestIssueLabelEffectsPropagatesRevisionConstructionFailure(t *testing.T) {
	t.Parallel()
	request := issueLabelWriteRequest{issueID: "FID-invalid-revision", label: "core", actor: standaloneActor{id: "developer", typeName: "agent"}, spec: issueLabelOperationSpec{effectOp: "insert", direction: "added", eventName: "issue-label-added"}}
	if _, err := issueLabelEffects(request, "ACT-test", "EVT-test", "2026-08-15T00:00:00Z", 0); err == nil {
		t.Fatal("issueLabelEffects() error = nil, want invalid observed revision rejection")
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

func assertIssueLabelRevisionEffect(t *testing.T, effect *journal.Effect, issueID string, prior, next int) {
	t.Helper()
	if effect == nil {
		t.Fatal("missing Issue revision effect")
	}
	if effect.PK["id"] != issueID || effect.Values["revision"] != float64(next) || effect.Prior["revision"] != float64(prior) || effect.Preconditions["revision"] != float64(prior) {
		t.Fatalf("Issue revision effect = %#v, want %q revision %d -> %d", effect, issueID, prior, next)
	}
}
