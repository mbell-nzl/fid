//go:build integration

package integration

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

// The completeness guard binds the real-process rows to the Source positions
// independently enumerated from review record --help.
func p24ReviewSourceStates() []string {
	return []string{"matching", "mismatching", "omitted-existing", "omitted-missing"}
}

func TestP24ReviewSourceMatrixCompleteness(t *testing.T) {
	t.Parallel()
	want := []string{"matching", "mismatching", "omitted-existing", "omitted-missing"}
	if !reflect.DeepEqual(p24ReviewSourceStates(), want) {
		t.Fatalf("Review Source states = %v, want %v", p24ReviewSourceStates(), want)
	}
}

func p24ReviewReplayMatches(original, replay map[string]any) bool {
	expected := p23Clone(original)
	expected["result"].(map[string]any)["replayed"] = true
	return reflect.DeepEqual(expected, replay)
}

func TestP24ReviewReplayEvidenceRejectsChangedHistory(t *testing.T) {
	t.Parallel()
	original := map[string]any{"event": map[string]any{"id": "EV-original"}, "result": map[string]any{"replayed": false, "review": map[string]any{"id": "REV-original"}}}
	replay := p23Clone(original)
	replay["result"].(map[string]any)["replayed"] = true
	if !p24ReviewReplayMatches(original, replay) {
		t.Fatal("unchanged replay rejected")
	}
	for _, changed := range []string{"event", "result"} {
		bad := p23Clone(replay)
		if changed == "event" {
			bad["event"].(map[string]any)["id"] = "EV-wrong"
		} else {
			bad["result"].(map[string]any)["review"].(map[string]any)["id"] = "REV-wrong"
		}
		if p24ReviewReplayMatches(original, bad) {
			t.Errorf("replay checker accepted changed %s", changed)
		}
	}
}

func p24ReviewAdapter(t *testing.T) (*p22FIDRecoveryAdapter, map[string]any, map[string]any) {
	t.Helper()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "matrix-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "matrix-reviewer")
	run := p23StartRun(t, a, "actor-denied", claim, "matrix-reviewer")
	return a, claim, run
}

func p24ReviewSources(value map[string]any) []map[string]any {
	coverage := value["coverage"].(map[string]any)
	return []map[string]any{
		coverage["items"].([]any)[0].(map[string]any)["source"].(map[string]any),
		coverage["exclusions"].([]any)[0].(map[string]any)["source"].(map[string]any),
		value["evidence_refs"].([]any)[0].(map[string]any),
	}
}

func p24ReviewProjectInput(t *testing.T, a *p22FIDRecoveryAdapter, run map[string]any) map[string]any {
	t.Helper()
	input := p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "pass", nil, nil)
	source := func(path string) map[string]any {
		return map[string]any{"source_type": "project-file", "display_label": path, "path": path, "section": nil, "external_id": nil}
	}
	input["coverage"] = map[string]any{"items": []any{map[string]any{"source": source("base.txt"), "assessed": true}},
		"exclusions": []any{map[string]any{"source": source(".gitignore"), "reason": "outside review scope"}}}
	input["evidence_refs"] = []any{source("base.txt")}
	return input
}

func p24ReviewAssertFailure(t *testing.T, a *p22FIDRecoveryAdapter, response p22ConformanceResponse, code string, details map[string]any, before int) {
	t.Helper()
	assertP22ProcessFailure(t, "review.record", code, 4, response)
	envelope := p23Envelope(t, response.Stderr)
	failure := envelope["error"].(map[string]any)
	if !reflect.DeepEqual(failure["details"], details) || failure["subject"] != nil {
		t.Fatalf("failure = %#v, want closed details %#v and null subject", failure, details)
	}
	if _, present := envelope["event"]; present || p21JournalEventCount(t, a.workspace) != before {
		t.Fatal("refusal emitted or appended an Event")
	}
}

func TestP24ReviewSourceMatrix(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	stdout, stderr, err := closureRunCommand(t.TempDir(), binary, "review", "record", "--help")
	if err != nil || stderr != "" {
		t.Fatalf("descriptor help: %v / %s", err, stderr)
	}
	positions := []string{}
	for _, line := range strings.Split(stdout, "\n") {
		path := strings.Split(strings.TrimSpace(line), ": ")[0]
		if strings.HasSuffix(path, ".source_type") {
			positions = append(positions, strings.TrimSuffix(path, ".source_type"))
		}
	}
	want := []string{"input.coverage.items[i].source", "input.coverage.exclusions[i].source", "input.evidence_refs[i]"}
	if !reflect.DeepEqual(positions, want) {
		t.Fatalf("descriptor Source positions = %v, want %v", positions, want)
	}
	for position, path := range positions {
		for _, state := range p24ReviewSourceStates() {
			t.Run(path+"/"+state, func(t *testing.T) {
				t.Parallel()
				a, _, run := p24ReviewAdapter(t)
				input := p24ReviewProjectInput(t, a, run)
				source := p24ReviewSources(input)[position]
				content, err := os.ReadFile(filepath.Join(a.workspace, source["path"].(string)))
				if err != nil {
					t.Fatal(err)
				}
				digest := fmt.Sprintf("sha256:%x", sha256.Sum256(content))
				wrong := "sha256:" + p22Zeros(64)
				if state == "matching" {
					source["digest"] = digest
				} else if state == "mismatching" {
					source["digest"] = wrong
				} else if state == "omitted-missing" {
					source["path"] = "missing.md"
				}
				observation := p24JudgmentExpectedRepository(t, a)
				before := p21JournalEventCount(t, a.workspace)
				response := p23ReviewRecord(t, a, "actor-denied", "source-matrix", input)
				if state == "mismatching" || state == "omitted-missing" {
					var expected, observed any
					if state == "mismatching" {
						expected, observed = wrong, digest
					}
					p24ReviewAssertFailure(t, a, response, "source-digest-mismatch", map[string]any{
						"display_label": source["display_label"], "path": source["path"], "expected_digest": expected, "observed_digest": observed}, before)
					return
				}
				p23AssertReviewSuccess(t, "Source matrix", response)
				if got := p24ReviewSources(p23ReviewFromResponse(t, response))[position]["digest"]; got != digest {
					t.Fatalf("stored %s digest = %v, want independently hashed %s", path, got, digest)
				}
				p24ReviewAssertRepositories(t, a, response, observation)
				if p21JournalEventCount(t, a.workspace) != before+1 {
					t.Fatal("accepted matrix row did not append one Event")
				}
			})
		}
	}
}

func p24ReviewAssertRepositories(t *testing.T, a *p22FIDRecoveryAdapter, response p22ConformanceResponse, observation map[string]any) {
	t.Helper()
	review := p23ReviewFromResponse(t, response)
	for i, repository := range []any{review["repository"], p24ReviewSources(review)[0]["repository"], p24ReviewSources(review)[1]["repository"], p24ReviewSources(review)[2]["repository"]} {
		if !reflect.DeepEqual(repository, observation) {
			t.Fatalf("stored repository %d = %#v, want independent observation %#v", i, repository, observation)
		}
	}
	event := p24EventByID(t, a.workspace, p21PathString(t, p23Envelope(t, response.Stdout), "event", "id"))
	if !reflect.DeepEqual(p24JSONMap(t, event.Repository), observation) {
		t.Fatalf("Event repository = %#v, want %#v", event.Repository, observation)
	}
}

// This fixture has exactly one untracked regular file. Git proves that fact;
// the test hashes its known image and the canonical dirty entries independently.
func p24ReviewDirtyRepository(t *testing.T, a *p22FIDRecoveryAdapter, path string, content []byte) map[string]any {
	t.Helper()
	if status := p22GitOutput(t, a.workspace, "status", "--porcelain", "--untracked-files=all", "--", ".", ":!.fid"); status != "?? "+path {
		t.Fatalf("dirty fixture status = %q, want one untracked %s", status, path)
	}
	root, err := filepath.Abs(a.workspace)
	if err != nil {
		t.Fatal(err)
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(config.LedgerDir(a.workspace), "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	head := p22GitOutput(t, a.workspace, "rev-parse", "HEAD")
	entries := []any{map[string]any{"path": path, "head": nil, "index": nil, "worktree": map[string]any{"mode": "100644", "digest": fmt.Sprintf("sha256:%x", sha256.Sum256(content))}}}
	payload, err := json.Marshal(map[string]any{"schema": "fid.repository-dirty.v1", "entries": entries})
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{"base_oid": head, "capture_source": "git", "dirty": true, "dirty_digest": fmt.Sprintf("sha256:%x", sha256.Sum256(payload)),
		"head_oid": head, "project_id": projectConfig.ProjectID, "tree_oid": p22GitOutput(t, a.workspace, "rev-parse", "HEAD^{tree}"),
		"worktree_id": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root)))}
}

func TestP24ReviewReplayUnavailableAndRestoredControls(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"head-advanced", "not-a-repository", "unmerged-paths", "deleted-source"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			a, claim, run := p24ReviewAdapter(t)
			input := p24ReviewProjectInput(t, a, run)
			// This project-owned file is ignored by Git. Its absence isolates
			// Source observation from the repository expectation check.
			p22Write(t, a.workspace, ".fid/coverage.md", "review source\n")
			for _, position := range []int{0, 2} {
				p24ReviewSources(input)[position]["path"] = ".fid/coverage.md"
			}
			observation := p24JudgmentExpectedRepository(t, a)
			original := p23ReviewRecord(t, a, "actor-denied", "before-unavailable", input)
			p23AssertReviewSuccess(t, "normal Review", original)
			p24ReviewAssertRepositories(t, a, original, observation)
			eventID := p21PathString(t, p23Envelope(t, original.Stdout), "event", "id")
			originalEvent := p24EventByID(t, a.workspace, eventID)
			p23CloseAndRelease(t, a, "actor-denied", run["id"].(string), claim["id"].(string), "before-control")
			controlClaim := p23AcquireClaim(t, a, "actor-denied", "review", "fresh-control")
			controlRun := p23StartRun(t, a, "actor-denied", controlClaim, "fresh-control")
			control := p23Clone(input)
			control["run_id"], control["expected_run_version"] = controlRun["id"], controlRun["version"]
			gitBackup := filepath.Join(t.TempDir(), "git-backup")
			switch state {
			case "head-advanced":
				p22Write(t, a.workspace, "base.txt", "advanced control HEAD\n")
				p22Git(t, a.workspace, "commit", "-qam", "advance control HEAD")
			case "not-a-repository":
				if err := os.Rename(filepath.Join(a.workspace, ".git"), gitBackup); err != nil {
					t.Fatal(err)
				}
			case "unmerged-paths":
				p24MakeUnmerged(t, a.workspace)
			case "deleted-source":
				if err := os.Remove(filepath.Join(config.LedgerDir(a.workspace), "coverage.md")); err != nil {
					t.Fatal(err)
				}
			}
			before := p21JournalEventCount(t, a.workspace)
			replay := p23ReviewRecord(t, a, "actor-denied", "before-unavailable", input)
			p23AssertReviewSuccess(t, "unavailable replay", replay)
			if !p24ReviewReplayMatches(p23Envelope(t, original.Stdout), p23Envelope(t, replay.Stdout)) ||
				!reflect.DeepEqual(originalEvent, p24EventByID(t, a.workspace, eventID)) || p21JournalEventCount(t, a.workspace) != before {
				t.Fatal("retry changed original result/Event or appended an Event")
			}
			refused := p23ReviewRecord(t, a, "actor-denied", "fresh-control-key", control)
			if state == "head-advanced" {
				p24ReviewAssertFailure(t, a, refused, "repository-mismatch", map[string]any{
					"expected": control["expected_repository"], "observed": p24JudgmentExpectedRepository(t, a)}, before)
			} else if state == "deleted-source" {
				source := p24ReviewSources(control)[0]
				p24ReviewAssertFailure(t, a, refused, "source-digest-mismatch", map[string]any{"display_label": source["display_label"],
					"path": source["path"], "expected_digest": nil, "observed_digest": nil}, before)
			} else {
				detail := "repository has unmerged paths"
				if state == "not-a-repository" {
					detail = "repository root: exit status 128"
				}
				p24ReviewAssertFailure(t, a, refused, "repository-unavailable", map[string]any{"reason": state, "detail": detail}, before)
				failure := p23Envelope(t, refused.Stderr)["error"].(map[string]any)
				if failure["category"] != "precondition" || failure["retryable"] != true {
					t.Fatalf("unavailable metadata = %#v", failure)
				}
			}
			switch state {
			case "head-advanced":
				p22Git(t, a.workspace, "reset", "--hard", observation["head_oid"].(string))
			case "not-a-repository":
				if err := os.Rename(gitBackup, filepath.Join(a.workspace, ".git")); err != nil {
					t.Fatal(err)
				}
			case "unmerged-paths":
				p22Git(t, a.workspace, "merge", "--abort")
				p22Git(t, a.workspace, "reset", "--hard", observation["head_oid"].(string))
			case "deleted-source":
				p22Write(t, a.workspace, ".fid/coverage.md", "review source\n")
			}
			// Exact same fresh Run, key and request now succeeds. No stale
			// expectation or duplicate Review can masquerade as a control.
			restored := p23ReviewRecord(t, a, "actor-denied", "fresh-control-key", control)
			p23AssertReviewSuccess(t, "restored observable control", restored)
			p24ReviewAssertRepositories(t, a, restored, p24JudgmentExpectedRepository(t, a))
			if p21JournalEventCount(t, a.workspace) != before+1 {
				t.Fatal("restored control did not append exactly one Event")
			}
		})
	}
}

func TestP24ReviewRecordResolvesSourcesAndReplaysWithoutObservation(t *testing.T) {
	t.Parallel()
	actions := []string{"claim.acquire", "claim.release", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.close", "run.observe", "run.start"}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "review-source-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "review-source")
	run := p23StartRun(t, a, "actor-denied", claim, "review-source")
	p22Write(t, a.workspace, "excluded.md", "excluded material\n")
	input := p23ReviewInput(a, p21PathString(t, run, "id"), p23Number(t, run["version"]), "pass", nil, nil)
	head := p22GitOutput(t, a.workspace, "rev-parse", "HEAD")
	dirty := p22GitOutput(t, a.workspace, "status", "--porcelain") != ""
	input["expected_repository"] = map[string]any{"head_oid": head, "dirty": dirty}
	basePath, exclusionPath := "base.txt", "excluded.md"
	baseBytes, err := os.ReadFile(filepath.Join(a.workspace, basePath))
	if err != nil {
		t.Fatal(err)
	}
	baseDigest := fmt.Sprintf("sha256:%x", sha256.Sum256(baseBytes))
	excludedBytes, err := os.ReadFile(filepath.Join(a.workspace, exclusionPath))
	if err != nil {
		t.Fatal(err)
	}
	excludedDigest := fmt.Sprintf("sha256:%x", sha256.Sum256(excludedBytes))
	baseSource := map[string]any{"source_type": "project-file", "display_label": "Covered", "path": basePath, "section": nil, "external_id": nil}
	exclusionSource := map[string]any{"source_type": "project-file", "display_label": "Excluded", "path": exclusionPath, "section": nil, "external_id": nil, "digest": excludedDigest}
	input["coverage"] = map[string]any{"items": []any{map[string]any{"source": baseSource, "assessed": true}}, "exclusions": []any{map[string]any{"source": exclusionSource, "reason": "outside review scope"}}}
	input["evidence_refs"] = []any{baseSource}
	for _, test := range []struct {
		name, path string
		change     func(map[string]any)
	}{
		{"top repository", "input.repository", func(v map[string]any) { v["repository"] = a.repository }},
		{"item repository", "input.coverage.items[0].source.repository", func(v map[string]any) {
			v["coverage"].(map[string]any)["items"].([]any)[0].(map[string]any)["source"].(map[string]any)["repository"] = a.repository
		}},
		{"exclusion repository", "input.coverage.exclusions[0].source.repository", func(v map[string]any) {
			v["coverage"].(map[string]any)["exclusions"].([]any)[0].(map[string]any)["source"].(map[string]any)["repository"] = a.repository
		}},
		{"evidence repository", "input.evidence_refs[0].repository", func(v map[string]any) {
			v["evidence_refs"].([]any)[0].(map[string]any)["repository"] = a.repository
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			removed := p23Clone(input)
			test.change(removed)
			before := p21JournalEventCount(t, a.workspace)
			result := p23ReviewRecord(t, a, "actor-denied", "review-removed-"+test.name, removed)
			assertP22ProcessFailure(t, "review.record", "invalid-argument", 2, result)
			var failure map[string]any
			if err := json.Unmarshal([]byte(result.Stderr), &failure); err != nil {
				t.Fatal(err)
			}
			if got := failure["error"].(map[string]any)["details"].(map[string]any)["argument"]; got != test.path || p21JournalEventCount(t, a.workspace) != before {
				t.Fatalf("removed repository member = %v, want %s and no Event", got, test.path)
			}
		})
	}
	wrongDigest := "sha256:" + p22Zeros(64)
	for _, test := range []struct {
		name, observed string
		change         func(map[string]any)
	}{
		{"item", baseDigest, func(v map[string]any) {
			v["coverage"].(map[string]any)["items"].([]any)[0].(map[string]any)["source"].(map[string]any)["digest"] = wrongDigest
		}},
		{"exclusion", excludedDigest, func(v map[string]any) {
			v["coverage"].(map[string]any)["exclusions"].([]any)[0].(map[string]any)["source"].(map[string]any)["digest"] = wrongDigest
		}},
		{"evidence", baseDigest, func(v map[string]any) { v["evidence_refs"].([]any)[0].(map[string]any)["digest"] = wrongDigest }},
	} {
		t.Run(test.name+" digest mismatch", func(t *testing.T) {
			wrong := p23Clone(input)
			test.change(wrong)
			before := p21JournalEventCount(t, a.workspace)
			response := p23ReviewRecord(t, a, "actor-denied", "review-source-wrong-"+test.name, wrong)
			assertP22ProcessFailure(t, "review.record", "source-digest-mismatch", 4, response)
			var failure map[string]any
			if err := json.Unmarshal([]byte(response.Stderr), &failure); err != nil {
				t.Fatal(err)
			}
			details := failure["error"].(map[string]any)["details"].(map[string]any)
			if details["expected_digest"] != wrongDigest || details["observed_digest"] != test.observed || p21JournalEventCount(t, a.workspace) != before {
				t.Fatalf("%s mismatch details = %#v; want both digests and no Event", test.name, details)
			}
		})
	}
	before := p21JournalEventCount(t, a.workspace)
	response := p23ReviewRecord(t, a, "actor-denied", "review-source-record", input)
	p23AssertReviewSuccess(t, "Review with resolved Sources", response)
	review := p23ReviewFromResponse(t, response)
	storedItem := review["coverage"].(map[string]any)["items"].([]any)[0].(map[string]any)["source"].(map[string]any)
	storedExclusion := review["coverage"].(map[string]any)["exclusions"].([]any)[0].(map[string]any)["source"].(map[string]any)
	storedEvidence := review["evidence_refs"].([]any)[0].(map[string]any)
	for _, check := range []struct {
		source map[string]any
		digest string
	}{{storedItem, baseDigest}, {storedExclusion, excludedDigest}, {storedEvidence, baseDigest}} {
		if check.source["digest"] != check.digest || !reflect.DeepEqual(check.source["repository"], review["repository"]) {
			t.Fatalf("resolved Source = %#v, want digest %s and observed Repository", check.source, check.digest)
		}
	}
	if repository := review["repository"].(map[string]any); repository["head_oid"] != head || repository["dirty"] != dirty || repository["tree_oid"] != p22GitOutput(t, a.workspace, "rev-parse", "HEAD^{tree}") {
		t.Fatalf("stored repository differs from git observation: %#v", repository)
	}
	p24ReviewAssertRepositories(t, a, response, p24ReviewDirtyRepository(t, a, exclusionPath, excludedBytes))
	if p21JournalEventCount(t, a.workspace) != before+1 {
		t.Fatal("accepted Review did not append exactly one Event")
	}
	p22Write(t, a.workspace, "base.txt", "advanced Review source\n")
	p22Git(t, a.workspace, "add", "base.txt", "excluded.md")
	p22Git(t, a.workspace, "commit", "-qm", "advance Review HEAD")
	before = p21JournalEventCount(t, a.workspace)
	replay := p23ReviewRecord(t, a, "actor-denied", "review-source-record", input)
	p23AssertReviewSuccess(t, "HEAD-advanced replay", replay)
	if !p24ReviewReplayMatches(p23Envelope(t, response.Stdout), p23Envelope(t, replay.Stdout)) || p21JournalEventCount(t, a.workspace) != before {
		t.Fatal("HEAD-advanced retry did not replay without Event")
	}
	newKey := p23ReviewRecord(t, a, "actor-denied", "review-head-control", input)
	assertP22ProcessFailure(t, "review.record", "repository-mismatch", 4, newKey)
	gitBackup := filepath.Join(t.TempDir(), "git-backup")
	if err := os.Rename(filepath.Join(a.workspace, ".git"), gitBackup); err != nil {
		t.Fatal(err)
	}
	replay = p23ReviewRecord(t, a, "actor-denied", "review-source-record", input)
	p23AssertReviewSuccess(t, "unobservable replay", replay)
	if !p24ReviewReplayMatches(p23Envelope(t, response.Stdout), p23Envelope(t, replay.Stdout)) || p21JournalEventCount(t, a.workspace) != before {
		t.Fatal("unobservable retry did not replay without Event")
	}
	unavailable := p23ReviewRecord(t, a, "actor-denied", "review-unobservable-control", input)
	assertP22ProcessFailure(t, "review.record", "repository-unavailable", 4, unavailable)
	if err := os.Rename(gitBackup, filepath.Join(a.workspace, ".git")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(a.workspace, basePath)); err != nil {
		t.Fatal(err)
	}
	replay = p23ReviewRecord(t, a, "actor-denied", "review-source-record", input)
	p23AssertReviewSuccess(t, "deleted Source replay", replay)
	if !p24ReviewReplayMatches(p23Envelope(t, response.Stdout), p23Envelope(t, replay.Stdout)) || p21JournalEventCount(t, a.workspace) != before {
		t.Fatal("deleted Source retry did not replay without Event")
	}
	deleted := p23Clone(input)
	deleted["expected_repository"] = map[string]any{"head_oid": p22GitOutput(t, a.workspace, "rev-parse", "HEAD"), "dirty": true}
	missing := p23ReviewRecord(t, a, "actor-denied", "review-deleted-control", deleted)
	assertP22ProcessFailure(t, "review.record", "source-digest-mismatch", 4, missing)
	var missingFailure map[string]any
	if err := json.Unmarshal([]byte(missing.Stderr), &missingFailure); err != nil {
		t.Fatal(err)
	}
	missingDetails := missingFailure["error"].(map[string]any)["details"].(map[string]any)
	if missingDetails["expected_digest"] != nil || missingDetails["observed_digest"] != nil {
		t.Fatalf("deleted omitted Source details = %#v, want both digests null", missingDetails)
	}
	if p21JournalEventCount(t, a.workspace) != before {
		t.Fatal("replay controls appended Event")
	}
}

// TestP23ReviewRecordRealProcessJourney keeps the VAL-015 real-process
// evidence at the built-binary boundary. The Run-only history direction is
// deliberately absent: DEC-239 confines that otherwise unreachable state to
// internal/review's isolated durable-history query fixture.
func TestP23ReviewRecordRealProcessJourney(t *testing.T) {
	t.Parallel()
	actions := []string{
		"claim.acquire", "claim.release", "envelope.activate", "envelope.create",
		"evidence.link", "gate.register", "review.record", "run.close", "run.observe", "run.start",
	}
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{}, actions, "mutation")
	a.task13ReplayCorpus = true

	t.Run("mutating_claim_actor_is_rejected_without_an_event", func(t *testing.T) {
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, a.actor, "builder-conflict", p23ReviewInput(a, a.runID, a.runVersion, "changes-required", []any{p23Finding(a.issueID, "FND-builder", "major", true)}, nil))
		assertP22ProcessFailure(t, "review.record", "review-not-independent", 4, response)
		for _, want := range []string{`"type":"claim"`, `"id":"` + a.claimID + `"`, `"conflicting_actor_id":"` + a.actor + `"`} {
			if !strings.Contains(response.Stderr, want) {
				t.Fatalf("review independence rejection = %q, want %s", response.Stderr, want)
			}
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("mutating Claim rejection appended Event(s): %d -> %d", before, after)
		}
	})

	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "builder-before-review")

	firstClaim := p23AcquireClaim(t, a, "actor-denied", "review", "reviewer-first")
	firstRun := p23StartRun(t, a, "actor-denied", firstClaim, "reviewer-first")
	firstInput := p23ReviewInput(a, p21PathString(t, firstRun, "id"), p23Number(t, firstRun["version"]), "changes-required", []any{
		p23Finding(a.issueID, "FND-prior", "major", true),
		p23Finding(a.issueID, "FND-prior-two", "major", true),
	}, nil)
	t.Run("dirty_expectation_mismatch_echoes_only_the_expected_pair", func(t *testing.T) {
		input := p23Clone(firstInput)
		delete(input, "repository")
		input["expected_repository"] = map[string]any{"head_oid": a.repository["head_oid"], "dirty": !a.repository["dirty"].(bool)}
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, "actor-denied", "dirty-mismatch", input)
		assertP22ProcessFailure(t, "review.record", "repository-mismatch", 4, response)
		var envelope map[string]any
		if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
			t.Fatal(err)
		}
		details := envelope["error"].(map[string]any)["details"].(map[string]any)
		if !reflect.DeepEqual(details["expected"], input["expected_repository"]) {
			t.Fatalf("expected echo = %#v, want %#v", details["expected"], input["expected_repository"])
		}
		if !reflect.DeepEqual(details["observed"], p24JSONMap(t, a.repository)) {
			t.Fatalf("observed repository = %#v, want %#v", details["observed"], a.repository)
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("dirty mismatch appended Event: %d -> %d", before, after)
		}
	})

	t.Run("repository_mismatch_is_typed_and_non_mutating", func(t *testing.T) {
		mismatch := p23Clone(firstInput)
		repository := p23Clone(mismatch["expected_repository"].(map[string]any))
		repository["head_oid"] = strings.Repeat("0", 40)
		mismatch["expected_repository"] = repository
		before := p21JournalEventCount(t, a.workspace)
		response := p23ReviewRecord(t, a, "actor-denied", "repository-mismatch", mismatch)
		assertP22ProcessFailure(t, "review.record", "repository-mismatch", 4, response)
		var envelope map[string]any
		if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
			t.Fatal(err)
		}
		details := envelope["error"].(map[string]any)["details"].(map[string]any)
		if len(details) != 2 || !reflect.DeepEqual(details["expected"], repository) || !reflect.DeepEqual(details["observed"], p24JSONMap(t, a.repository)) {
			t.Fatalf("repository mismatch details = %#v, want closed expected/observed shape", details)
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("repository mismatch appended Event(s): %d -> %d", before, after)
		}
	})

	beforeFirst := p21JournalEventCount(t, a.workspace)
	first := p23ReviewRecord(t, a, "actor-denied", "reviewer-first", firstInput)
	p23AssertReviewSuccess(t, "first independent Review", first)
	activity := p22FaultDB(t, a.workspace)
	var reviewActivityCount int
	p22FaultOK(t, activity.QueryRow(`SELECT count(*) FROM activity_log WHERE event='review-recorded' AND journal_event_id=?`,
		p21PathString(t, p23Envelope(t, first.Stdout), "event", "id")).Scan(&reviewActivityCount))
	p22FaultOK(t, activity.Close())
	if reviewActivityCount != 1 {
		t.Fatalf("review.record activity rows = %d, want one review-recorded row from its Event", reviewActivityCount)
	}
	if after := p21JournalEventCount(t, a.workspace); after != beforeFirst+1 {
		t.Fatalf("first Review Event count = %d, want %d", after, beforeFirst+1)
	}
	firstReview := p23ReviewFromResponse(t, first)
	humanReview, humanErr := closureRun(t, a.workspace, a.binary, "--project", a.workspace, "review", "get", firstReview["id"].(string))
	if humanErr != "" || !strings.Contains(humanReview, "Purpose: Inspect independent Review") || !strings.Contains(humanReview, "Independence: distinct-actor") || !strings.Contains(humanReview, "Findings by severity:") {
		t.Fatalf("Review Human projection = %q / %q", humanReview, humanErr)
	}
	resumeAfterFirst := p21Run(t, a.binary, a.workspace, "--json", "resume", a.issueID, "--repository", a.workspace)
	firstProjection := p21Path(t, resumeAfterFirst, "result", "resume").(map[string]any)
	if got := firstProjection["reviews"].([]any); len(got) != 1 || !reflect.DeepEqual(got[0], firstReview) {
		t.Fatalf("resume Reviews = %#v, want exact recorded Review %#v", got, firstReview)
	}
	if got := firstProjection["findings"].([]any); len(got) != 2 {
		t.Fatalf("resume open Findings = %#v, want both blocking Findings", got)
	}
	missingReview := a.p22Run(t, "review", "get", "REV-missing")
	assertP22ProcessFailure(t, "review.get", "review-not-found", 3, missingReview)

	t.Run("replay_duplicate_and_digest_conflict_preserve_history", func(t *testing.T) {
		before := p21JournalEventCount(t, a.workspace)
		replay := p23ReviewRecord(t, a, "actor-denied", "reviewer-first", firstInput)
		p23AssertReviewSuccess(t, "review replay", replay)
		if p21Path(t, p23Envelope(t, replay.Stdout), "result", "replayed") != true {
			t.Fatalf("review replay = %q, want replayed:true", replay.Stdout)
		}
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("review replay appended Event(s): %d -> %d", before, after)
		}

		conflictInput := p23Clone(firstInput)
		conflictInput["profile"] = "review-alternate"
		response := p23ReviewRecord(t, a, "actor-denied", "reviewer-first", conflictInput)
		assertP22ProcessFailure(t, "review.record", "idempotency-conflict", 4, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("idempotency conflict appended Event(s): %d -> %d", before, after)
		}

		response = p23ReviewRecord(t, a, "actor-denied", "same-run", firstInput)
		assertP22ProcessFailure(t, "review.record", "duplicate-entity", 4, response)
		if after := p21JournalEventCount(t, a.workspace); after != before {
			t.Fatalf("second Review for Run appended Event(s): %d -> %d", before, after)
		}
	})

	p23CloseAndRelease(t, a, "actor-denied", p21PathString(t, firstRun, "id"), p21PathString(t, firstClaim, "id"), "reviewer-first")
	reworkClaim := p23AcquireClaim(t, a, a.actor, "mutation", "builder-rework")
	reworkRun := p23StartRun(t, a, a.actor, reworkClaim, "builder-rework")
	p23CloseAndRelease(t, a, a.actor, p21PathString(t, reworkRun, "id"), p21PathString(t, reworkClaim, "id"), "builder-rework")

	secondClaim := p23AcquireClaim(t, a, "actor-denied", "review", "reviewer-second")
	secondRun := p23StartRun(t, a, "actor-denied", secondClaim, "reviewer-second")
	beforeSecond := p21JournalEventCount(t, a.workspace)
	second := p23ReviewRecord(t, a, "actor-denied", "reviewer-second", p23ReviewInput(a, p21PathString(t, secondRun, "id"), p23Number(t, secondRun["version"]), "pass", []any{}, []any{"FND-prior", "FND-prior-two"}))
	p23AssertReviewSuccess(t, "same Reviewer post-rework Review", second)
	p23AssertActivity(t, a.workspace, second, "review-recorded", a.issueID, 1)
	p23AssertActivity(t, a.workspace, second, "finding-resolved", a.issueID, 2)
	if after := p21JournalEventCount(t, a.workspace); after != beforeSecond+1 {
		t.Fatalf("second Review plus two resolutions changed Event count %d -> %d, want one atomic Event", beforeSecond, after)
	}
	secondReview := p23ReviewFromResponse(t, second)
	beforeTypeMismatch := p21JournalEventCount(t, a.workspace)
	wrongEvidenceType := p23RunInput(t, a, "evidence", "link", "actor-recovery", "review-on-gate-obligation", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "review", "receipt_id": nil, "review_id": secondReview["id"], "human_outcome_id": nil, "required_artifact_ids": []string{}})
	assertP22ProcessFailure(t, "evidence.link", "evidence-type-not-accepted", 4, wrongEvidenceType)
	if p21JournalEventCount(t, a.workspace) != beforeTypeMismatch {
		t.Fatal("wrong Evidence type appended an Event")
	}
	listed := a.p22Run(t, "review", "list", "--issue", a.issueID)
	p23AssertReviewSuccess(t, "review list by Issue", listed)
	items := p21Path(t, p23Envelope(t, listed.Stdout), "result", "items").([]any)
	if len(items) != 2 || items[0].(map[string]any)["id"] != firstReview["id"] && items[1].(map[string]any)["id"] != firstReview["id"] {
		t.Fatalf("review list = %#v, want both recorded Reviews", items)
	}
	bounded := a.p22Run(t, "review", "list", "--issue", a.issueID, "--limit", "1")
	assertP22ProcessFailure(t, "review.list", "result-limit-exceeded", 6, bounded)

	firstRead := p23ReviewGet(t, a, firstReview["id"].(string))
	secondRead := p23ReviewGet(t, a, secondReview["id"].(string))
	if firstRead["id"] != firstReview["id"] || firstRead["run_id"] != firstReview["run_id"] || firstRead["verdict"] != "changes-required" || !reflect.DeepEqual(firstRead["repository"], firstReview["repository"]) {
		t.Fatalf("first review readback = %#v, want immutable recorded binding %#v", firstRead, firstReview)
	}
	findings, ok := firstRead["findings"].([]any)
	if !ok || len(findings) != 2 {
		t.Fatalf("first review findings = %#v, want two projected Findings", firstRead["findings"])
	}
	for _, value := range findings {
		finding := value.(map[string]any)
		if finding["state"] != "resolved" || len(finding["dispositions"].([]any)) != 1 {
			t.Fatalf("resolved Finding = %#v, want one resolved disposition", finding)
		}
	}
	if !reflect.DeepEqual(secondRead, secondReview) || !reflect.DeepEqual(secondRead["resolved_finding_ids"], []any{"FND-prior", "FND-prior-two"}) {
		t.Fatalf("second review readback = %#v, want exact record result %#v with both resolved Finding IDs", secondRead, secondReview)
	}
}

func p23AssertActivity(t *testing.T, workspace string, response p22ConformanceResponse, name, issueID string, want int) {
	t.Helper()
	p23AssertActivitySubject(t, workspace, response, name, "issue", issueID, want)
}

func p23AssertActivitySubject(t *testing.T, workspace string, response p22ConformanceResponse, name, refType, refID string, want int) {
	t.Helper()
	conn := p22FaultDB(t, workspace)
	defer func() { p22FaultOK(t, conn.Close()) }()
	eventID := p21PathString(t, p23Envelope(t, response.Stdout), "event", "id")
	var total, matching int
	p22FaultOK(t, conn.QueryRow(`SELECT COUNT(*), COALESCE(SUM(CASE WHEN ref_type=? AND ref_id=? THEN 1 ELSE 0 END),0) FROM activity_log WHERE journal_event_id=? AND event=?`, refType, refID, eventID, name).Scan(&total, &matching))
	if total != want || matching != want {
		t.Fatalf("%s activity for %s = %d rows, %d with %s %s; want %d", name, eventID, total, matching, refType, refID, want)
	}
	rows, err := conn.Query(`SELECT event, ref_type FROM activity_log WHERE journal_event_id=?`, eventID)
	p22FaultOK(t, err)
	defer func() { p22FaultOK(t, rows.Close()) }()
	accepted := map[string]bool{
		"review-recorded": true, "finding-resolved": true, "finding-addressed": true,
		"human-outcome-recorded": true, "concern-recorded": true, "concern-disposed": true,
		"deviation-recorded": true, "requirement-referenced": true, "requirement-linked": true,
		"decision-referenced": true,
	}
	for rows.Next() {
		var event, refType string
		p22FaultOK(t, rows.Scan(&event, &refType))
		if !accepted[event] {
			t.Fatalf("Event %s wrote unaccepted ActivityEvent name %q", eventID, event)
		}
		if refType != "issue" && refType != "requirement" && refType != "decision" {
			t.Fatalf("Event %s wrote invalid ActivityEvent %q on %q", eventID, event, refType)
		}
	}
	p22FaultOK(t, rows.Err())
}

func requiresAdditionalJudgmentActors(actions []string) bool {
	for _, action := range actions {
		if action == "review.record" || action == "concern.dispose" {
			return true
		}
	}
	return false
}

func containsReviewRecordAction(actions []string) bool {
	return containsAction(actions, "review.record")
}

func containsAction(actions []string, want string) bool {
	for _, action := range actions {
		if action == want {
			return true
		}
	}
	return false
}

func p23ReviewInput(a *p22FIDRecoveryAdapter, runID string, version int, verdict string, findings, resolves []any) map[string]any {
	if findings == nil {
		findings = []any{}
	}
	if resolves == nil {
		resolves = []any{}
	}
	source := p23Clone(a.measurement)
	delete(source, "repository")
	source["source_type"] = "external-reference"
	source["path"] = nil
	source["external_id"] = "p22-corpus-coverage"
	return map[string]any{"run_id": runID, "expected_run_version": version, "profile": "review-default", "expected_repository": map[string]any{"head_oid": a.repository["head_oid"], "dirty": a.repository["dirty"]}, "coverage": map[string]any{"items": []any{map[string]any{"source": source, "assessed": true}}, "exclusions": []any{}}, "verdict": verdict, "findings": findings, "resolves_finding_ids": resolves, "evidence_refs": []any{}, "receipt_ids": []any{}}
}

func p23Finding(issueID, id, severity string, blocking bool) map[string]any {
	return map[string]any{"id": id, "subject": map[string]any{"type": "issue", "id": issueID}, "scope": map[string]any{"paths": []any{}}, "severity": severity, "summary": "real-process Review finding", "blocking": blocking}
}

func p23ReviewRecord(t *testing.T, a *p22FIDRecoveryAdapter, actor, key string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	return p23RunInput(t, a, "review", "record", actor, key, input)
}

func p23RunInput(t *testing.T, a *p22FIDRecoveryAdapter, group, command, actor, key string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	args := []string{"--project", a.workspace, "--json", group, command, "--actor", actor, "--input", p21Input(t, a.workspace, "p23-"+key, input), "--idempotency-key", key}
	stdout, stderr, err := closureRunCommand(a.workspace, a.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	if a.task13ReplayCorpus && response.ExitCode == 0 {
		before := p21JournalEventCount(t, a.workspace)
		replayedStdout, replayedStderr, replayErr := closureRunCommand(a.workspace, a.binary, args...)
		if replayErr != nil || replayedStderr != "" || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("%s.%s same-key replay changed state or failed: error=%v stderr=%q", group, command, replayErr, replayedStderr)
		}
		original := p23Envelope(t, response.Stdout)
		replayed := p23Envelope(t, replayedStdout)
		originalResult := original["result"].(map[string]any)
		originalResult["replayed"] = true
		if !reflect.DeepEqual(original, replayed) {
			t.Fatalf("%s.%s replay changed typed snapshot: first=%#v replay=%#v", group, command, original, replayed)
		}
		if conflictInput := p23ConflictInput(group+"."+command, input); conflictInput != nil {
			conflictArgs := []string{"--project", a.workspace, "--json", group, command, "--actor", actor, "--input", p21Input(t, a.workspace, "p23-"+key+"-digest-conflict", conflictInput), "--idempotency-key", key}
			conflictStdout, conflictStderr, conflictErr := closureRunCommand(a.workspace, a.binary, conflictArgs...)
			conflict := p22ConformanceResponse{Stdout: conflictStdout, Stderr: conflictStderr}
			if conflictErr != nil {
				conflict.ExitCode = processExitCode(t, conflictErr)
			}
			assertP22ProcessFailure(t, group+"."+command, "idempotency-conflict", 4, conflict)
			if p21JournalEventCount(t, a.workspace) != before {
				t.Fatalf("%s.%s digest conflict appended an Event", group, command)
			}
		}
	}
	return response
}

func p23ConflictInput(operation string, input map[string]any) map[string]any {
	mutated := p23Clone(input)
	switch operation {
	case "review.record":
		mutated["profile"] = "review-alternate"
	case "finding.address", "concern.dispose":
		mutated["reason"] = input["reason"].(string) + " (different)"
	case "concern.record":
		mutated["summary"] = input["summary"].(string) + " (different)"
	case "deviation.record", "human-outcome.record":
		mutated["reasoning"] = input["reasoning"].(string) + " (different)"
	case "decision.record-reference":
		if input["outcome"] == "approved" {
			mutated["outcome"] = "rejected"
		} else {
			mutated["outcome"] = "approved"
		}
	case "requirement.record-reference":
		if input["status"] == "active" {
			mutated["status"] = "withdrawn"
		} else {
			mutated["status"] = "active"
		}
	case "requirement.link":
		if input["relation"] == "implements" {
			mutated["relation"] = "verifies"
		} else {
			mutated["relation"] = "implements"
		}
	case "claim.acquire":
		mutated["ttl_seconds"] = p23NumberValue(input["ttl_seconds"]) + 1
	case "claim.release":
		mutated["reason"] = "failed"
	case "envelope.complete":
		mutated["expected_revision"] = p23NumberValue(input["expected_revision"]) + 1
	default:
		return nil
	}
	return mutated
}

func p23NumberValue(value any) int {
	switch n := value.(type) {
	case int:
		return n
	case float64:
		return int(n)
	}
	return 0
}

func p23AcquireClaim(t *testing.T, a *p22FIDRecoveryAdapter, actor, mode, key string) map[string]any {
	t.Helper()
	response := p23RunInput(t, a, "claim", "acquire", actor, key, map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": mode, "ttl_seconds": 3600})
	p23AssertReviewSuccess(t, "acquire "+mode+" Claim", response)
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "claim").(map[string]any)
}

func p23StartRun(t *testing.T, a *p22FIDRecoveryAdapter, actor string, claim map[string]any, key string) map[string]any {
	t.Helper()
	response := p23RunInput(t, a, "run", "start", actor, key, map[string]any{"issue_id": a.issueID, "envelope_id": p23EnvelopeID(t, a), "envelope_revision": 1, "claim_id": claim["id"], "claim_version": claim["version"], "executor": map[string]any{"id": "executor-" + key, "source": "activated-harness"}, "obligation_id": nil, "expected_repository": p24ExpectedRepository(a.repository)})
	p23AssertReviewSuccess(t, "start Run", response)
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "run").(map[string]any)
}

func p23CloseAndRelease(t *testing.T, a *p22FIDRecoveryAdapter, actor, runID, claimID, key string) {
	t.Helper()
	run := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "run", "get", runID)
	version := p23Number(t, p21Path(t, run, "result", "run", "version"))
	executorID := p21PathString(t, run, "result", "run", "executor", "id")
	observed := p23RunInput(t, a, "run", "observe", actor, key+"-observe", map[string]any{"run_id": runID, "expected_version": version, "liveness": map[string]any{"executor_id": executorID, "state": "terminated", "source_ref": "process:p23-review", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	p23AssertReviewSuccess(t, "observe Run", observed)
	observedEnvelope := p23Envelope(t, observed.Stdout)
	close := p23RunInput(t, a, "run", "close", actor, key+"-close", map[string]any{"run_id": runID, "expected_version": p21Path(t, observedEnvelope, "result", "run", "version"), "observation_id": p21PathString(t, observedEnvelope, "result", "observation", "id"), "expected_repository": p24ExpectedRepository(a.repository)})
	p23AssertReviewSuccess(t, "close Run", close)
	release := p23RunInput(t, a, "claim", "release", actor, key+"-release", map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"})
	p23AssertReviewSuccess(t, "release Claim", release)
}

func p23ReviewGet(t *testing.T, a *p22FIDRecoveryAdapter, reviewID string) map[string]any {
	t.Helper()
	response := a.p22Run(t, "review", "get", reviewID)
	p23AssertReviewSuccess(t, "review get", response)
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "review").(map[string]any)
}

func p23ReviewFromResponse(t *testing.T, response p22ConformanceResponse) map[string]any {
	t.Helper()
	return p21Path(t, p23Envelope(t, response.Stdout), "result", "review").(map[string]any)
}

func p23Envelope(t *testing.T, stdout string) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("decode Review response: %v; stdout=%q", err, stdout)
	}
	return result
}

func p23AssertReviewSuccess(t *testing.T, operation string, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 0 || response.Stderr != "" {
		t.Fatalf("%s = %#v, want stdout-only success", operation, response)
	}
}

func p23EnvelopeID(t *testing.T, a *p22FIDRecoveryAdapter) string {
	t.Helper()
	run := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "run", "get", a.runID)
	return p21PathString(t, run, "result", "run", "envelope_id")
}

func p23Number(t *testing.T, value any) int {
	t.Helper()
	number, ok := value.(float64)
	if !ok {
		t.Fatalf("number = %#v, want float64 JSON number", value)
	}
	return int(number)
}

func p23Clone(value map[string]any) map[string]any {
	payload, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var clone map[string]any
	if err := json.Unmarshal(payload, &clone); err != nil {
		panic(err)
	}
	return clone
}
