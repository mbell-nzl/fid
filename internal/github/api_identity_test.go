package github

import (
	"context"
	"testing"
)

func TestMergeCommitCIRejectsConflictingExecutionIdentity(t *testing.T) {
	t.Parallel()
	for _, field := range []struct {
		name   string
		mutate func(*WorkflowRun)
	}{
		{"workflow ID", func(run *WorkflowRun) { run.WorkflowID = 1000 }},
		{"head SHA", func(run *WorkflowRun) { run.HeadSHA = headSHA }},
		{"event", func(run *WorkflowRun) { run.Event = "pull_request" }},
	} {
		t.Run(field.name, func(t *testing.T) {
			t.Parallel()
			for _, layout := range []struct {
				name      string
				firstSize int
			}{
				{"same page", 1},
				{"across pages", 100},
			} {
				t.Run(layout.name, func(t *testing.T) {
					t.Parallel()
					first := make([]WorkflowRun, layout.firstSize)
					for i := range first {
						first[i] = sampleRun(int64(i+1), 1, int64(i+10), "completed", "success")
					}
					second := first[0]
					second.Attempt = 2
					field.mutate(&second)
					total := len(first) + 1
					pages := [][]byte{ciPage(t, total, append(first, second))}
					if layout.firstSize == 100 {
						pages = [][]byte{ciPage(t, total, first), ciPage(t, total, []WorkflowRun{second})}
					}
					calls := 0
					client := clientForTest(t, func(ctx context.Context, args, env []string) Result {
						if calls >= len(pages) {
							t.Fatal("unexpected extra page")
						}
						assertRead(t, ctx, args, env, apiArgs(ciEndpoint(calls+1), ciProjection))
						data := pages[calls]
						calls++
						return Result{Stdout: data}
					})
					got, err := client.MergeCommitCI(context.Background(), mergeSHA)
					if got.Verdict != "" || got.Commit != "" || len(got.Runs) != 0 || len(got.NotGreen) != 0 {
						t.Errorf("inconsistent execution returned CI facts: verdict=%q commit=%q observations=%d not-green=%d", got.Verdict, got.Commit, len(got.Runs), len(got.NotGreen))
					}
					assertUnavailable(t, err, MalformedResponse)
					if calls != len(pages) {
						t.Fatalf("calls = %d, want %d", calls, len(pages))
					}
				})
			}
		})
	}
}
