package main

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/checkrecord"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func laneFloorFixture(t *testing.T, lane string) *statusFixture {
	t.Helper()
	f := requiredCheckFixture(t, "['always']")
	startSQL(t, f, `UPDATE issues SET lane='`+lane+`' WHERE id='FID-9f13aa'`)
	path := filepath.Join(f.workspace, ".fid.toml")
	raw := string(mustReadCommandFile(t, path))
	raw = strings.Replace(raw, "[lanes.routine]\nchecks = []", "[lanes.routine]\npaths=['**/*.md']\nchecks = []", 1)
	raw = strings.Replace(raw, "[lanes.assured]\nchecks = []", "[lanes.assured]\npaths=['secure/**']\nchecks = ['always','go']", 1)
	raw = strings.Replace(raw, "[checks.always]", "[checks.always]\ndefined_by=['tools/**']", 1)
	writeBackupTestFile(t, path, raw)
	writeBackupTestFile(t, filepath.Join(f.workspace, "secure/old.md"), "old\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Floor policy and base")
	identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
	return f
}

func TestLaneFloorChecksAndAtomicHandoff(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, lane, path, want string
		names                  []string
	}{
		{"no changes", "routine", "", "routine", []string{}},
		{"routine change", "routine", "task.md", "routine", []string{}},
		{"never lower", "assured", "task.md", "assured", []string{"always"}},
		{"standard fallback", "routine", "task.go", "standard", []string{"always"}},
		{"assured path", "routine", "secure/task.go", "assured", []string{"always", "go"}},
		{"defined by", "standard", "tools/build.go", "assured", []string{"always", "go"}},
		{"hard floor", "standard", ".github/workflows/test.yml", "assured", []string{"always"}},
		{"unsupported UTF-8", "standard", "odd\xff.md", "assured", []string{"always", "go"}},
		{"unsupported newline", "standard", "odd\n.md", "assured", []string{"always", "go"}},
		{"deleted", "routine", "delete", "assured", []string{"always"}},
		{"renamed from", "routine", "rename", "assured", []string{"always"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := laneFloorFixture(t, tc.lane)
			switch tc.path {
			case "delete":
				identityGit(t, f.workspace, "rm", "secure/old.md")
			case "rename":
				identityGit(t, f.workspace, "mv", "secure/old.md", "plain.md")
			case "":
			default:
				writeBackupTestFile(t, filepath.Join(f.workspace, tc.path), "task\n")
			}
			if tc.path != "" {
				runGitForTest(t, f.workspace, "add", ".")
				runGitForTest(t, f.workspace, "commit", "-m", "Task change")
			}
			got := requiredNames(t, checkRequired(t, f, "floor-checks"))
			if !slices.Equal(got, tc.names) {
				t.Fatalf("checks=%v, want %v", got, tc.names)
			}
			// Check executes the effective lane without editing the task.
			detail := f.run("--json", "show", "FID-9f13aa")
			if !strings.Contains(detail.stdout, `"lane":"`+tc.lane+`"`) {
				t.Fatalf("check changed task lane: %s", detail.stdout)
			}
			lifecycleResult(t, f.run(append([]string{"--json"}, handoffArgs("floor-handoff")...)...), "task.handoff")
			var record handoffrecord.Record
			rows := f.run("sql", "SELECT record_json FROM handoffs WHERE task_id='FID-9f13aa'")
			if rows.err != nil {
				t.Fatal(rows.err)
			}
			var data []map[string]string
			if err := json.Unmarshal([]byte(rows.stdout), &data); err != nil {
				t.Fatalf("sql=%s: %v", rows.stdout, err)
			}
			if err := json.Unmarshal([]byte(data[0]["record_json"]), &record); err != nil {
				t.Fatal(err)
			}
			floor, err := policy.ComputeFloor(mustLoadFloorPolicy(t, f), record.Paths())
			if err != nil {
				t.Fatal(err)
			}
			comparison, err := policy.Compare(f.workspace, record.TaskID)
			if err != nil || !slices.Equal(record.Paths(), comparison.Paths) {
				t.Fatalf("receipt lost Git path bytes: %q / %q / %v", record.Paths(), comparison.Paths, err)
			}
			if record.Lane != tc.want || record.Floor == nil || checkrecord.Digest(record.Floor) != checkrecord.Digest(floor) {
				t.Fatalf("receipt=%+v floor=%+v", record, floor)
			}
			shown := f.run("--json", "show", "FID-9f13aa")
			var result struct {
				Result struct {
					Issue struct {
						Lane string `json:"lane"`
					} `json:"issue"`
					Floor *policy.Floor `json:"floor"`
				} `json:"result"`
			}
			if err := json.Unmarshal([]byte(shown.stdout), &result); err != nil {
				t.Fatal(err)
			}
			if result.Result.Issue.Lane != tc.want || result.Result.Floor == nil {
				t.Fatalf("show=%s", shown.stdout)
			}
			if tc.want == "assured" && (!strings.Contains(shown.stdout, `"approval":{"state":"missing"}`) || !strings.Contains(shown.stdout, `"no_other_model_review":true`)) {
				t.Fatalf("raised Assured lane needs approval: %s", shown.stdout)
			}
			if tc.want != tc.lane {
				text := f.run("show", "FID-9f13aa")
				if !strings.Contains(strings.Join(strings.Fields(text.stdout), " "), "raised from "+strings.ToUpper(tc.lane[:1])+tc.lane[1:]) || !strings.Contains(text.stdout, "floor") {
					t.Fatalf("show did not explain raise: %s", text.stdout)
				}
			}
		})
	}
}

func mustLoadFloorPolicy(t *testing.T, f *statusFixture) policy.Effective {
	t.Helper()
	loaded, err := policy.Load(f.workspace)
	if err != nil {
		t.Fatal(err)
	}
	return loaded.Effective
}

func TestLaneFloorSelectsConfiguredHigherLaneEvenWhenPlannedLaneIsUnset(t *testing.T) {
	t.Parallel()
	for _, empty := range []bool{false, true} {
		t.Run(map[bool]string{false: "required", true: "empty"}[empty], func(t *testing.T) {
			t.Parallel()
			f := laneFloorFixture(t, "routine")
			path := filepath.Join(f.workspace, ".fid.toml")
			raw := string(mustReadCommandFile(t, path))
			raw = strings.Replace(raw, "paths=['**/*.md']\nchecks = []", "paths=['**/*.md']", 1)
			if empty {
				raw = strings.Replace(raw, "checks = ['always','go']", "checks = []", 1)
			}
			writeBackupTestFile(t, path, raw)
			policyCommit(t, f.workspace)
			identityGit(t, f.workspace, "push", "-q", "origin", "HEAD:trunk")
			writeBackupTestFile(t, filepath.Join(f.workspace, "secure/task.go"), "task\n")
			runGitForTest(t, f.workspace, "add", ".")
			runGitForTest(t, f.workspace, "commit", "-m", "Assured task")
			names := requiredNames(t, checkRequired(t, f, "configured-floor"))
			if empty && len(names) != 0 || !empty && !slices.Equal(names, []string{"always", "go"}) {
				t.Fatalf("checks=%v", names)
			}
		})
	}
}

func TestLaneFloorShowGolden(t *testing.T) {
	t.Parallel()
	f := laneFloorFixture(t, "routine")
	for _, path := range []string{"secure/a.md", "secure/b.md", "secure/c.md"} {
		writeBackupTestFile(t, filepath.Join(f.workspace, path), "task\n")
	}
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Assured paths")
	requiredNames(t, checkRequired(t, f, "show-floor"))
	lifecycleResult(t, f.run(append([]string{"--json"}, handoffArgs("show-floor-handoff")...)...), "task.handoff")
	for _, unavailable := range []string{"known", "comparison-unavailable", "policy-unavailable"} {
		if unavailable == "comparison-unavailable" {
			writeBackupTestFile(t, filepath.Join(f.workspace, ".git", "shallow"), strings.TrimSpace(identityGit(t, f.workspace, "rev-parse", "HEAD"))+"\n")
		}
		if unavailable == "policy-unavailable" {
			identityGit(t, f.workspace, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
		}
		// The immutable floor remains explainable even without local origin metadata.
		for _, jsonMode := range []bool{false, true} {
			args, ext := []string{"show", "FID-9f13aa"}, ".txt"
			if jsonMode {
				args = append([]string{"--json"}, args...)
				ext = ".json"
			}
			before := startEventCount(t, f.workspace)
			got := f.run(args...)
			if got.err != nil || startEventCount(t, f.workspace) != before {
				t.Fatalf("read-only show=%#v", got)
			}
			if jsonMode && unavailable != "known" {
				var envelope struct {
					Result struct {
						Warnings []struct {
							Kind string `json:"kind"`
						} `json:"warnings"`
					} `json:"result"`
				}
				if err := json.Unmarshal([]byte(got.stdout), &envelope); err != nil {
					t.Fatal(err)
				}
				if len(envelope.Result.Warnings) != 1 || envelope.Result.Warnings[0].Kind != unavailable {
					t.Fatalf("wrong warning: %s", got.stdout)
				}
			}
			handoffGolden(t, f, "floor-"+unavailable+ext, got.stdout)
		}
	}
}

func TestShowFloorBeforeLaneSet(t *testing.T) {
	t.Parallel()
	f := laneFloorFixture(t, "routine")
	writeBackupTestFile(t, filepath.Join(f.workspace, "secure/task.go"), "task\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Unplanned task change")
	startSQL(t, f, `UPDATE issues SET lane=NULL WHERE id='FID-9f13aa'`)
	got := f.run("--json", "show", "FID-9f13aa")
	var envelope struct {
		Result struct {
			Floor *policy.Floor `json:"floor"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(got.stdout), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.Result.Floor == nil || envelope.Result.Floor.Raise("") != "assured" {
		t.Fatalf("show lost known floor: %s", got.stdout)
	}
}

func TestHandoffFloorRefusalDoesNotRaiseLane(t *testing.T) {
	t.Parallel()
	f := laneFloorFixture(t, "routine")
	writeBackupTestFile(t, filepath.Join(f.workspace, "secure/task.go"), "task\n")
	runGitForTest(t, f.workspace, "add", ".")
	runGitForTest(t, f.workspace, "commit", "-m", "Unverified Assured change")
	before := startEventCount(t, f.workspace)
	got := f.run(append([]string{"--json"}, handoffArgs("missing-floor-checks")...)...)
	assertJSONFailureCode(t, got.stderr, "task.handoff", "checks-unsatisfied")
	shown := f.run("--json", "show", "FID-9f13aa")
	if startEventCount(t, f.workspace) != before || !strings.Contains(shown.stdout, `"lane":"routine"`) || !strings.Contains(shown.stdout, `"handoff":null`) {
		t.Fatalf("refusal changed task: %s", shown.stdout)
	}
}
