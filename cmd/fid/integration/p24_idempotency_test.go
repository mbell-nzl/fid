//go:build integration

package integration

import (
	"reflect"
	"testing"
)

func TestP24IssueCreateConflictNamesChangedTitle(t *testing.T) {
	t.Parallel()
	binary, workspace := integrationFIDBinary(t), t.TempDir()
	leanFlagFixture(t, binary, workspace)
	p21Run(t, binary, workspace, "--json", "issue", "create", "--kind", "task", "--title", "original title", "--actor", "agent:developer", "--idempotency-key", "p22-corpus-issue")
	before := p21JournalEventCount(t, workspace)
	stdout, stderr, err := closureRunCommand(workspace, binary, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "different title", "--actor", "agent:developer", "--idempotency-key", "p22-corpus-issue")
	if err == nil || processExitCode(t, err) != 4 || stdout != "" {
		t.Fatalf("issue.create conflict = %v stdout=%q stderr=%q", err, stdout, stderr)
	}
	if got := p21Path(t, p23Envelope(t, stderr), "error", "details", "differing_fields"); !reflect.DeepEqual(got, []any{"--title"}) {
		t.Fatalf("issue.create differing_fields = %#v", got)
	}
	if got := p21JournalEventCount(t, workspace); got != before {
		t.Fatalf("issue.create conflict appended Event: %d -> %d", before, got)
	}
}
