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

	"github.com/mbell-nzl/fid/internal/config"
)

type p24HumanFixture struct {
	Authoring    bool     `json:"authoring,omitempty"`
	Name         string   `json:"name"`
	Operation    string   `json:"operation"`
	Code         string   `json:"code,omitempty"`
	Workspace    string   `json:"workspace"`
	CWD          string   `json:"cwd"`
	HumanArgs    []string `json:"human_args"`
	JSONArgs     []string `json:"json_args"`
	HumanCommand string   `json:"human_command"`
	JSONCommand  string   `json:"json_command"`
	ExpectedExit int      `json:"expected_exit"`
	BinarySHA256 string   `json:"binary_sha256"`
	HumanStdout  string   `json:"human_stdout"`
	HumanStderr  string   `json:"human_stderr"`
	JSONStdout   string   `json:"json_stdout"`
	JSONStderr   string   `json:"json_stderr"`
}
type p24HumanManifest struct {
	Version      int               `json:"version"`
	Binary       string            `json:"binary"`
	BinarySHA256 string            `json:"binary_sha256"`
	Fixtures     []p24HumanFixture `json:"fixtures"`
}

func TestP24HumanFixtures(t *testing.T) {
	t.Parallel()
	root := os.Getenv("FID_P24_HUMAN_FIXTURE_ROOT")
	if root == "" {
		root = t.TempDir()
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		t.Fatal(err)
	}
	fresh, err := os.MkdirTemp(root, "fixtures-")
	if err != nil {
		t.Fatal(err)
	}
	fresh, err = filepath.Abs(fresh)
	if err != nil {
		t.Fatal(err)
	}
	p24GenerateHumanFixtures(t, fresh)
	data, err := os.ReadFile(filepath.Join(fresh, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest p24HumanManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Fixtures) < 21 {
		t.Fatalf("manifest has %d fixtures, want five denials and four views in four states", len(manifest.Fixtures))
	}
	p24AssertHumanManifest(t, manifest)
	foundStart := false
	for _, row := range manifest.Fixtures {
		foundStart = foundStart || row.Name == "take-control-start"
	}
	if !foundStart {
		t.Fatal("manifest missing take-control-start workspace")
	}
	for _, name := range []string{"authoring-run-close", "authoring-decision-reference", "authoring-third-leaf", "judgment-concern.get", "judgment-deviation.get", "judgment-human-outcome.get", "judgment-resume"} {
		found := false
		for _, row := range manifest.Fixtures {
			found = found || row.Name == name
		}
		if !found {
			t.Errorf("manifest missing %s workspace", name)
		}
	}
	t.Logf("Retained Human fixture manifest: %s", filepath.Join(fresh, "manifest.json"))
}

func TestP24StopHumanRealCorpus(t *testing.T) {
	t.Parallel()
	codes, conditions := map[string]bool{}, map[string]bool{}
	for _, row := range p24RealStopRows {
		if codes[row.code] || conditions[row.condition] {
			t.Fatal("duplicate real Stop code or Condition")
		}
		codes[row.code], conditions[row.condition] = true, true
		t.Run(row.code, func(t *testing.T) {
			t.Parallel()
			p24RealStopCase(t, row)
		})
	}
	if len(codes) != 10 || len(conditions) != 10 {
		t.Fatalf("real stop.raise corpus has %d codes and %d Conditions, want 10 of each", len(codes), len(conditions))
	}
	for condition := range p23ConditionWireShapes {
		if !conditions[condition] {
			t.Errorf("missing real Condition %s", condition)
		}
	}
}

type p24RealStopRow struct {
	code, condition, details, next string
}

var p24RealStopRows = []p24RealStopRow{
	{"executor-liveness-unknown", "executor-terminated", `{"run_id":"RUN","executor_id":"executor-corpus-prior","observation_id":null}`, "fid run observe"},
	{"recovery-state-conflict", "run-recovered", `{"run_id":"RUN","observation_digest":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","paths":["base.txt"],"reasons":["changed"]}`, "fid run recover"},
	{"activation-incomplete", "activation-recovered", `{"activation_id":"ACT","point_of_no_return":true,"components":["kit"]}`, "fid activation recover"},
	{"capability-unavailable", "capability-available", `{"capability_id":"evidence.receipts.v1","available":false,"owner_phase":26}`, "wait for the capability evidence.receipts.v1"},
	{"budget-exhausted", "budget-observed", `{"envelope_id":"ENV","budget_name":"tokens","unit":"input-tokens","maximum":10,"consumed":10,"remaining":0}`, "fid run observe-budget"},
	{"operator-stop", "work-ended", `{"reason":"operator approval","reference":"DEC-human"}`, "fid run close or fid run abandon"},
	{"review-rework-exhausted", "review-recorded", `{"envelope_id":"ENV","budget_name":"reviews","unit":"review-cycles","maximum":1,"consumed":1,"remaining":0,"last_review_id":"REV-human"}`, "fid review record"},
	{"budget-data-unavailable", "finding-resolved", `{"envelope_id":"ENV","budget_name":"tokens","unit":"input-tokens","maximum":10,"consumed":null,"remaining":null}`, "fid review record"},
	{"human-outcome-required", "human-outcome-recorded", `{"envelope_id":"ENV","obligation_id":"OBL-human","outcome_class":"security"}`, "fid human-outcome record"},
	{"protected-decision-required", "decision-satisfied", `{"envelope_id":"ENV","decision_class":"security","required_capabilities":["decision.record-reference"],"affected_actions":["envelope.complete"]}`, "fid decision record-reference then fid human-outcome record"},
}

func p24RealStopCase(t *testing.T, row p24RealStopRow) {
	t.Helper()
	var execution p22ExecutionFixture
	if row.condition == "decision-satisfied" {
		execution.ProtectedDecisions = []any{map[string]any{"ordinal": 0, "decision_class": "security", "required_capabilities": []string{"decision.record-reference"}, "affected_actions": []string{"envelope.complete"}}}
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 86400, []string{}, p24FixtureActions, execution)
	subject := map[string]any{"type": "run", "id": a.runID}
	var reviewRun map[string]any
	if row.condition == "review-recorded" || row.condition == "finding-resolved" {
		subject = map[string]any{"type": "envelope", "id": p23EnvelopeID(t, a)}
		p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "corpus-builder")
		claim := p23AcquireClaim(t, a, "actor-denied", "review", "corpus-review")
		reviewRun = p23StartRun(t, a, "actor-denied", claim, "corpus-review")
		if row.condition == "finding-resolved" {
			response := p23ReviewRecord(t, a, "actor-denied", "corpus-finding", p23ReviewInput(a, reviewRun["id"].(string), p23Number(t, reviewRun["version"]), "changes-required", []any{p23Finding(a.issueID, "FND-human", "major", true)}, nil))
			p23AssertReviewSuccess(t, "corpus finding", response)
			p23CloseAndRelease(t, a, "actor-denied", reviewRun["id"].(string), claim["id"].(string), "corpus-first-review")
			claim = p23AcquireClaim(t, a, "actor-denied", "review", "corpus-second-review")
			reviewRun = p23StartRun(t, a, "actor-denied", claim, "corpus-second-review")
		}
	}
	if row.condition == "activation-recovered" || row.condition == "capability-available" {
		subject = map[string]any{"type": "activation", "id": a.activation["id"]}
	}
	condition := p22ConditionFields(row.condition)
	condition["subject"] = subject
	switch row.condition {
	case "budget-observed":
		condition["budget_name"] = "tokens"
	case "capability-available":
		condition["capability_id"] = "evidence.receipts.v1"
	case "finding-resolved":
		condition["finding_id"] = "FND-human"
	case "human-outcome-recorded":
		condition["outcome_class"] = "security"
	case "decision-satisfied":
		condition["decision_class"] = "security"
	}
	detailText := strings.NewReplacer(`"RUN"`, `"`+a.runID+`"`, `"ENV"`, `"`+p23EnvelopeID(t, a)+`"`, `"ACT"`, `"`+a.activation["id"].(string)+`"`).Replace(row.details)
	var details map[string]any
	if err := json.Unmarshal([]byte(detailText), &details); err != nil {
		t.Fatal(err)
	}
	blocked := []string{"gate.run"}
	if strings.HasPrefix(row.code, "budget-") {
		blocked = []string{"gate.run", "run.recover", "run.start"}
	}
	raised := p22AcceptedCall(t, a, "real-corpus-stop", "stop.raise", map[string]any{"subject": subject, "code": row.code, "blocked_actions": blocked, "details": details, "resolution_condition": condition, "resolver_capability": "stop.resolve"})
	stop := p21Path(t, raised, "result", "stop").(map[string]any)
	attention := p21Path(t, raised, "result", "attention").(map[string]any)
	p24AssertRealStopViews(t, a, row, stop, attention)
	fact := p24RealStopFact(t, a, row.condition, reviewRun)
	if p24EventByID(t, a.workspace, fact).EventType == "" {
		t.Fatal("resolution fact is not a real journal Event")
	}
	listed := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "stop", "list")
	stop = p24StopByID(t, listed, stop["id"].(string))
	if stop["resolution"] == nil {
		resolved := p23RunInput(t, a, "stop", "resolve", a.actor, "real-corpus-resolve", map[string]any{"stop_id": stop["id"], "fact_event_id": fact})
		p23AssertReviewSuccess(t, "resolve real corpus Stop", resolved)
		stop = p21Path(t, p23Envelope(t, resolved.Stdout), "result", "stop").(map[string]any)
	} // Judgment Events may have atomically disposed the Stop and Attention.
	if p21PathString(t, stop, "resolution", "fact_event_id") != fact {
		t.Fatalf("wrong qualifying resolution Event: %#v", stop)
	}
	// Before explicit Attention disposition, resume must render its resolved linked Stop.
	attentionList := p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "attention", "list")
	for _, raw := range p21Path(t, attentionList, "result", "items").([]any) {
		item := raw.(map[string]any)
		if item["id"] == attention["id"] {
			attention = item
		}
	}
	p24AssertRealStopViews(t, a, row, stop, attention)
	if attention["status"] == "open" {
		resolved := p23RunInput(t, a, "attention", "resolve", a.actor, "real-corpus-attention", map[string]any{"attention_id": attention["id"], "expected_version": attention["version"], "fact_event_ids": []string{p21PathString(t, stop, "resolution", "event_id")}})
		p23AssertReviewSuccess(t, "resolve real corpus Attention", resolved)
		attention = p21Path(t, p23Envelope(t, resolved.Stdout), "result", "attention").(map[string]any)
	}
	p24AssertRealStopViews(t, a, row, stop, attention)
}

func p24RealStopFact(t *testing.T, a *p22FIDRecoveryAdapter, condition string, reviewRun map[string]any) string {
	t.Helper()
	var response map[string]any
	switch condition {
	case "activation-recovered", "capability-available":
		response = p21Run(t, a.binary, a.workspace, "--project", a.workspace, "--json", "activation", "finalize", "--input", filepath.Join(a.workspace, "activation-finalize.json"), "--idempotency-key", "p22-corpus-activation-finalize")
	case "executor-terminated", "work-ended", "run-recovered":
		response = p22AcceptedCall(t, a, "corpus-terminated", "run.observe", nil)
		if condition == "work-ended" {
			response = p22AcceptedCall(t, a, "corpus-ended", "run.close", nil)
		} else if condition == "run-recovered" {
			p22AcceptedCall(t, a, "corpus-inspect", "run.inspect", nil)
			p22AcceptedCall(t, a, "corpus-failed", "run.close", nil)
			released := p23RunInput(t, a, "claim", "release", a.actor, "corpus-release", map[string]any{"claim_id": a.claimID, "expected_version": a.claimVersion, "reason": "failed"})
			p23AssertReviewSuccess(t, "release failed Run Claim", released)
			a.claimVersion = p22Number(p21Path(t, p23Envelope(t, released.Stdout), "result", "claim", "version"))
			response = p22AcceptedCall(t, a, "corpus-recovered", "run.recover", nil)
		}
	case "budget-observed":
		response = p22AcceptedCall(t, a, "corpus-budget", "run.observe-budget", map[string]any{"value": 0})
	case "review-recorded", "finding-resolved":
		var findings []any
		if condition == "finding-resolved" {
			findings = []any{"FND-human"}
		}
		result := p23ReviewRecord(t, a, "actor-denied", "corpus-review-fact", p23ReviewInput(a, reviewRun["id"].(string), p23Number(t, reviewRun["version"]), "pass", nil, findings))
		p23AssertReviewSuccess(t, "real Review fact", result)
		response = p23Envelope(t, result.Stdout)
	case "human-outcome-recorded", "decision-satisfied":
		if condition == "decision-satisfied" {
			result := p23RecordDecision(t, a, "corpus-security", "security", "approved", []string{"decision.record-reference"}, p23DecisionSource(a, true))
			p23AssertReviewSuccess(t, "real Decision fact", result)
		}
		result := p23RecordHumanOutcome(t, a, "corpus-security", a.actor, "security", "approved")
		p23AssertReviewSuccess(t, "real Human Outcome fact", result)
		response = p23Envelope(t, result.Stdout)
	default:
		t.Fatalf("missing real fact for %s", condition)
	}
	return p21PathString(t, response, "event", "id")
}

func p24StopByID(t *testing.T, envelope map[string]any, id string) map[string]any {
	t.Helper()
	for _, raw := range p21Path(t, envelope, "result", "items").([]any) {
		stop := raw.(map[string]any)
		if stop["id"] == id {
			return stop
		}
	}
	t.Fatalf("Stop %s missing", id)
	return nil
}

func p24AssertRealStopViews(t *testing.T, a *p22FIDRecoveryAdapter, row p24RealStopRow, stop, attention map[string]any) {
	t.Helper()
	labels := map[string]string{"run_id": "Run", "executor_id": "Executor", "observation_id": "Observation", "observation_digest": "Observation digest", "paths": "Paths", "reasons": "Reasons", "activation_id": "Activation", "point_of_no_return": "Point of no return", "components": "Components", "capability_id": "Capability", "available": "Available", "owner_phase": "Owner phase", "envelope_id": "Envelope", "budget_name": "Budget", "unit": "Unit", "maximum": "Maximum", "consumed": "Consumed", "remaining": "Remaining", "reason": "Reason", "reference": "Reference", "last_review_id": "Last Review", "obligation_id": "Obligation", "outcome_class": "Outcome class", "decision_class": "Decision class", "required_capabilities": "Required capabilities", "affected_actions": "Affected actions"}
	for _, args := range [][]string{{"stop", "list"}, {"attention", "list"}, {"resume", a.issueID, "--repository", a.workspace}} {
		before := p21JournalEventCount(t, a.workspace)
		out, stderr, err := closureRunCommand(a.workspace, a.binary, append([]string{"--project", a.workspace}, args...)...)
		if err != nil || stderr != "" || p21JournalEventCount(t, a.workspace) != before {
			t.Fatalf("real corpus view %v failed or wrote state: %v %s", args, err, stderr)
		}
		start := strings.Index(out, "Stop "+stop["id"].(string))
		if args[0] == "resume" && stop["resolution"] != nil {
			// The accepted resume projection contains open Stops and open
			// Attention. A resolved Stop remains visible through its open
			// Attention until disposition; recovered predecessor Runs and
			// atomically disposed judgment items can leave this projection.
			projection := p21Run(t, a.binary, a.workspace, append([]string{"--project", a.workspace, "--json"}, args...)...)
			visible := false
			for _, raw := range p21Path(t, projection, "result", "resume", "attention").([]any) {
				if raw.(map[string]any)["id"] == attention["id"] {
					visible = true
				}
			}
			if !visible {
				if start >= 0 || strings.Contains(out, "Attention "+attention["id"].(string)) {
					t.Fatal("Human resume retained a disposed item absent from its canonical projection")
				}
				continue
			}
		}
		if start < 0 {
			t.Fatalf("%s %v missing Stop %s: %s", row.code, args, stop["id"], out)
		}
		text := out[start:]
		// Bound assertions to this Stop, excluding guidance for other Stops.
		if end := strings.Index(text[5:], "\nStop "); end >= 0 {
			text = text[:end+5]
		}
		if end := strings.Index(text, "\nAttention "); end >= 0 {
			text = text[:end]
		}
		wants := []string{"Code: " + row.code, "Blocked actions: " + p24RealStopValue(stop["blocked_actions"]), "Resolver capability: stop.resolve", "Subject: " + p21PathString(t, stop, "subject", "type") + ":" + p21PathString(t, stop, "subject", "id"), "Resolution condition: " + row.condition}
		for key, value := range stop["details"].(map[string]any) {
			wants = append(wants, labels[key]+": "+p24RealStopValue(value))
		}
		if stop["resolution"] == nil {
			wants = append(wants, "Status: open", "Next: "+row.next+"; after the fact: fid stop resolve, then fid attention resolve where no lifecycle Event disposed it.")
		} else {
			wants = append(wants, "Status: resolved", "Fact Event: "+p21PathString(t, stop, "resolution", "fact_event_id"), "Resolution Event: "+p21PathString(t, stop, "resolution", "event_id"))
			if strings.Contains(text, "Next:") {
				t.Fatalf("resolved Stop has Next guidance: %s", text)
			}
		}
		for _, want := range wants {
			if !strings.Contains(text, want) {
				t.Errorf("%s %v missing %q: %s", row.code, args, want, text)
			}
		}
		if strings.Contains(out, "{") || strings.Contains(out, "<nil>") {
			t.Fatalf("%s %v leaked raw JSON/nil: %s", row.code, args, out)
		}
		if args[0] != "stop" {
			want := "Attention " + attention["id"].(string) + "\n  Status: " + attention["status"].(string) + "\n  Code: " + row.code + "\n  Linked Stops: " + stop["id"].(string)
			if !strings.Contains(out, want) {
				t.Errorf("%s %v missing Attention state: %s", row.code, args, out)
			}
			if attention["status"] == "resolved" && !strings.Contains(out, "Resolution Event: "+attention["resolution_event_id"].(string)) {
				t.Error("missing Attention resolution Event")
			}
		}
	}
}

func p24RealStopValue(value any) string {
	if value == nil {
		return "unknown"
	}
	if items, ok := value.([]any); ok {
		parts := make([]string, len(items))
		for i, item := range items {
			parts[i] = fmt.Sprint(item)
		}
		return strings.Join(parts, ", ")
	}
	return fmt.Sprint(value)
}

var p24FixtureActions = []string{"attention.resolve", "claim.acquire", "claim.expire", "claim.release", "decision.record-reference", "envelope.activate", "envelope.create", "gate.register", "gate.run", "human-outcome.record", "issue.transition", "review.record", "run.close", "run.observe", "run.observe-budget", "run.recover", "run.start", "stop.raise", "stop.resolve"}

func p24GenerateHumanFixtures(t *testing.T, root string) {
	t.Helper()
	contents, err := os.ReadFile(integrationFIDBinary(t))
	if err != nil {
		t.Fatal(err)
	}
	manifest := p24HumanManifest{Version: 1, Binary: filepath.Join(root, "fid"), BinarySHA256: fmt.Sprintf("sha256:%x", sha256.Sum256(contents))}
	if err := os.WriteFile(manifest.Binary, contents, 0700); err != nil {
		t.Fatal(err)
	}
	p24HumanDenialFixtures(t, root, &manifest)
	for _, dirty := range []bool{false, true} {
		p24HumanViewFixtures(t, root, &manifest, dirty)
	}
	start := p24NewTakeControlStart(t, manifest.Binary, root)
	p24AddHumanFixture(t, &manifest, "take-control-start", "issue.transition", "protected-decision-required", start.workspace, start.blockedArgs, 5)
	p24HumanJudgmentFixtures(t, root, &manifest)
	p24HumanAuthoringFixtures(t, root, &manifest)
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), append(data, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
}

func p24HumanDenialFixtures(t *testing.T, root string, manifest *p24HumanManifest) {
	t.Helper()
	workspace := p22HumanDirectory(t, root, "denials")
	a := p22NewFIDRecoveryAdapterAt(t, manifest.Binary, workspace, "human", 86400, []string{}, p24FixtureActions)
	p22AcceptedCall(t, a, "human-conflict", "run.observe-budget", map[string]any{"value": 0})
	fixture := p22ConformanceCase{Name: "human-conflict", Operation: "run.observe-budget"}
	args, input := a.p22BuildArgs(t, fixture)
	input["value"] = 1
	args = append(p22WithoutJSON(args), "--input", p21Input(t, workspace, "conflict", input), "--idempotency-key", "p22-corpus-human-conflict")
	p24AddHumanFixture(t, manifest, "idempotency-conflict", "run.observe-budget", "idempotency-conflict", workspace, args, 4)

	a.p22Observe(t, "terminated")
	args, input = a.p22BuildArgs(t, p22ConformanceCase{Operation: "run.close"})
	input["expected_repository"] = map[string]any{"head_oid": p22Zeros(40), "dirty": false}
	args = append(p22WithoutJSON(args), "--input", p21Input(t, workspace, "mismatch", input), "--idempotency-key", "human-mismatch")
	p24AddHumanFixture(t, manifest, "repository-mismatch", "run.close", "repository-mismatch", workspace, args, 4)

	input = p24ObservedGateInputWithoutObservation(a)
	input["expected_repository"].(map[string]any)["dirty"] = "yes"
	args = []string{"--project", workspace, "gate", "run", "--input", p21Input(t, workspace, "nested", input), "--actor", a.actor, "--idempotency-key", "human-nested"}
	p24AddHumanFixture(t, manifest, "invalid-argument", "gate.run", "invalid-argument", workspace, args, 2)

	args = []string{"--project", workspace, "requirement", "record-reference", "--input", p21Input(t, workspace, "digest", map[string]any{"id": "REQ-human", "source": p23DecisionSource(a, false), "status": "active"}), "--actor", a.actor, "--idempotency-key", "human-digest"}
	p24AddHumanFixture(t, manifest, "source-digest-mismatch", "requirement.record-reference", "source-digest-mismatch", workspace, args, 4)

	checkpoint := p23ProtectedDecisionFixture(t, manifest.Binary, root)
	p24AddHumanFixture(t, manifest, "protected-decision-required", "issue.transition", "protected-decision-required", checkpoint.Project, p22WithoutJSON(checkpoint.Steps[0].Args), 5)
}

func p24HumanViewFixtures(t *testing.T, root string, manifest *p24HumanManifest, dirty bool) {
	t.Helper()
	state := "clean"
	if dirty {
		state = "dirty"
	}
	workspace := p22HumanDirectory(t, root, state+"-same")
	a := p22NewFIDRecoveryAdapterAt(t, manifest.Binary, workspace, "human", 86400, []string{}, p24FixtureActions)
	p23CloseAndRelease(t, a, a.actor, a.runID, a.claimID, "human-builder")
	claim := p23AcquireClaim(t, a, "actor-denied", "review", "human-review")
	run := p23StartRun(t, a, "actor-denied", claim, "human-review")
	runID := p21PathString(t, run, "id")
	if dirty {
		p22Write(t, workspace, "base.txt", "uncommitted Human work\n")
	}
	a.repository = p23CurrentRepository(t, a)
	input := p23ReviewInput(a, runID, p23Number(t, run["version"]), "pass", nil, nil)
	review := p23ReviewRecord(t, a, "actor-denied", "human-review", input)
	p23AssertReviewSuccess(t, "Human fixture Review", review)
	reviewID := p23ReviewFromResponse(t, review)["id"].(string)
	p22AcceptedCall(t, a, "human-stop", "stop.raise", nil)
	copied := p22HumanDirectory(t, root, state+"-different")
	copyParityFixture(t, workspace, copied)
	for _, same := range []bool{true, false} {
		repository, relation := workspace, "same"
		if !same {
			repository, relation = copied, "different"
		}
		viewArgs := [][]string{
			{"stop", "list"}, {"attention", "list"},
			{"resume", a.issueID, "--repository", repository},
			{"review", "get", reviewID}, {"review", "list", "--issue", a.issueID},
			{"run", "inspect", runID, "--repository", repository},
			{"review", "record", "--input", filepath.Join(workspace, "human-review.json"), "--actor", "actor-denied", "--idempotency-key", "human-review"},
		}
		// Replay the exact accepted input, retaining its own key and original facts.
		viewArgs[6][3] = p21Input(t, workspace, "human-review-replay", input)
		viewArgs[6][7] = "human-review"
		for _, leaf := range viewArgs {
			operation := strings.Join(leaf[:2], ".")
			if leaf[0] == "resume" {
				operation = "resume"
			}
			project := workspace
			if leaf[0] == "review" && leaf[1] != "record" || leaf[0] == "stop" || leaf[0] == "attention" {
				project = repository
			}
			args := append([]string{"--project", project}, leaf...)
			if operation == "run.inspect" && !same {
				before := p21JournalEventCount(t, workspace)
				row := p24AddHumanFixture(t, manifest, state+"-"+relation+"-"+operation, operation, "recovery-state-conflict", workspace, args, 4)
				failure := p23Envelope(t, row.JSONStderr)
				want := map[string]any{"observation_digest": nil, "paths": []any{}, "reasons": []any{"repository identity changed"}, "run_id": runID}
				if !reflect.DeepEqual(p21Path(t, failure, "error", "details"), want) ||
					!reflect.DeepEqual(p21Path(t, failure, "error", "subject"), map[string]any{"type": "run", "id": runID}) ||
					p21PathString(t, failure, "error", "category") != "conflict" ||
					p21PathString(t, failure, "error", "message") != "recovery state conflict" ||
					p21PathString(t, failure, "operation") != operation ||
					p21Path(t, failure, "error", "retryable") != false ||
					p21Path(t, failure, "ok") != false ||
					p21PathString(t, failure, "contract_version") != "fid-cli.v1" {
					t.Fatalf("%s differs from DEC-287 refusal: %s", row.Name, row.JSONStderr)
				}
				wantHuman := "error: recovery state conflict\nreason: recovery-state-conflict; observation_digest=null, paths=, reasons=repository identity changed, run_id=" + runID + "\n"
				if row.HumanStderr != wantHuman || p21JournalEventCount(t, workspace) != before {
					t.Fatalf("%s refusal bytes or read-only state changed: %s", row.Name, row.HumanStderr)
				}
				continue
			}
			row := p24AddHumanFixture(t, manifest, state+"-"+relation+"-"+operation, operation, "", workspace, args, 0)
			if leaf[0] == "resume" || leaf[0] == "run" || leaf[0] == "review" {
				p24AssertShortHumanRepository(t, row.HumanStdout, p24FixtureHEAD(t, workspace), dirty)
			}
			if leaf[0] == "resume" || leaf[0] == "run" {
				if !strings.Contains(row.HumanStdout, relation+" worktree") {
					t.Fatalf("%s comparison missing: %s", row.Name, row.HumanStdout)
				}
			}
		}
	}
}

func p24AddHumanFixture(t *testing.T, manifest *p24HumanManifest, name, operation, code, workspace string, args []string, expected int) p24HumanFixture {
	t.Helper()
	row := p24HumanFixture{Name: name, Operation: operation, Code: code, Workspace: workspace, CWD: workspace, HumanArgs: args, JSONArgs: append([]string{"--json"}, args...), ExpectedExit: expected, BinarySHA256: manifest.BinarySHA256}
	row.HumanCommand = p24FixtureCommand(manifest.Binary, row.HumanArgs)
	row.JSONCommand = p24FixtureCommand(manifest.Binary, row.JSONArgs)
	row.JSONStdout, row.JSONStderr = p24FixtureRun(t, manifest.Binary, row, true)
	row.HumanStdout, row.HumanStderr = p24FixtureRun(t, manifest.Binary, row, false)
	if code != "" {
		failure := p23Envelope(t, row.JSONStderr)
		if p21PathString(t, failure, "error", "code") != code {
			t.Fatalf("%s wrong denial: %s", name, row.JSONStderr)
		}
		if row.HumanStderr != p24ExpectedFailureText(t, failure, operation, containsAction(args, "--input")) {
			t.Fatalf("%s Human/JSON mismatch: %s", name, row.HumanStderr)
		}
		if code == "invalid-argument" && p21PathString(t, failure, "error", "details", "argument") != "input.expected_repository.dirty" {
			t.Fatal("nested path missing")
		}
		if code == "protected-decision-required" && !strings.Contains(row.HumanStderr, "blocked by:") {
			t.Fatal("blocked by missing")
		}
	} else if row.HumanStdout == "" || strings.Contains(row.HumanStdout, "{") {
		t.Fatalf("%s empty or raw JSON view", name)
	}
	manifest.Fixtures = append(manifest.Fixtures, row)
	return row
}

func p24FixtureRun(t *testing.T, binary string, row p24HumanFixture, jsonMode bool) (string, string) {
	t.Helper()
	args := row.HumanArgs
	if jsonMode {
		args = row.JSONArgs
	}
	out, stderr, err := closureRunCommand(row.CWD, binary, args...)
	exit := 0
	if err != nil {
		exit = processExitCode(t, err)
	}
	if exit != row.ExpectedExit {
		t.Fatalf("%s json=%t exit=%d want=%d: %s %s", row.Name, jsonMode, exit, row.ExpectedExit, out, stderr)
	}
	if exit != 0 && out != "" || exit == 0 && stderr != "" {
		t.Fatalf("%s invalid streams: %q %q", row.Name, out, stderr)
	}
	return out, stderr
}

func p24FixtureCommand(binary string, args []string) string {
	parts := append([]string{binary}, args...)
	for i, part := range parts {
		parts[i] = "'" + strings.ReplaceAll(part, "'", "'\"'\"'") + "'"
	}
	return strings.Join(parts, " ")
}

func p24AssertShortHumanRepository(t *testing.T, text, head string, dirty bool) {
	t.Helper()
	state := "clean"
	if dirty {
		state = "dirty"
	}
	if !strings.Contains(text, "HEAD "+head[:12]+"; "+state) {
		t.Fatalf("missing short %s repository: %s", state, text)
	}
	for _, forbidden := range []string{"worktree_id", "tree_oid", "dirty_digest", "capture_source", "<nil>", head} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("repository leaked %q: %s", forbidden, text)
		}
	}
}

func p24FixtureHEAD(t *testing.T, workspace string) string {
	t.Helper()
	out, _ := closureRun(t, workspace, "git", "rev-parse", "HEAD")
	return strings.TrimSpace(out)
}

func p24AssertHumanManifest(t *testing.T, manifest p24HumanManifest) {
	t.Helper()
	data, err := os.ReadFile(manifest.Binary)
	if err != nil || fmt.Sprintf("sha256:%x", sha256.Sum256(data)) != manifest.BinarySHA256 {
		t.Fatalf("binary digest mismatch: %v", err)
	}
	seen := map[string]bool{}
	for _, row := range manifest.Fixtures {
		if seen[row.Name] || row.CWD == "" || row.Workspace == "" || row.BinarySHA256 != manifest.BinarySHA256 || row.HumanCommand != p24FixtureCommand(manifest.Binary, row.HumanArgs) || row.JSONCommand != p24FixtureCommand(manifest.Binary, row.JSONArgs) {
			t.Fatalf("incomplete/duplicate manifest row %s", row.Name)
		}
		seen[row.Name] = true
		if row.Authoring {
			humanOut, humanErr, jsonOut, jsonErr := p24RehearseAuthoring(t, manifest.Binary, row)
			if humanOut != row.HumanStdout || humanErr != row.HumanStderr || p21PathString(t, p23Envelope(t, jsonOut), "operation") != row.Operation || p21Path(t, p23Envelope(t, jsonOut), "result", "replayed") != true || jsonErr != row.JSONStderr {
				t.Fatalf("authoring fixture %s drifted", row.Name)
			}
			continue
		}
		out, stderr := p24FixtureRun(t, manifest.Binary, row, false)
		if out != row.HumanStdout || stderr != row.HumanStderr {
			t.Fatalf("retained fixture %s Human drifted", row.Name)
		}
		out, stderr = p24FixtureRun(t, manifest.Binary, row, true)
		if row.ExpectedExit == 0 {
			if row.Operation == "resume" {
				continue
			} // The lease projection is live; byte parity is checked at generation.
			if !reflect.DeepEqual(p23Envelope(t, out), p23Envelope(t, row.JSONStdout)) {
				t.Fatalf("retained fixture %s JSON drifted", row.Name)
			}
		} else if out != row.JSONStdout || stderr != row.JSONStderr {
			t.Fatalf("retained fixture %s JSON drifted", row.Name)
		}
	}
	for _, name := range []string{"idempotency-conflict", "repository-mismatch", "invalid-argument", "protected-decision-required", "source-digest-mismatch"} {
		if !seen[name] {
			t.Errorf("missing denial %s", name)
		}
	}
	for _, state := range []string{"clean", "dirty"} {
		for _, relation := range []string{"same", "different"} {
			for _, view := range []string{"stop.list", "attention.list", "resume", "review.get", "review.list", "review.record", "run.inspect"} {
				name := state + "-" + relation + "-" + view
				if !seen[name] {
					t.Errorf("missing view %s", name)
				}
			}
		}
	}
}

func p24HumanJudgmentFixtures(t *testing.T, root string, manifest *p24HumanManifest) {
	t.Helper()
	pairs := p23GenerateLeafPairs(t, manifest.Binary, root)
	workspace := pairs[0].CWD
	// The shared retained generator provides every judgment list/detail pair.
	for _, pair := range pairs {
		p24AddHumanFixture(t, manifest, "judgment-"+pair.Operation, pair.Operation, "", workspace, pair.HumanArgs, 0)
	}
	for family, id := range map[string]string{"concern": "CON-leaf", "deviation": "DEV-leaf", "human-outcome": "OUT-leaf"} {
		p24AddHumanFixture(t, manifest, "judgment-"+family+".get", family+".get", "", workspace, []string{family, "get", id}, 0)
	}
	issues := p21Run(t, manifest.Binary, workspace, "--json", "issue", "list")
	issueID := p21Path(t, issues, "result", "issues").([]any)[0].(map[string]any)["id"].(string)
	p24AddHumanFixture(t, manifest, "judgment-resume", "resume", "", workspace, []string{"resume", issueID, "--repository", workspace}, 0)
}

func p24HumanAuthoringFixtures(t *testing.T, root string, manifest *p24HumanManifest) {
	t.Helper()
	for _, kind := range []string{"run-close", "decision-reference", "third-leaf"} {
		workspace := p22HumanDirectory(t, root, "authoring-"+kind)
		actions := append(append([]string{}, p24FixtureActions...), "concern.record")
		a := p22NewFIDRecoveryAdapterAt(t, manifest.Binary, workspace, "human", 86400, []string{}, actions)
		var parts []string
		var input map[string]any
		switch kind {
		case "run-close":
			a.p22Observe(t, "terminated")
			parts, input = a.p22BuildArgs(t, p22ConformanceCase{Operation: "run.close"})
			parts = p22WithoutJSON(parts)
		case "decision-reference":
			parts = []string{"decision", "record-reference", "--actor", a.actor}
			source := p23DecisionSource(a, false)
			delete(source, "digest")
			input = map[string]any{"id": "DREF-authoring", "class": "security", "source": source, "authority": map[string]any{"actor_id": a.actor, "required_capabilities": []string{"decision.record-reference"}}, "outcome": "approved", "scope": map[string]any{"issue_id": a.issueID}}
		case "third-leaf":
			parts = []string{"concern", "record", "--actor", a.actor}
			input = map[string]any{"id": "CON-authoring", "subject": map[string]any{"type": "issue", "id": a.issueID}, "scope": map[string]any{"paths": []string{}}, "severity": "minor", "summary": "Human authored concern", "reasoning": "Human chooses another input leaf", "references": []any{}}
		}
		path := p21Input(t, root, "authoring-"+kind, input)
		args := append(parts, "--input", path, "--idempotency-key", "authoring-"+kind)
		operation := map[string]string{"run-close": "run.close", "decision-reference": "decision.record-reference", "third-leaf": "concern.record"}[kind]
		row := p24HumanFixture{Authoring: true, Name: "authoring-" + kind, Operation: operation, Workspace: workspace, CWD: workspace, HumanArgs: args, JSONArgs: append([]string{"--json"}, args...), ExpectedExit: 0, BinarySHA256: manifest.BinarySHA256}
		row.HumanCommand, row.JSONCommand = p24FixtureCommand(manifest.Binary, row.HumanArgs), p24FixtureCommand(manifest.Binary, row.JSONArgs)
		row.HumanStdout, row.HumanStderr, row.JSONStdout, row.JSONStderr = p24RehearseAuthoring(t, manifest.Binary, row)
		manifest.Fixtures = append(manifest.Fixtures, row)
	}
}

// Rehearse in the recorded physical worktree, then restore its pre-action
// workspace. The Human receives a valid first request, not a replay-only key.
func p24RehearseAuthoring(t *testing.T, binary string, row p24HumanFixture) (string, string, string, string) {
	t.Helper()
	scratch, err := os.MkdirTemp(filepath.Dir(row.Workspace), "authoring-rehearsal-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.RemoveAll(scratch); err != nil {
			t.Error(err)
		}
	}()
	snapshot := filepath.Join(scratch, "before")
	copyParityFixture(t, config.LedgerDir(row.Workspace), snapshot)
	humanOut, humanErr := p24FixtureRun(t, binary, row, false)
	jsonOut, jsonErr := p24FixtureRun(t, binary, row, true)
	if err := os.Rename(config.LedgerDir(row.Workspace), filepath.Join(scratch, "after")); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(snapshot, config.LedgerDir(row.Workspace)); err != nil {
		t.Fatal(err)
	}
	return humanOut, humanErr, jsonOut, jsonErr
}
