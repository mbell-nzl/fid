package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/epicplan"
	"github.com/mbell-nzl/fid/internal/epicwork"
	"github.com/mbell-nzl/fid/internal/reviewwork"
	"github.com/mbell-nzl/fid/internal/workrecord"
	fidmigrations "github.com/mbell-nzl/fid/migrations"
)

func TestEpicApprovalAcceptsOnlyItsPlanAndReplays(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	other, _ := lifecycleResult(t, f.run("--json", "issue", "create", "--kind", "epic", "--title", "Other", "--scheduling", "prioritised", "--actor", "human:mark", "--idempotency-key", "other"), "issue.create")
	otherID := other["issue"].(map[string]any)["id"].(string)
	create := func(kind, parent, key string) string {
		t.Helper()
		args := []string{"--json", "issue", "create", "--kind", kind, "--title", key, "--actor", "agent:builder", "--idempotency-key", key}
		if kind == "task" {
			args = append(args, "--epic", parent, "--lane", "standard")
		} else {
			args = append(args, "--parent", parent)
		}
		r, _ := lifecycleResult(t, f.run(args...), "issue.create")
		return r["issue"].(map[string]any)["id"].(string)
	}
	task := create("task", epic, "task")
	subtask := create("subtask", task, "subtask")
	outside := create("task", otherID, "outside")
	o := epicwork.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "human:mark", Reason: "Plan approved", IdempotencyKey: "approve", Confirm: func(p epicplan.Plan) (bool, error) {
		if len(p.Tasks) != 2 {
			t.Fatalf("preview=%+v", p)
		}
		for _, task := range p.Tasks {
			if !task.Accept {
				t.Fatal("preview omitted acceptance")
			}
		}
		return true, nil
	}}
	r, err := epicwork.Approve(o)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Accepted) != 2 {
		t.Fatalf("accepted=%v", r.Accepted)
	}
	before := startEventCount(t, f.workspace)
	lifecycleResult(t, f.run("--json", "accept", task, "--actor", "human:mark"), "task.accept")
	if startEventCount(t, f.workspace) != before {
		t.Fatal("accepting approved work wrote a second decision")
	}
	for _, rebuild := range []bool{false, true} {
		if rebuild {
			if got := f.run("rebuild"); got.err != nil {
				t.Fatalf("rebuild=%#v", got)
			}
		}
		q, err := reviewwork.OpenProjection(f.workspace)
		if err != nil {
			t.Fatal(err)
		}
		unaccepted, err := workrecord.Unaccepted(q)
		_ = q.Close()
		if err != nil {
			t.Fatal(err)
		}
		if unaccepted[task] != "" || unaccepted[subtask] != "" || unaccepted[outside] == "" {
			t.Fatalf("acceptances=%v", unaccepted)
		}
		o.Confirm = func(epicplan.Plan) (bool, error) { t.Fatal("retry prompted"); return false, nil }
		retry, err := epicwork.Approve(o)
		if err != nil || retry.DecisionID != r.DecisionID || retry.EventID != r.EventID {
			t.Fatalf("retry=%+v,%v", retry, err)
		}
	}
}

func TestEpicApprovalReleasesLockAndRejectsChangedPreview(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	f.clock = time.Now
	_, err := epicwork.Approve(epicwork.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "human:mark", Reason: "Approve", Confirm: func(epicplan.Plan) (bool, error) {
		_, err := epicwork.Change(epicwork.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "agent:other", IdempotencyKey: "concurrent", Validation: epicplan.Validation{Claim: "Concurrent validation", Verifier: "review:code"}}, "add")
		if err != nil {
			t.Fatalf("confirmation held the write lock: %v", err)
		}
		return true, nil
	}})
	if err == nil || !strings.Contains(err.Error(), "changed after preview") {
		t.Fatalf("changed preview=%v", err)
	}
	q, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = q.Close() }()
	a, err := epicplan.Current(q, epic)
	if err != nil || a.State != "missing" {
		t.Fatalf("approval=%+v,%v", a, err)
	}
}

func TestEpicApprovalTerminalAndPersonBoundary(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		actor    string
		terminal bool
		answer   string
		want     bool
	}{{"agent:builder", true, "", false}, {"service:ci", true, "", false}, {"human:mark", false, "", false}, {"human:mark", true, "wrong\n", true}, {"human:mark", true, "", true}} {
		t.Run(tc.actor+tc.answer, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			cmd := newRootCommandAt(f.workspace)
			old, _, _ := cmd.Find([]string{"epic"})
			cmd.RemoveCommand(old)
			cmd.AddCommand(newEpicCommandWithTerminal(func(any) bool { return tc.terminal }))
			var out, errout bytes.Buffer
			cmd.SetOut(&out)
			cmd.SetErr(&errout)
			cmd.SetIn(strings.NewReader(tc.answer))
			cmd.SetArgs([]string{"--json", "epic", "approve", epic, "--actor", tc.actor, "--reason", "Approve"})
			err := cmd.Execute()
			if (err == nil) != tc.want {
				t.Fatalf("approve=%v,%s", err, errout.String())
			}
			if tc.terminal && tc.actor == "human:mark" && !strings.Contains(errout.String(), "VALIDATIONS:") {
				t.Fatal("missing preview")
			}
		})
	}
}

func TestEpicApprovalStalesForEveryChange(t *testing.T) {
	t.Parallel()
	for _, action := range []string{"add", "update", "remove"} {
		t.Run(action, func(t *testing.T) {
			t.Parallel()
			f, epic := newEpicCancellationFixture(t)
			f.clock = time.Now
			base := epicwork.Options{WorkspaceRoot: f.workspace, MigrationsFS: fidmigrations.FS, EpicID: epic, Actor: "agent:builder", IdempotencyKey: "first", Validation: epicplan.Validation{Claim: "Original", Verifier: "review:code"}}
			if _, err := epicwork.Change(base, "add"); err != nil {
				t.Fatal(err)
			}
			approveTestEpic(t, f, epic)
			base.IdempotencyKey = "edit"
			base.ValidationID = "V1"
			base.Validation.Claim = "Changed"
			if _, err := epicwork.Change(base, action); err != nil {
				t.Fatal(err)
			}
			q, err := reviewwork.OpenProjection(f.workspace)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = q.Close() }()
			a, err := epicplan.Current(q, epic)
			if err != nil || a.State != "stale" {
				t.Fatalf("approval=%+v,%v", a, err)
			}
			approveTestEpic(t, f, epic)
			a, err = epicplan.Current(q, epic)
			if err != nil || a.State != "current" {
				t.Fatalf("renewed approval=%+v,%v", a, err)
			}
		})
	}
}

type brokenEpicInput struct{}

func (brokenEpicInput) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func TestEpicApprovalReportsInputError(t *testing.T) {
	t.Parallel()
	f, epic := newEpicCancellationFixture(t)
	cmd := newRootCommandAt(f.workspace)
	old, _, _ := cmd.Find([]string{"epic"})
	cmd.RemoveCommand(old)
	cmd.AddCommand(newEpicCommandWithTerminal(func(any) bool { return true }))
	var out, errout bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errout)
	cmd.SetIn(brokenEpicInput{})
	cmd.SetArgs([]string{"--json", "epic", "approve", epic, "--actor", "human:mark", "--reason", "Approve"})
	if err := cmd.Execute(); err == nil || !strings.Contains(errout.String(), "unexpected EOF") {
		t.Fatalf("input error=%v %s", err, errout.String())
	}
}
