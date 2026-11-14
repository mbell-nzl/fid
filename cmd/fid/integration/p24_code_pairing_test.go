//go:build integration

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/mbell-nzl/fid/internal/cli"
	"github.com/mbell-nzl/fid/internal/gittest"
)

type p24CodePair struct {
	root, operation string
	args            []string
}
type p24CodeFactory func(*testing.T, string) p24CodePair

// The corpus is enumerated against RegisteredCodes below. A row invokes the
// real fid binary, with no injected mapper errors or simulated responses.
var p24CodeRows = map[string]p24CodeFactory{
	"hold-not-current":              p24LifecycleCodePair,
	"worktree-busy":                 p24LifecycleCodePair,
	"already-held":                  p24LifecycleCodePair,
	"unsupported-kind":              p24LifecycleCodePair,
	"internal-error":                p24InternalPair,
	"journal-replay-stopped":        p24JournalPair,
	"lock-unavailable":              p24LockPair,
	"projection-divergence":         p24ProjectionPair,
	"import-target-conflict":        p24ImportPair,
	"activation-conflict":           p24ActivationCodePair,
	"activation-version-conflict":   p24ActivationCodePair,
	"activation-invalid":            p24ActivationCodePair,
	"activation-direction-conflict": p24ActivationCodePair,
	"activation-incomplete":         p24ActivationCodePair,
	"claim-conflict":                p24ClaimCodePair,
	"claim-stale":                   p24ClaimCodePair,
	"claim-expired":                 p24ClaimCodePair,
	"run-inactive":                  p24ClaimCodePair,
	"authority-unknown-actor":       p24GateCodePair,
	"authority-unavailable":         p24GateCodePair,
	"executor-policy-unsupported":   p24GateCodePair,
	"repository-unavailable":        p24GateCodePair,

	"claim-mode-invalid":         p24ReviewClaimPair,
	"obligation-unsatisfied":     p24ObligationCodePair,
	"human-outcome-required":     p24ObligationCodePair,
	"finding-not-found":          p24JudgmentCodePair,
	"concern-not-found":          p24JudgmentCodePair,
	"deviation-not-found":        p24JudgmentCodePair,
	"concern-already-disposed":   p24JudgmentCodePair,
	"source-digest-mismatch":     p24JudgmentCodePair,
	"authority-missing":          p24OutcomeCodePair,
	"outcome-not-found":          p24OutcomeCodePair,
	"outcome-not-permitted":      p24OutcomeCodePair,
	"review-verdict-invalid":     p24ReviewCodePair,
	"review-not-independent":     p24ReviewCodePair,
	"review-changes-required":    p24ReviewCodePair,
	"evidence-type-not-accepted": p24ReviewCodePair,
	"finding-already-disposed":   p24ReviewCodePair,
	"result-limit-exceeded":      p24JudgmentCodePair,

	"invalid-argument":      p24SimplePair("stop.raise", "stop", "raise", "--bogus"),
	"operation-unavailable": p24SimplePair("fid", "isue"),
	"project-not-found": func(t *testing.T, _ string) p24CodePair {
		return p24CodePair{gittest.NewWorkspace(t), "activity.list", []string{"activity", "list"}}
	},
	"project-already-initialized": p24SimplePair("project.init", "init", "--prefix", "P24"),
	"issue-not-found":             p24SimplePair("issue.get", "issue", "get", "FID-missing"),
	"label-not-found":             p24IssueCodePair,
	"state-conflict":              p24IssueCodePair,
	"hierarchy-invalid":           p24IssueCodePair,
	"transition-not-permitted":    p24IssueCodePair,
	"link-not-found":              p24IssueCodePair,
	"unsupported-relation":        p24IssueCodePair,
	"dependency-cycle":            p24IssueCodePair,
	"workflow-not-found":          p24SimplePair("workflow.assign", "workflow", "assign", "--name", "missing", "--project-wide", "--actor", "agent:a", "--idempotency-key", "pair"),
	"requirement-not-found":       p24SimplePair("requirement.get", "requirement", "get", "REQ-missing"),
	"capability-unavailable":      p24SimplePair("maintenance.vacuum", "maintenance", "vacuum", "--actor", "agent:a"),
	"activation-not-found":        p24SimplePair("activation.status", "activation", "status", "ACT-missing"),
	"envelope-not-found":          p24SimplePair("envelope.get", "envelope", "get", "ENV-missing"),
	"run-not-found":               p24SimplePair("run.get", "run", "get", "RUN-missing"),
	"receipt-not-found":           p24SimplePair("receipt.get", "receipt", "get", "RCP-missing"),
	"evidence-not-found":          p24SimplePair("evidence.get", "evidence", "get", "EVD-missing"),
	"review-not-found":            p24SimplePair("review.get", "review", "get", "REV-missing"),
	"archive-invalid":             p24ArchivePair,
	"claim-not-found":             p24RecoveryCodePair,
	"authority-denied":            p24RecoveryCodePair,
	"run-version-conflict":        p24RecoveryCodePair,
	"executor-liveness-unknown":   p24RecoveryCodePair,
	"executor-still-live":         p24RecoveryCodePair,
	"recovery-state-conflict":     p24RecoveryCodePair,
	"observation-not-found":       p24RecoveryCodePair,
	"recovery-already-linked":     p24RecoveryCodePair,
	"duplicate-entity":            p24RecoveryCodePair,
	"stop-not-found":              p24RecoveryCodePair,
	"attention-not-found":         p24RecoveryCodePair,
	"attention-version-conflict":  p24RecoveryCodePair,
	"resolution-invalid":          p24RecoveryCodePair,
	"stop-unresolved":             p24RecoveryCodePair,
	"budget-sample-conflict":      p24RecoveryCodePair,
	"repository-mismatch":         p24RecoveryCodePair,
	"registration-not-found":      p24RecoveryCodePair,
}

// Only exclusions proven by a current producer/caller audit belong here.
// Missing executable rows are failures, never silently classified unreachable.
var p24CodeExclusions = map[string]string{}

var p24SpecialCodePairs = map[string]func(*testing.T){
	"policy-unavailable":          TestPolicyUnavailableHumanAndJSON,
	"gate-unsuccessful":           TestP24HumanGateTerminalFailure,
	"budget-data-unavailable":     p24BudgetTerminalPair,
	"budget-exhausted":            p24BudgetTerminalPair,
	"protected-decision-required": TestP24HumanProtectedDecisionFailure,
	"review-rework-exhausted":     TestP23StopFailureHumanAndJSONParity,
}

// idempotency-conflict is a defined, reachable code that RegisteredCodes does
// not enumerate today. Keep its next-table process coverage explicitly, without
// changing that protected inventory or pretending it is a registry member.
var p24SupplementalPairs = map[string]p24CodeFactory{"idempotency-conflict": p24RecoveryCodePair}

func TestP24RegisteredCodePairingCompleteness(t *testing.T) {
	t.Parallel()
	registered := map[string]bool{}
	for _, code := range cli.RegisteredCodes() {
		registered[string(code)] = true
	}
	for code := range p24CodeRows {
		if !registered[code] {
			t.Errorf("stale pairing row %s", code)
		}
	}
	for code := range p24SpecialCodePairs {
		if !registered[code] || p24CodeRows[code] != nil {
			t.Errorf("invalid special pairing row %s", code)
		}
	}
	for code, reason := range p24CodeExclusions {
		if !registered[code] || (p24CodeRows[code] != nil || p24SpecialCodePairs[code] != nil) || strings.TrimSpace(reason) == "" {
			t.Errorf("invalid exclusion %s: %s", code, reason)
		}
	}
	for code := range registered {
		if p24CodeRows[code] == nil && p24SpecialCodePairs[code] == nil && p24CodeExclusions[code] == "" {
			t.Errorf("registered code has no process pair or pinned exclusion: %s", code)
		}
	}
}

func TestP24RegisteredCodePairs(t *testing.T) {
	t.Parallel()
	for code, test := range p24SpecialCodePairs {
		t.Run(code, test)
	}
	rows := map[string]p24CodeFactory{}
	for code, prepare := range p24CodeRows {
		rows[code] = prepare
	}
	for code, prepare := range p24SupplementalPairs {
		rows[code] = prepare
	}
	for code, prepare := range rows {
		t.Run(code, func(t *testing.T) {
			t.Parallel()
			pair := prepare(t, code)
			binary := integrationFIDBinary(t)
			args := append([]string{"--project", pair.root}, pair.args...)
			if pair.operation == "export.import" {
				args = pair.args
			}
			out, stderr, err := closureRunCommand(pair.root, binary, append([]string{"--json"}, args...)...)
			if err == nil || out != "" {
				t.Fatalf("JSON %s: %v %q %q", code, err, out, stderr)
			}
			payload := p23Envelope(t, stderr)
			f := payload["error"].(map[string]any)
			if payload["operation"] != pair.operation || f["code"] != code {
				t.Fatalf("wrong producer for %s: %#v", code, payload)
			}
			status := processExitCode(t, err)
			for _, r := range f["message"].(string) {
				if unicode.IsControl(r) {
					t.Fatal("control character in JSON message")
				}
			}
			out, stderr, err = closureRunCommand(pair.root, binary, args...)
			want := p24ExpectedFailureText(t, payload, pair.operation, strings.Contains(strings.Join(args, " "), "--input") || pair.operation == "stop.raise")
			if err == nil || processExitCode(t, err) != status || out != "" || stderr != want {
				t.Fatalf("Human %s: %v stdout=%q stderr=%q want=%q", code, err, out, stderr, want)
			}
		})
	}
}

func p24SimplePair(operation string, args ...string) p24CodeFactory {
	return func(t *testing.T, _ string) p24CodePair {
		t.Helper()
		root := gittest.NewWorkspace(t)
		p21Run(t, integrationFIDBinary(t), root, "init", "--prefix", "P24", "--json")
		return p24CodePair{root, operation, args}
	}
}

func p24IssueCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	root, binary := gittest.NewWorkspace(t), integrationFIDBinary(t)
	p21Run(t, binary, root, "init", "--prefix", "P24", "--json")
	create := func(title string) string {
		t.Helper()
		v := p21Run(t, binary, root, "--json", "issue", "create", "--kind", "task", "--title", title, "--actor", "agent:a", "--idempotency-key", title)
		return p21PathString(t, v, "result", "issue", "id")
	}
	id := create("first")
	args, operation := []string{}, ""
	actor := []string{"--actor", "agent:a", "--idempotency-key", "pair"}
	switch code {
	case "label-not-found":
		operation = "issue.label.add"
		args = []string{"issue", "label", "add", id, "missing"}
	case "state-conflict":
		operation = "issue.update"
		args = []string{"issue", "update", id, "--title", "changed", "--expected-revision", "999"}
	case "hierarchy-invalid":
		operation = "issue.create"
		args = []string{"issue", "create", "--kind", "subtask", "--title", "child"}
	case "transition-not-permitted":
		operation = "issue.transition"
		args = []string{"issue", "transition", id, "--to", "qa_pass", "--expected-revision", "1"}
	case "link-not-found":
		operation = "link.remove"
		args = []string{"link", "remove", id, "blocks", create("second"), "--from-type", "issue", "--to-type", "issue"}
	case "unsupported-relation":
		operation = "link.add"
		args = []string{"link", "add", id, "implements", create("second"), "--from-type", "issue", "--to-type", "issue"}
	case "dependency-cycle":
		peer := create("second")
		p21Run(t, binary, root, "--json", "link", "add", id, "blocks", peer, "--from-type", "issue", "--to-type", "issue", "--actor", "agent:a", "--idempotency-key", "edge")
		operation = "link.add"
		args = []string{"link", "add", peer, "blocks", id, "--from-type", "issue", "--to-type", "issue"}
	}
	return p24CodePair{root, operation, append(args, actor...)}
}

func p24ArchivePair(t *testing.T, _ string) p24CodePair {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, "invalid.zip")
	if err := os.WriteFile(path, []byte("invalid ZIP"), 0600); err != nil {
		t.Fatal(err)
	}
	return p24CodePair{root, "export.validate", []string{"export", "validate", path}}
}

func p24RecoveryCodePair(t *testing.T, code string) p24CodePair {
	t.Helper()
	a := p24RunAdapter(t)
	fixture := p22ConformanceCase{Name: "pair", Operation: "run.recover"}
	actor := a.actor
	switch code {
	case "claim-not-found":
		fixture.Operation = "claim.expire"
		fixture.Input = map[string]any{"claim_id": "CLM-missing"}
	case "authority-denied":
		fixture.Operation = "run.observe"
		actor = "actor-denied"
	case "run-version-conflict":
		fixture.Operation = "run.observe"
		fixture.Input = map[string]any{"expected_version": 999}
	case "executor-liveness-unknown", "executor-still-live":
		state := "unknown"
		if code == "executor-still-live" {
			state = "live"
		}
		a.p22Observe(t, state)
		a.p22Expire(t)
	case "recovery-state-conflict":
		a.p22Observe(t, "terminated")
		fixture.Operation = "run.checkpoint"
		fixture.Input = map[string]any{"reconciliation_digest": "sha256:" + p22Zeros(64)}
	case "observation-not-found", "duplicate-entity":
		a.p22Observe(t, "terminated")
		a.p22Expire(t)
		if code == "observation-not-found" {
			fixture.Input = map[string]any{"observation_id": "OBS-missing"}
		} else {
			fixture.Input = map[string]any{"executor": map[string]any{"id": a.executorID, "source": "activated-harness"}}
		}
	case "recovery-already-linked":
		a.p22Observe(t, "terminated")
		a.p22Expire(t)
		prior := a.runID
		p22AcceptedCall(t, a, "pair-first-recovery", "run.recover", nil)
		read := p21Run(t, a.binary, a.workspace, "--json", "run", "get", prior)
		fixture.Input = map[string]any{"prior_run_id": prior, "expected_run_version": p21Path(t, read, "result", "run", "version")}
	case "stop-not-found":
		fixture.Operation = "stop.resolve"
		fixture.Input = map[string]any{"stop_id": "STP-missing", "fact_event_id": "EVT-missing"}
	case "attention-not-found":
		fixture.Operation = "attention.resolve"
		fixture.Input = map[string]any{"attention_id": "ATT-missing", "expected_version": 1, "fact_event_ids": []string{"EVT-missing"}}
	case "attention-version-conflict", "resolution-invalid":
		p22AcceptedCall(t, a, "pair-stop", "stop.raise", nil)
		if code == "attention-version-conflict" {
			fixture.Operation = "attention.resolve"
			fixture.Input = map[string]any{"expected_version": 999, "fact_event_ids": []string{"EVT-missing"}}
		} else {
			fixture.Operation = "stop.resolve"
			fixture.Input = map[string]any{"fact_event_id": "EVT-missing"}
		}
	case "stop-unresolved":
		a.p22Observe(t, "terminated")
		a.p22Expire(t)
		p22AcceptedCall(t, a, "pair-blocked", "stop.raise", map[string]any{"blocked_actions": []string{"run.recover"}})
	case "budget-sample-conflict":
		p22AcceptedCall(t, a, "pair-budget", "run.observe-budget", nil)
		fixture.Operation = "run.observe-budget"
		fixture.Input = map[string]any{"source_sample_id": "pair-budget", "value": 4}
	case "repository-mismatch":
		fixture.Operation = "run.close"
		a.p22Observe(t, "terminated")
		fixture.Input = map[string]any{"expected_repository": map[string]any{"head_oid": p22Zeros(40), "dirty": false}}
	case "registration-not-found":
		return p24CodePair{a.workspace, "gate.resolve", []string{"gate", "resolve", "REG-missing", "--version", "1"}}
	case "idempotency-conflict":
		p22AcceptedCall(t, a, "pair-conflict", "run.observe-budget", nil)
		fixture.Name = "pair-conflict"
		fixture.Operation = "run.observe-budget"
		fixture.Input = map[string]any{"value": 4}
	}
	args, input := a.p22BuildArgs(t, fixture)
	for key, value := range fixture.Input {
		input[key] = value
	}
	human := []string{}
	for i := 0; i < len(args); i++ {
		if args[i] == "--json" {
			continue
		}
		if args[i] == "--actor" {
			human = append(human, args[i], actor)
			i++
			continue
		}
		human = append(human, args[i])
	}
	human = append(human, "--input", p21Input(t, a.workspace, "pair-input", input), "--idempotency-key", "p22-corpus-"+fixture.Name)
	return p24CodePair{a.workspace, fixture.Operation, human}
}

// Independent section-9 oracle. It reads JSON values and never calls the
// production renderer or BoundedMessage. Unit fixtures separately pin the cuts.
func p24ExpectedFailureText(t *testing.T, payload map[string]any, operation string, input bool) string {
	t.Helper()
	f := payload["error"].(map[string]any)
	d := f["details"].(map[string]any)
	code := f["code"].(string)
	join := func(key, separator string) string {
		values, _ := d[key].([]any)
		parts := []string{}
		for _, v := range values {
			parts = append(parts, fmt.Sprint(v))
		}
		return strings.Join(parts, separator)
	}
	detail := ""
	switch code {
	case "invalid-argument":
		detail = "argument=" + fmt.Sprint(d["argument"]) + ", allowed=" + join("allowed", "|")
	case "idempotency-conflict":
		detail = "differing fields: " + join("differing_fields", ",")
	case "transition-not-permitted":
		detail = join("unsatisfied", ",")
	case "obligation-unsatisfied", "stop-unresolved", "resolution-invalid":
		parts := []string{}
		for _, raw := range d["reasons"].([]any) {
			parts = append(parts, raw.(map[string]any)["code"].(string))
		}
		detail = strings.Join(parts, ",")
	case "repository-mismatch":
		state := func(key string) string {
			v := d[key].(map[string]any)
			head := v["head_oid"].(string)
			if len(head) > 12 {
				head = head[:12]
			}
			status := "clean"
			if v["dirty"] == true {
				status = "dirty"
			}
			return head + " " + status
		}
		detail = "expected HEAD " + state("expected") + ", observed HEAD " + state("observed")
	default:
		keys := []string{}
		for key := range d {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		parts := []string{}
		for _, key := range keys {
			value := d[key]
			text := ""
			switch value := value.(type) {
			case map[string]any:
				continue
			case []any:
				// The accepted activation failure schema declares this as an object
				// array even when JSON contains an empty array; section 9 omits it.
				if (key == "component_errors" && strings.HasPrefix(code, "activation-")) || (code == "projection-divergence" && (key == "warnings" || key == "errors")) {
					continue
				}
				scalars := []string{}
				nested := false
				for _, v := range value {
					if _, ok := v.(map[string]any); ok {
						nested = true
						break
					}
					scalars = append(scalars, fmt.Sprint(v))
				}
				if nested {
					continue
				}
				text = strings.Join(scalars, ",")
			case nil:
				text = "null"
			default:
				text = fmt.Sprint(value)
			}
			parts = append(parts, key+"="+text)
		}
		detail = strings.Join(parts, ", ")
	}
	reason := "reason: " + code
	if detail != "" {
		reason += "; " + detail
	}
	reason = p24PairReasonCut(reason)
	text := "error: " + f["message"].(string) + "\n" + reason + "\n"
	ids := append([]string{}, strings.FieldsFunc(join("stop_ids", ","), func(r rune) bool { return r == ',' })...)
	ids = append(ids, strings.FieldsFunc(join("attention_ids", ","), func(r rune) bool { return r == ',' })...)
	if len(ids) > 0 {
		text += "blocked by: " + strings.Join(ids, ",") + "\n"
	}
	next := map[string]string{
		"hold-not-current":            "fid next",
		"repository-mismatch":         "check HEAD and uncommitted changes, then correct expected_repository",
		"source-digest-mismatch":      "correct or omit digest for " + fmt.Sprint(d["display_label"]),
		"idempotency-conflict":        "resend the original input, or use a new --idempotency-key",
		"stop-unresolved":             "read each Stop's qualifying resolution: fid stop list --status open",
		"review-rework-exhausted":     "an activated Human records an approved Human Outcome of class review-rework-extension: fid human-outcome record (template: fid human-outcome record --example)",
		"protected-decision-required": "record the " + fmt.Sprint(d["decision_class"]) + " Decision Reference and an approving Human Outcome: fid decision record-reference, fid human-outcome record (templates: --example)",
		"human-outcome-required":      "an activated Human records a " + fmt.Sprint(d["outcome_class"]) + " Human Outcome, then link it: fid human-outcome record, fid evidence link (templates: --example)",
		"obligation-unsatisfied":      "see per-Obligation reasons: fid evidence evaluate --input <file> (template: fid evidence evaluate --example)",
	}[code]
	if code == "invalid-argument" && input {
		leaf := "fid " + strings.ReplaceAll(operation, ".", " ")
		next = "see " + leaf + " --help; template: " + leaf + " --example"
	}
	if next != "" {
		text += "next: " + next + "\n"
	}
	return text
}

func p24PairReasonCut(text string) string {
	var escaped strings.Builder
	for _, r := range text {
		switch r {
		case '\n':
			escaped.WriteString(`\n`)
		case '\r':
			escaped.WriteString(`\r`)
		case '\t':
			escaped.WriteString(`\t`)
		default:
			if unicode.IsControl(r) {
				fmt.Fprintf(&escaped, `\u%04x`, r)
			} else {
				escaped.WriteRune(r)
			}
		}
	}
	value := escaped.String()
	if utf8.RuneCountInString(value) <= 200 {
		return value
	}
	tokens := regexp.MustCompile(`\\(?:[abfnrtv\\"']|x[0-9a-fA-F]{2}|u[0-9a-fA-F]{4}|U[0-9a-fA-F]{8})|.`).FindAllString(value, -1)
	var result strings.Builder
	count := 0
	for _, token := range tokens {
		size := utf8.RuneCountInString(token)
		if count+size > 200 {
			break
		}
		result.WriteString(token)
		count += size
	}
	return result.String() + "…"
}
