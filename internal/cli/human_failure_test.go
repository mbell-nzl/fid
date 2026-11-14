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

		{"conflict", CodeIdempotencyConflict, map[string]any{"differing_fields": []string{"--actor", "source"}}, "artifact.verify", "fid artifact verify", "reason: idempotency-conflict; differing fields: --actor,source\nnext: resend the original input, or use a new --idempotency-key"},
		{"transition", CodeTransitionNotPermitted, map[string]any{"unsatisfied": []string{"first", "second"}}, "issue.transition", "fid issue transition", "reason: transition-not-permitted; first,second"},

		{"input", CodeInvalidArgument, map[string]any{"argument": "input.a", "allowed": []string{"a", "b"}, "value": "ignored"}, "artifact.verify", "fid artifact verify", "reason: invalid-argument; argument=input.a, allowed=a|b\nnext: see fid artifact verify --help; template: fid artifact verify --example"},
		{"no input", CodeInvalidArgument, map[string]any{"argument": "--title", "allowed": []string{"string"}}, "issue.create", "fid issue create", "reason: invalid-argument; argument=--title, allowed=string"},

		{"controls", CodeIssueNotFound, map[string]any{"id": "ISS-\n\r\t\x1b"}, "issue.get", "fid issue get", `reason: issue-not-found; id=ISS-\n\r\t\u001b`},
		{"show not found", CodeIssueNotFound, map[string]any{"id": "FID-nope"}, "issue.show", "fid show", "reason: issue-not-found; id=FID-nope\nnext: fid issue list"},
		{"get not found", CodeIssueNotFound, map[string]any{"id": "FID-nope"}, "issue.get", "fid issue get", "reason: issue-not-found; id=FID-nope"},
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
