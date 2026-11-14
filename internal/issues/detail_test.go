package issues

import (
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestGetReturnsCentralizedDetailShape(t *testing.T) {
	t.Parallel()

	t.Run("empty issue returns stable empty sections", func(t *testing.T) {
		t.Parallel()

		workspace := newIssueWorkspace(t)

		created, err := Create(CreateOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testIssueMigrationsDir(t),
			Kind:          "task",
			Title:         "Render empty detail",
			Actor:         testActorForRole("developer"), IdempotencyKey: testIdempotencyKey(),
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := Get(GetOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testIssueMigrationsDir(t),
			ID:            created.ID,
		})
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if got.ParentTitle != nil {
			t.Fatalf("Get() parent_title = %#v, want nil", got.ParentTitle)
		}
		if got.Labels == nil || len(got.Labels) != 0 {
			t.Fatalf("Get() labels = %#v, want empty array", got.Labels)
		}
		if got.AcceptanceCriteria == nil || len(got.AcceptanceCriteria) != 0 {
			t.Fatalf("Get() acceptance_criteria = %#v, want empty array", got.AcceptanceCriteria)
		}
		if got.Links == nil || len(got.Links) != 0 {
			t.Fatalf("Get() links = %#v, want empty array", got.Links)
		}
		if got.TriageStatus != nil {
			t.Fatalf("Get() triage_status = %#v, want nil", got.TriageStatus)
		}
	})

	t.Run("populated issue returns structured sections", func(t *testing.T) {
		t.Parallel()

		workspace := newIssueWorkspace(t)
		conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
		defer closeIssueTestDB(t, conn)

		seedIssueDetailFixture(t, workspace)

		got, err := Get(GetOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testIssueMigrationsDir(t),
			ID:            "FID-a1b2c3",
		})
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if got.ParentTitle == nil || *got.ParentTitle != "Core persistence layer" {
			t.Fatalf("Get() parent_title = %#v, want Core persistence layer", got.ParentTitle)
		}
		if want := []string{"core", "journal"}; !equalStrings(got.Labels, want) {
			t.Fatalf("Get() labels = %#v, want %#v", got.Labels, want)
		}
		if len(got.AcceptanceCriteria) != 2 {
			t.Fatalf("Get() acceptance_criteria len = %d, want 2", len(got.AcceptanceCriteria))
		}
		if got.AcceptanceCriteria[0].ID != "AC-d4e5f6" {
			t.Fatalf("Get() first acceptance criterion id = %q, want AC-d4e5f6", got.AcceptanceCriteria[0].ID)
		}
		if got.AcceptanceCriteria[0].Text != "Journal file is flushed to disk before returning success" ||
			got.AcceptanceCriteria[1].Text != "Partial write leaves journal intact" {
			t.Fatalf("Get() acceptance criteria = %#v, want retained ordered text", got.AcceptanceCriteria)
		}
		if got.TriageStatus == nil || *got.TriageStatus != "blocked" {
			t.Fatalf("Get() triage_status = %#v, want blocked", got.TriageStatus)
		}
		if len(got.Links) != 2 || got.Links[0].Relation != "blocks" || got.Links[1].Relation != "relates_to" {
			t.Fatalf("Get() links = %#v, want blocks + relates_to", got.Links)
		}
		if got.Links[0].TargetTitle == nil || *got.Links[0].TargetTitle != "Schema migration runner" {
			t.Fatalf("Get() first link target_title = %#v, want Schema migration runner", got.Links[0].TargetTitle)
		}
	})
}

func TestIssueParentTitle(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	if got, err := issueParentTitle(conn, nil); err != nil {
		t.Fatalf("issueParentTitle(nil) error = %v", err)
	} else if got != nil {
		t.Fatalf("issueParentTitle(nil) = %#v, want nil", got)
	}

	missingID := "FID-missing"
	if got, err := issueParentTitle(conn, &missingID); err != nil {
		t.Fatalf("issueParentTitle(missing) error = %v", err)
	} else if got != nil {
		t.Fatalf("issueParentTitle(missing) = %#v, want nil", got)
	}

	if _, err := conn.Exec(`
		INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, "FID-parent01", "epic", "Parent title", "todo", "epic-default", "tl", "2026-06-09T09:00:00Z", "2026-06-09T09:00:00Z"); err != nil {
		t.Fatalf("insert parent issue error = %v", err)
	}

	parentID := "FID-parent01"
	got, err := issueParentTitle(conn, &parentID)
	if err != nil {
		t.Fatalf("issueParentTitle(existing) error = %v", err)
	}
	if got == nil || *got != "Parent title" {
		t.Fatalf("issueParentTitle(existing) = %#v, want Parent title", got)
	}
}

func TestNormalizeDateOnlyPointer(t *testing.T) {
	t.Parallel()

	validDateOnly := " 2026-06-30 "
	validRFC3339 := "2026-06-30T15:04:05Z"
	blank := "   "
	malformed := "tomorrow"

	tests := []struct {
		name  string
		input *string
		want  *string
	}{
		{
			name:  "nil stays nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "blank string becomes nil",
			input: &blank,
			want:  nil,
		},
		{
			name:  "date only trims whitespace",
			input: &validDateOnly,
			want:  stringPointer("2026-06-30"),
		},
		{
			name:  "rfc3339 value is normalized to date only",
			input: &validRFC3339,
			want:  stringPointer("2026-06-30"),
		},
		{
			name:  "malformed value is preserved after trimming",
			input: &malformed,
			want:  stringPointer("tomorrow"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeDateOnlyPointer(tt.input)
			switch {
			case tt.want == nil && got != nil:
				t.Fatalf("normalizeDateOnlyPointer() = %#v, want nil", got)
			case tt.want != nil && got == nil:
				t.Fatal("normalizeDateOnlyPointer() = nil, want value")
			case tt.want != nil && *got != *tt.want:
				t.Fatalf("normalizeDateOnlyPointer() = %q, want %q", *got, *tt.want)
			}
		})
	}
}

func seedIssueDetailFixture(t *testing.T, workspace string) {
	t.Helper()
	appendIssueJournalEvent(t, workspace, journal.Event{EventSchema: journal.EventSchemaV1, EventType: "issue.detail-fixture", Actor: journal.Actor{Role: "developer"}, Command: journal.Command{Argv: []string{"fid", "test"}, CWD: ".", FIDVersion: "test"}, Effects: []journal.Effect{
		{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-e1f2g3"}, Values: map[string]any{"kind": "epic", "title": "Core persistence layer", "status": "todo", "workflow_id": "epic-default", "owner_role": "developer", "created_at": "2026-05-25T08:00:00Z", "updated_at": "2026-05-25T08:00:00Z"}},
		{Op: "insert", Table: "epic_metadata", PK: map[string]any{"issue_id": "FID-e1f2g3"}, Values: map[string]any{"scheduling": "sequential"}},
		{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-a1b2c3"}, Values: map[string]any{"kind": "task", "parent_id": "FID-e1f2g3", "title": "Implement journal append with flush", "workflow_id": "task-default", "status": "in_progress", "priority": 1, "sequence": 3, "owner_role": "developer", "created_at": "2026-05-25T09:00:00Z", "updated_at": "2026-05-27T10:20:00Z"}},
		{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-7c3e9f"}, Values: map[string]any{"kind": "task", "title": "Schema migration runner", "status": "todo", "workflow_id": "task-default", "owner_role": "developer", "created_at": "2026-05-25T09:05:00Z", "updated_at": "2026-05-25T09:05:00Z"}},
		{Op: "insert", Table: "labels", PK: map[string]any{"name": "core"}, Values: map[string]any{}},
		{Op: "insert", Table: "labels", PK: map[string]any{"name": "journal"}, Values: map[string]any{}},
		{Op: "insert", Table: "issue_labels", PK: map[string]any{"issue_id": "FID-a1b2c3", "label": "core"}, Values: map[string]any{}},
		{Op: "insert", Table: "issue_labels", PK: map[string]any{"issue_id": "FID-a1b2c3", "label": "journal"}, Values: map[string]any{}},
		{Op: "insert", Table: "acceptance_criteria", PK: map[string]any{"id": "AC-d4e5f6"}, Values: map[string]any{"issue_id": "FID-a1b2c3", "ordinal": 1, "text": "Journal file is flushed to disk before returning success"}},
		{Op: "insert", Table: "acceptance_criteria", PK: map[string]any{"id": "AC-e5f6g7"}, Values: map[string]any{"issue_id": "FID-a1b2c3", "ordinal": 2, "text": "Partial write leaves journal intact"}},
		{Op: "insert", Table: "triage_events", PK: map[string]any{"id": "TRG-000401"}, Values: map[string]any{"epic_id": "FID-e1f2g3", "issue_id": "FID-a1b2c3", "triager_role": "qa_lead", "source": "manual", "source_ref_id": "REV-000401", "verdict": "BLOCK", "reasoning_md": "project-docs/triage/TRG-000401.md", "created_at": "2026-05-27T10:19:00Z"}},
		{Op: "insert", Table: "links", PK: map[string]any{"from_type": "issue", "from_id": "FID-a1b2c3", "relation": "blocks", "to_type": "issue", "to_id": "FID-7c3e9f"}, Values: map[string]any{"created_at": "2026-05-27T10:18:00Z"}},
		{Op: "insert", Table: "links", PK: map[string]any{"from_type": "issue", "from_id": "FID-a1b2c3", "relation": "relates_to", "to_type": "issue", "to_id": "FID-7c3e9f"}, Values: map[string]any{"created_at": "2026-05-27T10:19:00Z"}},
	}})
	if _, err := journal.EnsureProjectionFresh(workspace, testIssueMigrationsDir(t)); err != nil {
		t.Fatal(err)
	}
}

// Kernel-only detail assertions leave with Concerns in 1.5c2.
