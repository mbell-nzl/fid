package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/policy"
	"github.com/mbell-nzl/fid/internal/policysnapshot"
)

func policyWorkspace(t *testing.T) (string, string) {
	t.Helper()
	root := gittest.NewWorkspace(t)
	runFIDCommand(t, root, "init", "--prefix", "FID")
	_, created := runFIDCommand(t, root, "issue", "create", "--kind", "task", "--title", "Policy task", "--actor", "human:mark")
	return root, parseCreatedIssueID(t, created)
}

func policyOrigin(t *testing.T, root string) {
	t.Helper()
	origin := filepath.Join(t.TempDir(), "origin.git")
	identityGit(t, root, "init", "-q", "--bare", "--initial-branch=trunk", origin)
	identityGit(t, root, "checkout", "-q", "-b", "trunk")
	policyCommit(t, root)
	identityGit(t, root, "remote", "add", "origin", origin)
	identityGit(t, root, "push", "-q", "-u", "origin", "trunk")
	identityGit(t, root, "remote", "set-head", "origin", "trunk")
}

func policyCommit(t *testing.T, root string) string {
	t.Helper()
	identityGit(t, root, "add", ".fid.toml")
	identityGit(t, root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false",
		"commit", "-q", "--no-verify", "--allow-empty", "-m", "policy")
	return strings.TrimSpace(identityGit(t, root, "rev-parse", "HEAD"))
}

func appendPolicy(t *testing.T, root, text string) {
	t.Helper()
	path := filepath.Join(root, ".fid.toml")
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	if _, err := file.WriteString(text); err != nil {
		t.Fatal(err)
	}
}

type policySnapshotRow struct {
	Sequence  int
	Ref       string
	Commit    sql.NullString
	Digest    string
	Effective string
}

func policySnapshots(t *testing.T, root string) []policySnapshotRow {
	t.Helper()
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	rows, err := conn.Query("SELECT sequence, source_ref, source_commit, digest, effective_json FROM policy_snapshots ORDER BY sequence")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	result := []policySnapshotRow{}
	for rows.Next() {
		var row policySnapshotRow
		if err := rows.Scan(&row.Sequence, &row.Ref, &row.Commit, &row.Digest, &row.Effective); err != nil {
			t.Fatal(err)
		}
		result = append(result, row)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestV7WorkingTreeEditsToPolicyChangeNothingUntilMerged(t *testing.T) {
	t.Parallel()
	root, task := policyWorkspace(t)
	appendPolicy(t, root, "\n[checks.unit]\nrun = ['go', 'test', './...']\ntimeout = '1m'\n[lanes.standard]\nchecks = ['unit']\n")
	policyOrigin(t, root)
	decodeNext(t, executeCommandForTest(root, "--json", "next"), "next.get")
	if got := policySnapshots(t, root); len(got) != 0 {
		t.Fatalf("read-only next wrote the first snapshot: %#v", got)
	}
	claim := func() {
		t.Helper()
		decodeNext(t, executeCommandForTest(root, "--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
	}
	claim()
	first := policySnapshots(t, root)
	if len(first) != 1 || first[0].Ref != "origin/trunk" || first[0].Commit.String != strings.TrimSpace(identityGit(t, root, "rev-parse", "origin/trunk")) {
		t.Fatalf("first snapshot = %#v", first)
	}
	var effective policy.Effective
	if err := json.Unmarshal([]byte(first[0].Effective), &effective); err != nil || effective.Lanes["standard"].TokenBudget != 0 || !effective.Lanes["standard"].ChecksSet || len(effective.Checks["unit"].Run) != 3 {
		t.Fatalf("effective policy = %#v, %v", effective, err)
	}
	appendPolicy(t, root, "token_budget = 42\n")
	identityGit(t, root, "add", ".fid.toml") // index edits have no authority either
	claim()
	identityGit(t, root, "checkout", "-q", "-b", "local-policy")
	commit := policyCommit(t, root)
	claim()
	if got := policySnapshots(t, root); len(got) != 1 || got[0] != first[0] {
		t.Fatalf("working tree/index/local commit changed policy: %#v", got)
	}
	// Pushing elsewhere doesn't refresh this clone's origin/trunk. Only a
	// subsequent caller-controlled fetch makes the merged policy authoritative.
	other := gittest.AddWorktree(t, root)
	identityGit(t, other, "push", "-q", "origin", commit+":refs/heads/trunk")
	// Push itself updates the shared remote-tracking ref, so put it back to
	// model an origin changed by a different clone, before this clone fetches.
	identityGit(t, root, "update-ref", "refs/remotes/origin/trunk", first[0].Commit.String)
	claim()
	if got := policySnapshots(t, root); len(got) != 1 {
		t.Fatalf("fid fetched policy: %#v", got)
	}
	identityGit(t, root, "fetch", "-q", "origin")
	// Read-only commands do not snapshot even a newly authoritative policy.
	beforeReads := readAllPolicyJournal(t, root)
	for _, args := range [][]string{{"next"}, {"status"}, {"show", task}} {
		result := executeCommandForTest(root, append([]string{"--json"}, args...)...)
		if result.err != nil {
			t.Fatalf("read %v: %#v", args, result)
		}
	}
	if afterReads := readAllPolicyJournal(t, root); afterReads != beforeReads {
		t.Fatal("read-only commands wrote an event after the effective digest changed")
	}
	if got := policySnapshots(t, root); len(got) != 1 {
		t.Fatalf("read-only commands wrote a snapshot: %#v", got)
	}
	claim() // also checks policy when returning an existing live hold
	second := policySnapshots(t, root)
	if len(second) != 2 || second[1].Commit.String != commit || second[1].Digest == first[0].Digest {
		t.Fatalf("snapshot after fetch = %#v", second)
	}
	if err := json.Unmarshal([]byte(second[1].Effective), &effective); err != nil || effective.Lanes["standard"].TokenBudget != 42 || len(effective.ActorRules["human"]) != 2 {
		t.Fatalf("new effective policy = %#v, %v", effective, err)
	}
	claim()
	before := executeCommandForTest(root, "sql", "SELECT * FROM policy_snapshots ORDER BY sequence")
	runFIDCommand(t, root, "rebuild")
	after := executeCommandForTest(root, "sql", "SELECT * FROM policy_snapshots ORDER BY sequence")
	if before.err != nil || after.err != nil || before.stdout != after.stdout || len(policySnapshots(t, root)) != 2 {
		t.Fatalf("rebuild changed snapshots: before=%#v after=%#v", before, after)
	}
	// Reformatting the authoritative file and a new source commit leave the
	// effective digest alone, so don't produce another snapshot.
	appendPolicy(t, root, "# same effective policy\n")
	policyCommit(t, root)
	identityGit(t, root, "push", "-q", "origin", "HEAD:trunk")
	claim()
	if got := policySnapshots(t, root); len(got) != 2 {
		t.Fatalf("formatting-only origin change wrote a snapshot: %#v", got)
	}
	// Journal order, including the initial snapshot immediately before taking.
	data := readAllPolicyJournal(t, root)
	if strings.Count(data, `"event_type":"policy.snapshotted"`) != 2 || strings.Index(data, `"event_type":"policy.snapshotted"`) > strings.Index(data, `"event_type":"hold.taken"`) {
		t.Fatalf("wrong snapshot/hold journal order: %s", data)
	}
}

func readAllPolicyJournal(t *testing.T, root string) string {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(config.LedgerDir(root), "journal", "*.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	var data strings.Builder
	for _, path := range paths {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		data.Write(contents)
	}
	return data.String()
}

func TestPolicySourceDefaultsAndRefusals(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"no origin", "no file on origin", "missing origin HEAD", "invalid origin policy"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root, _ := policyWorkspace(t)
			if name != "no origin" {
				policyOrigin(t, root)
			}
			var invalidCommit string
			switch name {
			case "no file on origin":
				identityGit(t, root, "rm", "-q", ".fid.toml")
				identityGit(t, root, "-c", "user.name=fid test", "-c", "user.email=fid@example.test", "-c", "commit.gpgsign=false", "commit", "-q", "-m", "remove policy")
				identityGit(t, root, "push", "-q", "origin", "trunk")
			case "missing origin HEAD":
				identityGit(t, root, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
			case "invalid origin policy":
				path := filepath.Join(root, ".fid.toml")
				valid, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				appendPolicy(t, root, "\n[lanes.assured]\nauto_merge = true\n")
				invalidCommit = policyCommit(t, root)
				identityGit(t, root, "push", "-q", "origin", "trunk")
				if err := os.WriteFile(path, valid, 0o600); err != nil {
					t.Fatal(err)
				}
			}
			before := readAllPolicyJournal(t, root)
			result := executeCommandForTest(root, "--json", "next", "--claim", "--actor", "agent:builder")
			if name == "missing origin HEAD" || name == "invalid origin policy" {
				want := "origin's default branch is unknown"
				if name == "invalid origin policy" {
					want = invalidCommit
				}
				if result.err == nil || !strings.Contains(result.stderr, `"code":"policy-unavailable"`) || !strings.Contains(result.stderr, want) {
					t.Fatalf("refusal = %#v, want %q", result, want)
				}
				if after := readAllPolicyJournal(t, root); before != after {
					t.Fatal("refused policy wrote to journal")
				}
				return
			}
			decodeNext(t, result, "next.claim")
			got := policySnapshots(t, root)
			if len(got) != 1 || got[0].Ref != "built-in defaults" || got[0].Commit.Valid {
				t.Fatalf("defaults snapshot = %#v", got)
			}
		})
	}
}

func TestPolicySnapshotRecordsBuiltInUpgradesAndReturnsToEarlierDigests(t *testing.T) {
	t.Parallel()
	root, _ := policyWorkspace(t)
	loaded, err := policy.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	oldDigest := loaded.Effective.Digest()
	// The boundary accepts an effective policy already resolved by fid. Use
	// that seam to represent a built-in rule upgrade without altering globals.
	err = journal.WithLock(root, journal.LockOwner{Command: "test policy upgrade"}, func(lock *journal.WriteLock) error {
		conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
		if err != nil {
			return err
		}
		defer closeCommandDB(t, conn)
		for _, version := range []int{1, 1, 2, 2, 1, 1} {
			loaded.Effective.RulesVersion = version
			if err := policysnapshot.Ensure(lock, conn, loaded, time.Now()); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	snapshots := policySnapshots(t, root)
	if len(snapshots) != 3 || snapshots[0].Digest != oldDigest || snapshots[1].Digest == oldDigest || snapshots[2].Digest != oldDigest {
		t.Fatalf("rule upgrade A→B→A snapshots = %#v", snapshots)
	}
	runFIDCommand(t, root, "rebuild")
	for index, replayed := range policySnapshots(t, root) {
		if replayed != snapshots[index] {
			t.Fatalf("upgrade snapshot changed during replay: %#v, want %#v", replayed, snapshots[index])
		}
	}
}

func TestNextClaimUpgradesAFreshPrePolicyProjection(t *testing.T) {
	t.Parallel()
	root, task := policyWorkspace(t)
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	// Reproduce a fully caught-up ledger left by the previous binary:
	// migration 20, projection format 5, with no missing journal Events.
	for _, statement := range []string{
		"DROP TABLE policy_snapshots",
		"UPDATE schema_migrations SET version = 20",
		"UPDATE projection_meta SET projection_schema_version = 5",
	} {
		if _, err := conn.Exec(statement); err != nil {
			closeCommandDB(t, conn)
			t.Fatal(err)
		}
	}
	closeCommandDB(t, conn)
	result, _ := decodeNext(t, executeCommandForTest(root, "--json", "next", "--claim", "--actor", "agent:builder"), "next.claim")
	if result["task"].(map[string]any)["id"] != task || len(policySnapshots(t, root)) != 1 {
		t.Fatalf("upgraded next = %#v, want the original task and its policy snapshot", result)
	}
	conn, err = fiddb.Open(filepath.Join(config.LedgerDir(root), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer closeCommandDB(t, conn)
	var migration, projection int
	if err := conn.QueryRow("SELECT version FROM schema_migrations").Scan(&migration); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRow("SELECT projection_schema_version FROM projection_meta WHERE key = 'default'").Scan(&projection); err != nil {
		t.Fatal(err)
	}
	if migration != 22 || projection != fiddb.ProjectionSchemaVersion {
		t.Fatalf("upgraded migration/projection = %d/%d", migration, projection)
	}
}
