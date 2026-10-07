package cli

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestHumanFailureLines(t *testing.T) {
	t.Parallel()
	rows := []struct {
		name                  string
		code                  Code
		details               map[string]any
		operation, leaf, want string
	}{
		{"scalars", CodeStateConflict, map[string]any{"z": true, "a": []any{"one", 2}, "nested": map[string]any{"ignored": true}, "null": nil}, "issue.update", "fid issue update", "reason: state-conflict; a=one,2, null=null, z=true"},
		{"empty nested array", CodeActivationInvalid, map[string]any{"status": "draft", "expected_statuses": []string{"staged"}, "component_errors": []map[string]any{}}, "activation.commit", "fid activation commit", "reason: activation-invalid; expected_statuses=staged, status=draft"},
		{"nullable object", CodeReviewReworkExhausted, map[string]any{"event": nil}, "claim.acquire", "fid claim acquire", "reason: review-rework-exhausted; event=null\nnext: an activated Human records an approved Human Outcome of class review-rework-extension: fid human-outcome record (template: fid human-outcome record --example)"},
		{"repository", CodeRepositoryMismatch, map[string]any{"expected": map[string]any{"head_oid": "abcdef1234567890", "dirty": false}, "observed": map[string]any{"head_oid": "7654321fedcba987", "dirty": true}}, "run.start", "fid run start", "reason: repository-mismatch; expected HEAD abcdef123456 clean, observed HEAD 7654321fedcb dirty\nnext: check HEAD and uncommitted changes, then correct expected_repository"},
		{"conflict", CodeIdempotencyConflict, map[string]any{"differing_fields": []string{"--actor", "source"}}, "run.start", "fid run start", "reason: idempotency-conflict; differing fields: --actor,source\nnext: resend the original input, or use a new --idempotency-key"},
		{"transition", CodeTransitionNotPermitted, map[string]any{"unsatisfied": []string{"first", "second"}}, "issue.transition", "fid issue transition", "reason: transition-not-permitted; first,second"},
		{"obligation", CodeObligationUnsatisfied, map[string]any{"reasons": []map[string]any{{"code": "missing"}, {"code": "stale"}}}, "envelope.complete", "fid envelope complete", "reason: obligation-unsatisfied; missing,stale\nnext: see per-Obligation reasons: fid evidence evaluate --input <file> (template: fid evidence evaluate --example)"},
		{"stops", CodeStopUnresolved, map[string]any{"reasons": []any{map[string]any{"code": "operator-stop"}}, "stop_ids": []string{"STP-1", "STP-2"}, "attention_ids": []string{"ATT-1"}}, "run.start", "fid run start", "reason: stop-unresolved; operator-stop\nblocked by: STP-1,STP-2,ATT-1\nnext: read each Stop's qualifying resolution: fid stop list --status open"},
		{"resolution", CodeResolutionInvalid, map[string]any{"reasons": []map[string]any{{"code": "bad-fact"}}}, "stop.resolve", "fid stop resolve", "reason: resolution-invalid; bad-fact"},
		{"input", CodeInvalidArgument, map[string]any{"argument": "input.a", "allowed": []string{"a", "b"}, "value": "ignored"}, "stop.raise", "fid stop raise", "reason: invalid-argument; argument=input.a, allowed=a|b\nnext: see fid stop raise --help; template: fid stop raise --example"},
		{"no input", CodeInvalidArgument, map[string]any{"argument": "--title", "allowed": []string{"string"}}, "issue.create", "fid issue create", "reason: invalid-argument; argument=--title, allowed=string"},
		{"source", CodeSourceDigestMismatch, map[string]any{"display_label": "design"}, "decision.record-reference", "fid decision record-reference", "reason: source-digest-mismatch; display_label=design\nnext: correct or omit digest for design"},
		{"review", CodeReviewReworkExhausted, map[string]any{}, "claim.acquire", "fid claim acquire", "reason: review-rework-exhausted\nnext: an activated Human records an approved Human Outcome of class review-rework-extension: fid human-outcome record (template: fid human-outcome record --example)"},
		{"decision", CodeProtectedDecisionRequired, map[string]any{"decision_class": "security"}, "issue.transition", "fid issue transition", "reason: protected-decision-required; decision_class=security\nnext: record the security Decision Reference and an approving Human Outcome: fid decision record-reference, fid human-outcome record (templates: --example)"},
		{"outcome", CodeHumanOutcomeRequired, map[string]any{"outcome_class": "release"}, "envelope.complete", "fid envelope complete", "reason: human-outcome-required; outcome_class=release\nnext: an activated Human records a release Human Outcome, then link it: fid human-outcome record, fid evidence link (templates: --example)"},
		{"controls", CodeIssueNotFound, map[string]any{"id": "ISS-\n\r\t\x1b"}, "issue.get", "fid issue get", `reason: issue-not-found; id=ISS-\n\r\t\u001b`},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			var out bytes.Buffer
			if err := WriteHumanFailure(&out, row.operation, row.leaf, Failure{Code: row.code, Message: `escaped\nmessage`, Details: row.details}, row.name == "input"); err != nil {
				t.Fatal(err)
			}
			if want := "error: escaped\\nmessage\n" + row.want + "\n"; out.String() != want {
				t.Fatalf("got %q want %q", out.String(), want)
			}
		})
	}
}

func TestHumanFailureReasonCut(t *testing.T) {
	t.Parallel()
	prefix := "reason: issue-not-found; id="
	for _, size := range []int{200, 201} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			t.Parallel()
			var out bytes.Buffer
			value := strings.Repeat("é", size-len(prefix))
			if err := WriteHumanFailure(&out, "issue.get", "fid issue get", Failure{Code: CodeIssueNotFound, Message: "missing", Details: map[string]any{"id": value}}, false); err != nil {
				t.Fatal(err)
			}
			want := prefix + value
			if size > 200 {
				want = string([]rune(want)[:200]) + "…"
			}
			if out.String() != "error: missing\n"+want+"\n" {
				t.Fatalf("got %q want %q", out.String(), want)
			}
		})
	}
	var out bytes.Buffer
	value := strings.Repeat("a", 199-len(prefix)) + "\x1b"
	if err := WriteHumanFailure(&out, "issue.get", "fid issue get", Failure{Code: CodeIssueNotFound, Message: "missing", Details: map[string]any{"id": value}}, false); err != nil {
		t.Fatal(err)
	}
	if want := "error: missing\n" + prefix + strings.Repeat("a", 199-len(prefix)) + "…\n"; out.String() != want {
		t.Fatalf("split escape: %q want %q", out.String(), want)
	}
}
