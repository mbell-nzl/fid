package main

import (
	"encoding/json"
	"strings"
	"testing"
)

// showJSON is the parts of fid --json show's result this test reads.
type showJSON struct {
	Issue struct {
		Priority           *int `json:"priority"`
		AcceptanceCriteria []struct {
			Text string `json:"text"`
		} `json:"acceptance_criteria"`
	} `json:"issue"`
	Epic *struct {
		ID       string `json:"id"`
		Progress struct {
			Done  int `json:"done"`
			Total int `json:"total"`
		} `json:"progress"`
	} `json:"epic"`
	Blocks []struct {
		ID string `json:"id"`
	} `json:"blocks"`
	Children []struct {
		ID string `json:"id"`
	} `json:"children"`
}

// Show reads work fid wrote through its journal, rebuilding a missing
// projection first: an epic, a started task under it with a criterion and
// a subtask, and a task it blocks.
func TestShowShowsTheJournalsWork(t *testing.T) {
	t.Parallel()
	fixture := newJournalStatusFixture(t)
	create := func(args ...string) string {
		t.Helper()
		_, stderr := runFIDCommand(t, fixture.workspace, append([]string{"issue", "create", "--actor", "human:mark"}, args...)...)
		return parseCreatedIssueID(t, stderr)
	}
	epic := create("--kind", "epic", "--title", "Operating experience", "--scheduling", "prioritised")
	task := create("--kind", "task", "--title", "Export validation", "--epic", epic, "--priority", "1")
	subtask := create("--kind", "subtask", "--title", "Validate member digests", "--parent", task)
	blocked := create("--kind", "task", "--title", "Retry a rejected backup push", "--epic", epic)
	runFIDCommand(t, fixture.workspace, "issue", "ac", "add", task, "--ordinal", "1", "--text", "A truncated archive is refused", "--actor", "human:mark")
	runFIDCommand(t, fixture.workspace, "link", "add", task, "blocks", blocked, "--from-type", "issue", "--to-type", "issue", "--actor", "human:mark")
	fixture.start(task)
	fixture.removeProjection()

	result := fixture.run("--json", "show", task)
	var envelope struct {
		Result showJSON `json:"result"`
	}
	if result.err != nil || result.stderr != "" || json.Unmarshal([]byte(result.stdout), &envelope) != nil {
		t.Fatalf("show --json = %v, stderr %q:\n%s", result.err, result.stderr, result.stdout)
	}
	got := envelope.Result
	if got.Epic == nil || got.Epic.ID != epic || got.Epic.Progress.Done != 0 || got.Epic.Progress.Total != 2 ||
		len(got.Blocks) != 1 || got.Blocks[0].ID != blocked || len(got.Children) != 1 || got.Children[0].ID != subtask ||
		got.Issue.Priority == nil || *got.Issue.Priority != 1 ||
		len(got.Issue.AcceptanceCriteria) != 1 || got.Issue.AcceptanceCriteria[0].Text != "A truncated archive is refused" {
		t.Fatalf("show --json doesn't hold the journal's work:\n%s", result.stdout)
	}
	text := fixture.run("show", task)
	for _, want := range []string{"● Export validation", "◆ " + epic, "○ " + blocked, "○ " + subtask, "1  A truncated archive is refused"} {
		if text.err != nil || !strings.Contains(text.stdout, want) {
			t.Fatalf("show = %v, text doesn't show %q:\n%s", text.err, want, text.stdout)
		}
	}
}
