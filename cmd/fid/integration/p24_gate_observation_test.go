//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func p24ObservedGateInput(t *testing.T, a *p22FIDRecoveryAdapter) map[string]any {
	t.Helper()
	return map[string]any{"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "expected_repository": p24ExpectedRepository(p22CorpusRepository(t, a.workspace, a.projectID))}
}

func TestP24GateAndEvaluateUnsafeAndReplay(t *testing.T) {
	t.Parallel()
	for _, state := range []string{"advanced-HEAD", "unmerged", "unborn-HEAD"} {
		t.Run(state, func(t *testing.T) {
			t.Parallel()
			a := p24RunAdapter(t)
			p22AcceptedCall(t, a, "budget", "run.observe-budget", map[string]any{"value": 0})
			input := p24ObservedGateInput(t, a)
			first := p23RunInput(t, a, "gate", "run", a.actor, "original", input)
			if first.ExitCode != 0 {
				t.Fatalf("original: %#v", first)
			}
			switch state {
			case "advanced-HEAD":
				p22Write(t, a.workspace, "base.txt", "advanced\n")
				p22Git(t, a.workspace, "commit", "-qam", "advance")
			case "unmerged":
				p24MakeUnmerged(t, a.workspace)
			case "unborn-HEAD":
				unbornGitHead(t, a.workspace)
			}
			before, events := p22FaultProjection(t, a.workspace), p21JournalEventCount(t, a.workspace)
			replayed := p23RunInput(t, a, "gate", "run", a.actor, "original", input)
			if replayed.ExitCode != 0 || replayed.Stderr != "" || p24GateTerminalResponseBytes(t, replayed.Stdout) != strings.Replace(p24GateTerminalResponseBytes(t, first.Stdout), `"replayed":false`, `"replayed":true`, 1) {
				t.Fatalf("replay: %#v", replayed)
			}
			if !reflect.DeepEqual(p22FaultProjection(t, a.workspace), before) || p21JournalEventCount(t, a.workspace) != events {
				t.Fatal("replay changed persistence")
			}
			if state == "advanced-HEAD" {
				return
			}
			for _, leaf := range []string{"gate.run", "evidence.evaluate"} {
				var response p22ConformanceResponse
				if leaf == "gate.run" {
					response = p23RunInput(t, a, "gate", "run", a.actor, "new-key", input)
				} else {
					response = a.p22Run(t, "evidence", "evaluate", "--input", p21Input(t, a.workspace, "evaluate", map[string]any{"obligation_id": nil}))
				}
				assertP22ProcessFailure(t, leaf, "repository-unavailable", 4, response)
				reason := "not-a-repository"
				if state == "unmerged" {
					reason = "unmerged-paths"
				}
				details := p21Path(t, p23Envelope(t, response.Stderr), "error", "details").(map[string]any)
				if len(details) != 2 || details["reason"] != reason || details["detail"] == "" {
					t.Fatalf("unavailable details: %#v", details)
				}
				if !reflect.DeepEqual(p22FaultProjection(t, a.workspace), before) || p21JournalEventCount(t, a.workspace) != events {
					t.Fatal("unsafe request changed persistence")
				}
			}
		})
	}
}

func TestP24ReceiptFreshnessRealGuardedActions(t *testing.T) {
	t.Parallel()
	for _, roundtrip := range []bool{false, true} {
		t.Run(strconv.FormatBool(roundtrip), func(t *testing.T) {
			t.Parallel()
			a := p24RunAdapter(t)
			p22AcceptedCall(t, a, "budget", "run.observe-budget", map[string]any{"value": 0})
			gate := p23RunInput(t, a, "gate", "run", a.actor, "fresh-gate", p24ObservedGateInput(t, a))
			if gate.ExitCode != 0 {
				t.Fatalf("gate: %#v", gate)
			}
			receipt := p24GateTerminal(t, gate.Stdout)["receipt_id"]
			link := p23RunInput(t, a, "evidence", "link", a.actor, "fresh-link", map[string]any{"obligation_id": "OBL-p22-corpus", "evidence_type": "receipt", "receipt_id": receipt, "review_id": nil, "human_outcome_id": nil, "required_artifact_ids": []string{}})
			if link.ExitCode != 0 {
				t.Fatalf("link: %#v", link)
			}
			original, err := os.ReadFile(filepath.Join(a.workspace, "base.txt"))
			if err != nil {
				t.Fatal(err)
			}
			evaluate := func(fresh string) {
				t.Helper()
				response := a.p22Run(t, "evidence", "evaluate", "--input", p21Input(t, a.workspace, "evaluate", map[string]any{"obligation_id": "OBL-p22-corpus"}))
				if response.ExitCode != 0 {
					t.Fatalf("evaluate: %#v", response)
				}
				value := p21Path(t, p23Envelope(t, response.Stdout), "result", "evaluation").(map[string]any)
				if value["freshness"] != fresh {
					t.Fatalf("freshness: %#v", value)
				}
				if fresh == "stale" && (len(value["reasons"].([]any)) != 1 || value["reasons"].([]any)[0].(map[string]any)["code"] != "evidence-stale") {
					t.Fatalf("stale reason: %#v", value)
				}
			}
			evaluate("fresh")
			for i, state := range []string{"in_progress", "dev_complete", "qa_review"} {
				p21Run(t, a.binary, a.workspace, "--json", "issue", "transition", a.issueID, "--to", state, "--actor", "agent:"+a.actor, "--expected-revision", strconv.Itoa(i+1), "--idempotency-key", state)
			}
			transition := func(key string) p22ConformanceResponse {
				return a.p22Run(t, "issue", "transition", a.issueID, "--to", "qa_pass", "--actor", "agent:"+a.actor, "--expected-revision", "4", "--idempotency-key", key)
			}
			complete := func(key string) p22ConformanceResponse {
				return p23RunInput(t, a, "envelope", "complete", a.actor, key, map[string]any{"envelope_id": p23EnvelopeID(t, a), "expected_revision": 1})
			}
			if roundtrip {
				p22Write(t, a.workspace, "base.txt", "edited\n")
				evaluate("stale")
				denied := transition("stale-terminal")
				assertP22ProcessFailure(t, "issue.transition", "transition-not-permitted", 4, denied)
				if !reflect.DeepEqual(p21Path(t, p23Envelope(t, denied.Stderr), "error", "details", "unsatisfied"), []any{"evidence-stale"}) {
					t.Fatalf("terminal denial lacks freshness reason: %s", denied.Stderr)
				}
				p22Write(t, a.workspace, "base.txt", string(original))
				evaluate("fresh")
			}
			if response := transition("fresh-terminal"); response.ExitCode != 0 {
				t.Fatalf("fresh terminal: %#v", response)
			}
			if roundtrip {
				p22Write(t, a.workspace, "base.txt", "edited\n")
				denied := complete("stale-complete")
				assertP22ProcessFailure(t, "envelope.complete", "obligation-unsatisfied", 4, denied)
				if !strings.Contains(denied.Stderr, `"code":"evidence-stale"`) {
					t.Fatal("completion denial lacks freshness reason")
				}
				p22Write(t, a.workspace, "base.txt", string(original))
				evaluate("fresh")
			}
			if response := complete("fresh-complete"); response.ExitCode != 0 {
				t.Fatalf("fresh completion: %#v", response)
			}
		})
	}
}
