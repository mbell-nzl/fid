package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/doneeval"
)

// doneChild uses the built fid binary, a strict gh stub and a git wrapper which
// rejects fetch. Only the child receives PATH and the stub's fixture variables.
func doneChild(t *testing.T, root string, facts doneeval.GitHubFacts, args ...string) commandExecutionResult {
	t.Helper()
	dir := t.TempDir()
	script, err := os.ReadFile(filepath.Join(commandRepositoryRoot(), "cmd", "fid", "testdata", "gh_stub.py"))
	if err != nil {
		t.Fatal(err)
	}
	writeDoneExecutable(t, filepath.Join(dir, "gh"), script)
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	writeDoneExecutable(t, filepath.Join(dir, "git"), []byte("#!/bin/sh\nfor arg do\n if [ \"$arg\" = fetch ]; then echo 'fetch is forbidden' >&2; exit 97; fi\ndone\nexec '"+git+"' \"$@\"\n"))
	pr := facts.PR
	var mergedAt, merge any
	if pr.Merged {
		mergedAt = pr.MergedAt
		merge = map[string]any{"oid": pr.MergeCommit}
	}
	requests := map[string]any{
		"pr view 63 --repo mbell-nzl/fid --json number,url,state,mergedAt,mergeCommit,headRefOid,headRefName,baseRefName,baseRefOid,isCrossRepository,body":                                                                                                                    map[string]any{"body": map[string]any{"number": 63, "url": pr.URL, "state": pr.State, "mergedAt": mergedAt, "mergeCommit": merge, "headRefOid": pr.HeadSHA, "headRefName": "task", "baseRefName": pr.BaseRef, "baseRefOid": pr.BaseSHA, "isCrossRepository": false, "body": pr.Body}},
		"api --hostname github.com --method GET repos/mbell-nzl/fid/git/commits/" + pr.HeadSHA + " --jq {sha, tree: .tree.sha}":                                                                                                                                                map[string]any{"body": map[string]any{"sha": pr.HeadSHA, "tree": facts.HeadTree}},
		"api --hostname github.com --method GET repos/mbell-nzl/fid/actions/runs?head_sha=" + pr.MergeCommit + "&per_page=100&page=1 --jq {total_count, runs: [.workflow_runs[] | {id, run_attempt, workflow_id, name, path, event, head_sha, status, conclusion, html_url}]}": map[string]any{"body": map[string]any{"total_count": len(facts.CI.Runs), "runs": facts.CI.Runs}},
	}
	raw, err := json.Marshal(map[string]any{"requests": requests})
	if err != nil {
		t.Fatal(err)
	}
	fixture := filepath.Join(dir, "responses.json")
	if err := os.WriteFile(fixture, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.CommandContext(t.Context(), binary, args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), runAsFIDVariable+"=1", "PATH="+dir+string(os.PathListSeparator)+os.Getenv("PATH"), "FID_TEST_GH_RESPONSES="+fixture, "FID_TEST_GH_LOG="+filepath.Join(dir, "requests.log"))
	if testing.CoverMode() != "" {
		cmd.Env = append(cmd.Env, "GOCOVERDIR="+t.TempDir())
	}
	var out, errout bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errout
	err = cmd.Run()
	log, logerr := os.ReadFile(filepath.Join(dir, "requests.log"))
	if logerr != nil || len(strings.Split(strings.TrimSpace(string(log)), "\n")) != 3 {
		t.Fatalf("GitHub calls=%s/%v; stderr=%s", log, logerr, errout.String())
	}
	return commandExecutionResult{stdout: out.String(), stderr: errout.String(), err: err}
}
func writeDoneExecutable(t *testing.T, path string, script []byte) {
	t.Helper()
	if err := os.WriteFile(path, script, 0o700); err != nil {
		t.Fatal(err)
	}
}
