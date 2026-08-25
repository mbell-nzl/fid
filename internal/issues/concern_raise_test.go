package issues

import "testing"

func TestRaiseConcernCreatesOpenConcernAndFlipsDerivedFlags(t *testing.T) {
	assertJudgmentUnavailable(t, "human.outcomes.v1", func() error {
		_, err := RaiseConcern(RaiseConcernOptions{})
		return err
	})
}

func TestRaiseConcernRejectsBeforeAppend(t *testing.T) {
	assertJudgmentUnavailable(t, "human.outcomes.v1", func() error {
		_, err := RaiseConcern(RaiseConcernOptions{IssueID: "FID-1", Summary: "deferred judgment"})
		return err
	})
}

func seedConcernRaiseIssueAndCompletion(t *testing.T, workspace string) (string, string) {
	t.Helper()

	created, err := Create(CreateOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testIssueMigrationsDir(t),
		Kind:          "task",
		Title:         "Concern history task",
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
