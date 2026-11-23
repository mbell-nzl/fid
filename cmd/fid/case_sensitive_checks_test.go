package main

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/handoffrecord"
	"github.com/mbell-nzl/fid/internal/policy"
)

func TestCaseSensitiveCheckNamesThroughClaimAndHandoff(t *testing.T) {
	t.Parallel()
	root := gittest.NewWorkspace(t)
	runFIDCommand(t, root, "init", "--prefix", "FID")
	runFIDCommand(t, root, "issue", "create", "--kind", "task", "--title", "Case-sensitive checks", "--lane", "standard", "--actor", "human:mark")
	appendPolicy(t, root, "\n[checks.build]\nrun=['true']\ntimeout='1m'\n[checks.Build]\nrun=['true']\ntimeout='2m'\n[lanes.standard]\nchecks=['build','Build']\n[review]\ntest_adequacy=[]\n")
	policyOrigin(t, root)
	runGitForTest(t, root, "checkout", "-q", "-b", "task")
	run := func(args ...string) map[string]any {
		t.Helper()
		fields, _ := lifecycleResult(t, executeCommandForTest(root, append([]string{"--json"}, args...)...), map[string]string{"next": "next.claim", "start": "task.start", "check": "check.execute", "handoff": "task.handoff"}[args[0]])
		return fields
	}
	claimed := run("next", "--claim", "--actor", "agent:builder", "--model", "codex/gpt-6.1-sol")
	hold := claimed["hold"].(map[string]any)["id"].(string)
	run("start", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "start")
	checked := executeCommandForTest(root, "--json", "check", "--hold", hold, "--actor", "agent:builder", "--idempotency-key", "checks")
	if got := requiredNames(t, checked); !slices.Equal(got, []string{"Build", "build"}) {
		t.Fatalf("counting checks = %v", got)
	}
	run("handoff", "--hold", hold, "--actor", "agent:builder", "--pr", "https://github.com/mbell-nzl/fid/pull/101", "--idempotency-key", "handoff")
	read := func() handoffrecord.Record {
		t.Helper()
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
		if err != nil {
			t.Fatal(err)
		}
		defer closeCommandDB(t, conn)
		var raw string
		if err := conn.QueryRow(`SELECT record_json FROM handoffs`).Scan(&raw); err != nil {
			t.Fatal(err)
		}
		var record handoffrecord.Record
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			t.Fatal(err)
		}
		return record
	}
	snapshots := policySnapshots(t, root)
	var effective policy.Effective
	if len(snapshots) != 1 {
		t.Fatalf("snapshots = %v, want one", snapshots)
	}
	if err := json.Unmarshal([]byte(snapshots[0].Effective), &effective); err != nil || len(effective.Checks) != 2 || effective.Checks["build"].Timeout != time.Minute || effective.Checks["Build"].Timeout != 2*time.Minute {
		t.Fatalf("case-sensitive policy = %#v: %v", effective.Checks, err)
	}
	before := read()
	if len(before.DefinitionDigests) != 2 || before.DefinitionDigests["build"] == "" || before.DefinitionDigests["Build"] == "" || before.DefinitionDigests["build"] == before.DefinitionDigests["Build"] {
		t.Fatalf("definition digests = %v", before.DefinitionDigests)
	}
	task21Read(t, root, "--json", "rebuild")
	if after := read(); !reflect.DeepEqual(before, after) || !reflect.DeepEqual(snapshots, policySnapshots(t, root)) {
		t.Fatal("rebuild changed policy or hand-off receipt")
	}
}
