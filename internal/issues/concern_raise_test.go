package issues

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

func TestRaiseConcernCreatesOpenConcernAndFlipsDerivedFlags(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Raise concern against completion",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	before, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(before) error = %v", err)
	}
	if before.ConcernsRaised {
		t.Fatal("Get(before).ConcernsRaised = true, want false")
	}
	if before.OpenConcernsCount != 0 {
		t.Fatalf("Get(before).OpenConcernsCount = %d, want 0", before.OpenConcernsCount)
	}

	completion, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	result, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              created.ID,
		ActorRole:            "developer",
		RaisedInCompletionID: completion.ID,
		Summary:              "Potential race in rebuild parity",
		BodyMDPath:           stringPointer("project-docs/concerns/CON-new.md"),
	})
	if err != nil {
		t.Fatalf("RaiseConcern() error = %v", err)
	}

	if !strings.HasPrefix(result.ID, "CON-") {
		t.Fatalf("RaiseConcern().ID = %q, want CON- prefix", result.ID)
	}
	if result.IssueID != created.ID {
		t.Fatalf("RaiseConcern().IssueID = %q, want %q", result.IssueID, created.ID)
	}
	if result.RaisedInCompletionID != completion.ID {
		t.Fatalf("RaiseConcern().RaisedInCompletionID = %q, want %q", result.RaisedInCompletionID, completion.ID)
	}
	if result.RaisedByRole != "developer" {
		t.Fatalf("RaiseConcern().RaisedByRole = %q, want developer", result.RaisedByRole)
	}
	if result.Summary != "Potential race in rebuild parity" {
		t.Fatalf("RaiseConcern().Summary = %q, want summary", result.Summary)
	}
	if result.BodyMDPath == nil || *result.BodyMDPath != "project-docs/concerns/CON-new.md" {
		t.Fatalf("RaiseConcern().BodyMDPath = %#v, want body path", result.BodyMDPath)
	}
	if result.Status != "open" {
		t.Fatalf("RaiseConcern().Status = %q, want open", result.Status)
	}

	got, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            created.ID,
	})
	if err != nil {
		t.Fatalf("Get(after) error = %v", err)
	}
	if !got.ConcernsRaised {
		t.Fatal("Get(after).ConcernsRaised = false, want true")
	}
	if got.OpenConcernsCount != 1 {
		t.Fatalf("Get(after).OpenConcernsCount = %d, want 1", got.OpenConcernsCount)
	}
	if len(got.Concerns) != 1 {
		t.Fatalf("Get(after).Concerns len = %d, want 1", len(got.Concerns))
	}

	concern := got.Concerns[0]
	if concern.ID != result.ID {
		t.Fatalf("Get(after).Concerns[0].ID = %q, want %q", concern.ID, result.ID)
	}
	if concern.RaisedInCompletionID != completion.ID {
		t.Fatalf("Get(after).Concerns[0].RaisedInCompletionID = %q, want %q", concern.RaisedInCompletionID, completion.ID)
	}
	if concern.RaisedByRole != "developer" {
		t.Fatalf("Get(after).Concerns[0].RaisedByRole = %q, want developer", concern.RaisedByRole)
	}
	if concern.Summary != "Potential race in rebuild parity" {
		t.Fatalf("Get(after).Concerns[0].Summary = %q, want summary", concern.Summary)
	}
	if concern.BodyMDPath == nil || *concern.BodyMDPath != "project-docs/concerns/CON-new.md" {
		t.Fatalf("Get(after).Concerns[0].BodyMDPath = %#v, want body path", concern.BodyMDPath)
	}
	if concern.Status != "open" {
		t.Fatalf("Get(after).Concerns[0].Status = %q, want open", concern.Status)
	}

	listed, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("List() len = %d, want 1", len(listed))
	}
	if !listed[0].ConcernsRaised {
		t.Fatal("List()[0].ConcernsRaised = false, want true")
	}
	if listed[0].OpenConcernsCount != 1 {
		t.Fatalf("List()[0].OpenConcernsCount = %d, want 1", listed[0].OpenConcernsCount)
	}
}

func TestRaiseConcernRejectsBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(t *testing.T, workspace string) RaiseConcernOptions
		wantErr string
		verify  func(t *testing.T, workspace string)
	}{
		{
			name: "unknown role",
			prepare: func(t *testing.T, workspace string) RaiseConcernOptions {
				issueID, completionID := seedConcernRaiseIssueAndCompletion(t, workspace)
				return RaiseConcernOptions{
					WorkspaceRoot:        workspace,
					MigrationsDir:        testIssueMigrationsDir(t),
					IssueID:              issueID,
					ActorRole:            "intern",
					RaisedInCompletionID: completionID,
					Summary:              "Unknown role reject",
				}
			},
			wantErr: `role "intern" not in project profile`,
		},
		{
			name: "missing issue",
			prepare: func(t *testing.T, workspace string) RaiseConcernOptions {
				_, completionID := seedConcernRaiseIssueAndCompletion(t, workspace)
				return RaiseConcernOptions{
					WorkspaceRoot:        workspace,
					MigrationsDir:        testIssueMigrationsDir(t),
					IssueID:              "FID-missing",
					ActorRole:            "developer",
					RaisedInCompletionID: completionID,
					Summary:              "Missing issue reject",
				}
			},
			wantErr: `issue "FID-missing" not found`,
		},
		{
			name: "missing completion",
			prepare: func(t *testing.T, workspace string) RaiseConcernOptions {
				issueID, _ := seedConcernRaiseIssueAndCompletion(t, workspace)
				return RaiseConcernOptions{
					WorkspaceRoot:        workspace,
					MigrationsDir:        testIssueMigrationsDir(t),
					IssueID:              issueID,
					ActorRole:            "developer",
					RaisedInCompletionID: "TC-missing",
					Summary:              "Missing completion reject",
				}
			},
			wantErr: `raise concern: completion "TC-missing" not found`,
		},
		{
			name: "completion belongs to another issue",
			prepare: func(t *testing.T, workspace string) RaiseConcernOptions {
				firstIssueID, _ := seedConcernRaiseIssueAndCompletion(t, workspace)
				otherIssueID, otherCompletionID := seedConcernRaiseIssueAndCompletion(t, workspace)
				return RaiseConcernOptions{
					WorkspaceRoot:        workspace,
					MigrationsDir:        testIssueMigrationsDir(t),
					IssueID:              firstIssueID,
					ActorRole:            "developer",
					RaisedInCompletionID: otherCompletionID,
					Summary:              "Cross issue completion reject",
					BodyMDPath:           stringPointer(otherIssueID),
				}
			},
			wantErr: `does not belong to issue`,
			verify: func(t *testing.T, workspace string) {
				t.Helper()

				conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
				defer closeIssueTestDB(t, conn)
				assertConcernCount(t, conn, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newIssueWorkspace(t)
			options := tt.prepare(t, workspace)
			journalPath := journalPathForWorkspace(t, workspace)
			linesBefore := countJournalLines(journalPath)

			_, err := RaiseConcern(options)
			if err == nil {
				t.Fatalf("RaiseConcern() error = nil, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("RaiseConcern() error = %q, want substring %q", err.Error(), tt.wantErr)
			}
			if got := countJournalLines(journalPath); got != linesBefore {
				t.Fatalf("journal line count after rejection = %d, want unchanged %d", got, linesBefore)
			}

			conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
			defer closeIssueTestDB(t, conn)
			assertConcernCount(t, conn, 0)

			if tt.verify != nil {
				tt.verify(t, workspace)
			}
		})
	}
}

func seedConcernRaiseIssueAndCompletion(t *testing.T, workspace string) (string, string) {
	t.Helper()

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Concern raise task",
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	completion, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       created.ID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	return created.ID, completion.ID
}

func assertConcernCount(t *testing.T, conn *sql.DB, want int) {
	t.Helper()

	var count int
	if err := conn.QueryRow(`SELECT COUNT(*) FROM concerns`).Scan(&count); err != nil {
		t.Fatalf("count concerns: %v", err)
	}
	if count != want {
		t.Fatalf("concerns rows = %d, want %d", count, want)
	}
}
