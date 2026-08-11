package main

import "testing"

func TestV01ReadCommandsRejectRetiredFormats(t *testing.T) {
	workspace := t.TempDir()
	runInitCommand(t, workspace, "init")
	_, stderr := runFIDCommand(t, workspace, "issue", "create", "--kind", "task", "--title", "Read format retirement", "--by", "developer")
	issueID := parseCreatedIssueID(t, stderr)
	for _, args := range [][]string{{"issue", "get", issueID}, {"issue", "list"}, {"ready"}} {
		args = append(args, "--format", "json")
		if result := executeCommandForTest(workspace, args...); result.err == nil {
			t.Fatalf("retired format accepted for %v", args[:2])
		}
	}
}
