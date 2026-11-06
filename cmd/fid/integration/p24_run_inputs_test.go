//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/gittest"
)

func p24RunInput(t *testing.T, a *p22FIDRecoveryAdapter, leaf, key string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "run", leaf, a.actor, key, input)
}

func p24RunAdapter(t *testing.T) *p22FIDRecoveryAdapter {
	t.Helper()
	return p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, p22FinalActions)
}

func p24TerminalInput(t *testing.T, a *p22FIDRecoveryAdapter, leaf string) map[string]any {
	t.Helper()
	input := map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "observation_id": a.observation, "expected_repository": p24ExpectedRepository(a.repository)}
	if leaf == "abandon" {
		input["expected_claim_version"] = a.claimVersion
		input["authority"] = map[string]any{"mode": "ordinary", "actor_id": a.actor, "reason": "stop work", "reference": nil, "grant_ref": nil}
	}
	return input
}

func TestP24RunTerminalExpectationObservationAndUnsafe(t *testing.T) {
	t.Parallel()
	for _, leaf := range []string{"close", "abandon"} {
		for _, kind := range []string{"match", "HEAD", "dirty", "unmerged", "unborn-HEAD"} {
			t.Run(leaf+"/"+kind, func(t *testing.T) {
				t.Parallel()
				a := p24RunAdapter(t)
				a.p22Observe(t, "terminated")
				p22Write(t, a.workspace, "base.txt", "committed after start\n")
				p22Git(t, a.workspace, "commit", "-qam", "advance Run HEAD")
				p22Write(t, a.workspace, "base.txt", "edited after commit\n")
				a.repository = p24RunDirtyRepository(t, a.workspace, "base.txt")
				input := p24TerminalInput(t, a, leaf)
				expected := input["expected_repository"].(map[string]any)
				if kind == "HEAD" {
					expected["head_oid"] = strings.Repeat("f", 40)
				}
				if kind == "dirty" {
					expected["dirty"] = false
				}
				if kind == "unmerged" {
					p22Git(t, a.workspace, "checkout", "--", "base.txt")
					p24MakeUnmerged(t, a.workspace)
				}
				if kind == "unborn-HEAD" {
					gittest.UnbornHead(t, a.workspace)
				}
				before, events := p22FaultProjection(t, a.workspace), p21JournalEventCount(t, a.workspace)
				response := p24RunInput(t, a, leaf, "terminal", input)
				if kind == "match" {
					if response.ExitCode != 0 {
						t.Fatalf("terminal: %#v", response)
					}
					record := p21Path(t, p23Envelope(t, response.Stdout), "result", "run", "terminal_repository").(map[string]any)
					head := p22GitOutput(t, a.workspace, "rev-parse", "HEAD")
					tree := p22GitOutput(t, a.workspace, "rev-parse", "HEAD^{tree}")
					root, _ := filepath.Abs(a.workspace)
					if record["head_oid"] != head || record["base_oid"] != head || record["tree_oid"] != tree || record["dirty"] != true || p22GitOutput(t, a.workspace, "status", "--porcelain", "--", "base.txt") == "" || record["worktree_id"] != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root))) || record["project_id"] != a.projectID || record["capture_source"] != "git" {
						t.Fatalf("terminal observation: %#v", record)
					}
					if !reflect.DeepEqual(record, a.repository) {
						t.Fatalf("terminal = %#v, want independent repository %#v", record, a.repository)
					}
					p24AssertRunEventRepository(t, a, record)
					event := p24EventByID(t, a.workspace, p21PathString(t, p23Envelope(t, response.Stdout), "event", "id"))
					if !reflect.DeepEqual(p24JSONMap(t, event.Repository), a.repository) {
						t.Fatalf("Event repository = %#v, want independent repository %#v", event.Repository, a.repository)
					}
					return
				}
				code := "repository-mismatch"
				if kind == "unmerged" || kind == "unborn-HEAD" {
					code = "recovery-state-conflict"
				}
				assertP22ProcessFailure(t, "run."+leaf, code, 4, response)
				if code == "repository-mismatch" {
					details := p21Path(t, p23Envelope(t, response.Stderr), "error", "details")
					observed := cloneP23Repository(a.repository)
					raw, _ := json.Marshal(observed)
					if err := json.Unmarshal(raw, &observed); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(details, map[string]any{"expected": expected, "observed": observed}) {
						t.Fatalf("mismatch details: %#v", details)
					}
				}
				if events != p21JournalEventCount(t, a.workspace) || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
					t.Fatal("terminal rejection changed durable state")
				}
			})
		}
	}
}

// Require the entire dirty fixture to be one unstaged tracked edit, and derive
// its identity from Git and config rather than the production observer.
func p24RunDirtyRepository(t *testing.T, workspace, path string) map[string]any {
	t.Helper()
	gitBytes := func(args ...string) []byte {
		t.Helper()
		data, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).Output()
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	git := func(args ...string) string {
		t.Helper()
		return strings.TrimSuffix(string(gitBytes(args...)), "\n")
	}
	if status := git("status", "--porcelain", "--untracked-files=all", "--", "."); status != " M "+path {
		t.Fatalf("dirty fixture status = %q, want one tracked edit %s", status, path)
	}
	root, err := filepath.Abs(git("rev-parse", "--show-toplevel"))
	if err != nil {
		t.Fatal(err)
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(workspace), "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	head := git("rev-parse", "HEAD")
	image := func(mode string, data []byte) map[string]any {
		return map[string]any{"mode": mode, "digest": fmt.Sprintf("sha256:%x", sha256.Sum256(data))}
	}
	headEntry := strings.Fields(git("ls-tree", "HEAD", "--", path))
	indexEntry := strings.Fields(git("ls-files", "-s", "--", path))
	if len(headEntry) != 4 || headEntry[1] != "blob" || len(indexEntry) != 4 || indexEntry[2] != "0" {
		t.Fatalf("unsupported tracked fixture: HEAD %v / index %v", headEntry, indexEntry)
	}
	info, err := os.Lstat(filepath.Join(root, path))
	if err != nil || !info.Mode().IsRegular() {
		t.Fatalf("worktree fixture must be regular: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	mode := "100644"
	if info.Mode().Perm()&0o111 != 0 {
		mode = "100755"
	}
	// Hash the actual HEAD, index and worktree bytes, with sorted map keys
	// from encoding/json instead of production observation/canonicalization.
	entries := []any{map[string]any{"path": path,
		"head":     image(headEntry[0], gitBytes("cat-file", "blob", headEntry[2])),
		"index":    image(indexEntry[0], gitBytes("cat-file", "blob", indexEntry[1])),
		"worktree": image(mode, content)}}
	payload, err := json.Marshal(map[string]any{"schema": "fid.repository-dirty.v1", "entries": entries})
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{"base_oid": head, "capture_source": "git", "dirty": true, "dirty_digest": fmt.Sprintf("sha256:%x", sha256.Sum256(payload)),
		"head_oid": head, "project_id": projectConfig.ProjectID, "tree_oid": git("rev-parse", "HEAD^{tree}"),
		"worktree_id": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root)))}
}

func p24AssertRunEventRepository(t *testing.T, a *p22FIDRecoveryAdapter, want map[string]any) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Join(config.LedgerDir(a.workspace), "journal"))
	if err != nil {
		t.Fatal(err)
	}
	var last map[string]any
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(config.LedgerDir(a.workspace), "journal", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
			if err := json.Unmarshal([]byte(line), &last); err != nil {
				t.Fatal(err)
			}
		}
	}
	if !reflect.DeepEqual(last["repository"], want) {
		t.Fatalf("Event repository = %#v, want %#v", last["repository"], want)
	}
}

func p24BudgetInput(a *p22FIDRecoveryAdapter) map[string]any {
	return map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "budget_name": "tokens", "source_sample_id": "measurement", "sample_kind": "baseline", "previous_sample_id": nil, "value": 1, "measurement_ref": map[string]any{"source_type": "project-file", "display_label": "measurement", "path": "project-docs/p22-corpus.md", "section": nil, "external_id": nil}, "observed_at": time.Now().UTC().Format(time.RFC3339Nano)}
}

func TestP24RunBudgetSourceInputMatrix(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"matching", "mismatching", "omitted", "missing", "null-unborn-HEAD", "null-unmerged", "unmerged", "unborn-HEAD"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			a := p24RunAdapter(t)
			input := p24BudgetInput(a)
			source := input["measurement_ref"].(map[string]any)
			digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("corpus source\n")))
			if kind == "matching" {
				source["digest"] = digest
			}
			if kind == "mismatching" {
				source["digest"] = "sha256:" + strings.Repeat("0", 64)
			}
			if kind == "missing" {
				source["path"] = "missing.md"
			}
			if kind == "null-unborn-HEAD" || kind == "null-unmerged" {
				input["measurement_ref"] = nil
			}
			if kind == "unmerged" || kind == "null-unmerged" {
				p24MakeUnmerged(t, a.workspace)
			}
			if kind == "unborn-HEAD" || kind == "null-unborn-HEAD" {
				gittest.UnbornHead(t, a.workspace)
			}
			before, events := p22FaultProjection(t, a.workspace), p21JournalEventCount(t, a.workspace)
			response := p24RunInput(t, a, "observe-budget", "budget", input)
			switch kind {
			case "matching", "omitted", "null-unborn-HEAD", "null-unmerged":
				if response.ExitCode != 0 {
					t.Fatalf("budget: %#v", response)
				}
				sample := p21Path(t, p23Envelope(t, response.Stdout), "result", "sample").(map[string]any)
				if kind == "null-unborn-HEAD" || kind == "null-unmerged" {
					if sample["measurement_ref"] != nil {
						t.Fatal("null measurement changed")
					}
					return
				}
				stored := sample["measurement_ref"].(map[string]any)
				want := p24JudgmentExpectedRepository(t, a)
				if stored["digest"] != digest || !reflect.DeepEqual(stored["repository"], want) {
					t.Fatalf("stored Source: %#v", stored)
				}
				p24AssertRunEventRepository(t, a, want)
			default:
				code := "source-digest-mismatch"
				if kind == "unmerged" || kind == "unborn-HEAD" {
					code = "repository-unavailable"
				}
				assertP22ProcessFailure(t, "run.observe-budget", code, 4, response)
				details := p21Path(t, p23Envelope(t, response.Stderr), "error", "details").(map[string]any)
				if code == "source-digest-mismatch" {
					expected, observed := source["digest"], any(digest)
					if kind == "missing" {
						observed = nil
					}
					if !reflect.DeepEqual(details, map[string]any{"display_label": "measurement", "path": source["path"], "expected_digest": expected, "observed_digest": observed}) {
						t.Fatalf("Source mismatch details: %#v", details)
					}
				} else {
					reason := "not-a-repository"
					if kind == "unmerged" {
						reason = "unmerged-paths"
					}
					if !p23SameKeys(details, "reason detail") || details["reason"] != reason || details["detail"] == "" {
						t.Fatalf("unavailable details: %#v", details)
					}
				}
				if events != p21JournalEventCount(t, a.workspace) || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
					t.Fatal("budget rejection changed durable state")
				}
			}
		})
	}
}

func TestP24RunReplayBeforeObservationAndSourceReads(t *testing.T) {
	t.Parallel()
	for _, leaf := range []string{"start", "close", "abandon", "observe-budget"} {
		states := []string{"HEAD-and-file", "unmerged", "unborn-HEAD"}
		if leaf == "observe-budget" {
			states = append(states, "deleted-source")
		}
		for _, state := range states {
			t.Run(leaf+"/"+state, func(t *testing.T) {
				t.Parallel()
				a := p24RunAdapter(t)
				key := "replay-original"
				var input map[string]any
				var first p22ConformanceResponse
				switch leaf {
				case "start":
					data, err := os.ReadFile(filepath.Join(a.workspace, "p22-corpus-run.json"))
					if err != nil {
						t.Fatal(err)
					}
					if err := json.Unmarshal(data, &input); err != nil {
						t.Fatal(err)
					}
					key = "p22-corpus-run-start"
					// The constructor committed this exact request; capture its replay before changing the worktree.
					first = p24RunInput(t, a, leaf, key, input)
				case "close", "abandon":
					a.p22Observe(t, "terminated")
					input = p24TerminalInput(t, a, leaf)
					first = p24RunInput(t, a, leaf, key, input)
				default:
					input = p24BudgetInput(a)
					first = p24RunInput(t, a, leaf, key, input)
				}
				if first.ExitCode != 0 {
					t.Fatalf("original: %#v", first)
				}
				original := p23Envelope(t, first.Stdout)
				// Fresh controls preserve all prerequisites, so observation is the only failing condition.
				control := cloneP23Repository(input)
				if leaf == "observe-budget" {
					control["expected_version"] = a.runVersion + 1
					control["source_sample_id"] = "control-sample"
					control["sample_kind"] = "delta"
					control["previous_sample_id"] = p21Path(t, original, "result", "sample", "id")
				} else {
					if leaf == "start" {
						p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "prepare-control")
					}
					if leaf == "close" {
						released := p23RunInput(t, a, "claim", "release", a.actor, "prepare-control-release", map[string]any{"claim_id": a.claimID, "expected_version": a.claimVersion, "reason": "completed"})
						if released.ExitCode != 0 {
							t.Fatalf("release: %#v", released)
						}
					}
					claim := p23AcquireClaim(t, a, a.actor, "mutation", "control-claim")
					if leaf == "start" {
						control["claim_id"], control["claim_version"] = claim["id"], claim["version"]
						control["executor"] = map[string]any{"id": "executor-control", "source": "activated-harness"}
					} else {
						run := p23StartRun(t, a, a.actor, claim, "control-run")
						a.runID, a.runVersion, a.executorID = run["id"].(string), p23Number(t, run["version"]), p21PathString(t, run, "executor", "id")
						a.claimID, a.claimVersion = claim["id"].(string), p23Number(t, claim["version"])
						observed := p24RunInput(t, a, "observe", "control-observe", map[string]any{"run_id": a.runID, "expected_version": a.runVersion, "liveness": map[string]any{"executor_id": a.executorID, "state": "terminated", "source_ref": "process:p24-control", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
						if observed.ExitCode != 0 {
							t.Fatalf("control observation: %#v", observed)
						}
						a.p22Advance(t, "run.observe", observed.Stdout)
						control = p24TerminalInput(t, a, leaf)
					}
				}
				switch state {
				case "HEAD-and-file":
					p22Write(t, a.workspace, "base.txt", "new commit\n")
					p22Git(t, a.workspace, "commit", "-qam", "advance HEAD")
					p22Write(t, a.workspace, "project-docs/p22-corpus.md", "changed source\n")
				case "unmerged":
					p24MakeUnmerged(t, a.workspace)
				case "unborn-HEAD":
					gittest.UnbornHead(t, a.workspace)
				case "deleted-source":
					if err := os.Remove(filepath.Join(a.workspace, "project-docs/p22-corpus.md")); err != nil {
						t.Fatal(err)
					}
				}
				before, events := p22FaultProjection(t, a.workspace), p21JournalEventCount(t, a.workspace)
				retry := p24RunInput(t, a, leaf, key, input)
				if retry.ExitCode != 0 || !p24ReviewReplayMatches(original, p23Envelope(t, retry.Stdout)) {
					t.Fatalf("replay: %#v", retry)
				}
				if state != "HEAD-and-file" {
					response := p24RunInput(t, a, leaf, "new-key-control", control)
					code := "repository-unavailable"
					if leaf == "close" || leaf == "abandon" {
						code = "recovery-state-conflict"
					}
					if state == "deleted-source" {
						code = "source-digest-mismatch"
					}
					assertP22ProcessFailure(t, "run."+leaf, code, 4, response)
				}
				if events != p21JournalEventCount(t, a.workspace) || !reflect.DeepEqual(before, p22FaultProjection(t, a.workspace)) {
					t.Fatal("replay/control changed durable state")
				}
			})
		}
	}
}

func TestP24RunChangedFakeFIDInputParity(t *testing.T) {
	t.Parallel()
	// Derive the changed intersection from the shared adapters' operation corpus.
	changed := map[string]string{"run.start": "repository", "run.close": "repository", "run.abandon": "repository", "run.observe-budget": "measurement_ref.repository"}
	leaves := map[string]bool{}
	for _, fixture := range p22RecoveryFixtures(t) {
		if _, ok := changed[fixture.Operation]; ok {
			leaves[fixture.Operation] = true
		}
	}
	if len(leaves) != 3 {
		t.Fatalf("shared changed Run dispatch = %v, want close/abandon/budget", leaves)
	}
	for operation := range leaves {
		t.Run(operation, func(t *testing.T) {
			t.Parallel()
			a := p24RunAdapter(t)
			fake := p22NewFakeRecoveryAdapter()
			a.p22Observe(t, "terminated")
			fake.model.liveness, fake.model.observation = "terminated", a.observation
			fake.model.run["version"] = a.runVersion
			_, input := a.p22BuildArgs(t, p22ConformanceCase{Operation: operation, Name: "parity"})
			removed := p23Clone(input)
			if operation == "run.observe-budget" {
				removed["measurement_ref"].(map[string]any)["repository"] = a.repository
			} else {
				removed["repository"] = a.repository
			}
			before := p21JournalEventCount(t, a.workspace)
			for _, request := range []struct {
				key   string
				input map[string]any
				exit  int
			}{{"removed", removed, 2}, {"accepted", input, 0}} {
				fid := p24RunInput(t, a, strings.TrimPrefix(operation, "run."), request.key, request.input)
				model := fake.model.execute(t, operation, request.key, a.actor, request.input, false)
				if fid.ExitCode != request.exit || model.ExitCode != fid.ExitCode {
					t.Fatalf("%s parity: fid=%#v fake=%#v", request.key, fid, model)
				}
				if request.exit != 0 {
					if p21JournalEventCount(t, a.workspace) != before {
						t.Fatal("removed member appended an Event")
					}
					if p21Path(t, p23Envelope(t, fid.Stderr), "error", "details", "argument") != "input."+changed[operation] {
						t.Fatalf("removed path: %#v", fid)
					}
				}
			}
		})
	}
}
