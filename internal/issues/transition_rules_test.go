package issues

import (
	"path/filepath"
	"slices"
	"testing"

	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/workflow"
)

func TestEvaluateTransitionRuleConcernsAddressedUsesRealConcernWrites(t *testing.T) {
	t.Parallel()

	workspace := newIssueWorkspace(t)
	issueID, raiseCompletionID := seedConcernRaiseIssueAndCompletion(t, workspace)

	firstConcern, err := RaiseConcern(RaiseConcernOptions{
		WorkspaceRoot:        workspace,
		MigrationsDir:        testIssueMigrationsDir(t),
		IssueID:              issueID,
		ActorRole:            "developer",
		RaisedInCompletionID: raiseCompletionID,
		Summary:              "Open concern one",
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
		Summary:              "Open concern two",
	})
	if err != nil {
		t.Fatalf("RaiseConcern(second) error = %v", err)
	}

	conn, err := fiddb.Open(filepath.Join(workspace, ".fid", "fid.sqlite"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer closeIssueTestDB(t, conn)

	current, err := loadIssueTransitionState(conn, issueID)
	if err != nil {
		t.Fatalf("loadIssueTransitionState() error = %v", err)
	}

	openResult, err := evaluateTransitionRule(conn, current, workflow.RuleReference{Name: workflow.RuleConcernsAddressed})
	if err != nil {
		t.Fatalf("evaluateTransitionRule(open) error = %v", err)
	}
	if openResult.Satisfied {
		t.Fatal("evaluateTransitionRule(open).Satisfied = true, want false")
	}
	if len(openResult.Unmet) != 1 {
		t.Fatalf("evaluateTransitionRule(open).Unmet len = %d, want 1", len(openResult.Unmet))
	}
	if got := openResult.Unmet[0].Rule; got != workflow.RuleConcernsAddressed {
		t.Fatalf("evaluateTransitionRule(open).Unmet[0].Rule = %q, want %q", got, workflow.RuleConcernsAddressed)
	}
	wantOpenConcernIDs := []string{firstConcern.ID, secondConcern.ID}
	slices.Sort(wantOpenConcernIDs)
	if got := openResult.Unmet[0].EntityIDs; !slices.Equal(got, wantOpenConcernIDs) {
		t.Fatalf("evaluateTransitionRule(open).Unmet[0].EntityIDs = %#v, want %#v", got, wantOpenConcernIDs)
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
	if _, err := AddressConcern(AddressConcernOptions{
		WorkspaceRoot:           workspace,
		MigrationsDir:           testIssueMigrationsDir(t),
		ConcernID:               firstConcern.ID,
		ActorRole:               "developer",
		AddressedInCompletionID: addressCompletion.ID,
	}); err != nil {
		t.Fatalf("AddressConcern(first) error = %v", err)
	}

	partialResult, err := evaluateTransitionRule(conn, current, workflow.RuleReference{Name: workflow.RuleConcernsAddressed})
	if err != nil {
		t.Fatalf("evaluateTransitionRule(partial) error = %v", err)
	}
	if partialResult.Satisfied {
		t.Fatal("evaluateTransitionRule(partial).Satisfied = true, want false")
	}
	if got := partialResult.Unmet[0].EntityIDs; len(got) != 1 || got[0] != secondConcern.ID {
		t.Fatalf("evaluateTransitionRule(partial).Unmet[0].EntityIDs = %#v, want [%q]", got, secondConcern.ID)
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

	clearedResult, err := evaluateTransitionRule(conn, current, workflow.RuleReference{Name: workflow.RuleConcernsAddressed})
	if err != nil {
		t.Fatalf("evaluateTransitionRule(cleared) error = %v", err)
	}
	if !clearedResult.Satisfied {
		t.Fatalf("evaluateTransitionRule(cleared).Satisfied = false, want true (unmet=%#v)", clearedResult.Unmet)
	}
	if len(clearedResult.Unmet) != 0 {
		t.Fatalf("evaluateTransitionRule(cleared).Unmet = %#v, want empty", clearedResult.Unmet)
	}
}
