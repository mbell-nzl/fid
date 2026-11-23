package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/defaultconflict"
	"github.com/mbell-nzl/fid/internal/doneeval"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/reviewwork"
)

func approveConflictReviews(t *testing.T, f *statusFixture, task, lane string) {
	t.Helper()
	for _, stream := range []string{"primary", "separate"} {
		if stream == "separate" && lane != "assured" {
			continue
		}
		in := reviewInput()
		in["stream"] = stream
		in["reviewer"].(map[string]any)["session_id"] = stream
		lifecycleResult(t, recordReviewTask(t, f, in, "approve-"+stream, "human:courier", task), "review.record")
	}
}

func moveConflictDefault(t *testing.T, f *statusFixture, conflict bool) string {
	t.Helper()
	other := gittest.AddWorktree(t, f.workspace)
	runGitForTest(t, other, "checkout", "--detach", "origin/trunk")
	path := "unrelated.txt"
	if conflict {
		path = "notes.md"
	}
	writeBackupTestFile(t, filepath.Join(other, path), "default branch change\n")
	runGitForTest(t, other, "add", ".")
	runGitForTest(t, other, "commit", "-m", "Move default")
	commit := strings.TrimSpace(identityGit(t, other, "rev-parse", "HEAD"))
	runGitForTest(t, other, "update-ref", "refs/remotes/origin/trunk", commit)
	return commit
}

func TestApprovedDefaultConflictSelection(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name               string
		approved, conflict bool
		missing            string
		want               string
	}{
		{"conflict", true, true, "", "rework"},
		{"clean", true, false, "", "none"},
		{"recorded merge", true, true, "", "none"},
		{"review due", false, true, "", "none"},
		{"separate due", true, true, "", "none"},
		{"missing default commit", true, true, "default", "none"},
		{"missing ref", true, true, "ref", "none"},
		{"missing head", true, true, "head", "none"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			lane := "standard"
			if tc.name == "separate due" {
				lane = "assured"
			}
			f, task := canonicalReviewFixtureWithLane(t, lane)
			if tc.approved {
				approveConflictReviews(t, f, task, "standard")
			}
			commit := moveConflictDefault(t, f, tc.conflict)
			if tc.name == "recorded merge" {
				startSQL(t, f, fmt.Sprintf(`INSERT INTO merge_observations(id,task_id,handoff_id,merge_commit,record_json) SELECT 'MRG-abcdef',task_id,id,'%s','{}' FROM handoffs WHERE task_id='%s'`, strings.Repeat("a", 40), task))
			}
			if tc.missing == "ref" {
				runGitForTest(t, f.workspace, "update-ref", "-d", "refs/remotes/origin/trunk")
			}
			if tc.missing == "head" || tc.missing == "default" {
				head := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))
				if tc.missing == "default" {
					head = commit
				}
				objects := strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "--path-format=absolute", "--git-path", "objects"))
				path := filepath.Join(objects, head[:2], head[2:])
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.WriteFile(path, raw, 0o444); err != nil {
						t.Error(err)
					}
				})
			}
			before := startEventCount(t, f.workspace)
			got, _ := lifecycleResult(t, f.run("--json", "next", "--explain"), "next.get")
			if got["kind"] != tc.want {
				t.Fatalf("next kind=%v, want %s: %v", got["kind"], tc.want, got)
			}
			explained := f.run("next", "--explain")
			if explained.err != nil {
				t.Fatal(explained)
			}
			if tc.missing != "" && !strings.Contains(explained.stdout, "unavailable") {
				t.Fatalf("missing local evidence unexplained: %s", explained.stdout)
			}
			if tc.want == "rework" {
				evidence, ok := got["rework"].(map[string]any)
				if !ok || evidence["cause"] != "default-conflict" || evidence["default_commit"] != commit || fmt.Sprint(evidence["conflicting_paths"]) != "[notes.md]" {
					t.Fatalf("conflict evidence=%v", got)
				}
				if !strings.Contains(explained.stdout, "default-conflict") {
					t.Fatalf("cause unexplained: %s", explained.stdout)
				}
				status, _ := lifecycleResult(t, f.run("--json", "status"), "project.status")
				if !strings.Contains(fmt.Sprint(status["needs_you"]), "no longer merges into origin/trunk") || !strings.Contains(fmt.Sprint(status["needs_you"]), "fid next --claim") {
					t.Fatalf("status still asks for merge: %v", status)
				}
				original := f.workspace
				f.workspace = gittest.AddWorktree(t, original)
				other, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:other", "--model", "codex/gpt-6.1-sol"), "next.claim")
				if other["kind"] != "none" {
					t.Fatalf("another worktree got update: %v", other)
				}
				f.workspace = original
			}
			if startEventCount(t, f.workspace) != before {
				t.Fatal("read or other worktree appended update")
			}
		})
	}
}

func TestApprovedDefaultConflictClaimAndFreshness(t *testing.T) {
	t.Parallel()
	for _, unchanged := range []bool{true, false} {
		t.Run(fmt.Sprint("unchanged=", unchanged), func(t *testing.T) {
			t.Parallel()
			f, task := canonicalReviewFixtureWithLane(t, "assured")
			approveConflictReviews(t, f, task, "assured")
			if got := approveForTest(f, task, task+"\n", true); got.err != nil {
				t.Fatal(got)
			}
			moveConflictDefault(t, f, true)
			_, revision := taskState(t, f, task)
			got, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
			if got["kind"] != "rework" {
				t.Fatalf("approved conflict not returned: %v", got)
			}
			if got["task"].(map[string]any)["status"] != "in_progress" {
				t.Fatalf("claimed task state=%v", got["task"])
			}
			hold := got["hold"].(map[string]any)
			if hold["replaces"] == nil {
				t.Fatal("missing kept reservation link")
			}
			if state, rev := taskState(t, f, task); state != "in_progress" || rev != revision+1 {
				t.Fatalf("update state=%s revision=%v", state, rev)
			}
			if unchanged {
				// A second default move reverts its conflict. Merging that tip preserves
				// the original patch, so the new hand-off must retain both approvals.
				other := gittest.AddWorktree(t, f.workspace)
				runGitForTest(t, other, "checkout", "--detach", "origin/trunk")
				runGitForTest(t, other, "revert", "--no-edit", "HEAD")
				runGitForTest(t, other, "update-ref", "refs/remotes/origin/trunk", "HEAD")
				runGitForTest(t, f.workspace, "merge", "--no-edit", "origin/trunk")
			} else {
				// Resolve the add/add conflict explicitly, as the builder would.
				runGitForTest(t, f.workspace, "merge", "-s", "ours", "--no-edit", "origin/trunk")
			}
			lifecycleResult(t, f.run("--json", "check", "--hold", hold["id"].(string), "--actor", "agent:builder", "--idempotency-key", "update-check"), "check.execute")
			reviewHandoffTask(t, f, "updated-handoff", task)
			for _, rebuild := range []bool{false, true} {
				if rebuild {
					if _, err := journal.RebuildWorkspace(f.workspace, issueMigrationsDir()); err != nil {
						t.Fatal(err)
					}
				}
				rounds, err := reviewwork.List(f.workspace, task)
				if err != nil || len(rounds) != 2 || rounds[0].Current != unchanged || rounds[1].Current != unchanged {
					t.Fatalf("round currency=%v/%v", rounds, err)
				}
				shown, _ := lifecycleResult(t, f.run("--json", "show", task), "issue.show")
				want := "stale"
				if unchanged {
					want = "current"
				}
				if shown["approval"].(map[string]any)["state"] != want {
					t.Fatalf("approval currency=%v", shown["approval"])
				}
			}
		})
	}
}

func TestApprovedTaskCannotReturnForMinorOrCapCharge(t *testing.T) {
	t.Parallel()
	// No CHANGES rounds are allowed: a default conflict still permits an update.
	f, task := canonicalReviewFixtureWithClock(t, "standard", false, time.Now, func(f *statusFixture) {
		trustedShowPolicy(t, f, func(raw string) string {
			return strings.Replace(raw, "[lanes.standard]", "[lanes.standard]\nrework_cap=0", 1)
		})
	})
	input := reviewInput()
	input["findings"] = []any{reviewFinding("minor")}
	lifecycleResult(t, recordReviewTask(t, f, input, "minor-approval", "human:courier", task), "review.record")
	before := readAllPolicyJournal(t, f.workspace)
	refused, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if refused["kind"] != "none" || readAllPolicyJournal(t, f.workspace) != before {
		t.Fatalf("approved minor returned to builder: %v", refused)
	}
	moveConflictDefault(t, f, true)
	claimed, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if claimed["kind"] != "rework" {
		t.Fatalf("conflict charged a CHANGES round: %v", claimed)
	}
	again, event := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if again["kind"] != "resume" || event != nil {
		t.Fatalf("update blocked by cap: %v", again)
	}
}

func TestDefaultConflictGolden(t *testing.T) {
	t.Parallel()
	f := reviewFixture(t)
	approveConflictReviews(t, f, "FID-9f13aa", "standard")
	moveConflictDefault(t, f, true)
	for _, tc := range []struct {
		name string
		args []string
	}{
		{"next", []string{"next", "--explain"}},
		{"status", []string{"status"}},
		{"claimed", []string{"--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"}},
	} {
		modes := []bool{false, true}
		if tc.name == "claimed" {
			modes = []bool{true}
		}
		for _, jsonMode := range modes {
			args, ext := tc.args, ".txt"
			if jsonMode {
				args, ext = append([]string{"--json"}, args...), ".json"
			}
			got := f.run(args...)
			if got.err != nil {
				t.Fatal(got)
			}
			handoffGolden(t, f, filepath.Join("..", "default-conflict", tc.name+ext), reworkGoldenOutput(f, got.stdout))
		}
	}
}

func TestDefaultConflictAppendRefusalAndRecovery(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	approveConflictReviews(t, f, task, "standard")
	moveConflictDefault(t, f, true)
	conn, err := reviewwork.OpenProjection(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	snapshot, err := doneeval.GatherReviews(f.workspace, conn, task, nil)
	if err != nil {
		t.Fatal(err)
	}
	evidence, err := defaultconflict.Inspect(f.workspace, snapshot.Handoff.ID, snapshot.Handoff.Head)
	if err != nil {
		t.Fatal(err)
	}
	record := defaultconflict.Record{Evidence: *evidence, Binding: snapshot.Fingerprint}
	effects, err := journal.DefaultConflictEffects(conn, task, snapshot.Handoff.HoldID, "DFC-abcdef", record, f.clock())
	if err != nil {
		t.Fatal(err)
	}
	event := journal.Event{EventSchema: 2, EventType: "task.default-conflict", CommandName: "next.claim", OccurredAt: f.clock(), ActorV2: journal.EventActor{ActorType: "agent", ActorID: "builder", Producer: "fid-standalone", ProvenanceRef: "test"}, Aggregate: journal.Aggregate{ID: "DFC-abcdef", Type: "default-conflict", Sequence: 1}, Effects: effects}
	before := readAllPolicyJournal(t, f.workspace)
	err = journal.WithLock(f.workspace, journal.LockOwner{Command: "test"}, func(lock *journal.WriteLock) error {
		malformed := event
		malformed.Effects = effects[:2]
		if _, err := journal.CommitEvent(lock, conn, malformed); err == nil {
			t.Fatal("split return appended")
		}
		if readAllPolicyJournal(t, f.workspace) != before {
			t.Fatal("refused append changed journal")
		}
		var state string
		var stamp any
		if err := conn.QueryRow(`SELECT i.status,h.rework_at FROM issues i JOIN holds h ON h.task_id=i.id WHERE h.id=?`, snapshot.Handoff.HoldID).Scan(&state, &stamp); err != nil || state != "in_review" || stamp != nil {
			t.Fatalf("refused append changed state: %s/%v/%v", state, stamp, err)
		}
		// Simulate interruption after the atomic return, before hold.taken.
		_, err := journal.CommitEvent(lock, conn, event)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	got, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if got["kind"] != "rework" || got["rework"] == nil || got["rework"].(map[string]any)["cause"] != "default-conflict" {
		t.Fatalf("return interrupted before hold lost its cause: %v", got)
	}
}

func TestDefaultConflictUpgradesExistingProjection(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	approveConflictReviews(t, f, task, "standard")
	moveConflictDefault(t, f, true)
	// Model the previous release's complete projection and matching journal.
	startSQL(t, f, `DROP TABLE default_conflicts; UPDATE schema_migrations SET version=37; UPDATE projection_meta SET projection_schema_version=19`)
	got, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if got["kind"] != "rework" {
		t.Fatalf("existing projection cannot return conflict: %v", got)
	}
}

func TestDefaultConflictRoutineWithoutRequiredReviews(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixtureWithLane(t, "routine")
	moveConflictDefault(t, f, true)
	got, _ := lifecycleResult(t, f.run("--json", "next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol"), "next.claim")
	if got["kind"] != "rework" || got["rework"].(map[string]any)["cause"] != "default-conflict" {
		t.Fatalf("Routine conflict cannot return without a review round: %v", got)
	}
	listed, _ := lifecycleResult(t, f.run("--json", "review", "list", task), "review.list")
	if len(listed["rounds"].([]any)) != 0 {
		t.Fatalf("Routine update manufactured a review round: %v", listed)
	}
}

func TestDefaultConflictRetainsAllowanceGuidance(t *testing.T) {
	t.Parallel()
	f, task := canonicalReviewFixture(t)
	approveConflictReviews(t, f, task, "standard")
	const guidance = "Keep the documented compatibility behavior"
	lifecycleResult(t, f.run("--json", "rework", "allow", task, "--rounds", "1", "--actor", "human:mark", "--reason", guidance), "task.rework-allow")
	moveConflictDefault(t, f, true)
	got := f.run("next", "--explain")
	if got.err != nil || !strings.Contains(got.stdout, guidance) {
		t.Fatalf("default update lost existing allowance guidance: %#v", got)
	}
	if !strings.Contains(got.stdout, "fid show "+task) {
		t.Fatalf("default update lost the task inspection action: %s", got.stdout)
	}
}
