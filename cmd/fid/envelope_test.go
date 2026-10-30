package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
	"github.com/mbell-nzl/fid/internal/project"
)

const testEvidenceObligationID = "OBL-1"

// Retained fixtures that use caller snapshots keep their local observation as
// the input to later Run/Gate/Evidence leaves. Observing fixtures still capture
// their current worktree. Neither kind shares observed history across roots.
var retainedRepositoryInputs sync.Map

func TestP24EnvelopeReplayAfterSourceChangeOrDeletion(t *testing.T) {
	t.Parallel()
	for _, position := range p24EnvelopeSourcePositions(t) {
		for _, state := range []string{"changed", "deleted"} {
			t.Run(position+"/"+state, func(t *testing.T) {
				t.Parallel()
				workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
				input := cloneEnvelopeInput(t, envelopeCreateInput(t, workspace, issueID, registration))
				source := p24EnvelopeSource(input, position)
				delete(source, "digest")
				source["path"] = "referenced.json"
				path := filepath.Join(workspace, "referenced.json")
				if err := os.WriteFile(path, []byte("original\n"), 0o600); err != nil {
					t.Fatal(err)
				}
				first := p24BindingMutation(t, workspace, "envelope create", input, "original")
				if first.err != nil {
					t.Fatalf("original: %#v", first)
				}
				if state == "deleted" {
					if err := os.Remove(path); err != nil {
						t.Fatal(err)
					}
				} else if err := os.WriteFile(path, []byte("changed\n"), 0o600); err != nil {
					t.Fatal(err)
				}
				before := snapshotWorkspacePersistence(t, workspace)
				replayed := p24BindingMutation(t, workspace, "envelope create", input, "original")
				if replayed.err != nil {
					t.Fatalf("replay: %#v", replayed)
				}
				want := decodeCommandJSONStructure(t, first.stdout).(map[string]any)
				got := decodeCommandJSONStructure(t, replayed.stdout).(map[string]any)
				if got["result"].(map[string]any)["replayed"] != true {
					t.Fatalf("not replayed: %#v", got)
				}
				got["result"].(map[string]any)["replayed"] = false
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("Source replay changed: %#v / %#v", got, want)
				}
				assertClaimStateUnchanged(t, workspace, before)
				control := cloneEnvelopeInput(t, input)
				control["obligations"].([]any)[0].(map[string]any)["id"] = "OBL-control"
				newKey := p24BindingMutation(t, workspace, "envelope create", control, "new-key")
				if state == "deleted" {
					if cli.ExitCode(newKey.err) != 4 || !strings.Contains(newKey.stderr, `"code":"source-digest-mismatch"`) || !strings.Contains(newKey.stderr, `"observed_digest":null`) || !strings.Contains(newKey.stderr, `"expected_digest":null`) {
						t.Fatalf("deleted Source control: %#v", newKey)
					}
					assertClaimStateUnchanged(t, workspace, before)
					if err := os.WriteFile(path, []byte("changed\n"), 0o600); err != nil {
						t.Fatal(err)
					}
					newKey = p24BindingMutation(t, workspace, "envelope create", control, "new-key")
				}
				if newKey.err != nil {
					t.Fatalf("viable Source control: %#v", newKey)
				}
				envelope := decodeCanonicalResult(t, newKey.stdout)["envelope"].(map[string]any)
				if p24EnvelopeSource(envelope, position)["digest"] != fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("changed\n"))) {
					t.Fatal("new-key control did not read changed Source")
				}
			})
		}
	}
}

func TestP24BindingReplayBeforeObservation(t *testing.T) {
	t.Parallel()
	for _, leaf := range []string{"gate register", "envelope create", "claim acquire", "envelope complete"} {
		for _, state := range []string{"head-change", "non-git", "unmerged"} {
			t.Run(leaf+"/"+state, func(t *testing.T) {
				t.Parallel()
				workspace, input, control := p24BindingReplayFixture(t, leaf)
				originalRepository := p24IndependentBindingRepository(t, workspace)
				first := p24BindingMutation(t, workspace, leaf, input, "original")
				if first.err != nil {
					t.Fatalf("original: %#v", first)
				}
				p24AssertBindingEventRepository(t, workspace, originalRepository)
				if leaf == "claim acquire" {
					claim := claimFromResponse(t, first.stdout)
					released := p24BindingMutation(t, workspace, "claim release", map[string]any{"claim_id": claim["id"], "expected_version": 1, "reason": "completed"}, "release")
					if released.err != nil {
						t.Fatalf("release for viable new-key control: %#v", released)
					}
				}
				git := func(args ...string) {
					t.Helper()
					if data, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).CombinedOutput(); err != nil {
						t.Fatalf("git %v: %v (%s)", args, err, data)
					}
				}
				restore := func() {}
				switch state {
				case "head-change":
					git("commit", "--allow-empty", "-qm", "new HEAD")
				case "non-git":
					if err := os.Rename(filepath.Join(workspace, ".git"), filepath.Join(workspace, ".git-hidden")); err != nil {
						t.Fatal(err)
					}
					restore = func() {
						if err := os.Rename(filepath.Join(workspace, ".git-hidden"), filepath.Join(workspace, ".git")); err != nil {
							t.Fatal(err)
						}
					}
				case "unmerged":
					blob, err := exec.Command("git", "-C", workspace, "rev-parse", "HEAD:.gitignore").Output()
					if err != nil {
						t.Fatal(err)
					}
					cmd := exec.Command("git", "-C", workspace, "update-index", "--index-info")
					cmd.Stdin = strings.NewReader(fmt.Sprintf("0 %s\t.gitignore\n100644 %s 1\t.gitignore\n100644 %s 2\t.gitignore\n", strings.Repeat("0", 40), strings.TrimSpace(string(blob)), strings.TrimSpace(string(blob))))
					if data, err := cmd.CombinedOutput(); err != nil {
						t.Fatalf("unmerged fixture: %v (%s)", err, data)
					}
					restore = func() { git("reset", "-q", "HEAD", "--", ".gitignore") }
				}
				before := snapshotWorkspacePersistence(t, workspace)
				replayed := p24BindingMutation(t, workspace, leaf, input, "original")
				if replayed.err != nil {
					t.Fatalf("replay %s: %#v", state, replayed)
				}
				originalJSON := decodeCommandJSONStructure(t, first.stdout).(map[string]any)
				replayJSON := decodeCommandJSONStructure(t, replayed.stdout).(map[string]any)
				if replayJSON["result"].(map[string]any)["replayed"] != true {
					t.Fatalf("not a replay: %#v", replayJSON)
				}
				replayJSON["result"].(map[string]any)["replayed"] = false
				if !reflect.DeepEqual(originalJSON, replayJSON) {
					t.Fatalf("replay changed original: %#v / %#v", originalJSON, replayJSON)
				}
				assertClaimStateUnchanged(t, workspace, before)
				newKey := p24BindingMutation(t, workspace, leaf, control, "new-key")
				if state != "head-change" {
					if cli.ExitCode(newKey.err) != 4 || newKey.stdout != "" {
						t.Fatalf("unavailable new-key control: %#v", newKey)
					}
					failure := decodeCommandJSONStructure(t, newKey.stderr).(map[string]any)["error"].(map[string]any)
					wantReason := "not-a-repository"
					if state == "unmerged" {
						wantReason = "unmerged-paths"
					}
					if failure["code"] != "repository-unavailable" || failure["category"] != "precondition" || failure["retryable"] != true || failure["subject"] != nil || len(failure["details"].(map[string]any)) != 2 || failure["details"].(map[string]any)["reason"] != wantReason {
						t.Fatalf("unavailable failure: %#v", failure)
					}
					assertClaimStateUnchanged(t, workspace, before)
					restore()
					newKey = p24BindingMutation(t, workspace, leaf, control, "new-key")
				}
				if newKey.err != nil {
					t.Fatalf("viable new-key control after restoring observation: %#v", newKey)
				}
			})
		}
	}
}

func p24BindingReplayFixture(t *testing.T, leaf string) (string, map[string]any, map[string]any) {
	t.Helper()
	if leaf == "gate register" {
		workspace := finalizedGateWorkspace(t, nil)
		input := gateRegisterInput(t, workspace)
		control := cloneEnvelopeInput(t, input)
		control["scope"] = map[string]any{"issue_id": "ISS-control"}
		return workspace, input, control
	}
	if leaf == "claim acquire" {
		workspace, issueID, envelopeID := claimEnvelopeForTest(t)
		input := map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"cmd"}}, "ttl_seconds": 3600}
		return workspace, input, cloneEnvelopeInput(t, input)
	}
	if leaf == "envelope create" {
		workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
		input := envelopeCreateInput(t, workspace, issueID, registration)
		control := cloneEnvelopeInput(t, input)
		control["obligations"].([]any)[0].(map[string]any)["id"] = "OBL-control"
		return workspace, input, control
	}
	workspace, issueID, envelopeID, _, _ := task18QualifiedEvidenceFixture(t, false, true)
	for index, state := range []string{"in_progress", "dev_complete", "qa_review", "qa_pass"} {
		result := executeCommandForTest(workspace, "--json", "issue", "transition", issueID, "--to", state, "--actor", "agent:actor-builder", "--expected-revision", fmt.Sprint(index+1), "--idempotency-key", "p24-"+state)
		if result.err != nil {
			t.Fatalf("terminal fixture: %#v", result)
		}
	}
	runID, claimID := task18RunAndClaim(t, workspace, issueID)
	closeRunForTest(t, workspace, runID, runRepositoryForTest(t, workspace), "p24-close")
	released := p24BindingMutation(t, workspace, "claim release", map[string]any{"claim_id": claimID, "expected_version": 1, "reason": "completed"}, "p24-release")
	if released.err != nil {
		t.Fatalf("release completion fixture: %#v", released)
	}
	registration := registerGateForTest(t, workspace, filepath.Join(workspace, "gate-register-envelope.json"), "p24-registration")
	controlCreate := envelopeCreateInput(t, workspace, issueID, registration)
	controlCreate["obligations"].([]any)[0].(map[string]any)["id"] = "OBL-control"
	controlCreate["obligations"].([]any)[0].(map[string]any)["guarded_actions"] = []string{"issue.transition-terminal-success"}
	grant := controlCreate["authority_grants"].([]any)[0].(map[string]any)
	grant["actions"] = []string{"envelope.activate", "envelope.complete"}
	grant["actor_selector"].(map[string]any)["required_capabilities"] = []string{"envelope.activate", "envelope.complete"}
	created := p24BindingMutation(t, workspace, "envelope create", controlCreate, "control-create")
	if created.err != nil {
		t.Fatalf("fresh completion control: %#v", created)
	}
	controlID := decodeCanonicalResult(t, created.stdout)["envelope"].(map[string]any)["id"].(string)
	activate := p24BindingMutation(t, workspace, "envelope activate", map[string]any{"envelope_id": controlID, "expected_revision": 1}, "control-activate")
	if activate.err != nil {
		t.Fatalf("fresh completion activation: %#v", activate)
	}
	return workspace, map[string]any{"envelope_id": envelopeID, "expected_revision": 1}, map[string]any{"envelope_id": controlID, "expected_revision": 1}
}

func TestP24EnvelopeAndClaimOriginalCopyAndImport(t *testing.T) {
	t.Parallel()
	for _, leaf := range []string{"envelope create", "claim acquire"} {
		t.Run(leaf, func(t *testing.T) {
			t.Parallel()
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := envelopeCreateInput(t, workspace, issueID, registration)
			if leaf == "claim acquire" {
				grant := input["authority_grants"].([]any)[0].(map[string]any)
				grant["actions"] = []string{"claim.acquire", "envelope.activate"}
				grant["actor_selector"].(map[string]any)["required_capabilities"] = []string{"claim.acquire", "envelope.activate"}
				created := p24BindingMutation(t, workspace, "envelope create", input, "create")
				if created.err != nil {
					t.Fatalf("create: %#v", created)
				}
				envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
				activated := p24BindingMutation(t, workspace, "envelope activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}, "activate")
				if activated.err != nil {
					t.Fatalf("activate: %#v", activated)
				}
				input = map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "scope": map[string]any{"paths": []string{"cmd"}}, "ttl_seconds": 3600}
			}
			copyRoot := t.TempDir()
			if err := copyWorkspaceTree(workspace, copyRoot); err != nil {
				t.Fatal(err)
			}
			// The copy is observable, with the same history but a different root.
			p24IndependentBindingRepository(t, copyRoot)
			before := snapshotWorkspacePersistence(t, copyRoot)
			refused := p24BindingMutation(t, copyRoot, leaf, input, "binding")
			if cli.ExitCode(refused.err) != 2 || refused.stdout != "" {
				t.Fatalf("plain copy: %#v", refused)
			}
			failure := decodeCommandJSONStructure(t, refused.stderr).(map[string]any)["error"].(map[string]any)
			message, allowed := "prior registration binding does not match envelope", "accepted envelope.create input"
			if leaf == "claim acquire" {
				message, allowed = "claim acquire binding is invalid", "accepted claim.acquire input"
			}
			if failure["code"] != "invalid-argument" || failure["message"] != message || !reflect.DeepEqual(failure["details"], map[string]any{"argument": "input", "value": nil, "allowed": []any{allowed}}) {
				t.Fatalf("plain copy failure: %#v", failure)
			}
			assertClaimStateUnchanged(t, copyRoot, before)
			archive := filepath.Join(t.TempDir(), "archive.zip")
			exported := executeCommandForTest(workspace, "export", "create", archive, "--json")
			if exported.err != nil {
				t.Fatalf("export: %#v", exported)
			}
			target := filepath.Join(t.TempDir(), "imported")
			if data, err := exec.Command("git", "clone", "-q", workspace, target).CombinedOutput(); err != nil {
				t.Fatalf("clone: %v (%s)", err, data)
			}
			imported := executeCommandForTest(target, "export", "import", archive, "--repository", target, "--idempotency-key", "import", "--json")
			if imported.err != nil {
				t.Fatalf("import: %#v", imported)
			}
			// Project sources are supplied by the destination checkout.
			if err := os.MkdirAll(filepath.Join(target, "project-docs/gates"), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(target, "project-docs/gates/unit.json"), []byte("gate policy fixture\n"), 0o600); err != nil {
				t.Fatal(err)
			}
			accepted := p24BindingMutation(t, target, leaf, input, "binding")
			if accepted.err != nil || accepted.stderr != "" {
				t.Fatalf("imported destination: %#v", accepted)
			}
			original := p24BindingMutation(t, workspace, leaf, input, "binding")
			if original.err != nil || original.stderr != "" {
				t.Fatalf("original control: %#v", original)
			}
			for _, result := range []struct {
				root     string
				response commandExecutionResult
			}{{workspace, original}, {target, accepted}} {
				repository := p24IndependentBindingRepository(t, result.root)
				if leaf == "claim acquire" {
					if claimFromResponse(t, result.response.stdout)["worktree_id"] != repository["worktree_id"] {
						t.Fatalf("Claim did not record local worktree: %#v", result.response)
					}
					// Claim Events retain the Envelope's observed scope, including imported history.
					envelopeID := input["envelope_id"].(string)
					get := executeCommandForTest(result.root, "--json", "envelope", "get", envelopeID)
					if get.err != nil {
						t.Fatalf("read scope: %#v", get)
					}
					repository = decodeCanonicalResult(t, get.stdout)["envelope"].(map[string]any)["repository_scope"].(map[string]any)
				}
				p24AssertBindingEventRepository(t, result.root, repository)
			}
		})
	}
}

func TestP24ImportedClaimStoredAuthority(t *testing.T) {
	t.Parallel()
	for _, action := range []string{"renew", "release", "expire"} {
		t.Run(action, func(t *testing.T) {
			t.Parallel()
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := envelopeCreateInput(t, workspace, issueID, registration)
			grant := input["authority_grants"].([]any)[0].(map[string]any)
			actions := []string{"claim.acquire", "claim.expire", "claim.release", "claim.renew", "envelope.activate"}
			grant["actions"] = actions
			grant["actor_selector"].(map[string]any)["required_capabilities"] = actions
			created := p24BindingMutation(t, workspace, "envelope create", input, "create")
			if created.err != nil {
				t.Fatalf("create: %#v", created)
			}
			envelopeID := assertEnvelopeShape(t, created.stdout, "draft")
			if activated := p24BindingMutation(t, workspace, "envelope activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}, "activate"); activated.err != nil {
				t.Fatalf("activate: %#v", activated)
			}
			archive := filepath.Join(t.TempDir(), "archive.zip")
			if exported := executeCommandForTest(workspace, "export", "create", archive, "--json"); exported.err != nil {
				t.Fatalf("export: %#v", exported)
			}
			target := filepath.Join(t.TempDir(), "imported")
			if data, err := exec.Command("git", "clone", "-q", workspace, target).CombinedOutput(); err != nil {
				t.Fatalf("clone: %v (%s)", err, data)
			}
			if imported := executeCommandForTest(target, "export", "import", archive, "--repository", target, "--idempotency-key", "import", "--json"); imported.err != nil {
				t.Fatalf("import: %#v", imported)
			}
			for _, root := range []string{workspace, target} {
				claim, _ := acquireGenericClaimForTest(t, root, issueID, envelopeID, "stored-authority-acquire")
				if claim["worktree_id"] != p24IndependentBindingRepository(t, root)["worktree_id"] {
					t.Fatalf("Claim did not record destination: %#v", claim)
				}
				request := map[string]any{"claim_id": claim["id"], "expected_version": 1}
				if action == "renew" {
					request["ttl_seconds"] = 120
				}
				if action == "release" {
					request["reason"] = "completed"
				}
				if root == target {
					copied := t.TempDir()
					if err := copyWorkspaceTree(target, copied); err != nil {
						t.Fatal(err)
					}
					before := p24BindingPersistenceSnapshot(t, copied)
					refused := p24BindingMutation(t, copied, "claim "+action, request, "stored-authority")
					if cli.ExitCode(refused.err) != 2 || refused.stdout != "" {
						t.Fatalf("copy exit: %#v", refused)
					}
					failure := decodeCommandJSONStructure(t, refused.stderr).(map[string]any)["error"].(map[string]any)
					if failure["code"] != "invalid-argument" || failure["message"] != "claim worktree binding is invalid" || !reflect.DeepEqual(failure["details"], map[string]any{"argument": "input", "value": nil, "allowed": []any{"accepted claim." + action + " input"}}) {
						t.Fatalf("copy failure: %#v", failure)
					}
					if !before.equal(p24BindingPersistenceSnapshot(t, copied)) {
						t.Fatal("copy changed persistence")
					}
				}
				mutated := p24BindingMutation(t, root, "claim "+action, request, "stored-authority")
				if mutated.err != nil || mutated.stderr != "" || stringifyClaimValue(claimFromResponse(t, mutated.stdout)["version"]) != "2" {
					t.Fatalf("%s in %s: %#v", action, root, mutated)
				}
				before := p24BindingPersistenceSnapshot(t, root)
				replayed := p24BindingMutation(t, root, "claim "+action, request, "stored-authority")
				if replayed.err != nil || !claimResponseReplayed(t, replayed.stdout) || !reflect.DeepEqual(claimFromResponse(t, replayed.stdout), claimFromResponse(t, mutated.stdout)) {
					t.Fatalf("replay %s: %#v", action, replayed)
				}
				if !before.equal(p24BindingPersistenceSnapshot(t, root)) {
					t.Fatal("replay changed persistence")
				}
				if action == "renew" {
					released := p24BindingMutation(t, root, "claim release", map[string]any{"claim_id": claim["id"], "expected_version": 2, "reason": "completed"}, "release-renewed")
					if released.err != nil || released.stderr != "" || stringifyClaimValue(claimFromResponse(t, released.stdout)["version"]) != "3" || claimFromResponse(t, released.stdout)["status"] != "released" {
						t.Fatalf("release renewed Claim: %#v", released)
					}
				}
			}
		})
	}
}

func TestP24EnvelopeSourcePositions(t *testing.T) {
	t.Parallel()
	positions := p24EnvelopeSourcePositions(t)
	for _, position := range positions {
		for _, state := range []string{"matching", "omitted", "mismatch", "missing", "removed-repository"} {
			t.Run(position+"/"+state, func(t *testing.T) {
				t.Parallel()
				workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
				input := cloneEnvelopeInput(t, envelopeCreateInput(t, workspace, issueID, registration))
				source := p24EnvelopeSource(input, position)
				digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("gate policy fixture\n")))
				switch state {
				case "matching":
					source["digest"] = digest
				case "omitted":
					delete(source, "digest")
				case "mismatch":
					source["digest"] = "sha256:" + strings.Repeat("0", 64)
				case "missing":
					delete(source, "digest")
					source["path"] = "missing.json"
				case "removed-repository":
					source["repository"] = testEnvelopeRepository(t, workspace)
				}
				observation := p24IndependentBindingRepository(t, workspace)
				before := snapshotWorkspacePersistence(t, workspace)
				result := p24BindingMutation(t, workspace, "envelope create", input, "source-position")
				if state == "matching" || state == "omitted" {
					if result.err != nil || result.stderr != "" {
						t.Fatalf("Source %s: %#v", state, result)
					}
					envelope := decodeCanonicalResult(t, result.stdout)["envelope"].(map[string]any)
					stored := p24EnvelopeSource(envelope, position)
					if stored["digest"] != digest || !sameJSONMap(stored["repository"].(map[string]any), observation) || !sameJSONMap(envelope["repository_scope"].(map[string]any), observation) {
						t.Fatalf("stored Source/scope differs from independent observation: %#v / %#v", stored, observation)
					}
					p24AssertBindingEventRepository(t, workspace, observation)
					return
				}
				assertClaimStateUnchanged(t, workspace, before)
				failure := decodeCommandJSONStructure(t, result.stderr).(map[string]any)["error"].(map[string]any)
				if state == "removed-repository" {
					if cli.ExitCode(result.err) != 2 || failure["code"] != "invalid-argument" || failure["details"].(map[string]any)["argument"] != "input."+position+".repository" {
						t.Fatalf("removed repository: %#v", result)
					}
					return
				}
				if cli.ExitCode(result.err) != 4 || failure["code"] != "source-digest-mismatch" {
					t.Fatalf("Source refusal: %#v", result)
				}
				wantExpected, wantObserved := source["digest"], any(digest)
				if state == "missing" {
					wantObserved = nil
				}
				want := map[string]any{"display_label": source["display_label"], "path": source["path"], "expected_digest": wantExpected, "observed_digest": wantObserved}
				if !reflect.DeepEqual(failure["details"], want) {
					t.Fatalf("Source details = %#v, want %#v", failure["details"], want)
				}
			})
		}
	}
}

func p24EnvelopeSourcePositions(t *testing.T) []string {
	t.Helper()
	var positions []string
	var walk func([]cli.InputField, string)
	walk = func(fields []cli.InputField, prefix string) {
		for _, field := range fields {
			if field.Name == "source_type" {
				positions = append(positions, prefix)
			}
			path := strings.TrimPrefix(prefix+"."+field.Name, ".")
			walk(field.Children, path)
			if field.Element != nil {
				walk(field.Element.Children, path+"[0]")
			}
		}
	}
	walk(inputDescriptorRegistry["envelope create"].Fields, "")
	want := []string{"risk_result.source", "sources[0]", "obligations[0].source"}
	if !reflect.DeepEqual(positions, want) {
		t.Fatalf("Source corpus positions = %v, want %v", positions, want)
	}
	return positions
}

func p24EnvelopeSource(input map[string]any, position string) map[string]any {
	switch position {
	case "risk_result.source":
		return input["risk_result"].(map[string]any)["source"].(map[string]any)
	case "sources[0]":
		return input["sources"].([]any)[0].(map[string]any)
	default:
		return input["obligations"].([]any)[0].(map[string]any)["source"].(map[string]any)
	}
}

func p24BindingMutation(t *testing.T, workspace, leaf string, input map[string]any, key string) commandExecutionResult {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.json")
	writeActivationInput(t, path, input)
	args := append([]string{"--project", workspace, "--json"}, strings.Split(leaf, " ")...)
	args = append(args, "--actor", "actor-builder", "--input", path, "--idempotency-key", key)
	return executeCommandForTest(workspace, args...)
}

func p24IndependentBindingRepository(t *testing.T, workspace string) map[string]any {
	t.Helper()
	git := func(args ...string) string {
		t.Helper()
		data, err := exec.Command("git", append([]string{"-C", workspace}, args...)...).Output()
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(data))
	}
	head, tree := git("rev-parse", "HEAD"), git("rev-parse", "HEAD^{tree}")
	if status := git("status", "--porcelain", "--untracked-files=all"); status != "" {
		t.Fatalf("independent fixture must be clean: %s", status)
	}
	projectConfig, err := config.ReadProjectConfig(filepath.Join(workspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(workspace)
	if err != nil {
		t.Fatal(err)
	}
	return map[string]any{"base_oid": head, "capture_source": "git", "dirty": false, "dirty_digest": nil, "head_oid": head, "project_id": projectConfig.ProjectID, "tree_oid": tree, "worktree_id": fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(root)))}
}

func p24BindingPersistenceSnapshot(t *testing.T, workspace string) workspacePersistenceSnapshot {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(workspace, ".fid/journal/*.jsonl"))
	if err != nil || len(paths) == 0 {
		t.Fatalf("journal inventory: %v / %v", paths, err)
	}
	snapshot := workspacePersistenceSnapshot{database: mustReadCommandFile(t, filepath.Join(workspace, ".fid/fid.sqlite"))}
	for _, path := range paths {
		snapshot.journal = append(snapshot.journal, mustReadCommandFile(t, path)...)
	}
	return snapshot
}

func p24AssertBindingEventRepository(t *testing.T, workspace string, repository map[string]any) {
	t.Helper()
	paths, err := filepath.Glob(filepath.Join(workspace, ".fid/journal/*.jsonl"))
	if err != nil || len(paths) == 0 {
		t.Fatalf("journal inventory: %v / %v", paths, err)
	}
	var event journal.Event
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range splitJournalLines(data) {
			var candidate journal.Event
			if err := json.Unmarshal([]byte(line), &candidate); err != nil {
				t.Fatal(err)
			}
			if event.EventID == "" || candidate.OccurredAt.After(event.OccurredAt) {
				event = candidate
			}
		}
	}
	encoded, err := json.Marshal(event.Repository)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, repository) {
		t.Fatalf("Event repository = %#v, want %#v", got, repository)
	}
}

func TestEnvelopeCreateObservesScopeAndSourcesWithoutCallerRepository(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedObservedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(nil))
	input := envelopeCreateInput(t, workspace, issueID, registration)
	delete(input, "repository_scope")
	source := input["sources"].([]any)[0].(map[string]any)
	delete(source, "repository")
	source["source_type"] = "external-reference"
	source["path"] = nil
	source["external_id"] = "gate-policy"
	source["digest"] = "sha256:" + strings.Repeat("0", 64)
	path := filepath.Join(workspace, "observed-envelope.json")
	writeActivationInput(t, path, input)
	created := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", path, "--idempotency-key", "observed-envelope")
	if created.err != nil {
		t.Fatalf("create observed Envelope: %v; stderr=%q", created.err, created.stderr)
	}
	envelope := decodeCommandJSONStructure(t, created.stdout).(map[string]any)["result"].(map[string]any)["envelope"].(map[string]any)
	if got, want := envelope["repository_scope"].(map[string]any)["worktree_id"], testEnvelopeRepository(t, workspace)["worktree_id"]; got != want {
		t.Fatalf("repository_scope.worktree_id = %v, want %v", got, want)
	}
	storedSource := envelope["sources"].([]any)[0].(map[string]any)
	if got, want := storedSource["repository"].(map[string]any)["worktree_id"], testEnvelopeRepository(t, workspace)["worktree_id"]; got != want {
		t.Fatalf("source.repository.worktree_id = %v, want %v", got, want)
	}
}

func TestEnvelopeCreateGetActivateRoundTripAndRejectsWidening(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	inputPath := filepath.Join(workspace, "envelope.json")
	input := envelopeCreateInput(t, workspace, issueID, registration)
	writeActivationInput(t, inputPath, input)

	created := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder",
		"--input", inputPath, "--idempotency-key", "envelope-create-1",
	)
	if created.err != nil || created.stderr != "" {
		t.Fatalf("envelope create = %#v, want stdout-only success", created)
	}
	envelopeID := assertEnvelopeShape(t, created.stdout, "draft")

	got := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "get", envelopeID)
	if got.err != nil || got.stderr != "" || assertEnvelopeShape(t, got.stdout, "draft") != envelopeID {
		t.Fatalf("envelope get = %#v, want exact draft envelope", got)
	}

	activatePath := filepath.Join(workspace, "envelope-activate.json")
	writeActivationInput(t, activatePath, map[string]any{"envelope_id": envelopeID, "expected_revision": 1})
	activated := executeCommandForTest(workspace,
		"--project", workspace, "--json", "envelope", "activate", "--actor", "actor-builder",
		"--input", activatePath, "--idempotency-key", "envelope-activate-1",
	)
	if activated.err != nil || activated.stderr != "" || assertEnvelopeShape(t, activated.stdout, "active") != envelopeID {
		t.Fatalf("envelope activate = %#v, want active exact envelope", activated)
	}

	for _, testCase := range []struct {
		name  string
		apply func(map[string]any)
	}{
		{"related issues", func(value map[string]any) { value["related_issue_ids"] = []string{"ISS-2"} }},
		{"duplicate obligation ID", func(value map[string]any) {
			value["obligations"] = append(value["obligations"].([]any), value["obligations"].([]any)[0])
		}},
		{"computed status", func(value map[string]any) {
			value["obligations"].([]any)[0].(map[string]any)["derived_status"] = "satisfied"
		}},
		{"other obligation type", func(value map[string]any) { value["obligations"].([]any)[0].(map[string]any)["type"] = "human-outcome" }},
		{"guard widening", func(value map[string]any) {
			value["obligations"].([]any)[0].(map[string]any)["guarded_actions"] = []string{"envelope.complete", "issue.transition-terminal-success", "run.close"}
		}},
		{"mutable budget", func(value map[string]any) { value["budgets"].([]any)[0].(map[string]any)["consumed"] = 1 }},
		{"budget output state empty", func(value map[string]any) { value["budgets"].([]any)[0].(map[string]any)["state"] = "" }},
		{"budget output state null", func(value map[string]any) { value["budgets"].([]any)[0].(map[string]any)["state"] = nil }},
		{"routine lane", func(value map[string]any) { value["routine"] = map[string]any{"classification": "external"} }},
		{"revision", func(value map[string]any) { value["revision"] = 2 }},
		{"cancellation", func(value map[string]any) { value["cancellation"] = map[string]any{} }},
		{"protected decision", func(value map[string]any) { value["protected_decision"] = "DEC-1" }},
		{"stop", func(value map[string]any) { value["stop"] = map[string]any{} }},
		{"attention", func(value map[string]any) { value["attention"] = map[string]any{} }},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			candidate := cloneEnvelopeInput(t, input)
			testCase.apply(candidate)
			writeActivationInput(t, inputPath, candidate)
			before := journalEventCount(t, workspace)
			result := executeCommandForTest(workspace,
				"--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder",
				"--input", inputPath, "--idempotency-key", "envelope-reject-"+strings.ReplaceAll(testCase.name, " ", "-"),
			)
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
				t.Fatalf("%s result = %#v, want typed closed-input rejection", testCase.name, result)
			}
			if after := journalEventCount(t, workspace); after != before {
				t.Fatalf("%s event count = %d, want no Event after %d", testCase.name, after, before)
			}
		})
	}
}

func TestEnvelopeCreateAcceptsTypedObligationsAndProtectedDecisions(t *testing.T) {
	t.Parallel()
	workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
	input := envelopeCreateInput(t, workspace, issueID, registration)
	source := input["sources"].([]any)[0]
	input["obligations"] = []any{
		input["obligations"].([]any)[0],
		map[string]any{"id": "OBL-review", "type": "review", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "review_profile": "standard", "accepted_evidence_types": []string{"review"}, "accepted_trust_sources": []string{"activated-actor"}, "freshness_target": "repository-tree", "artifact_policy": "metadata-sufficient", "independence": "distinct-actor", "required_lane": nil, "source": source},
		map[string]any{"id": "OBL-human", "type": "human-outcome", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "outcome_class": "release", "accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"}, "freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source},
	}
	input["protected_decisions"] = []any{map[string]any{"ordinal": 0, "decision_class": "release", "required_capabilities": []string{"envelope.complete"}, "affected_actions": []string{"envelope.complete"}}}
	inputPath := filepath.Join(workspace, "typed-envelope.json")
	writeActivationInput(t, inputPath, input)

	result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "typed-envelope")
	if result.err != nil || result.stderr != "" {
		t.Fatalf("typed envelope create = %#v, want success", result)
	}
	payload := decodeCommandJSONStructure(t, result.stdout).(map[string]any)
	envelope := payload["result"].(map[string]any)["envelope"].(map[string]any)
	obligations := envelope["obligations"].([]any)
	if len(obligations) != 3 || obligations[1].(map[string]any)["registration"] != nil || obligations[2].(map[string]any)["registration"] != nil {
		t.Fatalf("typed Envelope obligations = %#v, want three rows with null new-type registrations", obligations)
	}
	if decisions := envelope["protected_decisions"].([]any); len(decisions) != 1 {
		t.Fatalf("typed Envelope protected decisions = %#v, want one authoring row", decisions)
	}
}

func TestEnvelopeCreateCommandCoversObligationAndDecisionBoundaries(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name            string
		obligationCount int
		decisionCount   int
	}{
		{name: "one obligation and zero decisions", obligationCount: 1, decisionCount: 0},
		{name: "sixteen obligations and sixteen decisions", obligationCount: 16, decisionCount: 16},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := typedEnvelopeCreateInput(t, workspace, issueID, registration)
			input["obligations"] = commandObligationsForCount(t, input, issueID, test.obligationCount)
			input["protected_decisions"] = commandProtectedDecisionMaps(test.decisionCount)
			inputPath := filepath.Join(workspace, "envelope-boundary.json")
			writeActivationInput(t, inputPath, input)

			result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "envelope-boundary-"+strings.ReplaceAll(test.name, " ", "-"))
			if result.err != nil || result.stderr != "" {
				t.Fatalf("boundary create = %#v, want success", result)
			}
			payload := decodeCommandJSONStructure(t, result.stdout).(map[string]any)
			envelope := payload["result"].(map[string]any)["envelope"].(map[string]any)
			if got := len(envelope["obligations"].([]any)); got != test.obligationCount {
				t.Fatalf("obligation count = %d, want %d", got, test.obligationCount)
			}
			if got := len(envelope["protected_decisions"].([]any)); got != test.decisionCount {
				t.Fatalf("protected decision count = %d, want %d", got, test.decisionCount)
			}
		})
	}
}

func TestEnvelopeCreateCommandRejectsEveryTypedAuthoringNegativeWithoutPersistenceChange(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		mutate func(t *testing.T, input map[string]any, issueID string, registration gateRegistrationView)
	}{
		{name: "zero obligations", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["obligations"] = []any{}
		}},
		{name: "seventeen obligations", mutate: func(t *testing.T, input map[string]any, issueID string, registration gateRegistrationView) {
			input["obligations"] = commandObligationsForCount(t, input, issueID, 17)
		}},
		{name: "duplicate obligation ID", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			obligations := input["obligations"].([]any)
			obligations[1].(map[string]any)["id"] = obligations[0].(map[string]any)["id"]
		}},
		{name: "second gate-evidence obligation", mutate: func(t *testing.T, input map[string]any, issueID string, registration gateRegistrationView) {
			obligations := input["obligations"].([]any)
			input["obligations"] = append(obligations, commandObligationsForCount(t, input, issueID, 1)[0])
			input["obligations"].([]any)[3].(map[string]any)["id"] = "OBL-second-gate"
		}},
		{name: "duplicate review profile", mutate: func(_ *testing.T, input map[string]any, issueID string, _ gateRegistrationView) {
			obligations := input["obligations"].([]any)
			duplicate := commandReviewObligation(input["sources"].([]any)[0], issueID, "OBL-review-duplicate", "standard")
			input["obligations"] = append(obligations, duplicate)
		}},
		{name: "duplicate outcome class", mutate: func(_ *testing.T, input map[string]any, issueID string, _ gateRegistrationView) {
			obligations := input["obligations"].([]any)
			duplicate := commandHumanObligation(input["sources"].([]any)[0], issueID, "OBL-human-duplicate", "release")
			input["obligations"] = append(obligations, duplicate)
		}},
		{name: "non-primary subject", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["obligations"].([]any)[0].(map[string]any)["subject"] = map[string]any{"type": "issue", "id": "ISS-other"}
		}},
		{name: "seventeen protected decisions", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["protected_decisions"] = commandProtectedDecisionMaps(17)
		}},
		{name: "duplicate protected decision class", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["protected_decisions"] = []any{commandProtectedDecisionMap(0, "release", "envelope.complete"), commandProtectedDecisionMap(1, "release", "envelope.complete")}
		}},
		{name: "duplicate protected decision ordinal", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			input["protected_decisions"] = []any{commandProtectedDecisionMap(0, "release", "envelope.complete"), commandProtectedDecisionMap(0, "rollback", "envelope.complete")}
		}},
		{name: "unguarded protected decision action", mutate: func(_ *testing.T, input map[string]any, _ string, _ gateRegistrationView) {
			for _, obligation := range input["obligations"].([]any) {
				obligation.(map[string]any)["guarded_actions"] = []string{"envelope.complete"}
			}
			input["protected_decisions"] = []any{commandProtectedDecisionMap(0, "release", "issue.transition-terminal-success")}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := typedEnvelopeCreateInput(t, workspace, issueID, registration)
			test.mutate(t, input, issueID, registration)
			inputPath := filepath.Join(workspace, "envelope-rejected.json")
			writeActivationInput(t, inputPath, input)
			before := snapshotWorkspacePersistence(t, workspace)

			result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "envelope-rejected-"+strings.ReplaceAll(test.name, " ", "-"))
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
				t.Fatalf("rejected command = %#v, want invalid-argument with no stdout", result)
			}
			if after := snapshotWorkspacePersistence(t, workspace); !before.equal(after) {
				t.Fatalf("rejected %s changed journal or projection", test.name)
			}
		})
	}
}

func TestEnvelopeCreateCommandRejectsEveryForbiddenTypedFieldCombinationWithoutPersistenceChange(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{name: "gate missing registration", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["registration"] = nil
		}},
		{name: "gate has review profile", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["review_profile"] = "profile"
		}},
		{name: "gate has outcome class", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["outcome_class"] = "outcome"
		}},
		{name: "gate evidence type mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["accepted_evidence_types"] = []string{"review"}
		}},
		{name: "gate trust source mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["accepted_trust_sources"] = []string{"activated-actor"}
		}},
		{name: "gate freshness mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["freshness_target"] = "none"
		}},
		{name: "gate artifact policy mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["artifact_policy"] = "metadata-sufficient"
		}},
		{name: "gate independence mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "gate-evidence")["independence"] = "distinct-actor"
		}},
		{name: "review has registration", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["registration"] = map[string]any{"id": "REG-test", "version": 1}
		}},
		{name: "review missing profile", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["review_profile"] = nil
		}},
		{name: "review has outcome class", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["outcome_class"] = "outcome"
		}},
		{name: "review evidence type mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["accepted_evidence_types"] = []string{"receipt"}
		}},
		{name: "review trust source mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["accepted_trust_sources"] = []string{"fid-executed"}
		}},
		{name: "review freshness mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["freshness_target"] = "none"
		}},
		{name: "review artifact policy mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["artifact_policy"] = "required-available"
		}},
		{name: "review independence mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "review")["independence"] = "human"
		}},
		{name: "human has registration", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["registration"] = map[string]any{"id": "REG-test", "version": 1}
		}},
		{name: "human has review profile", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["review_profile"] = "profile"
		}},
		{name: "human missing outcome class", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["outcome_class"] = nil
		}},
		{name: "human evidence type mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["accepted_evidence_types"] = []string{"review"}
		}},
		{name: "human trust source mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["accepted_trust_sources"] = []string{"activated-actor"}
		}},
		{name: "human freshness mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["freshness_target"] = "repository-tree"
		}},
		{name: "human artifact policy mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["artifact_policy"] = "required-available"
		}},
		{name: "human independence mismatch", mutate: func(input map[string]any) {
			commandObligationByType(t, input, "human-outcome")["independence"] = "distinct-actor"
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace, issueID, registration := finalizedEnvelopeWorkspace(t)
			input := typedEnvelopeCreateInput(t, workspace, issueID, registration)
			test.mutate(input)
			inputPath := filepath.Join(workspace, "envelope-forbidden-field.json")
			writeActivationInput(t, inputPath, input)
			before := snapshotWorkspacePersistence(t, workspace)

			result := executeCommandForTest(workspace, "--project", workspace, "--json", "envelope", "create", "--actor", "actor-builder", "--input", inputPath, "--idempotency-key", "envelope-forbidden-field-"+strings.ReplaceAll(test.name, " ", "-"))
			if result.err == nil || result.stdout != "" || !strings.Contains(result.stderr, `"code":"invalid-argument"`) {
				t.Fatalf("forbidden typed field command = %#v, want invalid-argument with no stdout", result)
			}
			if after := snapshotWorkspacePersistence(t, workspace); !before.equal(after) {
				t.Fatalf("forbidden typed field %s changed journal or projection", test.name)
			}
		})
	}
}

func commandObligationByType(t *testing.T, input map[string]any, wantType string) map[string]any {
	t.Helper()
	for _, value := range input["obligations"].([]any) {
		obligation := value.(map[string]any)
		if obligation["type"] == wantType {
			return obligation
		}
	}
	t.Fatalf("typed input has no %q Obligation", wantType)
	return nil
}

func TestDecodeClosedEnvelopeInputRejectsAuthoredBudgetState(t *testing.T) {
	t.Parallel()
	for _, state := range []string{`""`, `null`} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			descriptor := inputDescriptorRegistry["envelope create"]
			example, err := cli.InputExample(descriptor)
			if err != nil {
				t.Fatal(err)
			}
			var input map[string]any
			if err := json.Unmarshal(example, &input); err != nil {
				t.Fatal(err)
			}
			p24ExpandDescriptorObjects(t, input, descriptor.Fields)
			var authored any
			if err := json.Unmarshal([]byte(state), &authored); err != nil {
				t.Fatal(err)
			}
			input["budgets"].([]any)[0].(map[string]any)["state"] = authored
			root := t.TempDir()
			path := filepath.Join(t.TempDir(), "authored-state.json")
			writeActivationInput(t, path, input)
			result := executeCommandForTest(root, "--json", "envelope", "create", "--input", path, "--actor", "ACT-1", "--idempotency-key", "authored-state")
			argument, allowed := "input", []string{"accepted envelope.create input"}
			value := any(nil)
			if state == "null" {
				argument, allowed, value = "input.budgets[0].state", []string{"string"}, "null"
			}
			if cli.ExitCode(result.err) != 2 || result.stdout != "" {
				t.Fatalf("authored state accepted: %#v", result)
			}
			assertAdoptedJSONFailure(t, result.stderr, "envelope.create", argument, value, allowed)
			entries, err := os.ReadDir(root)
			if err != nil || len(entries) != 0 {
				t.Fatalf("authored state changed root: %v / %v", entries, err)
			}
		})
	}
}

func finalizedEnvelopeWorkspace(t *testing.T) (string, string, gateRegistrationView) {
	return finalizedEnvelopeWorkspaceWithDefinition(t, nil)
}

func finalizedEnvelopeWorkspaceWithDefinition(t *testing.T, modifyDefinition func(map[string]any)) (string, string, gateRegistrationView) {
	t.Helper()
	return finalizedEnvelopeWorkspaceWithDefinitions(t, gateDefinitionsForTest(modifyDefinition))
}

func finalizedEnvelopeWorkspaceWithDefinitions(t *testing.T, definitions []any) (string, string, gateRegistrationView) {
	return finalizedEnvelopeWorkspaceWithRepository(t, definitions, false)
}

func finalizedObservedEnvelopeWorkspaceWithDefinitions(t *testing.T, definitions []any) (string, string, gateRegistrationView) {
	return finalizedEnvelopeWorkspaceWithRepository(t, definitions, true)
}

func finalizedEnvelopeWorkspaceWithRepository(t *testing.T, definitions []any, observeRepository bool) (string, string, gateRegistrationView) {
	t.Helper()
	return buildFinalizedEnvelopeWorkspaceWithRepository(t, definitions, observeRepository)
}

func buildFinalizedEnvelopeWorkspace(t *testing.T, definitions []any) (string, string, gateRegistrationView) {
	return buildFinalizedEnvelopeWorkspaceWithRepository(t, definitions, false)
}

func buildFinalizedEnvelopeWorkspaceWithRepository(t *testing.T, definitions []any, observeRepository bool) (string, string, gateRegistrationView) {
	t.Helper()
	// Reuse only facts produced before any path-bound repository observation.
	// Registration and every subsequent Envelope/Claim/Run use this clone's root.
	key := reusableWorkspaceFixtureKey(t, "unregistered-envelope", definitions)
	workspace, values := cloneReusableWorkspaceFixture(t, key, func() (string, []string) {
		root, issueID := buildUnregisteredEnvelopeWorkspace(t, definitions)
		return root, []string{issueID}
	})
	issueID := values[0]
	initializeObservedRepositoryForTest(t, workspace)
	registrationInput := gateRegisterInput(t, workspace)
	registrationInput["scope"] = map[string]any{"issue_id": issueID}
	registrationPath := filepath.Join(workspace, "gate-register-envelope.json")
	writeActivationInput(t, registrationPath, registrationInput)
	registration := registerGateForTest(t, workspace, registrationPath, "gate-register-envelope")
	if !observeRepository {
		retainedRepositoryInputs.Store(workspace, testEnvelopeRepository(t, workspace))
	}
	return workspace, issueID, registration
}

func buildUnregisteredEnvelopeWorkspace(t *testing.T, definitions []any) (string, string) {
	t.Helper()
	prepare := func() (string, []string) {
		workspace := clonePristineWorkspace(t)
		created := executeCommandForTest(workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Envelope issue", "--actor", "agent:developer")
		if created.err != nil {
			t.Fatalf("create issue: %v; stderr=%q", created.err, created.stderr)
		}
		var issue struct {
			Result struct {
				Issue struct {
					ID string `json:"id"`
				} `json:"issue"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(created.stdout), &issue); err != nil || issue.Result.Issue.ID == "" {
			t.Fatalf("decode created issue: %v; stdout=%q", err, created.stdout)
		}
		activationID := prepareActivationForTest(t, workspace)
		profileDigest := stageActivationComponentForTest(t, workspace, activationID, "authority-profile", "envelope-builders", []any{map[string]any{
			"actor_id": "actor-builder", "actor_type": "agent", "producer": "kit", "provenance_ref": "activation:profile/envelope-builders", "session_ref": nil,
			"capabilities": authorityCapabilities(),
		}, map[string]any{
			"actor_id": "actor-harness", "actor_type": "service", "producer": "activated-harness", "provenance_ref": "activation:profile/envelope-harness", "session_ref": nil,
			"capabilities": []string{"run.observe"},
		}})
		return workspace, []string{issue.Result.Issue.ID, activationID, profileDigest}
	}
	workspace, values := cloneReusableWorkspaceFixture(t, reusableWorkspaceFixtureKey(t, "prepared-envelope", false), prepare)
	issueID, activationID, profileDigest := values[0], values[1], values[2]
	gateDigest := stageActivationComponentForTest(t, workspace, activationID, "gate-definitions", "default-gates", definitions)
	observeActivationResourcesForTest(t, workspace, activationID, "new")
	digests := []string{profileDigest, gateDigest}
	sort.Strings(digests)
	commitPath := filepath.Join(workspace, "activation-commit-envelope.json")
	writeActivationInput(t, commitPath, map[string]any{"activation_id": activationID, "expected_status": "staged", "expected_version": activationVersionForTest(t, workspace, activationID), "component_digests": digests})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "commit", "--input", commitPath, "--idempotency-key", "activation-commit-envelope"); result.err != nil {
		t.Fatalf("commit activation: %v; stderr=%q", result.err, result.stderr)
	}
	finalizePath := filepath.Join(workspace, "activation-finalize-envelope.json")
	writeActivationInput(t, finalizePath, map[string]any{"activation_id": activationID, "expected_status": "committed", "expected_version": activationVersionForTest(t, workspace, activationID)})
	if result := executeCommandForTest(workspace, "--project", workspace, "--json", "activation", "finalize", "--input", finalizePath, "--idempotency-key", "activation-finalize-envelope"); result.err != nil {
		t.Fatalf("finalize activation: %v; stderr=%q", result.err, result.stderr)
	}
	return workspace, issueID
}

func gateDefinitionsForTest(modifyDefinition func(map[string]any)) []any {
	definitions := validGateDefinitions()
	if modifyDefinition != nil {
		modifyDefinition(definitions[0].(map[string]any))
	}
	return definitions
}

func authorityCapabilities() []string {
	return []string{
		"artifact.verify", "claim.acquire", "claim.expire", "claim.release", "claim.renew", "envelope.activate", "envelope.complete",
		"envelope.create", "evidence.link", "finding.address", "gate.register", "gate.run", "issue.transition", "review.record", "run.close", "run.observe", "run.start",
	}
}

func envelopeCreateInput(t *testing.T, workspace, issueID string, registration gateRegistrationView) map[string]any {
	t.Helper()
	sourcePath := filepath.Join(workspace, "project-docs", "gates", "unit.json")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sourcePath, []byte("gate policy fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	source := map[string]any{
		"digest":        fmt.Sprintf("sha256:%x", sha256.Sum256([]byte("gate policy fixture\n"))),
		"display_label": "Gate policy", "external_id": nil,
		"path": "project-docs/gates/unit.json", "section": nil, "source_type": "project-file",
	}
	return map[string]any{
		"outcome": "Prove the envelope value", "non_goals": []string{},
		"risk_result": map[string]any{"source": source, "result_digest": "sha256:" + strings.Repeat("1", 64)},
		"sources":     []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "protected_decisions": []any{},
		"authority_grants": []any{map[string]any{
			"actions":        []string{"envelope.activate"},
			"actor_selector": map[string]any{"actor_id": "actor-builder", "required_capabilities": []string{"envelope.activate"}},
			"scope":          map[string]any{"issue_id": issueID},
		}},
		"obligations": []any{map[string]any{
			"id": testEvidenceObligationID, "type": "gate-evidence", "applicability": "required",
			"guarded_actions":         []string{"envelope.complete", "issue.transition-terminal-success"},
			"subject":                 map[string]any{"type": "issue", "id": issueID},
			"registration":            map[string]any{"id": registration.ID, "version": registration.Version},
			"accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"},
			"freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source,
		}},
		"budgets": []any{map[string]any{
			"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3, "warning_threshold": nil,
			"exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil,
		}},
	}
}

func typedEnvelopeCreateInput(t *testing.T, workspace, issueID string, registration gateRegistrationView) map[string]any {
	t.Helper()
	input := envelopeCreateInput(t, workspace, issueID, registration)
	source := input["sources"].([]any)[0]
	input["obligations"] = []any{
		input["obligations"].([]any)[0],
		commandReviewObligation(source, issueID, "OBL-review", "standard"),
		commandHumanObligation(source, issueID, "OBL-human", "release"),
	}
	input["protected_decisions"] = []any{}
	return input
}

func commandObligationsForCount(t *testing.T, input map[string]any, issueID string, count int) []any {
	t.Helper()
	if count < 1 || count > 17 {
		t.Fatalf("command obligation count = %d, want 1..17", count)
	}
	source := input["sources"].([]any)[0]
	gate := input["obligations"].([]any)[0].(map[string]any)
	obligations := []any{gate}
	for index := 1; index < count; index++ {
		if index%2 == 1 {
			obligations = append(obligations, commandReviewObligation(source, issueID, "OBL-review-"+string(rune('a'+index)), "profile-"+string(rune('a'+index))))
			continue
		}
		obligations = append(obligations, commandHumanObligation(source, issueID, "OBL-human-"+string(rune('a'+index)), "outcome-"+string(rune('a'+index))))
	}
	return obligations
}

func commandReviewObligation(source any, issueID, id, profile string) map[string]any {
	return map[string]any{
		"id": id, "type": "review", "applicability": "required", "guarded_actions": []string{"envelope.complete"},
		"subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "review_profile": profile,
		"accepted_evidence_types": []string{"review"}, "accepted_trust_sources": []string{"activated-actor"},
		"freshness_target": "repository-tree", "artifact_policy": "metadata-sufficient", "independence": "distinct-actor", "required_lane": nil, "source": source,
	}
}

func commandHumanObligation(source any, issueID, id, outcomeClass string) map[string]any {
	return map[string]any{
		"id": id, "type": "human-outcome", "applicability": "required", "guarded_actions": []string{"envelope.complete"},
		"subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "outcome_class": outcomeClass,
		"accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"},
		"freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source,
	}
}

func commandProtectedDecisionMaps(count int) []any {
	decisions := make([]any, 0, count)
	for ordinal := 0; ordinal < count; ordinal++ {
		decisions = append(decisions, commandProtectedDecisionMap(ordinal, "decision-"+string(rune('a'+ordinal)), "envelope.complete"))
	}
	return decisions
}

func commandProtectedDecisionMap(ordinal int, decisionClass, action string) map[string]any {
	return map[string]any{
		"ordinal": ordinal, "decision_class": decisionClass, "required_capabilities": []string{"envelope.complete"}, "affected_actions": []string{action},
	}
}

func testEnvelopeRepository(t *testing.T, workspace string) map[string]any {
	t.Helper()
	projectConfig, err := config.ReadProjectConfig(filepath.Join(workspace, ".fid", "config.toml"))
	if err != nil {
		t.Fatalf("read project config: %v", err)
	}
	if _, err := os.Stat(filepath.Join(workspace, ".git")); err == nil {
		observation, err := project.ObserveRepository(workspace)
		if err != nil {
			t.Fatalf("observe test repository: %v", err)
		}
		return map[string]any{"base_oid": observation.HeadOID, "capture_source": "git", "dirty": observation.Dirty, "dirty_digest": observation.DirtyDigest, "head_oid": observation.HeadOID, "project_id": projectConfig.ProjectID, "tree_oid": observation.TreeOID, "worktree_id": observation.WorktreeID}
	}
	return map[string]any{
		"base_oid": nil, "capture_source": "git", "dirty": false, "dirty_digest": nil,
		"head_oid": strings.Repeat("a", 40), "project_id": projectConfig.ProjectID, "tree_oid": strings.Repeat("b", 40), "worktree_id": "worktree-1",
	}
}

func initializeObservedRepositoryForTest(t *testing.T, workspace string) {
	t.Helper()
	for _, arguments := range [][]string{{"init", "-q"}, {"config", "user.email", "fid@example.test"}, {"config", "user.name", "fid test"}} {
		if output, err := exec.Command("git", append([]string{"-C", workspace}, arguments...)...).CombinedOutput(); err != nil {
			t.Fatalf("initialize observed repository %v: %v (%s)", arguments, err, output)
		}
	}
	if err := os.WriteFile(filepath.Join(workspace, ".gitignore"), []byte(".fid/\n*.json\n"), 0o600); err != nil {
		t.Fatalf("write observed repository ignore: %v", err)
	}
	for _, arguments := range [][]string{{"add", ".gitignore"}, {"commit", "-qm", "fixture"}} {
		if output, err := exec.Command("git", append([]string{"-C", workspace}, arguments...)...).CombinedOutput(); err != nil {
			t.Fatalf("commit observed repository %v: %v (%s)", arguments, err, output)
		}
	}
}

func assertEnvelopeShape(t *testing.T, output, wantStatus string) string {
	t.Helper()
	value := decodeCommandJSONStructure(t, output).(map[string]any)
	envelope := value["result"].(map[string]any)["envelope"].(map[string]any)
	if value["operation"] != "envelope.create" && value["operation"] != "envelope.get" && value["operation"] != "envelope.activate" {
		t.Fatalf("operation = %#v, want envelope operation", value["operation"])
	}
	if envelope["current_revision"] != json.Number("1") || envelope["status"] != wantStatus || len(envelope["related_issue_ids"].([]any)) != 0 {
		t.Fatalf("envelope = %#v, want revision 1 %s and no related issues", envelope, wantStatus)
	}
	obligation := envelope["obligations"].([]any)[0].(map[string]any)
	if obligation["type"] != "gate-evidence" || obligation["derived_status"] != "unsatisfied" ||
		strings.Join(anyStrings(t, obligation["guarded_actions"]), ",") != "envelope.complete,issue.transition-terminal-success" {
		t.Fatalf("obligation = %#v, want one gate-evidence dual-guard value", obligation)
	}
	budget := envelope["budgets"].([]any)[0].(map[string]any)
	if budget["consumed"] != json.Number("0") || budget["remaining"] != json.Number("3") || budget["state"] != "available" {
		t.Fatalf("budget = %#v, want derived counters and available state", budget)
	}
	return envelope["id"].(string)
}

func cloneEnvelopeInput(t *testing.T, input map[string]any) map[string]any {
	t.Helper()
	payload, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal envelope input: %v", err)
	}
	var cloned map[string]any
	if err := json.Unmarshal(payload, &cloned); err != nil {
		t.Fatalf("unmarshal envelope input: %v", err)
	}
	return cloned
}

func anyStrings(t *testing.T, value any) []string {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("value = %T, want []any", value)
	}
	result := make([]string, len(items))
	for index, item := range items {
		result[index], ok = item.(string)
		if !ok {
			t.Fatalf("item %d = %T, want string", index, item)
		}
	}
	return result
}

func journalEventCount(t *testing.T, workspace string) int {
	t.Helper()
	return len(splitJournalLines(mustReadCommandFile(t, mustJournalPath(t, workspace))))
}
