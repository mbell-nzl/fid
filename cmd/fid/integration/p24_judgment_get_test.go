//go:build integration

package integration

import (
	"encoding/json"
	"errors"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

func TestP24JudgmentGetRecords(t *testing.T) {
	t.Parallel()
	a := p24FailureJudgmentAdapter(t)
	subject := map[string]any{"type": "issue", "id": a.issueID}
	inputs := []struct {
		family, member, id string
		input              map[string]any
	}{
		{"concern", "concern", "CON-1", map[string]any{"id": "CON-1", "subject": subject, "scope": map[string]any{"paths": []string{}}, "severity": "major", "summary": "Concern summary", "reasoning": "Recorded concern reasoning", "references": []any{p23DecisionSource(a, true)}}},
		{"deviation", "deviation", "DEV-1", map[string]any{"id": "DEV-1", "subject": subject, "expected": "Expected behavior", "observed": "Observed behavior", "authority": map[string]any{"actor_id": a.actor, "reference": nil}, "disposition": "pending-decision", "reasoning": "Recorded deviation reasoning", "references": []any{p23DecisionSource(a, true)}, "supersedes_deviation_id": nil}},
		{"human-outcome", "outcome", "OUT-1", map[string]any{"id": "OUT-1", "class": "security", "subject": subject, "result": "approved", "scope": map[string]any{"issue_id": a.issueID}, "reasoning": "Recorded Human reasoning", "references": []any{p23DecisionSource(a, true)}}},
	}
	for _, row := range inputs {
		p23AssertReviewSuccess(t, row.family+" record", p23RunInput(t, a, row.family, "record", a.actor, "get-"+row.id, row.input))
	}
	before := p23PortableJournalEvents(t, a.workspace)
	for _, row := range inputs {
		t.Run(row.family, func(t *testing.T) {
			listOut, listErr, listExit := closureRunCommand(a.workspace, a.binary, "--json", row.family, "list")
			if listExit != nil || listErr != "" {
				t.Fatalf("list: %v %s", listExit, listErr)
			}
			items := p21Path(t, p23Envelope(t, listOut), "result", "items").([]any)
			if len(items) != 1 {
				t.Fatalf("list items = %#v, want the recorded item", items)
			}
			out, stderr, err := closureRunCommand(a.workspace, a.binary, "--json", row.family, "get", row.id)
			if err != nil || stderr != "" {
				t.Fatalf("get: %v stdout=%q stderr=%q", err, out, stderr)
			}
			envelope := p23Envelope(t, out)
			result := envelope["result"].(map[string]any)
			_, hasEvent := envelope["event"]
			if envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != true || !hasEvent || envelope["operation"] != row.family+".get" || envelope["event"] != nil || len(result) != 1 || !reflect.DeepEqual(result[row.member], items[0]) {
				t.Fatalf("get envelope = %#v, want same list item and event:null", envelope)
			}
			human, humanErr, err := closureRunCommand(a.workspace, a.binary, row.family, "get", row.id)
			if err != nil || humanErr != "" || !strings.HasPrefix(human, "Purpose: ") || !strings.Contains(human, row.id) || strings.Contains(human, `"result"`) {
				t.Fatalf("Human get: %v %q %q", err, human, humanErr)
			}
			// Globals select the project even from outside the workspace.
			selected, selectedErr, err := closureRunCommand(t.TempDir(), a.binary, "--project", a.workspace, "--json", row.family, "get", row.id)
			if err != nil || selectedErr != "" || !reflect.DeepEqual(p23Envelope(t, selected), envelope) {
				t.Fatalf("selected project get: %v %q %q", err, selected, selectedErr)
			}
		})
	}
	if after := p23PortableJournalEvents(t, a.workspace); !reflect.DeepEqual(after, before) {
		t.Fatal("get leaves changed journal events")
	}
}

func TestP24JudgmentGetRejections(t *testing.T) {
	t.Parallel()
	binary := integrationFIDBinary(t)
	workspace := t.TempDir()
	p21Run(t, binary, workspace, "--json", "init", "--prefix", "GET")
	before := p23PortableJournalEvents(t, workspace)
	for _, row := range []struct{ family, code, idLabel string }{
		{"concern", "concern-not-found", "<concern-id>"},
		{"deviation", "deviation-not-found", "<deviation-id>"},
		{"human-outcome", "outcome-not-found", "<outcome-id>"},
	} {
		t.Run(row.family, func(t *testing.T) {
			missingID := "ID-missing"
			jsonOut, jsonErr, jsonExit := p24GetRunExit(t, binary, workspace, "--json", row.family, "get", missingID)
			if jsonExit != 3 || jsonOut != "" {
				t.Fatalf("not found: exit=%d stdout=%q stderr=%q", jsonExit, jsonOut, jsonErr)
			}
			var failure map[string]any
			if err := json.Unmarshal([]byte(jsonErr), &failure); err != nil {
				t.Fatal(err)
			}
			detail := failure["error"].(map[string]any)
			if failure["contract_version"] != "fid-cli.v1" || failure["ok"] != false || failure["operation"] != row.family+".get" || detail["code"] != row.code || detail["category"] != "not-found" || detail["retryable"] != false || !reflect.DeepEqual(detail["subject"], map[string]any{"type": row.family, "id": missingID}) || !reflect.DeepEqual(detail["details"], map[string]any{"id": missingID}) {
				t.Fatalf("not-found envelope = %#v", failure)
			}
			humanOut, humanErr, humanExit := p24GetRunExit(t, binary, workspace, row.family, "get", missingID)
			if humanExit != 3 || humanOut != "" || !strings.HasPrefix(humanErr, "error: "+detail["message"].(string)+"\nreason: "+row.code) {
				t.Fatalf("not-found pairing: exit=%d stdout=%q stderr=%q", humanExit, humanOut, humanErr)
			}
			for _, grammar := range []struct {
				args     []string
				argument string
				value    any
				allowed  []string
			}{
				{nil, "arguments", nil, []string{row.idLabel}},
				{[]string{missingID, "extra"}, "arguments", "extra", []string{row.idLabel}},
				{[]string{missingID, "--bogus"}, "--bogus", "--bogus", []string{}},
			} {
				args := append([]string{row.family, "get"}, grammar.args...)
				for _, dir := range []string{workspace, t.TempDir()} {
					out, stderr, exit := p24GetRunExit(t, binary, dir, append([]string{"--json"}, args...)...)
					if exit != 2 || out != "" {
						t.Fatalf("grammar %v: exit=%d stdout=%q stderr=%q", args, exit, out, stderr)
					}
					var envelope map[string]any
					if err := json.Unmarshal([]byte(stderr), &envelope); err != nil {
						t.Fatal(err)
					}
					failure := envelope["error"].(map[string]any)
					wantDetails := map[string]any{"argument": grammar.argument, "value": grammar.value, "allowed": grammar.allowed}
					wantBytes, err := json.Marshal(wantDetails)
					if err != nil {
						t.Fatal(err)
					}
					var want map[string]any
					if err := json.Unmarshal(wantBytes, &want); err != nil {
						t.Fatal(err)
					}
					if envelope["contract_version"] != "fid-cli.v1" || envelope["ok"] != false || envelope["operation"] != row.family+".get" || failure["code"] != "invalid-argument" || failure["category"] != "invalid-input" || failure["retryable"] != false || failure["subject"] != nil || !reflect.DeepEqual(failure["details"], want) {
						t.Fatalf("grammar envelope = %#v, want details %#v", envelope, want)
					}
					hOut, hErr, hExit := p24GetRunExit(t, binary, dir, args...)
					if hExit != 2 || hOut != "" || !strings.HasPrefix(hErr, "error: "+failure["message"].(string)+"\nreason: invalid-argument;") {
						t.Fatalf("grammar pairing: exit=%d stdout=%q stderr=%q", hExit, hOut, hErr)
					}
				}
			}
		})
	}
	if after := p23PortableJournalEvents(t, workspace); !reflect.DeepEqual(after, before) {
		t.Fatal("get rejections changed journal events")
	}
}

func p24GetRunExit(t *testing.T, binary, workspace string, args ...string) (string, string, int) {
	t.Helper()
	out, stderr, err := closureRunCommand(workspace, binary, args...)
	if err == nil {
		return out, stderr, 0
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatal(err)
	}
	return out, stderr, exit.ExitCode()
}
