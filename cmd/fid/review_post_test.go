package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

//nolint:paralleltest // The fake gh owns this test process's PATH.
func TestReviewPostingPendingAndRetry(t *testing.T) {
	f, task := canonicalReviewFixture(t)
	dir := t.TempDir()
	log, offline := filepath.Join(dir, "comments"), filepath.Join(dir, "offline")
	requests := filepath.Join(dir, "requests")
	script := `#!/usr/bin/python3
import json, pathlib, sys
d=pathlib.Path(__file__).parent
with (d/'requests').open('a') as f: f.write('request\n')
if (d/'offline').exists(): sys.exit(1)
args=sys.argv[1:]
p=d/'stored.json'
comments=json.loads(p.read_text()) if p.exists() else []
if args[0]=='api':
    print(json.dumps(comments))
else:
    body=pathlib.Path(args[args.index('--body-file')+1]).read_text()
    with (d/'comments').open('a') as f: f.write(body)
    url='https://github.com/mbell-nzl/fid/issues/63#issuecomment-'+str(123+len(comments))
    comments.append({'body':body,'html_url':url})
    p.write_text(json.dumps(comments))
    print(url)
`
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	if err := os.WriteFile(offline, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	in := reviewInput()
	in["findings"] = []any{reviewFinding("minor")}
	first, _ := lifecycleResult(t, recordReviewTask(t, f, in, "post-round", "human:courier", task), "review.record")
	round := first["round"].(map[string]any)
	id := round["id"].(string)
	finding := round["findings"].([]any)[0].(map[string]any)["id"].(string)
	status := f.run("--json", "status")
	if status.err != nil || !strings.Contains(status.stdout, "review-post-pending") || !strings.Contains(status.stdout, id) {
		t.Fatalf("missing pending warning: %#v", status)
	}
	normalize := func(s string) string {
		s = f.placeholders(s)
		for _, prefix := range []string{"RND", "FRS", "HOF", "FID", "HLD"} {
			s = regexp.MustCompile(prefix+`-[0-9a-f]{6}`).ReplaceAllString(s, prefix+"-<new>")
		}
		return regexp.MustCompile(`\b(?:[0-9a-f]{64}|[0-9a-f]{40})\b`).ReplaceAllString(s, "<fingerprint-or-commit>")
	}
	fields, _ := lifecycleResult(t, status, "project.status")
	warningsJSON, err := json.Marshal(map[string]any{"warnings": fields["warnings"]})
	if err != nil {
		t.Fatal(err)
	}
	assertGolden(t, filepath.Join("review", "pending.json"), normalize(string(warningsJSON)))
	pendingText := f.run("status")
	if pendingText.err != nil {
		t.Fatal(pendingText)
	}
	warningText := ""
	lines := strings.Split(pendingText.stdout, "\n")
	for i, line := range lines {
		if strings.Contains(line, "PR comments pending") {
			warningText = line + "\n" + lines[i+1] + "\n"
		}
	}
	assertGolden(t, filepath.Join("review", "pending.txt"), normalize(warningText))
	for _, item := range fields["needs_you"].([]any) {
		if item.(map[string]any)["kind"] == "review-post-pending" {
			t.Fatal("pending post needs a person")
		}
	}
	attempts, err := os.ReadFile(requests)
	if err != nil || len(attempts) == 0 {
		t.Fatal("recording did not attempt posting")
	}
	pendingResult, _ := lifecycleResult(t, f.run("--json", "review", "post", id), "review.post")
	if len(pendingResult["pending"].([]any)) != 1 || pendingResult["pending"].([]any)[0] != id {
		t.Fatalf("offline retry=%v", pendingResult)
	}
	if err := os.Remove(offline); err != nil {
		t.Fatal(err)
	}
	lifecycleResult(t, f.run("--json", "review", "post", id), "review.post")
	body, err := os.ReadFile(log)
	if err != nil || !strings.Contains(string(body), "## Review: APPROVE") || !strings.Contains(string(body), "Retry conflicts lack coverage") || !strings.Contains(string(body), "test-adequacy") {
		t.Fatalf("comment not rendered from record: %s / %v", body, err)
	}
	assertGolden(t, filepath.Join("review", "comment.md"), normalize(string(body)))
	before := startEventCount(t, f.workspace)
	requestCount, _ := os.ReadFile(requests)
	lifecycleResult(t, f.run("--json", "review", "post", id), "review.post")
	got, _ := os.ReadFile(log)
	if string(got) != string(body) || startEventCount(t, f.workspace) != before {
		t.Fatal("retry posted twice")
	}
	status = f.run("--json", "status")
	if status.err != nil || strings.Contains(status.stdout, "review-post-pending") {
		t.Fatalf("warning not cleared: %#v", status)
	}
	for _, args := range [][]string{{"review", "list", task}, {"show", task}, {"next"}} {
		if got := f.run(args...); got.err != nil {
			t.Fatal(got)
		}
	}
	afterReads, _ := os.ReadFile(requests)
	if string(afterReads) != string(requestCount) {
		t.Fatal("retry/read contacted gh")
	}
	answer, _ := lifecycleResult(t, f.run("--json", "finding", "answer", finding, "--actor", "agent:builder", "--text", "Accepted explanation"), "finding.answer")
	responseID := answer["response"].(map[string]any)["id"].(string)
	if got := f.run("--json", "review", "post", responseID); got.err == nil {
		t.Fatal("review post accepted a response as its round")
	}
	got, _ = os.ReadFile(log)
	if !strings.Contains(string(got), "Accepted explanation") {
		t.Fatal("answer was not posted")
	}
	lifecycleResult(t, f.run("--json", "finding", "defer", finding, "--actor", "agent:builder", "--title", "Follow up"), "finding.defer")
	got, _ = os.ReadFile(log)
	if !strings.Contains(string(got), "Deferred") || !strings.Contains(string(got), "Follow up") {
		t.Fatal("deferral was not posted")
	}
}
