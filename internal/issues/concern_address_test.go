package issues

import (
	"database/sql"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddressConcernUpdatesConcernInPlaceAndClearsDerivedFlags(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	issueID, firstCompletionID := seedConcernRaiseIssueAndCompletion(t, workspace)

	firstConcern, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              issueID,
		ActorRole:            "developer",
		RaisedInCompletionID: firstCompletionID,
		Summary:              "First open concern",
	})
	if err != nil {
		t.Fatalf("RaiseConcern(first) error = %v", err)
	}

	secondCompletion, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete(second) error = %v", err)
	}

	secondConcern, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              issueID,
		ActorRole:            "developer",
		RaisedInCompletionID: secondCompletion.ID,
		Summary:              "Second open concern",
	})
	if err != nil {
		t.Fatalf("RaiseConcern(second) error = %v", err)
	}

	addressCompletion, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete(address first) error = %v", err)
	}

	result, err := AddressConcern(AddressConcernOptions{
		WorkspaceRoot:           workspace,
		MigrationsDir:           testIssueMigrationsDir(t),
		ConcernID:               firstConcern.ID,
		ActorRole:               "developer",
		AddressedInCompletionID: addressCompletion.ID,
	})
	if err != nil {
		t.Fatalf("AddressConcern(first) error = %v", err)
	}

	if result.ID != firstConcern.ID {
		t.Fatalf("AddressConcern().ID = %q, want %q", result.ID, firstConcern.ID)
	}
	if result.IssueID != issueID {
		t.Fatalf("AddressConcern().IssueID = %q, want %q", result.IssueID, issueID)
	}
	if result.AddressingCompletionID != addressCompletion.ID {
		t.Fatalf("AddressConcern().AddressingCompletionID = %q, want %q", result.AddressingCompletionID, addressCompletion.ID)
	}
	if result.Status != "addressed" {
		t.Fatalf("AddressConcern().Status = %q, want addressed", result.Status)
	}

	conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
	defer closeIssueTestDB(t, conn)

	assertConcernCount(t, conn, 2)

	storedFirst := readConcernRow(t, conn, firstConcern.ID)
	if storedFirst.Status != "addressed" {
		t.Fatalf("storedFirst.Status = %q, want addressed", storedFirst.Status)
	}
	if storedFirst.AddressedInCompletionID == nil || *storedFirst.AddressedInCompletionID != addressCompletion.ID {
		t.Fatalf("storedFirst.AddressedInCompletionID = %#v, want %q", storedFirst.AddressedInCompletionID, addressCompletion.ID)
	}
	if storedFirst.UpdatedAt == nil || firstConcern.UpdatedAt >= *storedFirst.UpdatedAt {
		t.Fatalf("storedFirst.UpdatedAt = %#v, want after %q", storedFirst.UpdatedAt, firstConcern.UpdatedAt)
	}

	storedSecond := readConcernRow(t, conn, secondConcern.ID)
	if storedSecond.Status != "open" {
		t.Fatalf("storedSecond.Status = %q, want open", storedSecond.Status)
	}

	afterFirstAddress, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            issueID,
	})
	if err != nil {
		t.Fatalf("Get(after first address) error = %v", err)
	}
	if !afterFirstAddress.ConcernsRaised {
		t.Fatal("Get(after first address).ConcernsRaised = false, want true")
	}
	if afterFirstAddress.OpenConcernsCount != 1 {
		t.Fatalf("Get(after first address).OpenConcernsCount = %d, want 1", afterFirstAddress.OpenConcernsCount)
	}

	listAfterFirstAddress, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("List(after first address) error = %v", err)
	}
	if len(listAfterFirstAddress) != 1 {
		t.Fatalf("List(after first address) len = %d, want 1", len(listAfterFirstAddress))
	}
	if !listAfterFirstAddress[0].ConcernsRaised {
		t.Fatal("List(after first address)[0].ConcernsRaised = false, want true")
	}
	if listAfterFirstAddress[0].OpenConcernsCount != 1 {
		t.Fatalf("List(after first address)[0].OpenConcernsCount = %d, want 1", listAfterFirstAddress[0].OpenConcernsCount)
	}

	finalAddressCompletion, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete(address second) error = %v", err)
	}

	if _, err := AddressConcern(AddressConcernOptions{
		WorkspaceRoot:           workspace,
		MigrationsDir:           testIssueMigrationsDir(t),
		ConcernID:               secondConcern.ID,
		ActorRole:               "developer",
		AddressedInCompletionID: finalAddressCompletion.ID,
	}); err != nil {
		t.Fatalf("AddressConcern(second) error = %v", err)
	}

	cleared, err := Get(GetOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		ID:            issueID,
	})
	if err != nil {
		t.Fatalf("Get(after final address) error = %v", err)
	}
	if cleared.ConcernsRaised {
		t.Fatal("Get(after final address).ConcernsRaised = true, want false")
	}
	if cleared.OpenConcernsCount != 0 {
		t.Fatalf("Get(after final address).OpenConcernsCount = %d, want 0", cleared.OpenConcernsCount)
	}

	listCleared, err := List(ListOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
	})
	if err != nil {
		t.Fatalf("List(after final address) error = %v", err)
	}
	if len(listCleared) != 1 {
		t.Fatalf("List(after final address) len = %d, want 1", len(listCleared))
	}
	if listCleared[0].ConcernsRaised {
		t.Fatal("List(after final address)[0].ConcernsRaised = true, want false")
	}
	if listCleared[0].OpenConcernsCount != 0 {
		t.Fatalf("List(after final address)[0].OpenConcernsCount = %d, want 0", listCleared[0].OpenConcernsCount)
	}
}

func TestAddressConcernRejectsBeforeAppend(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(t *testing.T, workspace string) AddressConcernOptions
		wantErr string
		verify  func(t *testing.T, workspace string, conn *sql.DB)
	}{
		{
			name: "unknown role",
			prepare: func(t *testing.T, workspace string) AddressConcernOptions {
				concernID, _, addressCompletionID := seedAddressConcernScenario(t, workspace)
				return AddressConcernOptions{
					WorkspaceRoot:           workspace,
					MigrationsDir:           testIssueMigrationsDir(t),
					ConcernID:               concernID,
					ActorRole:               "intern",
					AddressedInCompletionID: addressCompletionID,
				}
			},
			wantErr: `role "intern" not in project profile`,
		},
		{
			name: "missing concern",
			prepare: func(t *testing.T, workspace string) AddressConcernOptions {
				_, _, addressCompletionID := seedAddressConcernScenario(t, workspace)
				return AddressConcernOptions{
					WorkspaceRoot:           workspace,
					MigrationsDir:           testIssueMigrationsDir(t),
					ConcernID:               "CON-missing",
					ActorRole:               "developer",
					AddressedInCompletionID: addressCompletionID,
				}
			},
			wantErr: `address concern: concern "CON-missing" not found`,
		},
		{
			name: "already addressed concern",
			prepare: func(t *testing.T, workspace string) AddressConcernOptions {
				concernID, _, addressCompletionID := seedAddressConcernScenario(t, workspace)
				if _, err := AddressConcern(AddressConcernOptions{
					WorkspaceRoot:           workspace,
					MigrationsDir:           testIssueMigrationsDir(t),
					ConcernID:               concernID,
					ActorRole:               "developer",
					AddressedInCompletionID: addressCompletionID,
				}); err != nil {
					t.Fatalf("AddressConcern(seed addressed) error = %v", err)
				}
				return AddressConcernOptions{
					WorkspaceRoot:           workspace,
					MigrationsDir:           testIssueMigrationsDir(t),
					ConcernID:               concernID,
					ActorRole:               "developer",
					AddressedInCompletionID: addressCompletionID,
				}
			},
			wantErr: `address concern: concern`,
			verify: func(t *testing.T, workspace string, conn *sql.DB) {
				t.Helper()
				assertConcernCount(t, conn, 1)
			},
		},
		{
			name: "missing completion",
			prepare: func(t *testing.T, workspace string) AddressConcernOptions {
				concernID, _, _ := seedAddressConcernScenario(t, workspace)
				return AddressConcernOptions{
					WorkspaceRoot:           workspace,
					MigrationsDir:           testIssueMigrationsDir(t),
					ConcernID:               concernID,
					ActorRole:               "developer",
					AddressedInCompletionID: "TC-missing",
				}
			},
			wantErr: `address concern: completion "TC-missing" not found`,
		},
		{
			name: "completion belongs to another issue",
			prepare: func(t *testing.T, workspace string) AddressConcernOptions {
				concernID, issueID, _ := seedAddressConcernScenario(t, workspace)
				otherIssueID, otherCompletionID := seedConcernRaiseIssueAndCompletion(t, workspace)
				if otherIssueID == issueID {
					t.Fatalf("other issue id = %q, want different from %q", otherIssueID, issueID)
				}
				return AddressConcernOptions{
					WorkspaceRoot:           workspace,
					MigrationsDir:           testIssueMigrationsDir(t),
					ConcernID:               concernID,
					ActorRole:               "developer",
					AddressedInCompletionID: otherCompletionID,
				}
			},
			wantErr: `does not belong to issue`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace := newIssueWorkspace(t)
			options := tt.prepare(t, workspace)
			journalPath := journalPathForWorkspace(t, workspace)
			linesBefore := countJournalLines(journalPath)

			_, err := AddressConcern(options)
			if err == nil {
				t.Fatalf("AddressConcern() error = nil, want %q", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("AddressConcern() error = %q, want substring %q", err.Error(), tt.wantErr)
			}
			if got := countJournalLines(journalPath); got != linesBefore {
				t.Fatalf("journal line count after rejection = %d, want unchanged %d", got, linesBefore)
			}

			conn := openIssueDB(t, filepath.Join(workspace, ".fid", "fid.sqlite"))
			defer closeIssueTestDB(t, conn)
			if tt.verify != nil {
				tt.verify(t, workspace, conn)
			}
		})
	}
}

func seedAddressConcernScenario(t *testing.T, workspace string) (concernID, issueID, addressCompletionID string) {
	t.Helper()

	issueID, raiseCompletionID := seedConcernRaiseIssueAndCompletion(t, workspace)
	raised, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              issueID,
		ActorRole:            "developer",
		RaisedInCompletionID: raiseCompletionID,
		Summary:              "Address me",
	})
	if err != nil {
		t.Fatalf("RaiseConcern() error = %v", err)
	}

	addressCompletion, err := Complete(CompleteOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		IssueID:       issueID,
		ActorRole:     "developer",
	})
	if err != nil {
		t.Fatalf("Complete(address seed) error = %v", err)
	}

	return raised.ID, issueID, addressCompletion.ID
}

type storedConcernRow struct {
	Status                  string
	AddressedInCompletionID *string
	UpdatedAt               *string
}

func readConcernRow(t *testing.T, conn *sql.DB, concernID string) storedConcernRow {
	t.Helper()

	var (
		row                     storedConcernRow
		addressedInCompletionID sql.NullString
		updatedAt               sql.NullString
	)
	if err := conn.QueryRow(`
		SELECT status, addressed_in_completion_id, updated_at
		FROM concerns
		WHERE id = ?
	`, concernID).Scan(&row.Status, &addressedInCompletionID, &updatedAt); err != nil {
		t.Fatalf("query concern %q: %v", concernID, err)
	}
	row.AddressedInCompletionID = nullStringPointer(addressedInCompletionID)
	row.UpdatedAt = nullStringPointer(updatedAt)
	return row
}
