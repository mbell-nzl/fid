package reviewcomment

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mbell-nzl/fid/internal/reviewrecord"
	_ "modernc.org/sqlite"
)

func TestLegacyCommentBodies(t *testing.T) {
	t.Parallel()
	for name, id := range commentSubjects {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assertCommentBody(t, commentFixture(t), id, "v1_"+name)
		})
	}
}

var commentSubjects = map[string]string{"round": "RND-current", "none": "RND-empty", "answer": "FRS-answer", "deferral": "FRS-defer"}

func assertCommentBody(t *testing.T, conn *sql.DB, id, golden string) {
	t.Helper()
	s, err := ReadSubject(conn, id)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != id || s.TaskID != "FID-task" || s.HandoffID != "HOF-one" || s.PR != "https://github.com/owner/repo/pull/63" {
		t.Fatalf("subject metadata: %+v", s)
	}
	want, err := os.ReadFile(filepath.Join("testdata", golden+".md"))
	if err != nil {
		t.Fatal(err)
	}
	if s.Body != string(want) {
		t.Fatalf("body differs from %s:\ngot:\n%s\nwant:\n%s", golden, s.Body, want)
	}
}

func commentFixture(t *testing.T) *sql.DB {
	t.Helper()
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	conn.SetMaxOpenConns(1)
	_, err = conn.Exec(`CREATE TABLE review_rounds(id TEXT,task_id TEXT,handoff_id TEXT,record_json TEXT);
CREATE TABLE finding_responses(id TEXT,record_json TEXT);
CREATE TABLE handoffs(id TEXT,task_id TEXT,record_json TEXT);
INSERT INTO handoffs VALUES('HOF-one','FID-task','{"pr":"https://github.com/owner/repo/pull/63","lane":"standard"}');`)
	if err != nil {
		t.Fatal(err)
	}
	r := reviewrecord.Round{ID: "RND-current", TaskID: "FID-task", HandoffID: "HOF-one",
		Reviewer: reviewrecord.Reviewer{Actor: "agent:reviewer", Model: "openai/gpt-6.1-sol", Effort: "xhigh", SessionID: "codex-01a2"},
		Stream:   "primary", Profiles: []string{"code", "test-adequacy"}, Verdict: "CHANGES",
		Purpose: "Keep **reviewer Markdown** and `code`.", Verification: "Ran `go test`.\n\n```text\nPASS\n```"}
	r.Findings = []reviewrecord.Finding{
		{ID: "RND-current-1", FindingInput: reviewrecord.FindingInput{Severity: "major", Summary: "Check **replay**", Location: &reviewrecord.Location{File: "internal/status/approval.go", Line: 42}, Evidence: "A `receipt` failed.\n    More evidence.", Fix: "Keep *old* rendering."}},
		{ID: "RND-current-2", FindingInput: reviewrecord.FindingInput{Severity: "minor", Summary: "Explain the choice", Evidence: "No explanation.", Fix: "Add a reason."}},
		{ID: "RND-current-3", FindingInput: reviewrecord.FindingInput{Severity: "nit", Summary: "Polish", Evidence: "A label.", Fix: "Make it **bold**."}},
	}
	r.Dispositions = []reviewrecord.Disposition{
		{FindingID: "RND-earlier-1", Disposition: "resolved", Reason: "Now **covered**."},
		{FindingID: "RND-earlier-2", Disposition: "not-applicable", Reason: "Different path."},
		{FindingID: "RND-earlier-3", Disposition: "open", Reason: "Still applies."},
	}
	for _, empty := range []bool{false, true} {
		if empty {
			r.ID, r.Verdict, r.Findings, r.Dispositions = "RND-empty", "APPROVE", nil, nil
			r.Reviewer.Model, r.Reviewer.Effort = "", ""
		}
		raw, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Exec(`INSERT INTO review_rounds VALUES(?,?,?,?)`, r.ID, r.TaskID, r.HandoffID, string(raw)); err != nil {
			t.Fatal(err)
		}
	}
	a := reviewrecord.Response{ID: "FRS-answer", FindingID: "RND-current-2", TaskID: r.TaskID, HandoffID: r.HandoffID, Actor: "agent:builder", Text: "Fixed in `abc123`.\n\n**Details** kept as written."}
	for _, deferred := range []bool{false, true} {
		if deferred {
			a.ID, a.DeferredTask = "FRS-defer", &reviewrecord.DeferredTask{ID: "FID-later", Title: "Keep **this title**"}
		}
		raw, err := json.Marshal(a)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Exec(`INSERT INTO finding_responses VALUES(?,?)`, a.ID, string(raw)); err != nil {
			t.Fatal(err)
		}
	}
	return conn
}
