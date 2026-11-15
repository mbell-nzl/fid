package github

import (
	"fmt"
	"testing"
)

const mergeSHA = "4517e22de7eb09fad1741dab73a38b28419749a3"
const headSHA = "14b530c9ac45f9fdd705b1bef3ede2ccf9e1297c"
const treeSHA = "3ceee6a1141c0f7d49839bc2ff81e14b325c523f"

func sampleRun(id, attempt, workflow int64, status, conclusion string) WorkflowRun {
	return WorkflowRun{ID: id, Attempt: attempt, WorkflowID: workflow, Name: "ci", Path: ".github/workflows/ci.yml", Event: "push", HeadSHA: mergeSHA, Status: status, Conclusion: conclusion, URL: fmt.Sprintf("https://github.com/mbell-nzl/fid/actions/runs/%d", id)}
}

func TestCIVerdict(t *testing.T) {
	t.Parallel()
	success := sampleRun(1, 1, 10, "completed", "success")
	failed := sampleRun(1, 1, 10, "completed", "failure")
	queued := sampleRun(1, 1, 10, "queued", "")
	progress := sampleRun(1, 1, 10, "in_progress", "")
	skipped := sampleRun(1, 1, 10, "completed", "skipped")
	wrongHead := success
	wrongHead.HeadSHA = headSHA
	pullRequest := success
	pullRequest.Event = "pull_request"
	secondFailure := failed
	secondFailure.Attempt = 2
	secondSuccess := success
	secondSuccess.Attempt = 2
	other := success
	other.ID = 2
	other.WorkflowID = 20
	newer := success
	newer.ID = 2
	for _, tt := range []struct {
		name     string
		runs     []WorkflowRun
		verdict  Verdict
		ignored  []IgnoreReason
		failures int
	}{
		{"none", nil, Waiting, nil, 0},
		{"queued", []WorkflowRun{queued}, Waiting, []IgnoreReason{""}, 0},
		{"in progress", []WorkflowRun{progress}, Waiting, []IgnoreReason{""}, 0},
		{"success", []WorkflowRun{success}, Green, []IgnoreReason{""}, 0},
		{"failure", []WorkflowRun{failed}, NotGreen, []IgnoreReason{""}, 1},
		{"only skipped", []WorkflowRun{skipped}, Waiting, []IgnoreReason{Skipped}, 0},
		{"mixed workflows", []WorkflowRun{failed, other}, NotGreen, []IgnoreReason{"", ""}, 1},
		{"passing workflows", []WorkflowRun{success, other}, Green, []IgnoreReason{"", ""}, 0},
		{"failure after pass", []WorkflowRun{success, secondFailure}, NotGreen, []IgnoreReason{Superseded, ""}, 1},
		{"pass after failure", []WorkflowRun{secondSuccess, failed}, Green, []IgnoreReason{"", Superseded}, 0},
		{"new workflow execution", []WorkflowRun{secondFailure, newer}, Green, []IgnoreReason{Superseded, ""}, 0},
		{"wrong head", []WorkflowRun{wrongHead}, Waiting, []IgnoreReason{WrongHead}, 0},
		{"PR event", []WorkflowRun{pullRequest}, Waiting, []IgnoreReason{WrongEvent}, 0},
		{"skip and pass", []WorkflowRun{skipped, other}, Green, []IgnoreReason{Skipped, ""}, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := EvaluateCI(mergeSHA, tt.runs)
			if got.Verdict != tt.verdict || len(got.NotGreen) != tt.failures || len(got.Runs) != len(tt.runs) {
				t.Fatalf("CI = %#v", got)
			}
			for i, reason := range tt.ignored {
				if got.Runs[i].Ignored != reason {
					t.Errorf("workflow %d ignored=%q want %q", i, got.Runs[i].Ignored, reason)
				}
			}
		})
	}
	for _, conclusion := range []string{"cancelled", "timed_out", "action_required", "neutral", "stale", "startup_failure", "unknown_future_value"} {
		t.Run(conclusion, func(t *testing.T) {
			t.Parallel()
			got := EvaluateCI(mergeSHA, []WorkflowRun{sampleRun(1, 1, 10, "completed", conclusion)})
			if got.Verdict != NotGreen || len(got.NotGreen) != 1 {
				t.Fatalf("CI = %#v", got)
			}
		})
	}
}
