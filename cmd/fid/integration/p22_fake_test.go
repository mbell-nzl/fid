//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/project"
)

const p22CorpusDir = "testdata/p22_conformance"

type p22ConformanceCase struct {
	Name         string         `json:"name"`
	Operation    string         `json:"operation"`
	Kind         string         `json:"kind"`
	ResultFields []string       `json:"result_fields"`
	Event        bool           `json:"event"`
	FailureCode  string         `json:"failure_code"`
	ExitCode     int            `json:"exit_code"`
	Expected     map[string]any `json:"expected"`
	Input        map[string]any `json:"input,omitempty"`
	Reset        bool           `json:"reset,omitempty"`
}

type p22ConformanceResponse struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type p22RecoveryAdapter interface {
	Name() string
	Execute(t *testing.T, fixture p22ConformanceCase) p22ConformanceResponse
}

func TestP22DualAdapterRecoveryConformance(t *testing.T) {
	fixtures := p22RecoveryFixtures(t)
	for _, newAdapter := range []func(*testing.T) p22RecoveryAdapter{
		func(_ *testing.T) p22RecoveryAdapter { return p22NewFakeRecoveryAdapter() },
		func(t *testing.T) p22RecoveryAdapter { return p22NewFIDRecoveryAdapter(t) },
	} {
		adapter := newAdapter(t)
		t.Run(adapter.Name(), func(t *testing.T) {
			var priorRun string
			var denialStops, denialAttention any
			adapterScope := t
			for _, fixture := range fixtures {
				t.Run(fixture.Name, func(t *testing.T) {
					if fixture.Reset {
						adapter = newAdapter(adapterScope)
						priorRun = ""
						denialStops, denialAttention = nil, nil
					}
					response := adapter.Execute(t, fixture)
					if fixture.Kind == "denied" {
						first := assertP22LiveDenial(t, response, false)
						denialStops = p22ResultPath(first, "error.details.stop_ids")
						denialAttention = p22ResultPath(first, "error.details.attention_ids")
						replay := assertP22LiveDenial(t, adapter.Execute(t, fixture), true)
						replay["error"].(map[string]any)["details"].(map[string]any)["replayed"] = false
						if !p22SameResult(first, replay) {
							t.Fatalf("committed denial replay changed identities: %#v %#v", first, replay)
						}
						return
					}
					assertP22ConformanceResponse(t, fixture, response)
					if response.ExitCode != 0 {
						return
					}
					var original map[string]any
					if err := json.Unmarshal([]byte(response.Stdout), &original); err != nil {
						t.Fatal(err)
					}
					result := original["result"].(map[string]any)
					if fixture.Operation == "stop.raise" && priorRun != "" {
						stop := result["stop"].(map[string]any)
						attention := result["attention"].(map[string]any)
						if p22ResultPath(stop, "subject.id") != priorRun || !reflect.DeepEqual(attention["stop_ids"], []any{stop["id"]}) {
							t.Fatalf("Stop/Attention linkage = %#v; run=%q", result, priorRun)
						}
					}
					if fixture.Operation == "stop.resolve" && p22ResultPath(result, "stop.resolution.event_id") != p22ResultPath(original, "event.id") {
						t.Fatalf("Stop disposition does not bind its Event: %#v", original)
					}
					if fixture.Operation == "stop.list" && denialStops != nil {
						if !reflect.DeepEqual([]any{p22ResultPath(result, "items.0.id")}, denialStops) {
							t.Fatalf("listed Stop differs from committed denial: %#v", result)
						}
					}
					if fixture.Operation == "attention.list" && denialAttention != nil {
						if !reflect.DeepEqual([]any{p22ResultPath(result, "items.0.id")}, denialAttention) || !reflect.DeepEqual(p22ResultPath(result, "items.0.stop_ids"), denialStops) {
							t.Fatalf("listed Attention differs from committed denial: %#v", result)
						}
					}

					if fixture.Operation == "run.recover" {
						run := result["run"].(map[string]any)
						claim := result["claim"].(map[string]any)
						recovery := result["recovery"].(map[string]any)
						if run["id"] == priorRun || run["predecessor_run_id"] != priorRun || recovery["prior_run_id"] != priorRun || recovery["successor_run_id"] != run["id"] || run["claim_id"] != claim["id"] {
							t.Fatalf("recovery identity relationships = %#v; prior=%q", result, priorRun)
						}
					}
					if run, ok := result["run"].(map[string]any); ok {
						priorRun, _ = run["id"].(string)
					}
					if fixture.Event {
						assertP22Replay(t, fixture, original, adapter.Execute(t, fixture))
					}
				})
			}
		})
	}
}

func p22RecoveryFixtures(t *testing.T) []p22ConformanceCase {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(p22CorpusDir, "manifest.json"))
	if err != nil {
		t.Fatalf("read Phase 22 recovery conformance manifest: %v", err)
	}
	var manifest struct {
		Cases []string `json:"cases"`
	}
	if err := json.Unmarshal(contents, &manifest); err != nil {
		t.Fatalf("decode Phase 22 recovery conformance manifest: %v", err)
	}
	if len(manifest.Cases) != 41 {
		t.Fatalf("Phase 22 recovery conformance fixture count = %d, want 41", len(manifest.Cases))
	}
	fixtures := make([]p22ConformanceCase, 0, len(manifest.Cases))
	for _, path := range manifest.Cases {
		contents, err := os.ReadFile(filepath.Join(p22CorpusDir, path))
		if err != nil {
			t.Fatalf("read Phase 22 recovery fixture %s: %v", path, err)
		}
		var fixture p22ConformanceCase
		if err := json.Unmarshal(contents, &fixture); err != nil {
			t.Fatalf("decode Phase 22 recovery fixture %s: %v", path, err)
		}
		if fixture.Name == "" || fixture.Operation == "" || fixture.Kind == "" {
			t.Fatalf("Phase 22 recovery fixture %s is incomplete: %#v", path, fixture)
		}
		fixtures = append(fixtures, fixture)
	}
	return fixtures
}

func assertP22ConformanceResponse(t *testing.T, fixture p22ConformanceCase, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != fixture.ExitCode {
		t.Fatalf("%s %s exit = %d, want %d; stdout=%q stderr=%q", fixture.Operation, fixture.Kind, response.ExitCode, fixture.ExitCode, response.Stdout, response.Stderr)
	}
	if fixture.Kind == "accepted" {
		if response.Stderr != "" {
			t.Fatalf("%s accepted stderr = %q, want empty", fixture.Operation, response.Stderr)
		}
		var envelope map[string]any
		if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
			t.Fatalf("decode %s accepted stdout: %v; stdout=%q", fixture.Operation, err, response.Stdout)
		}
		if envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != true || envelope["operation"] != fixture.Operation {
			t.Fatalf("%s accepted envelope = %#v", fixture.Operation, envelope)
		}
		result, ok := envelope["result"].(map[string]any)
		if !ok {
			t.Fatalf("%s accepted result = %#v, want object", fixture.Operation, envelope["result"])
		}
		for path, want := range fixture.Expected {
			got := p22ResultPath(result, path)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("%s result %s = %#v, want %#v", fixture.Operation, path, got, want)
			}
		}
		for _, field := range fixture.ResultFields {
			if _, ok := result[field]; !ok {
				t.Fatalf("%s accepted result missing %q: %#v", fixture.Operation, field, result)
			}
		}
		if fixture.Event {
			if replayed, ok := result["replayed"].(bool); !ok || replayed {
				t.Fatalf("%s first delivery replayed = %#v, want boolean false", fixture.Operation, result["replayed"])
			}
		}
		for _, field := range []string{"run", "claim", "observation", "checkpoint", "recovery", "reconciliation"} {
			if value, exists := result[field]; exists {
				if _, ok := value.(map[string]any); !ok {
					t.Fatalf("%s %s = %#v, want typed object", fixture.Operation, field, value)
				}
			}
		}
		if fixture.Event {
			if _, ok := envelope["event"].(map[string]any); !ok {
				t.Fatalf("%s accepted Event = %#v, want Event object", fixture.Operation, envelope["event"])
			}
		} else if envelope["event"] != nil {
			t.Fatalf("%s accepted Event = %#v, want null", fixture.Operation, envelope["event"])
		}
		return
	}
	if response.Stdout != "" {
		t.Fatalf("%s unknown-field stdout = %q, want empty", fixture.Operation, response.Stdout)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode %s unknown-field stderr: %v; stderr=%q", fixture.Operation, err, response.Stderr)
	}
	errorValue, ok := envelope["error"].(map[string]any)
	if !ok || envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != fixture.Operation || errorValue["code"] != fixture.FailureCode {
		t.Fatalf("%s unknown-field envelope = %#v", fixture.Operation, envelope)
	}
}

func p22FixtureResponse(t *testing.T, envelope map[string]any, stdout bool, exitCode int) p22ConformanceResponse {
	t.Helper()
	encoded, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshal fake recovery response: %v", err)
	}
	value := string(encoded) + "\n"
	if stdout {
		return p22ConformanceResponse{Stdout: value, ExitCode: exitCode}
	}
	return p22ConformanceResponse{Stderr: value, ExitCode: exitCode}
}

type p22FIDRecoveryAdapter struct {
	activation                                                     map[string]any
	registrationID                                                 string
	registrationVersion                                            any
	binary                                                         string
	workspace                                                      string
	projectID                                                      string
	issueID                                                        string
	actor                                                          string
	repository                                                     map[string]any
	claimID                                                        string
	claimVersion                                                   int
	runID                                                          string
	runVersion                                                     int
	observation                                                    string
	digest                                                         string
	requests                                                       map[string]p22PreparedRequest
	executorID                                                     string
	stopID, attentionID, terminalEvent, dispositionEvent, sampleID string
	attentionVersion                                               int
	measurement                                                    map[string]any
}

type p22PreparedRequest struct {
	args  []string
	input map[string]any
}

func (adapter *p22FIDRecoveryAdapter) Name() string { return "fid" }

func p22NewFIDRecoveryAdapter(t *testing.T) *p22FIDRecoveryAdapter {
	t.Helper()
	return p22NewFIDRecoveryAdapterWithActorType(t, "agent")
}

func p22NewFIDRecoveryAdapterWithActorType(t *testing.T, actorType string) *p22FIDRecoveryAdapter {
	t.Helper()
	return p22NewFIDRecoveryAdapterWithLease(t, actorType, 60)
}

func p22NewFIDRecoveryAdapterWithLease(t *testing.T, actorType string, leaseSeconds int) *p22FIDRecoveryAdapter {
	t.Helper()
	return p22NewFIDRecoveryAdapterWithArtifacts(t, actorType, leaseSeconds, []string{})
}

func p22NewFIDRecoveryAdapterWithArtifacts(t *testing.T, actorType string, leaseSeconds int, artifactPaths []string) *p22FIDRecoveryAdapter {
	t.Helper()
	return p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), actorType, leaseSeconds, artifactPaths, nil)
}

// p22ExecutionFixture declares immutable policy through normal fixture creation.
type p22ExecutionFixture struct {
	Argv                []string
	ExtraBudgets        []any
	HumanObligations    bool
	ProtectedDecisions  []any
	ReviewObligation    bool
	ReviewCyclesMaximum int
	// CoexistingReviewObligation adds a review Obligation under its own ID
	// (OBL-p23-review) beside the gate-evidence Obligation instead of replacing
	// it in place the way ReviewObligation does.
	CoexistingReviewObligation bool
	// ReviewGuardsTerminalSuccess makes the coexisting review Obligation guard
	// issue.transition-terminal-success as well as envelope.complete (DEC-249:
	// an Obligation blocks only the guarded actions it declares). It has no
	// effect without CoexistingReviewObligation.
	ReviewGuardsTerminalSuccess bool
}

// Explicit locations keep generated fixtures' repository identities durable.
// A non-nil action set is used by the independent final admission matrix.
func p22NewFIDRecoveryAdapterAt(t *testing.T, binary, workspace, actorType string, leaseSeconds int, artifactPaths, actions []string, execution ...p22ExecutionFixture) *p22FIDRecoveryAdapter {
	return p22NewFIDRecoveryAdapterWithClaimModeAt(t, binary, workspace, actorType, leaseSeconds, artifactPaths, actions, "mutation", execution...)
}

func p22NewFIDRecoveryAdapterWithClaimModeAt(t *testing.T, binary, workspace, actorType string, leaseSeconds int, artifactPaths, actions []string, claimMode string, execution ...p22ExecutionFixture) *p22FIDRecoveryAdapter {
	t.Helper()
	p22Git(t, workspace, "init", "-q")
	p22Git(t, workspace, "config", "user.email", "p22-corpus@example.test")
	p22Git(t, workspace, "config", "user.name", "Phase 22 corpus")
	p22Write(t, workspace, "base.txt", "base\n")
	p22Write(t, workspace, ".gitignore", "*.json\n")
	p22Git(t, workspace, "add", "base.txt", ".gitignore")
	p22Git(t, workspace, "commit", "-qm", "base")
	init := p21Run(t, binary, workspace, "--project", workspace, "--json", "init", "--prefix", "P22")
	projectID := p21PathString(t, init, "result", "project", "id")
	issue := p21Run(t, binary, workspace, "--project", workspace, "--json", "issue", "create", "--kind", "task", "--title", "Phase 22 conformance", "--actor", "agent:developer", "--idempotency-key", "p22-corpus-issue")
	issueID := p21PathString(t, issue, "result", "issue", "id")
	actor := "actor-recovery"
	capabilities := []string{"artifact.verify", "attention.resolve", "claim.acquire", "claim.expire", "envelope.activate", "envelope.create", "gate.register", "gate.run", "run.abandon", "run.checkpoint", "run.close", "run.observe", "run.observe-budget", "run.recover", "run.start", "stop.raise", "stop.resolve"}
	if actorType == "human" {
		capabilities = append(capabilities, "run.force-recover")
		sort.Strings(capabilities)
	}
	grantActions := capabilities[1:]
	if actions != nil {
		capabilities = append([]string(nil), actions...)
		sort.Strings(capabilities)
		grantActions = capabilities
	}
	activationID := p21Activation(t, binary, workspace, "prepare", "p22-corpus-activation-prepare", map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "p22-corpus", "phase": 21})
	profiles := []any{
		map[string]any{"actor_id": "actor-denied", "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities},
		map[string]any{"actor_id": actor, "actor_type": actorType, "producer": "activated-harness", "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities},
	}
	if requiresAdditionalJudgmentActors(actions) {
		profiles = append(profiles, map[string]any{"actor_id": "actor-reviewer-two", "actor_type": "agent", "producer": "activated-harness", "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities})
	}
	// A fixture whose declared actions include human-outcome.record but whose
	// primary Actor is not itself human-typed needs a dedicated human Actor:
	// recording requires actor_type "human" in the finalized profile, which an
	// "agent" primary Actor can never satisfy regardless of capability.
	if actorType != "human" && containsAction(actions, "human-outcome.record") {
		profiles = append(profiles, map[string]any{"actor_id": "actor-human-outcome", "actor_type": "human", "producer": "activated-harness", "provenance_ref": "activation:p22-corpus", "session_ref": nil, "capabilities": capabilities})
	}
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].(map[string]any)["actor_id"].(string) < profiles[j].(map[string]any)["actor_id"].(string)
	})
	profileDigest := p21Component(t, binary, workspace, activationID, "authority-profile", "p22-corpus-authority", profiles)
	argv := []string{"/bin/true"}
	if len(execution) != 0 && len(execution[0].Argv) != 0 {
		argv = execution[0].Argv
	}
	gateDigest := p21Component(t, binary, workspace, activationID, "gate-definitions", "p22-corpus-gates", []any{map[string]any{"artifact_bytes": 1024, "artifact_paths": artifactPaths, "argv": argv, "credential_names": []string{}, "cwd": ".", "definition_id": "gate-p22-corpus", "definition_version": 1, "environment": []any{}, "excerpt_bytes": 0, "network": "allowed", "output_bytes": 1024, "redact_literals": []string{}, "sandbox": "none", "shell": false, "success_exit_codes": []int{0}, "timeout_ms": 1000, "write": "allowed"}})
	digests := []string{gateDigest, profileDigest}
	sort.Strings(digests)
	p21Activation(t, binary, workspace, "commit", "p22-corpus-activation-commit", map[string]any{"activation_id": activationID, "expected_status": "staged", "component_digests": digests})
	p21Activation(t, binary, workspace, "finalize", "p22-corpus-activation-finalize", map[string]any{"activation_id": activationID, "expected_status": "committed"})
	repository := p22CorpusRepository(t, workspace, projectID)
	registered := p21Run(t, binary, workspace, "--project", workspace, "--json", "gate", "register", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-gate-register", map[string]any{"definition_id": "gate-p22-corpus", "definition_version": 1, "worktree_id": repository["worktree_id"], "scope": map[string]any{"issue_id": issueID}, "repository": repository}), "--idempotency-key", "p22-corpus-gate-register")
	registrationID := p21PathString(t, registered, "result", "registration", "id")
	registrationVersion := p21Path(t, registered, "result", "registration", "version")
	source := map[string]any{"digest": "sha256:" + p22Zeros(64), "display_label": "Phase 22 corpus", "external_id": nil, "path": "project-docs/p22-corpus.md", "repository": repository, "section": nil, "source_type": "project-file"}
	authorityGrants := []any{map[string]any{"actions": grantActions, "actor_selector": map[string]any{"actor_id": actor, "required_capabilities": grantActions}, "scope": map[string]any{"issue_id": issueID}}}
	// The Review journey needs a second, independently authorized Actor. Keep the
	// standard recovery fixture's single-Actor authority otherwise.
	if containsReviewRecordAction(actions) {
		authorityGrants = append(authorityGrants, map[string]any{"actions": grantActions, "actor_selector": map[string]any{"actor_id": "actor-denied", "required_capabilities": grantActions}, "scope": map[string]any{"issue_id": issueID}})
		authorityGrants = append(authorityGrants, map[string]any{"actions": grantActions, "actor_selector": map[string]any{"actor_id": "actor-reviewer-two", "required_capabilities": grantActions}, "scope": map[string]any{"issue_id": issueID}})
		sort.Slice(authorityGrants, func(i, j int) bool {
			left := authorityGrants[i].(map[string]any)["actor_selector"].(map[string]any)["actor_id"].(string)
			right := authorityGrants[j].(map[string]any)["actor_selector"].(map[string]any)["actor_id"].(string)
			return left < right
		})
	}
	obligation := map[string]any{"id": "OBL-p22-corpus", "type": "gate-evidence", "applicability": "required", "guarded_actions": []string{"envelope.complete", "issue.transition-terminal-success"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": map[string]any{"id": registrationID, "version": registrationVersion}, "accepted_evidence_types": []string{"receipt"}, "accepted_trust_sources": []string{"fid-executed"}, "freshness_target": "repository-tree", "artifact_policy": "required-available", "independence": "fid-executed", "required_lane": nil, "source": source}
	if len(execution) != 0 && execution[0].ReviewObligation {
		obligation = map[string]any{"id": "OBL-p22-corpus", "type": "review", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "review_profile": "review-default", "accepted_evidence_types": []string{"review"}, "accepted_trust_sources": []string{"activated-actor"}, "freshness_target": "repository-tree", "artifact_policy": "metadata-sufficient", "independence": "distinct-actor", "required_lane": nil, "source": source}
	}
	obligations := []any{obligation}
	if len(execution) != 0 && execution[0].CoexistingReviewObligation {
		reviewGuardedActions := []string{"envelope.complete"}
		if execution[0].ReviewGuardsTerminalSuccess {
			reviewGuardedActions = []string{"envelope.complete", "issue.transition-terminal-success"}
		}
		obligations = append(obligations, map[string]any{"id": "OBL-p23-review", "type": "review", "applicability": "required", "guarded_actions": reviewGuardedActions, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "review_profile": "review-default", "accepted_evidence_types": []string{"review"}, "accepted_trust_sources": []string{"activated-actor"}, "freshness_target": "repository-tree", "artifact_policy": "metadata-sufficient", "independence": "distinct-actor", "required_lane": nil, "source": source})
	}
	if len(execution) != 0 && execution[0].HumanObligations {
		for _, outcomeClass := range []string{"release", "security"} {
			obligations = append(obligations, map[string]any{"id": "OBL-human-" + outcomeClass, "type": "human-outcome", "applicability": "required", "guarded_actions": []string{"envelope.complete"}, "subject": map[string]any{"type": "issue", "id": issueID}, "registration": nil, "outcome_class": outcomeClass, "accepted_evidence_types": []string{"human-outcome"}, "accepted_trust_sources": []string{"activated-human"}, "freshness_target": "none", "artifact_policy": "metadata-sufficient", "independence": "human", "required_lane": nil, "source": source})
		}
	}
	sort.Slice(obligations, func(i, j int) bool {
		return obligations[i].(map[string]any)["id"].(string) < obligations[j].(map[string]any)["id"].(string)
	})
	envelopeRequest := map[string]any{"outcome": "Recover one interrupted Run", "non_goals": []string{}, "risk_result": map[string]any{"source": source, "result_digest": "sha256:" + p22Ones(64)}, "repository_scope": repository, "sources": []any{source}, "primary_issue_id": issueID, "related_issue_ids": []string{}, "authority_grants": authorityGrants, "obligations": obligations, "budgets": []any{map[string]any{"name": "attempts", "scope": "envelope", "unit": "attempts", "maximum": 3, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}, map[string]any{"name": "tokens", "scope": "envelope", "unit": "input-tokens", "maximum": 10, "warning_threshold": nil, "exhaustion_action": "budget-exhausted", "consumed": nil, "remaining": nil}}}
	if len(execution) != 0 {
		if execution[0].ProtectedDecisions != nil {
			envelopeRequest["protected_decisions"] = execution[0].ProtectedDecisions
		}
		budgets := append(envelopeRequest["budgets"].([]any), execution[0].ExtraBudgets...)
		if execution[0].ReviewCyclesMaximum != 0 {
			budgets = append(budgets, map[string]any{"name": "review-cycles", "scope": "envelope", "unit": "review-cycles", "maximum": execution[0].ReviewCyclesMaximum, "warning_threshold": nil, "exhaustion_action": "review-rework-exhausted", "consumed": nil, "remaining": nil})
		}
		sort.Slice(budgets, func(i, j int) bool {
			return budgets[i].(map[string]any)["name"].(string) < budgets[j].(map[string]any)["name"].(string)
		})
		envelopeRequest["budgets"] = budgets
	}
	envelopeInput := p21Input(t, workspace, "p22-corpus-envelope", envelopeRequest)
	created := p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "create", "--actor", actor, "--input", envelopeInput, "--idempotency-key", "p22-corpus-envelope-create")
	envelopeID := p21PathString(t, created, "result", "envelope", "id")
	p21Run(t, binary, workspace, "--project", workspace, "--json", "envelope", "activate", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-envelope-activate", map[string]any{"envelope_id": envelopeID, "expected_revision": 1}), "--idempotency-key", "p22-corpus-envelope-activate")
	claim := p21Run(t, binary, workspace, "--project", workspace, "--json", "claim", "acquire", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-claim", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "worktree_id": repository["worktree_id"], "scope": map[string]any{"paths": []string{"base.txt"}}, "mode": claimMode, "ttl_seconds": leaseSeconds}), "--idempotency-key", "p22-corpus-claim-acquire")
	claimID := p21PathString(t, claim, "result", "claim", "id")
	run := p21Run(t, binary, workspace, "--project", workspace, "--json", "run", "start", "--actor", actor, "--input", p21Input(t, workspace, "p22-corpus-run", map[string]any{"issue_id": issueID, "envelope_id": envelopeID, "envelope_revision": 1, "claim_id": claimID, "claim_version": 1, "executor": map[string]any{"id": "executor-corpus-prior", "source": "activated-harness"}, "obligation_id": nil, "repository": repository}), "--idempotency-key", "p22-corpus-run-start")
	return &p22FIDRecoveryAdapter{activation: map[string]any{"id": activationID}, binary: binary, workspace: workspace, projectID: projectID, issueID: issueID, actor: actor, repository: repository, claimID: claimID, claimVersion: 1, runID: p21PathString(t, run, "result", "run", "id"), runVersion: 1, executorID: "executor-corpus-prior", measurement: source, registrationID: registrationID, registrationVersion: registrationVersion}
}

func p22CorpusRepository(t *testing.T, workspace, projectID string) map[string]any {
	t.Helper()
	observation, err := project.ObserveRepository(workspace)
	if err != nil {
		t.Fatalf("observe corpus repository: %v", err)
	}
	return map[string]any{"base_oid": observation.HeadOID, "capture_source": "git", "dirty": observation.Dirty, "dirty_digest": observation.DirtyDigest, "head_oid": observation.HeadOID, "project_id": projectID, "tree_oid": observation.TreeOID, "worktree_id": observation.WorktreeID}
}

func (adapter *p22FIDRecoveryAdapter) Execute(t *testing.T, fixture p22ConformanceCase) p22ConformanceResponse {
	t.Helper()
	args, input := adapter.p22Args(t, fixture)
	if input != nil {
		args = append(args, "--input", p21Input(t, adapter.workspace, fixture.Name, input), "--idempotency-key", "p22-corpus-"+fixture.Name)
	}
	stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	} else {
		response.ExitCode = 0
	}
	if response.ExitCode == 0 {
		adapter.p22Advance(t, fixture.Operation, stdout)
	}
	return response
}

func (adapter *p22FIDRecoveryAdapter) p22Args(t *testing.T, fixture p22ConformanceCase) ([]string, map[string]any) {
	t.Helper()
	if request, ok := adapter.requests[fixture.Name]; ok {
		return request.args, request.input
	}
	args, input := adapter.p22BuildArgs(t, fixture)
	for key, value := range fixture.Input {
		input[key] = value
	}
	if adapter.requests == nil {
		adapter.requests = map[string]p22PreparedRequest{}
	}
	adapter.requests[fixture.Name] = p22PreparedRequest{args: args, input: input}
	return args, input
}

func (adapter *p22FIDRecoveryAdapter) p22BuildArgs(t *testing.T, fixture p22ConformanceCase) ([]string, map[string]any) {
	t.Helper()
	if strings.HasPrefix(fixture.Operation, "activation.") {
		return adapter.p22ActivationArgs(t, fixture)
	}
	args := []string{"--project", adapter.workspace, "--json"}
	unknown := fixture.Kind == "unknown-field"
	var input map[string]any
	switch fixture.Operation {
	case "run.observe":
		args = append(args, "run", "observe", "--actor", adapter.actor)
		input = map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "liveness": map[string]any{"executor_id": adapter.executorID, "state": "terminated", "source_ref": "process:p22-corpus", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}}
	case "run.inspect":
		args = append(args, "run", "inspect", adapter.runID)
		if unknown {
			args = append(args, "--unknown-field")
		}
		return args, nil
	case "run.checkpoint":
		args = append(args, "run", "checkpoint", "--actor", adapter.actor)
		input = map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "reconciliation_digest": adapter.digest}
	case "claim.expire":
		args = append(args, "claim", "expire", "--actor", adapter.actor)
		input = map[string]any{"claim_id": adapter.claimID, "expected_version": adapter.claimVersion}
	case "run.recover":
		args = append(args, "run", "recover", "--actor", adapter.actor)
		input = map[string]any{"prior_run_id": adapter.runID, "expected_run_version": adapter.runVersion, "expected_claim_version": adapter.claimVersion, "observation_id": adapter.observation, "reconciliation_digest": adapter.digest, "adopted_paths": []string{}, "authority": map[string]any{"mode": "ordinary", "actor_id": adapter.actor, "reason": "trusted termination", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-corpus-successor", "source": "activated-harness"}, "ttl_seconds": 60}
	case "stop.list", "attention.list":
		group := strings.Split(fixture.Operation, ".")[0]
		args = append(args, group, "list", "--subject", "run:"+adapter.runID, "--limit", "10")
		if unknown {
			args = append(args, "--unknown-field")
		}
		return args, nil
	case "stop.raise":
		args = append(args, "stop", "raise", "--actor", adapter.actor)
		subject := map[string]any{"type": "run", "id": adapter.runID}
		input = map[string]any{"subject": subject, "code": "operator-stop", "blocked_actions": []string{"gate.run"}, "details": map[string]any{"reason": "conformance operator intervention", "reference": "policy:conformance"}, "resolution_condition": map[string]any{"type": "work-ended", "subject": subject}, "resolver_capability": "stop.resolve"}
	case "stop.resolve":
		args = append(args, "stop", "resolve", "--actor", adapter.actor)
		input = map[string]any{"stop_id": adapter.stopID, "fact_event_id": adapter.terminalEvent}
	case "attention.resolve":
		args = append(args, "attention", "resolve", "--actor", adapter.actor)
		input = map[string]any{"attention_id": adapter.attentionID, "expected_version": adapter.attentionVersion, "fact_event_ids": []string{adapter.dispositionEvent}}
	case "run.observe-budget":
		args = append(args, "run", "observe-budget", "--actor", adapter.actor)
		var previous any
		kind := "baseline"
		if adapter.sampleID != "" {
			previous, kind = adapter.sampleID, "delta"
		}
		input = map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "budget_name": "tokens", "source_sample_id": fixture.Name, "sample_kind": kind, "previous_sample_id": previous, "value": 3, "measurement_ref": adapter.measurement, "observed_at": time.Now().UTC().Format(time.RFC3339Nano)}
	case "gate.run":
		args = append(args, "gate", "run", "--actor", adapter.actor)
		input = map[string]any{"run_id": adapter.runID, "registration_id": adapter.registrationID, "registration_version": adapter.registrationVersion, "obligation_id": "OBL-p22-corpus", "repository": adapter.repository}
	case "run.close", "run.abandon":
		args = append(args, "run", strings.TrimPrefix(fixture.Operation, "run."), "--actor", adapter.actor)
		input = map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "observation_id": adapter.observation, "repository": adapter.repository}
		if fixture.Operation == "run.abandon" {
			input["expected_claim_version"] = adapter.claimVersion
			input["authority"] = map[string]any{"mode": "ordinary", "actor_id": adapter.actor, "reason": "intentional conformance cessation", "reference": nil, "grant_ref": nil}
		}
	default:
		t.Fatalf("unsupported Phase 22 conformance operation %q", fixture.Operation)
	}
	if unknown {
		input["unknown_field"] = true
	}
	return args, input
}

func (adapter *p22FIDRecoveryAdapter) p22Advance(t *testing.T, operation, stdout string) {
	t.Helper()
	var envelope map[string]any
	if err := json.Unmarshal([]byte(stdout), &envelope); err != nil {
		t.Fatalf("decode %s accepted response: %v; stdout=%q", operation, err, stdout)
	}
	result := p21Path(t, envelope, "result").(map[string]any)
	if result["replayed"] == true {
		return
	}
	if activation, ok := result["activation"].(map[string]any); ok {
		adapter.activation = activation
		return
	}
	switch operation {
	case "run.observe":
		adapter.observation = p21PathString(t, envelope, "result", "observation", "id")
		adapter.runVersion = int(p21Path(t, envelope, "result", "run", "version").(float64))
	case "run.inspect":
		adapter.digest = p21PathString(t, envelope, "result", "reconciliation", "digest")
	case "run.checkpoint":
		adapter.runVersion = int(p21Path(t, envelope, "result", "run", "version").(float64))
		if adapter.binary == "" {
			return
		}
		inspection := adapter.p22Run(t, "run", "inspect", adapter.runID)
		if inspection.ExitCode != 0 {
			t.Fatalf("inspect accepted checkpoint: %#v", inspection)
		}
		adapter.p22Advance(t, "run.inspect", inspection.Stdout)
	case "claim.expire":
		adapter.claimVersion = int(p21Path(t, envelope, "result", "claim", "version").(float64))
	case "run.recover":
		if p21PathString(t, envelope, "result", "recovery", "prior_run_id") != adapter.runID {
			t.Fatalf("recovery response = %#v, want prior Run %q", result, adapter.runID)
		}
		adapter.runID = p21PathString(t, envelope, "result", "run", "id")
		adapter.runVersion = p22Number(p21Path(t, envelope, "result", "run", "version"))
		adapter.claimID = p21PathString(t, envelope, "result", "claim", "id")
		adapter.claimVersion = p22Number(p21Path(t, envelope, "result", "claim", "version"))
		adapter.executorID = "executor-corpus-successor"
	case "run.observe-budget":
		adapter.sampleID = p21PathString(t, envelope, "result", "sample", "id")
		adapter.runVersion++
	case "stop.raise":
		adapter.stopID = p21PathString(t, envelope, "result", "stop", "id")
		adapter.attentionID = p21PathString(t, envelope, "result", "attention", "id")
		adapter.attentionVersion = p22Number(p21Path(t, envelope, "result", "attention", "version"))
	case "run.close", "run.abandon":
		adapter.terminalEvent = p21PathString(t, envelope, "event", "id")
		adapter.runVersion = p22Number(p21Path(t, envelope, "result", "run", "version"))
	case "stop.resolve":
		adapter.dispositionEvent = p21PathString(t, envelope, "event", "id")
	case "attention.resolve":
		adapter.attentionVersion = p22Number(p21Path(t, envelope, "result", "attention", "version"))
	}
}

func TestP22RealBinaryRecoveryFailureCorpus(t *testing.T) {
	for _, fixture := range []struct {
		name, operation, code string
		exit                  int
		run                   func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse
	}{
		{"not-found", "run.inspect", "run-not-found", 3, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			return adapter.p22Run(t, "run", "inspect", "RUN-missing")
		}},
		{"authority", "run.observe", "authority-denied", 5, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			return adapter.p22RunInput(t, "run", "observe", "actor-denied", map[string]any{"run_id": adapter.runID, "expected_version": 1, "liveness": map[string]any{"executor_id": "executor-corpus-prior", "state": "terminated", "source_ref": "process:p22-authority", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
		}},
		{"liveness", "run.recover", "executor-still-live", 4, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			adapter.p22Observe(t, "live")
			adapter.p22Expire(t)
			return adapter.p22Recover(t)
		}},
		{"revision", "run.observe", "run-version-conflict", 4, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			return adapter.p22RunInput(t, "run", "observe", adapter.actor, map[string]any{"run_id": adapter.runID, "expected_version": 99, "liveness": map[string]any{"executor_id": "executor-corpus-prior", "state": "terminated", "source_ref": "process:p22-revision", "observed_at": time.Now().UTC().Format(time.RFC3339Nano)}})
		}},
		{"conflict", "run.checkpoint", "recovery-state-conflict", 4, func(t *testing.T, adapter *p22FIDRecoveryAdapter) p22ConformanceResponse {
			adapter.p22Observe(t, "terminated")
			return adapter.p22RunInput(t, "run", "checkpoint", adapter.actor, map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "reconciliation_digest": "sha256:" + p22Zeros(64)})
		}},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			response := fixture.run(t, p22NewFIDRecoveryAdapter(t))
			assertP22ProcessFailure(t, fixture.operation, fixture.code, fixture.exit, response)
		})
	}
}

func (adapter *p22FIDRecoveryAdapter) p22Run(t *testing.T, command ...string) p22ConformanceResponse {
	t.Helper()
	args := append([]string{"--project", adapter.workspace, "--json"}, command...)
	stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func (adapter *p22FIDRecoveryAdapter) p22RunInput(t *testing.T, group, command, actor string, input map[string]any) p22ConformanceResponse {
	t.Helper()
	args := []string{"--project", adapter.workspace, "--json", group, command, "--actor", actor, "--input", p21Input(t, adapter.workspace, "p22-failure-"+command, input), "--idempotency-key", "p22-failure-" + command + "-" + actor}
	stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, args...)
	response := p22ConformanceResponse{Stdout: stdout, Stderr: stderr}
	if err != nil {
		response.ExitCode = processExitCode(t, err)
	}
	return response
}

func (adapter *p22FIDRecoveryAdapter) p22Observe(t *testing.T, state string) {
	t.Helper()
	response := adapter.p22RunInput(t, "run", "observe", adapter.actor, map[string]any{"run_id": adapter.runID, "expected_version": adapter.runVersion, "liveness": map[string]any{"executor_id": "executor-corpus-prior", "state": state, "source_ref": "process:p22-failure", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}})
	if response.ExitCode != 0 {
		t.Fatalf("observe %s = exit %d stdout=%q stderr=%q", state, response.ExitCode, response.Stdout, response.Stderr)
	}
	adapter.p22Advance(t, "run.observe", response.Stdout)
	response = adapter.p22Run(t, "run", "inspect", adapter.runID)
	if response.ExitCode != 0 {
		t.Fatalf("inspect after observe = exit %d stdout=%q stderr=%q", response.ExitCode, response.Stdout, response.Stderr)
	}
	adapter.p22Advance(t, "run.inspect", response.Stdout)
}

func (adapter *p22FIDRecoveryAdapter) p22Expire(t *testing.T) {
	t.Helper()
	response := adapter.p22RunInput(t, "claim", "expire", adapter.actor, map[string]any{"claim_id": adapter.claimID, "expected_version": adapter.claimVersion})
	if response.ExitCode != 0 {
		t.Fatalf("expire Claim = exit %d stdout=%q stderr=%q", response.ExitCode, response.Stdout, response.Stderr)
	}
	adapter.p22Advance(t, "claim.expire", response.Stdout)
}

func (adapter *p22FIDRecoveryAdapter) p22Recover(t *testing.T) p22ConformanceResponse {
	t.Helper()
	return adapter.p22RunInput(t, "run", "recover", adapter.actor, map[string]any{"prior_run_id": adapter.runID, "expected_run_version": adapter.runVersion, "expected_claim_version": adapter.claimVersion, "observation_id": adapter.observation, "reconciliation_digest": adapter.digest, "adopted_paths": []string{}, "authority": map[string]any{"mode": "ordinary", "actor_id": adapter.actor, "reason": "failure corpus", "reference": nil, "grant_ref": nil}, "executor": map[string]any{"id": "executor-corpus-failure-successor", "source": "activated-harness"}, "ttl_seconds": 60})
}

func assertP22ProcessFailure(t *testing.T, operation, code string, exit int, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != exit || response.Stdout != "" {
		t.Fatalf("%s failure = exit %d stdout=%q stderr=%q, want exit %d and no stdout", operation, response.ExitCode, response.Stdout, response.Stderr, exit)
	}
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatalf("decode %s failure: %v; stderr=%q", operation, err, response.Stderr)
	}
	errorValue, ok := envelope["error"].(map[string]any)
	if !ok || envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != operation || errorValue["code"] != code {
		t.Fatalf("%s failure envelope = %#v", operation, envelope)
	}
	if fmt.Sprint(errorValue["code"]) != code {
		t.Fatalf("%s failure code = %#v, want %q", operation, errorValue["code"], code)
	}
}
