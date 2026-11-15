package github

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"
)

const ciProjection = "{total_count, runs: [.workflow_runs[] | {id, run_attempt, workflow_id, name, path, event, head_sha, status, conclusion, html_url}]}"
const treeProjection = "{sha, tree: .tree.sha}"

func apiArgs(endpoint, projection string) []string {
	return []string{"api", "--hostname", "github.com", "--method", "GET", endpoint, "--jq", projection}
}

func TestCommitTree(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name      string
		data      []byte
		want      string
		malformed bool
	}{
		{"sample", fixture(t, "head-commit74.json"), treeSHA, false},
		{"wrong commit", []byte(fmt.Sprintf(`{"sha":%q,"tree":%q}`, mergeSHA, treeSHA)), "", true},
		{"missing tree", []byte(fmt.Sprintf(`{"sha":%q}`, headSHA)), "", true},
		{"invalid tree", []byte(fmt.Sprintf(`{"sha":%q,"tree":"no"}`, headSHA)), "", true},
		{"truncated", []byte(`{"sha":`), "", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := clientForTest(t, func(ctx context.Context, args, env []string) Result {
				assertRead(t, ctx, args, env, apiArgs("repos/mbell-nzl/fid/git/commits/"+headSHA, treeProjection))
				return Result{Stdout: tt.data}
			})
			got, err := client.CommitTree(context.Background(), headSHA)
			if tt.malformed {
				assertUnavailable(t, err, MalformedResponse)
			} else if err != nil || got != tt.want {
				t.Fatalf("tree = %q, %v", got, err)
			}
		})
	}
}

func ciPage(t *testing.T, total int, runs []WorkflowRun) []byte {
	t.Helper()
	if runs == nil {
		runs = []WorkflowRun{}
	}
	data, err := json.Marshal(struct {
		Total int           `json:"total_count"`
		Runs  []WorkflowRun `json:"runs"`
	}{total, runs})
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func ciEndpoint(page int) string {
	return fmt.Sprintf("repos/mbell-nzl/fid/actions/runs?head_sha=%s&per_page=100&page=%d", mergeSHA, page)
}

func TestMergeCommitCI(t *testing.T) {
	t.Parallel()
	success := sampleRun(43626136158, 1, 290273570, "completed", "success")
	failure := success
	failure.Conclusion = "failure"
	secondSuccess := success
	secondSuccess.Attempt = 2
	secondFailure := failure
	secondFailure.Attempt = 2
	queued := success
	queued.Status = "queued"
	queued.Conclusion = ""
	progress := success
	progress.Status = "in_progress"
	progress.Conclusion = ""
	skipped := success
	skipped.Conclusion = "skipped"
	other := success
	other.ID++
	other.WorkflowID++
	other.URL = "https://github.com/mbell-nzl/fid/actions/runs/43626136159"
	wrongHead := success
	wrongHead.HeadSHA = headSHA
	wrongEvent := success
	wrongEvent.Event = "pull_request"
	for _, tt := range []struct {
		name    string
		pages   [][]byte
		verdict Verdict
		ignored IgnoreReason
		runs    []WorkflowRun
		failed  []WorkflowRun
	}{
		{"captured sample", [][]byte{fixture(t, "rest-runs-merge74.json")}, Green, "", []WorkflowRun{success}, nil},
		{"none", [][]byte{ciPage(t, 0, nil)}, Waiting, "", nil, nil},
		{"queued", [][]byte{ciPage(t, 1, []WorkflowRun{queued})}, Waiting, "", []WorkflowRun{queued}, nil},
		{"in progress", [][]byte{ciPage(t, 1, []WorkflowRun{progress})}, Waiting, "", []WorkflowRun{progress}, nil},
		{"skipped only", [][]byte{ciPage(t, 1, []WorkflowRun{skipped})}, Waiting, Skipped, []WorkflowRun{ignoredRun(skipped, Skipped)}, nil},
		{"mixed workflows", [][]byte{ciPage(t, 2, []WorkflowRun{failure, other})}, NotGreen, "", []WorkflowRun{failure, other}, []WorkflowRun{failure}},
		{"multiple passing workflows", [][]byte{ciPage(t, 2, []WorkflowRun{success, other})}, Green, "", []WorkflowRun{success, other}, nil},
		{"pending without conclusion", [][]byte{ciMutate(t, queued, func(v map[string]any) { v["conclusion"] = nil })}, Waiting, "", []WorkflowRun{queued}, nil},
		{"rerun failed", [][]byte{ciPage(t, 2, []WorkflowRun{success, secondFailure})}, NotGreen, Superseded, []WorkflowRun{ignoredRun(success, Superseded), secondFailure}, []WorkflowRun{secondFailure}},
		{"rerun passed", [][]byte{ciPage(t, 2, []WorkflowRun{failure, secondSuccess})}, Green, Superseded, []WorkflowRun{ignoredRun(failure, Superseded), secondSuccess}, nil},
		{"wrong head", [][]byte{ciPage(t, 1, []WorkflowRun{wrongHead})}, Waiting, WrongHead, []WorkflowRun{ignoredRun(wrongHead, WrongHead)}, nil},
		{"PR event", [][]byte{ciPage(t, 1, []WorkflowRun{wrongEvent})}, Waiting, WrongEvent, []WorkflowRun{ignoredRun(wrongEvent, WrongEvent)}, nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calls := 0
			client := clientForTest(t, func(ctx context.Context, args, env []string) Result {
				if calls >= len(tt.pages) {
					t.Fatal("unexpected extra page")
				}
				assertRead(t, ctx, args, env, apiArgs(ciEndpoint(calls+1), ciProjection))
				data := tt.pages[calls]
				calls++
				return Result{Stdout: data}
			})
			got, err := client.MergeCommitCI(context.Background(), mergeSHA)
			if err != nil || got.Verdict != tt.verdict || got.Commit != mergeSHA || calls != len(tt.pages) {
				t.Fatalf("CI = %#v, %v, %d calls", got, err, calls)
			}
			if tt.ignored != "" && got.Runs[0].Ignored != tt.ignored {
				t.Fatalf("ignored = %#v", got.Runs)
			}
			if !slices.Equal(got.Runs, tt.runs) {
				t.Errorf("workflows = %#v, want %#v", got.Runs, tt.runs)
			}
			if !slices.Equal(got.NotGreen, tt.failed) {
				t.Errorf("not green = %#v, want %#v", got.NotGreen, tt.failed)
			}
		})
	}
}

func ignoredRun(run WorkflowRun, reason IgnoreReason) WorkflowRun {
	run.Ignored = reason
	return run
}

func TestCIPagination(t *testing.T) {
	t.Parallel()
	first := make([]WorkflowRun, 100)
	for i := range first {
		first[i] = sampleRun(int64(i+1), 1, int64(i+1), "completed", "success")
	}
	failed := sampleRun(101, 1, 101, "completed", "failure")
	for _, tt := range []struct {
		name   string
		second []byte
		result Result
		want   Verdict
		reason UnavailableReason
	}{
		{name: "last page failure", second: ciPage(t, 101, []WorkflowRun{failed}), want: NotGreen},
		{name: "changed total", second: ciPage(t, 102, []WorkflowRun{failed}), reason: MalformedResponse},
		{name: "missing page", second: ciPage(t, 101, nil), reason: MalformedResponse},
		{name: "duplicate execution", second: ciPage(t, 101, []WorkflowRun{first[0]}), reason: MalformedResponse},
		{name: "truncated last page", second: []byte(`{"total_count":101,"runs":[`), reason: MalformedResponse},
		{name: "failed last page", result: Result{ExitCode: 1, Stderr: []byte("HTTP 503")}, reason: APIError},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calls := 0
			client := clientForTest(t, func(ctx context.Context, args, env []string) Result {
				calls++
				assertRead(t, ctx, args, env, apiArgs(ciEndpoint(calls), ciProjection))
				switch calls {
				case 1:
					return Result{Stdout: ciPage(t, 101, first)}
				case 2:
					if tt.second != nil {
						return Result{Stdout: tt.second}
					}
					return tt.result
				default:
					t.Fatal("unexpected extra page")
					return Result{}
				}
			})
			got, err := client.MergeCommitCI(context.Background(), mergeSHA)
			if calls != 2 {
				t.Fatalf("calls = %d", calls)
			}
			if tt.reason != "" {
				assertUnavailable(t, err, tt.reason)
				if got.Verdict == Green {
					t.Fatal("partial CI was green")
				}
				return
			}
			if err != nil || got.Verdict != tt.want || len(got.Runs) != 101 || len(got.NotGreen) != 1 || got.NotGreen[0].ID != 101 {
				t.Fatalf("CI = %#v, %v", got, err)
			}
			if !slices.Equal(got.Runs, append(slices.Clone(first), failed)) || !slices.Equal(got.NotGreen, []WorkflowRun{failed}) {
				t.Fatalf("workflow facts changed across pages: CI = %#v", got)
			}
		})
	}
}

func TestCIIncompleteResponse(t *testing.T) {
	t.Parallel()
	success := sampleRun(1, 1, 10, "completed", "success")
	for _, tt := range []struct {
		name string
		data []byte
	}{
		{"missing total", []byte(`{"runs":[]}`)}, {"missing list", []byte(`{"total_count":0}`)},
		{"null list", []byte(`{"total_count":0,"runs":null}`)}, {"short page", ciPage(t, 2, []WorkflowRun{success})},
		{"negative count", ciPage(t, -1, nil)}, {"extra workflow", ciPage(t, 0, []WorkflowRun{success})},
		{"invalid ID", ciMutate(t, success, func(v map[string]any) { v["id"] = 0 })},
		{"missing attempt", ciMutate(t, success, func(v map[string]any) { delete(v, "run_attempt") })},
		{"missing status", ciMutate(t, success, func(v map[string]any) { delete(v, "status") })},
		{"missing conclusion", ciMutate(t, success, func(v map[string]any) { delete(v, "conclusion") })},
		{"missing workflow", ciMutate(t, success, func(v map[string]any) { delete(v, "workflow_id") })},
		{"missing URL", ciMutate(t, success, func(v map[string]any) { delete(v, "html_url") })},
		{"missing path", ciMutate(t, success, func(v map[string]any) { delete(v, "path") })},
		{"completed without conclusion", ciMutate(t, success, func(v map[string]any) { v["conclusion"] = nil })},
		{"queued with conclusion", ciMutate(t, success, func(v map[string]any) { v["status"] = "queued" })},
		{"unknown status", ciMutate(t, success, func(v map[string]any) { v["status"] = "unknown" })},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := clientForTest(t, func(context.Context, []string, []string) Result { return Result{Stdout: tt.data} })
			got, err := client.MergeCommitCI(context.Background(), mergeSHA)
			assertUnavailable(t, err, MalformedResponse)
			if got.Verdict == Green {
				t.Fatal("incomplete CI was green")
			}
		})
	}
}

func ciMutate(t *testing.T, run WorkflowRun, mutate func(map[string]any)) []byte {
	t.Helper()
	var value map[string]any
	data, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	mutate(value)
	data, err = json.Marshal(map[string]any{"total_count": 1, "runs": []any{value}})
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestInvalidReadInput(t *testing.T) {
	t.Parallel()
	client := clientForTest(t, func(context.Context, []string, []string) Result {
		t.Fatal("invalid input executed gh")
		return Result{}
	})
	for _, sha := range []string{"", "--help", "a/b", strings.Repeat("a", 39)} {
		t.Run(sha, func(t *testing.T) {
			t.Parallel()
			if _, err := client.CommitTree(context.Background(), sha); err == nil {
				t.Fatal("accepted invalid commit")
			}
			if _, err := client.MergeCommitCI(context.Background(), sha); err == nil {
				t.Fatal("accepted invalid commit")
			}
		})
	}
}

func TestCIRejectsInvalidWorkflowURL(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"nonsense", "https://github.com/other/repo/actions/runs/1", "https://github.com/mbell-nzl/fid/actions/runs/99"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			run := sampleRun(1, 1, 10, "completed", "success")
			run.URL = value
			client := clientForTest(t, func(context.Context, []string, []string) Result {
				return Result{Stdout: ciPage(t, 1, []WorkflowRun{run})}
			})
			_, err := client.MergeCommitCI(context.Background(), mergeSHA)
			assertUnavailable(t, err, MalformedResponse)
		})
	}
}

func TestDuplicateResponseKeys(t *testing.T) {
	t.Parallel()
	success := sampleRun(1, 1, 10, "completed", "success")
	data := strings.Replace(string(ciPage(t, 1, []WorkflowRun{success})), `"conclusion":"success"`, `"conclusion":"failure","conclusion":"success"`, 1)
	client := clientForTest(t, func(context.Context, []string, []string) Result { return Result{Stdout: []byte(data)} })
	_, err := client.MergeCommitCI(context.Background(), mergeSHA)
	assertUnavailable(t, err, MalformedResponse)
}

func TestCIAggregateOutputBound(t *testing.T) {
	t.Parallel()
	first := make([]WorkflowRun, 100)
	for i := range first {
		first[i] = sampleRun(int64(i+1), 1, int64(i+1), "completed", "success")
		first[i].Name = strings.Repeat("x", 70000)
	}
	last := sampleRun(101, 1, 101, "completed", "success")
	last.Name = strings.Repeat("x", 2<<20)
	pages := [][]byte{ciPage(t, 101, first), ciPage(t, 101, []WorkflowRun{last})}
	calls := 0
	client := clientForTest(t, func(context.Context, []string, []string) Result {
		if calls >= len(pages) {
			t.Fatal("unexpected extra page")
		}
		page := pages[calls]
		calls++
		return Result{Stdout: page}
	})
	got, err := client.MergeCommitCI(context.Background(), mergeSHA)
	assertUnavailable(t, err, OversizedOutput)
	if calls != 2 || got.Verdict == Green {
		t.Fatalf("CI=%#v, calls=%d", got, calls)
	}
}
