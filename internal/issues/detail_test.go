package issues

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
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
		if got.Concerns == nil || len(got.Concerns) != 0 {
			t.Fatalf("Get() concerns = %#v, want empty array", got.Concerns)
		}
		if got.Links == nil || len(got.Links) != 0 {
			t.Fatalf("Get() links = %#v, want empty array", got.Links)
		}
		if got.ConcernsRaised {
			t.Fatal("Get() concerns_raised = true, want false")
		}
		if got.OpenConcernsCount != 0 {
			t.Fatalf("Get() open_concerns_count = %d, want 0", got.OpenConcernsCount)
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

		seedIssueDetailFixture(t, conn)

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
		if !got.ConcernsRaised {
			t.Fatal("Get() concerns_raised = false, want true")
		}
		if got.OpenConcernsCount != 1 {
			t.Fatalf("Get() open_concerns_count = %d, want 1", got.OpenConcernsCount)
		}
		if got.TriageStatus == nil || *got.TriageStatus != "blocked" {
			t.Fatalf("Get() triage_status = %#v, want blocked", got.TriageStatus)
		}
		if len(got.Concerns) != 1 || got.Concerns[0].ID != "CON-000401" {
			t.Fatalf("Get() concerns = %#v, want one open concern", got.Concerns)
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

func seedIssueDetailFixture(t *testing.T, conn *sql.DB) {
	t.Helper()
	statements := []string{
		`INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
		 VALUES ('FID-e1f2g3','epic','Core persistence layer','todo','epic-default','developer','2026-05-25T08:00:00Z','2026-05-25T08:00:00Z')`,
		`INSERT INTO epic_metadata (issue_id, scheduling) VALUES ('FID-e1f2g3', 'sequential')`,
		`INSERT INTO issues (id, kind, parent_id, title, workflow_id, status, priority, sequence, owner_role, created_at, updated_at)
		 VALUES ('FID-a1b2c3','task','FID-e1f2g3','Implement journal append with flush','task-default','in_progress',1,3,'developer','2026-05-25T09:00:00Z','2026-05-27T10:20:00Z')`,
		`INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
		 VALUES ('FID-7c3e9f','task','Schema migration runner','todo','task-default','developer','2026-05-25T09:05:00Z','2026-05-25T09:05:00Z')`,
		`INSERT INTO labels (name) VALUES ('core'), ('journal')`,
		`INSERT INTO issue_labels (issue_id, label) VALUES ('FID-a1b2c3','core'), ('FID-a1b2c3','journal')`,
		`INSERT INTO acceptance_criteria (id, issue_id, ordinal, text)
		 VALUES ('AC-d4e5f6','FID-a1b2c3',1,'Journal file is flushed to disk before returning success'),
		        ('AC-e5f6g7','FID-a1b2c3',2,'Partial write leaves journal intact')`,
		`INSERT INTO concerns (id,subject_type,subject_id,actor_id,scope_json,severity,summary,reasoning,references_json,status,created_at,created_event_id)
		 VALUES ('CON-000401','issue','FID-a1b2c3','developer','{"paths":[]}','major','Need triage before closing','Needs review','[]','open','2026-05-27T10:18:00Z','EVT-concern-detail')`,
		`INSERT INTO triage_events (id,epic_id,issue_id,triager_role,source,source_ref_id,verdict,reasoning_md,created_at)
		 VALUES ('TRG-000401','FID-e1f2g3','FID-a1b2c3','qa_lead','manual','REV-000401','BLOCK','project-docs/triage/TRG-000401.md','2026-05-27T10:19:00Z')`,
		`INSERT INTO links (from_type,from_id,relation,to_type,to_id,created_at)
		 VALUES ('issue','FID-a1b2c3','blocks','issue','FID-7c3e9f','2026-05-27T10:18:00Z'),
		        ('issue','FID-a1b2c3','relates_to','issue','FID-7c3e9f','2026-05-27T10:19:00Z')`,
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			t.Fatalf("seed Issue detail with %q: %v", statement, err)
		}
	}
}
