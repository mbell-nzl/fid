//go:build integration

package integration

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestP22RecoveryHumanResultMatchesDurableJSON(t *testing.T) {
	for _, mode := range []string{"ordinary", "force"} {
		t.Run(mode, func(t *testing.T) {
			adapter, fixture := p22HumanRecoveryFixture(t, mode)
			response := adapter.Execute(t, fixture)
			if response.ExitCode != 0 {
				t.Fatalf("recover %s: %#v", mode, response)
			}
			var envelope map[string]any
			if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
				t.Fatal(err)
			}
			result := envelope["result"].(map[string]any)
			before := p21JournalEventCount(t, adapter.workspace)
			args, input := adapter.p22Args(t, fixture)
			humanArgs := p22WithoutJSON(args)
			humanArgs = append(humanArgs, "--input", p21Input(t, adapter.workspace, fixture.Name, input), "--idempotency-key", "p22-corpus-"+fixture.Name)
			stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, humanArgs...)
			if err != nil || stderr != "" {
				t.Fatalf("Human recovery replay: err=%v stderr=%q", err, stderr)
			}
			p22AssertHumanRecovery(t, stdout, result["recovery"].(map[string]any), result["run"].(map[string]any)["budgets"].([]any))
			if after := p21JournalEventCount(t, adapter.workspace); after != before {
				t.Fatalf("Human replay appended Events: %d -> %d", before, after)
			}
		})
	}
}

func TestP22ResumeHumanMatchesEligibleAndStoppedRecoveryState(t *testing.T) {
	adapter, fixture := p22HumanRecoveryFixture(t, "ordinary")
	response := adapter.Execute(t, fixture)
	if response.ExitCode != 0 {
		t.Fatalf("recover: %#v", response)
	}
	var recovered map[string]any
	if err := json.Unmarshal([]byte(response.Stdout), &recovered); err != nil {
		t.Fatal(err)
	}
	issueID := p21PathString(t, recovered, "result", "run", "issue_id")
	p22ExecuteHumanSetup(t, adapter, p22ConformanceCase{Name: "human-budget", Operation: "run.observe-budget", Input: map[string]any{"value": 0}})
	for _, blocked := range []bool{false, true} {
		if blocked {
			p22ExecuteHumanSetup(t, adapter, p22ConformanceCase{Name: "human-stop", Operation: "stop.raise"})
		}
		before := p21JournalEventCount(t, adapter.workspace)
		resume := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "resume", issueID, "--repository", adapter.workspace)
		stdout, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, "--project", adapter.workspace, "resume", issueID, "--repository", adapter.workspace)
		if err != nil || stderr != "" {
			t.Fatalf("Human resume: err=%v stderr=%q", err, stderr)
		}
		result := p21Path(t, recovered, "result").(map[string]any)
		p22AssertHumanRecovery(t, stdout, result["recovery"].(map[string]any), p21Path(t, resume, "result", "resume", "run", "budgets").([]any))
		p22AssertHumanResume(t, stdout, p21Path(t, resume, "result", "resume").(map[string]any), blocked)
		if after := p21JournalEventCount(t, adapter.workspace); after != before {
			t.Fatalf("resume appended Events: %d -> %d", before, after)
		}
	}
}

func TestP22ResumeHumanDistinguishesElapsedLeaseFromRecordedStatus(t *testing.T) {
	adapter := p22NewFIDRecoveryAdapterWithLease(t, "agent", 5)
	run := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "get", adapter.runID)
	issueID := p21PathString(t, run, "result", "run", "issue_id")
	before := p21JournalEventCount(t, adapter.workspace)
	initial := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "resume", issueID)
	expires := p21PathString(t, initial, "result", "resume", "claim", "expires_at")
	deadline, err := time.Parse(time.RFC3339Nano, expires)
	if err != nil {
		t.Fatal(err)
	}
	// Wait for the actual stored expiry, not a guessed scheduling interval.
	if remaining := time.Until(deadline); remaining > 0 {
		<-time.After(remaining + time.Millisecond)
	}
	after := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "resume", issueID)
	claim := p21Path(t, after, "result", "resume", "claim").(map[string]any)
	if claim["status"] != "active" || claim["effective_status"] != "expired" || claim["expires_at"] != expires {
		t.Fatalf("elapsed lease rewrote durable Claim facts: %#v", claim)
	}
	human, stderr, err := closureRunCommand(adapter.workspace, adapter.binary, "--project", adapter.workspace, "resume", issueID)
	want := fmt.Sprintf("Claim %s: effective expired; recorded active; expires %s", claim["id"], expires)
	if err != nil || stderr != "" || !strings.Contains(human, want) {
		t.Fatalf("Human lease differs from JSON: err=%v stderr=%q human=%q want=%q", err, stderr, human, want)
	}
	if count := p21JournalEventCount(t, adapter.workspace); count != before {
		t.Fatalf("expiry read appended an Event: %d -> %d", before, count)
	}
}

func p22AssertHumanResume(t *testing.T, human string, resume map[string]any, blocked bool) {
	t.Helper()
	for _, key := range []string{"stops", "attention", "budgets", "recoveries"} {
		if _, ok := resume[key].([]any); !ok {
			t.Fatalf("resume %s is not an explicit list: %#v", key, resume)
		}
	}
	for _, value := range []any{p22ResultPath(resume, "issue.id"), p22ResultPath(resume, "envelope.id"), p22ResultPath(resume, "run.status"), p22ResultPath(resume, "claim.effective_status")} {
		if value == nil || !strings.Contains(human, fmt.Sprint(value)) {
			t.Errorf("Human resume missing JSON fact %v: %q", value, human)
		}
	}
	findings, ok := resume["findings"].(map[string]any)
	if !ok || findings["available"] != false || !strings.Contains(human, "Findings: unavailable") {
		t.Errorf("Findings availability contradicts JSON: %#v / %q", findings, human)
	}
	actions, ok := resume["next_actions"].([]any)
	if !ok {
		t.Fatalf("next_actions is not a list: %#v", resume)
	}
	foundGate := false
	for _, item := range actions {
		action, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("action is not typed: %#v", item)
		}
		reasons := action["reasons"].([]any)
		subject := action["subject"].(map[string]any)
		state := "ready"
		if len(reasons) != 0 {
			state = "blocked"
		}
		want := fmt.Sprintf("Action %s %s:%s: %s", action["operation"], subject["type"], subject["id"], state)
		if !strings.Contains(human, want) {
			t.Errorf("Human action contradicts JSON %q: %q", want, human)
		}
		for _, required := range action["required_capabilities"].([]any) {
			if !strings.Contains(human, fmt.Sprint(required)) {
				t.Errorf("Human action missing required capability %v", required)
			}
		}
		for _, reason := range reasons {
			if code := reason.(map[string]any)["code"]; !strings.Contains(human, fmt.Sprint(code)) {
				t.Errorf("Human action missing reason %v", code)
			}
		}
		if action["operation"] == "gate.run" {
			foundGate = true
			if (len(reasons) != 0) != blocked {
				t.Errorf("gate eligibility blocked=%t, want %t: %#v", len(reasons) != 0, blocked, action)
			}
		}
	}
	if !foundGate {
		t.Error("resume omitted the current Gate action")
	}
	if blocked {
		for _, key := range []string{"stops", "attention"} {
			items := resume[key].([]any)
			if len(items) == 0 {
				t.Errorf("blocked resume omitted %s", key)
			}
			for _, item := range items {
				if id := item.(map[string]any)["id"]; !strings.Contains(human, fmt.Sprint(id)) {
					t.Errorf("Human resume omitted %s %v", key, id)
				}
			}
		}
	}
}

func p22HumanRecoveryFixture(t *testing.T, mode string) (*p22FIDRecoveryAdapter, p22ConformanceCase) {
	t.Helper()
	actorType := "agent"
	if mode == "force" {
		actorType = "human"
	}
	adapter := p22NewFIDRecoveryAdapterWithActorType(t, actorType)
	p22Write(t, adapter.workspace, "base.txt", "adopted work needs independent proof\n")
	liveness := "terminated"
	if mode == "force" {
		liveness = "live"
	}
	p22ExecuteHumanSetup(t, adapter, p22ConformanceCase{Name: "human-observe", Operation: "run.observe", Input: map[string]any{"liveness": map[string]any{"executor_id": adapter.executorID, "state": liveness, "source_ref": "process:human-parity", "observed_at": time.Now().UTC().Add(-time.Second).Format(time.RFC3339Nano)}}})
	p22ExecuteHumanSetup(t, adapter, p22ConformanceCase{Name: "human-inspect", Operation: "run.inspect"})
	if mode == "ordinary" {
		p22ExecuteHumanSetup(t, adapter, p22ConformanceCase{Name: "human-expire", Operation: "claim.expire"})
	}
	authority := map[string]any{"mode": mode, "actor_id": adapter.actor, "reason": "recover the observed work", "reference": nil, "grant_ref": nil}
	if mode == "force" {
		run := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "run", "get", adapter.runID)
		authority["reference"] = "HUM-parity-takeover"
		authority["grant_ref"] = map[string]any{"envelope_id": p21Path(t, run, "result", "run", "envelope_id"), "revision": 1, "ordinal": 0}
	}
	return adapter, p22ConformanceCase{Name: "human-recover", Operation: "run.recover", Input: map[string]any{"adopted_paths": []string{"base.txt"}, "authority": authority}}
}

func p22ExecuteHumanSetup(t *testing.T, adapter *p22FIDRecoveryAdapter, fixture p22ConformanceCase) {
	t.Helper()
	if response := adapter.Execute(t, fixture); response.ExitCode != 0 {
		t.Fatalf("Human setup %s: %#v", fixture.Name, response)
	}
}

func p22WithoutJSON(args []string) []string {
	result := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "--json" {
			result = append(result, arg)
		}
	}
	return result
}

func p22AssertHumanRecovery(t *testing.T, human string, recovery map[string]any, budgets []any) {
	t.Helper()
	authority := recovery["authority"].(map[string]any)
	for _, value := range []any{recovery["prior_run_id"], recovery["successor_run_id"], recovery["prior_claim_id"], recovery["successor_claim_id"], authority["mode"], authority["actor_id"], authority["reason"]} {
		if !strings.Contains(human, fmt.Sprint(value)) {
			t.Errorf("Human recovery missing JSON fact %v: %q", value, human)
		}
	}
	if !strings.Contains(human, "unproven work") {
		t.Errorf("Human recovery does not distinguish adoption from proof: %q", human)
	}
	for _, path := range recovery["adopted_paths"].([]any) {
		if !strings.Contains(human, fmt.Sprint(path)) {
			t.Errorf("Human recovery missing adopted path %v: %q", path, human)
		}
	}
	for _, budget := range budgets {
		value := budget.(map[string]any)
		remaining := "unknown"
		if value["remaining"] != nil {
			remaining = fmt.Sprint(value["remaining"])
		}
		want := fmt.Sprintf("Budget %s: remaining %s %s", value["name"], remaining, value["unit"])
		if !strings.Contains(human, want) {
			t.Errorf("Human recovery missing budget fact %q: %q", want, human)
		}
	}
	if authority["mode"] == "force" {
		for _, want := range []string{"known-live takeover", "residual filesystem risk", fmt.Sprint(authority["reference"])} {
			if !strings.Contains(human, want) {
				t.Errorf("Human force recovery missing %q: %q", want, human)
			}
		}
	} else if strings.Contains(human, "known-live takeover") {
		t.Errorf("ordinary recovery mislabeled force takeover: %q", human)
	}
}
