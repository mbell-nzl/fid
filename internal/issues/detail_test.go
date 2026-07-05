package issues

import (
	"database/sql"
	"errors"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetReturnsCentralizedDetailShape(t *testing.T) {
	t.Parallel()

	t.Run("empty issue returns stable empty sections and null due date", func(t *testing.T) {
		t.Parallel()

		workspace := newIssueWorkspace(t)

		created, err := Create(CreateOptions{
			WorkspaceRoot: workspace,
			MigrationsDir: testIssueMigrationsDir(t),
			Kind:          "task",
			Title:         "Render empty detail",
			ActorRole:     "developer",
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

		if got.DueDate != nil {
			t.Fatalf("Get() due_date = %#v, want nil", got.DueDate)
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
		if got.Completions == nil || len(got.Completions) != 0 {
			t.Fatalf("Get() completions = %#v, want empty array", got.Completions)
		}
		if got.Links == nil || len(got.Links) != 0 {
			t.Fatalf("Get() links = %#v, want empty array", got.Links)
		}
		if got.Claims == nil || len(got.Claims) != 0 {
			t.Fatalf("Get() claims = %#v, want empty array", got.Claims)
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

	t.Run("populated issue returns due date and structured sections", func(t *testing.T) {
		t.Parallel()

		workspace := newIssueWorkspace(t)
		conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
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

		if got.DueDate == nil || *got.DueDate != "2026-06-30" {
			t.Fatalf("Get() due_date = %#v, want 2026-06-30", got.DueDate)
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
		developerVerification, ok := got.AcceptanceCriteria[0].LatestVerifications["developer"].(map[string]any)
		if !ok {
			t.Fatalf("Get() latest_verifications[developer] = %#v, want map", got.AcceptanceCriteria[0].LatestVerifications["developer"])
		}
		if developerVerification["status"] != "pass" || developerVerification["verified_at"] != "2026-05-27T10:15:00Z" {
			t.Fatalf("Get() latest_verifications[developer] = %#v, want latest pass record", developerVerification)
		}
		if got.AcceptanceCriteria[1].LatestVerifications["developer"] != nil {
			t.Fatalf("Get() second criterion latest_verifications[developer] = %#v, want nil owner placeholder", got.AcceptanceCriteria[1].LatestVerifications["developer"])
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
		if len(got.Completions) != 2 || got.Completions[0].ID != "TC-000402" || got.Completions[1].ID != "TC-000401" {
			t.Fatalf("Get() completions = %#v, want latest-first completions", got.Completions)
		}
		if len(got.Links) != 2 || got.Links[0].Relation != "blocks" || got.Links[1].Relation != "implements" {
			t.Fatalf("Get() links = %#v, want blocks + implements", got.Links)
		}
		if got.Links[0].TargetTitle == nil || *got.Links[0].TargetTitle != "Schema migration runner" {
			t.Fatalf("Get() first link target_title = %#v, want Schema migration runner", got.Links[0].TargetTitle)
		}
		if len(got.Claims) != 1 || got.Claims[0].ClaimedBy != "developer" {
			t.Fatalf("Get() claims = %#v, want one live developer claim", got.Claims)
		}
		if len(got.Gates) != 2 {
			t.Fatalf("Get() gates = %#v, want two gate runs", got.Gates)
		}
		if got.Gates[0].ID != "GATE-000401" || got.Gates[0].AttemptNumber != 1 || got.Gates[0].Command != "go test ./..." {
			t.Fatalf("Get() first gate = %#v, want attempt 1 go test gate", got.Gates[0])
		}
		if got.Gates[0].AgentRunID == nil || *got.Gates[0].AgentRunID != "RUN-000401" {
			t.Fatalf("Get() first gate agent_run_id = %#v, want RUN-000401", got.Gates[0].AgentRunID)
		}
		if got.Gates[0].OutputExcerpt == nil || *got.Gates[0].OutputExcerpt != "ok" {
			t.Fatalf("Get() first gate output_excerpt = %#v, want ok", got.Gates[0].OutputExcerpt)
		}
		if got.Gates[1].ID != "GATE-000402" || got.Gates[1].AttemptNumber != 2 || got.Gates[1].Command != "golangci-lint run ./..." {
			t.Fatalf("Get() second gate = %#v, want attempt 2 golangci-lint gate", got.Gates[1])
		}
		if got.Gates[1].AgentRunID != nil {
			t.Fatalf("Get() second gate agent_run_id = %#v, want nil", got.Gates[1].AgentRunID)
		}
		if got.Gates[1].OutputExcerpt != nil {
			t.Fatalf("Get() second gate output_excerpt = %#v, want nil", got.Gates[1].OutputExcerpt)
		}
		if len(got.Completions) == 0 || len(got.Completions[0].GateResults) != 2 {
			t.Fatalf("Get() completions gate_results = %#v, want completion snapshot preserved separately", got.Completions)
		}
	})
}

func TestSectionValue(t *testing.T) {
	t.Parallel()

	issue := Issue{
		AcceptanceCriteria: []AcceptanceCriterion{{ID: "AC-1"}},
		Concerns:           []IssueConcern{{ID: "CON-1"}},
		Completions:        []IssueCompletion{{ID: "TC-1"}},
		Links:              []IssueLink{{Relation: "blocks", TargetID: "FID-2"}},
		Claims:             []IssueClaim{{RefID: "REQ-1", ClaimedBy: "developer"}},
		Gates:              []IssueGate{{ID: "GATE-1", Command: "go test ./..."}},
	}

	tests := []struct {
		name    string
		section string
		want    any
		wantErr string
	}{
		{
			name:    "acceptance section",
			section: "acceptance",
			want:    issue.AcceptanceCriteria,
		},
		{
			name:    "concerns section",
			section: "concerns",
			want:    issue.Concerns,
		},
		{
			name:    "completions section",
			section: "completions",
			want:    issue.Completions,
		},
		{
			name:    "links section",
			section: "links",
			want:    issue.Links,
		},
		{
			name:    "claims section",
			section: "claims",
			want:    issue.Claims,
		},
		{
			name:    "gates section",
			section: "gates",
			want:    issue.Gates,
		},
		{
			name:    "unknown section is rejected with supported names",
			section: " timeline ",
			wantErr: `unsupported issue section "timeline"; supported sections: acceptance, concerns, completions, links, claims, gates`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SectionValue(issue, tt.section)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("SectionValue() error = nil, want %q", tt.wantErr)
				}
				var unsupportedErr *UnsupportedSectionError
				if !errors.As(err, &unsupportedErr) {
					t.Fatalf("SectionValue() error = %T, want *UnsupportedSectionError", err)
				}
				if unsupportedErr.Section != "timeline" {
					t.Fatalf("UnsupportedSectionError.Section = %q, want timeline", unsupportedErr.Section)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("SectionValue() error = %q, want %q", err.Error(), tt.wantErr)
				}
				if !equalStrings(SupportedSections(), []string{"acceptance", "concerns", "completions", "links", "claims", "gates"}) {
					t.Fatalf("SupportedSections() = %#v", SupportedSections())
				}
				return
			}
			if err != nil {
				t.Fatalf("SectionValue() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SectionValue() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestIssueParentTitle(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
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

	statements := []struct {
		query string
		args  []any
	}{
		{
			query: `
				INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-e1f2g3",
				"epic",
				"Core persistence layer",
				"todo",
				"epic-default",
				"developer",
				"2026-05-25T08:00:00Z",
				"2026-05-25T08:00:00Z",
			},
		},
		{
			query: `
				INSERT INTO epic_metadata (issue_id, scheduling)
				VALUES (?, ?)
			`,
			args: []any{"FID-e1f2g3", "sequential"},
		},
		{
			query: `
				INSERT INTO issues (
					id, kind, parent_id, title, workflow_id, status, priority, sequence,
					spec_md_path, due_date, owner_role, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-a1b2c3",
				"task",
				"FID-e1f2g3",
				"Implement journal append with flush",
				"task-default",
				"in_progress",
				1,
				3,
				"project-docs/docs/phases/phase-1/tasks/FID-a1b2c3.md",
				"2026-06-30",
				"developer",
				"2026-05-25T09:00:00Z",
				"2026-05-27T10:20:00Z",
			},
		},
		{
			query: `
				INSERT INTO issues (id, kind, title, status, workflow_id, owner_role, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"FID-7c3e9f",
				"task",
				"Schema migration runner",
				"todo",
				"task-default",
				"developer",
				"2026-05-25T09:05:00Z",
				"2026-05-25T09:05:00Z",
			},
		},
		{
			query: `
				INSERT INTO requirements (id, type, summary, source_doc, source_section, status)
				VALUES (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"PRD-ISSUE-001",
				"PRD",
				"Implement journal append with flush",
				"01_PRD.md",
				"ISSUE",
				"active",
			},
		},
		{
			query: `INSERT INTO labels (name) VALUES (?), (?)`,
			args:  []any{"core", "journal"},
		},
		{
			query: `INSERT INTO issue_labels (issue_id, label) VALUES (?, ?), (?, ?)`,
			args:  []any{"FID-a1b2c3", "core", "FID-a1b2c3", "journal"},
		},
		{
			query: `
				INSERT INTO acceptance_criteria (id, issue_id, ordinal, text)
				VALUES (?, ?, ?, ?), (?, ?, ?, ?)
			`,
			args: []any{
				"AC-d4e5f6",
				"FID-a1b2c3",
				1,
				"Journal file is flushed to disk before returning success",
				"AC-e5f6g7",
				"FID-a1b2c3",
				2,
				"Partial write leaves journal intact",
			},
		},
		{
			query: `
				INSERT INTO task_completions (
					id, issue_id, attempt_number, developer_role, model, status, actual_tokens, gate_results, human_intervention, completed_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"TC-000401",
				"FID-a1b2c3",
				1,
				"developer",
				"gpt-5",
				"dev_complete",
				1440,
				`["go test ./..."]`,
				false,
				"2026-05-27T09:30:00Z",
				"TC-000402",
				"FID-a1b2c3",
				2,
				"developer",
				"gpt-5",
				"qa_pass",
				1550,
				`[{"cmd":"go test ./...","exit_code":0,"ran_at":"2026-05-27T10:10:00Z","output_excerpt":"ok"},{"cmd":"go test -race ./...","exit_code":0,"ran_at":"2026-05-27T10:12:00Z","output_excerpt":"race ok"}]`,
				false,
				"2026-05-27T10:20:00Z",
			},
		},
		{
			query: `
				INSERT INTO agent_runs (id, issue_id, role, model, started_at)
				VALUES (?, ?, ?, ?, ?)
			`,
			args: []any{
				"RUN-000401",
				"FID-a1b2c3",
				"developer",
				"gpt-5",
				"2026-05-27T09:00:00Z",
			},
		},
		{
			query: `
				INSERT INTO gate_runs (
					id, issue_id, attempt_number, agent_run_id, recorded_by_role, command, exit_code, output_excerpt, ran_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"GATE-000401",
				"FID-a1b2c3",
				1,
				"RUN-000401",
				"developer",
				"go test ./...",
				0,
				"ok",
				"2026-05-27T09:20:00Z",
				"GATE-000402",
				"FID-a1b2c3",
				2,
				nil,
				"developer",
				"golangci-lint run ./...",
				1,
				nil,
				"2026-05-27T10:25:00Z",
			},
		},
		{
			query: `
				INSERT INTO ac_verifications (id, ac_id, verifier_role, status, completion_id, verified_at)
				VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"ACV-older-dev",
				"AC-d4e5f6",
				"developer",
				"fail",
				"TC-000401",
				"2026-05-27T09:35:00Z",
				"ACV-latest-dev",
				"AC-d4e5f6",
				"developer",
				"pass",
				"TC-000402",
				"2026-05-27T10:15:00Z",
			},
		},
		{
			query: `
				INSERT INTO concerns (
					id, issue_id, raised_in_completion_id, raised_by_role, summary, body_md_path, status, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"CON-000401",
				"FID-a1b2c3",
				"TC-000402",
				"developer",
				"Need triage before closing",
				"project-docs/concerns/CON-000401.md",
				"open",
				"2026-05-27T10:18:00Z",
				"2026-05-27T10:18:00Z",
			},
		},
		{
			query: `
				INSERT INTO triage_events (
					id, epic_id, issue_id, triager_role, source, source_ref_id, verdict, reasoning_md, created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"TRG-000401",
				"FID-e1f2g3",
				"FID-a1b2c3",
				"qa_lead",
				"task_completion",
				"TC-000402",
				"BLOCK",
				"project-docs/triage/TRG-000401.md",
				"2026-05-27T10:19:00Z",
			},
		},
		{
			query: `
				INSERT INTO links (from_type, from_id, relation, to_type, to_id, created_at)
				VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"issue",
				"FID-a1b2c3",
				"blocks",
				"issue",
				"FID-7c3e9f",
				"2026-05-27T10:18:00Z",
				"issue",
				"FID-a1b2c3",
				"implements",
				"requirement",
				"PRD-ISSUE-001",
				"2026-05-27T10:19:00Z",
			},
		},
		{
			query: `
				INSERT INTO claims (ref_type, ref_id, claimed_by, claimed_at, ttl_seconds, note)
				VALUES (?, ?, ?, ?, ?, ?)
			`,
			args: []any{
				"issue",
				"FID-a1b2c3",
				"developer",
				"2026-05-27T09:50:00Z",
				1800,
				"Working the flush path",
			},
		},
	}

	for _, statement := range statements {
		if _, err := conn.Exec(statement.query, statement.args...); err != nil {
			t.Fatalf("Exec(%q) error = %v", statement.query, err)
		}
	}
}
