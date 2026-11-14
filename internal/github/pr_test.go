package github

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

const prURL = "https://github.com/mbell-nzl/fid/pull/74"
const prFields = "number,url,state,mergedAt,mergeCommit,headRefOid,headRefName,baseRefName,baseRefOid,isCrossRepository,body"

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func prFixture(t *testing.T, mutate func(map[string]any)) []byte {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(fixture(t, "pr74.json"), &value); err != nil {
		t.Fatal(err)
	}
	value["body"] = "Original PR body\n\n## Rollback\nRevert the commit.\n"
	if mutate != nil {
		mutate(value)
	}
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func clientForTest(t *testing.T, run Runner) *Client {
	t.Helper()
	ciCalls := 0
	audited := func(ctx context.Context, args, env []string) Result {
		switch {
		case len(args) > 0 && args[0] == "pr":
			assertRead(t, ctx, args, env, []string{"pr", "view", "74", "--repo", "mbell-nzl/fid", "--json", prFields})
		case len(args) > 5 && strings.Contains(args[5], "/git/commits/"):
			assertRead(t, ctx, args, env, apiArgs("repos/mbell-nzl/fid/git/commits/"+headSHA, treeProjection))
		default:
			ciCalls++
			assertRead(t, ctx, args, env, apiArgs(ciEndpoint(ciCalls), ciProjection))
		}
		return run(ctx, args, env)
	}
	client, err := New(Repository{Owner: "mbell-nzl", Name: "fid"}, "main", audited)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func assertRead(t *testing.T, ctx context.Context, args, env, want []string) {
	t.Helper()
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("argv = %#v, want %#v", args, want)
	}
	deadline, ok := ctx.Deadline()
	if !ok || time.Until(deadline) > DefaultTimeout {
		t.Fatalf("unbounded request context")
	}
	for _, arg := range args {
		if slices.Contains([]string{"fetch", "merge", "push", "POST", "PUT", "PATCH", "DELETE"}, arg) {
			t.Fatalf("write argv %q", arg)
		}
	}
	for _, value := range []string{"GH_PROMPT_DISABLED=1", "GH_PAGER=cat", "NO_COLOR=1", "GH_NO_UPDATE_NOTIFIER=1", "GH_REPO=mbell-nzl/fid"} {
		if !slices.Contains(env, value) {
			t.Errorf("missing %s", value)
		}
	}
}

func TestPRFacts(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, state   string
		merged, cross bool
	}{
		{"merged", "MERGED", true, false}, {"open", "OPEN", false, false}, {"closed unmerged", "CLOSED", false, false}, {"fork head", "MERGED", true, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			data := prFixture(t, func(value map[string]any) {
				value["state"] = tt.state
				value["isCrossRepository"] = tt.cross
				if !tt.merged {
					value["mergeCommit"] = nil
					value["mergedAt"] = nil
				}
				if tt.cross {
					value["headRepositoryOwner"] = map[string]any{"login": "fork-owner"}
				}
			})
			calls := 0
			client := clientForTest(t, func(ctx context.Context, args, env []string) Result {
				calls++
				assertRead(t, ctx, args, env, []string{"pr", "view", "74", "--repo", "mbell-nzl/fid", "--json", prFields})
				return Result{Stdout: data}
			})
			got, err := client.PR(context.Background(), prURL)
			if err != nil {
				t.Fatal(err)
			}
			if calls != 1 || got.Number != 74 || got.URL != prURL || got.State != tt.state || got.Merged != tt.merged || got.CrossRepository != tt.cross || got.HeadSHA != headSHA || got.HeadRef != "reset-1.5d2-drop-tables" || got.BaseRef != "main" || got.BaseSHA != "526091c64c69ac62b5bf2f5c30a2770ca7f88e72" || got.BaseRepository.String() != "mbell-nzl/fid" || got.Body != "Original PR body\n\n## Rollback\nRevert the commit.\n" {
				t.Fatalf("PR = %#v (%d calls)", got, calls)
			}
			if tt.merged && (got.MergeCommit != mergeSHA || got.MergedAt == nil || got.MergedAt.Format(time.RFC3339) != "2026-11-14T18:21:19Z") {
				t.Fatalf("merge facts = %#v", got)
			}
		})
	}
}

func TestPRMismatch(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, url string
		mutate    func(map[string]any)
		kind      MismatchKind
		calls     int
	}{
		{"linked elsewhere", "https://github.com/elsewhere/fid/pull/74", nil, RepositoryMismatch, 0},
		{"response base repository", prURL, func(v map[string]any) { v["url"] = "https://github.com/elsewhere/fid/pull/74" }, RepositoryMismatch, 1},
		{"base branch", prURL, func(v map[string]any) { v["baseRefName"] = "release" }, BranchMismatch, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			calls := 0
			client := clientForTest(t, func(ctx context.Context, args, env []string) Result {
				calls++
				assertRead(t, ctx, args, env, []string{"pr", "view", "74", "--repo", "mbell-nzl/fid", "--json", prFields})
				return Result{Stdout: prFixture(t, tt.mutate)}
			})
			_, err := client.PR(context.Background(), tt.url)
			var mismatch *MismatchError
			var unavailable *UnavailableError
			if !errors.As(err, &mismatch) || mismatch.Kind != tt.kind || errors.As(err, &unavailable) || calls != tt.calls {
				t.Fatalf("mismatch = %v (%d calls)", err, calls)
			}
		})
	}
}

func TestPRUnavailable(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		result Result
		reason UnavailableReason
	}{
		{"missing gh", Result{Err: exec.ErrNotFound}, MissingGH},
		{"authentication", Result{ExitCode: 4, Stderr: []byte("authentication required")}, Unauthenticated},
		{"network", Result{ExitCode: 1, Stderr: []byte("error connecting to api.github.com")}, APIError},
		{"not found", Result{ExitCode: 1, Stderr: []byte("gh: Not Found (HTTP 404)")}, PRNotFound},
		{"rate limit", Result{ExitCode: 1, Stderr: []byte("gh: API rate limit exceeded (HTTP 403)")}, RateLimited},
		{"429", Result{ExitCode: 1, Stderr: []byte("HTTP 429")}, RateLimited},
		{"timeout", Result{Err: context.DeadlineExceeded}, TimedOut},
		{"truncated JSON", Result{Stdout: []byte(`{"number":74`)}, MalformedResponse},
		{"trailing JSON", Result{Stdout: append(prFixture(t, nil), []byte(`{}`)...)}, MalformedResponse},
		{"null", Result{Stdout: []byte("null")}, MalformedResponse},
		{"oversized", Result{Stdout: []byte(strings.Repeat(" ", DefaultOutputLimit+1))}, OversizedOutput},
		{"runner output bound", Result{Err: ErrOutputLimit}, OversizedOutput},
		{"failed with good JSON", Result{ExitCode: 1, Stdout: prFixture(t, nil), Stderr: []byte("API failed")}, APIError},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := clientForTest(t, func(ctx context.Context, args, env []string) Result {
				assertRead(t, ctx, args, env, []string{"pr", "view", "74", "--repo", "mbell-nzl/fid", "--json", prFields})
				return tt.result
			})
			_, err := client.PR(context.Background(), prURL)
			assertUnavailable(t, err, tt.reason)
		})
	}
}

func TestPRInconsistentResponse(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"missing number", func(v map[string]any) { delete(v, "number") }},
		{"wrong number", func(v map[string]any) { v["number"] = 75 }},
		{"unknown state", func(v map[string]any) { v["state"] = "UNKNOWN" }},
		{"merged without commit", func(v map[string]any) { v["mergeCommit"] = nil }},
		{"closed with merged time", func(v map[string]any) { v["state"] = "CLOSED" }},
		{"missing body", func(v map[string]any) { delete(v, "body") }},
		{"missing cross repository", func(v map[string]any) { delete(v, "isCrossRepository") }},
		{"invalid head", func(v map[string]any) { v["headRefOid"] = "bad" }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := clientForTest(t, func(context.Context, []string, []string) Result { return Result{Stdout: prFixture(t, tt.mutate)} })
			_, err := client.PR(context.Background(), prURL)
			assertUnavailable(t, err, MalformedResponse)
		})
	}
}

func assertUnavailable(t *testing.T, err error, reason UnavailableReason) {
	t.Helper()
	var unavailable *UnavailableError
	if !errors.As(err, &unavailable) || unavailable.Reason != reason || !strings.HasPrefix(err.Error(), "GitHub unavailable:") {
		t.Fatalf("failure = %v, want reason %s", err, reason)
	}
}
