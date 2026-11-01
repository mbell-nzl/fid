//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var p22FinalActions = strings.Fields(`artifact.verify attention.resolve claim.acquire claim.expire claim.release claim.renew envelope.activate envelope.complete envelope.create evidence.link gate.register gate.run issue.transition run.abandon run.checkpoint run.close run.force-recover run.observe run.observe-budget run.recover run.start stop.raise stop.resolve`)

func TestP22FinalActivatedActionSetAcceptsProfilesAndGrants(t *testing.T) {
	t.Parallel()
	if len(p22FinalActions) != 23 {
		t.Fatalf("independent action count = %d, want 23", len(p22FinalActions))
	}
	doc, err := os.ReadFile(filepath.Join("testdata", "p22_public_interface.md"))
	if err != nil {
		t.Fatal(err)
	}
	section := p21Section(t, string(doc), "### 2.1 Activated Capabilities And Grant Actions", "## 3. New And Replaced Leaves")
	var documented []string
	for _, row := range regexp.MustCompile("(?m)^\\| `([^`]+)` \\|").FindAllStringSubmatch(section, -1) {
		documented = append(documented, row[1])
	}
	sort.Strings(documented)
	if !reflect.DeepEqual(documented, p22FinalActions) {
		t.Fatalf("documented actions = %v, want independent %v", documented, p22FinalActions)
	}
	a := p22NewFIDRecoveryAdapterAt(t, integrationFIDBinary(t), t.TempDir(), "human", 3600, []string{}, p22FinalActions)
	run := p21Run(t, a.binary, a.workspace, "--json", "run", "get", a.runID)
	envelope := p21Run(t, a.binary, a.workspace, "--json", "envelope", "get", p21PathString(t, run, "result", "run", "envelope_id"))
	want := make([]any, len(p22FinalActions))
	for i, value := range p22FinalActions {
		want[i] = value
	}
	if got := p22ResultPath(envelope, "result.envelope.authority_grants.0.actions"); !reflect.DeepEqual(got, want) {
		t.Fatalf("stored accepted grant = %#v, want all23 %#v", got, want)
	}
	// Every rejected request has otherwise valid content; reuse the unchanged
	// prepared activation and draft input to isolate the action-domain rejection.
	prepared := p22NewFIDActivationAdapter(t, "prepared").(*p22FIDRecoveryAdapter)
	for _, action := range strings.Fields(`unknown.action run.finish activation.prepare activation.component activation.commit activation.abort activation.finalize activation.resource activation.recover activation.recover-complete export.import run.inspect attention.list`) {
		t.Run(action, func(t *testing.T) {
			before := p21JournalEventCount(t, prepared.workspace)
			profile := p22ActivationProfile()
			profile[0].(map[string]any)["capabilities"] = []string{action}
			input := map[string]any{"activation_id": prepared.activation["id"], "expected_version": prepared.activation["version"], "component_type": "authority-profile", "component_id": "invalid-action-profile", "component_version": 1, "expected_prior_digest": nil, "content": profile}
			response := prepared.p22Run(t, "activation", "component", "--input", p21Input(t, prepared.workspace, "reject-"+action, input), "--idempotency-key", "reject-profile-"+action)
			p22AssertActionRejection(t, response)
			if p21JournalEventCount(t, prepared.workspace) != before {
				t.Fatal("rejected profile appended Event")
			}
			data, err := os.ReadFile(filepath.Join(a.workspace, "p22-corpus-envelope.json"))
			if err != nil {
				t.Fatal(err)
			}
			var draft map[string]any
			if err := json.Unmarshal(data, &draft); err != nil {
				t.Fatal(err)
			}
			draft["authority_grants"].([]any)[0].(map[string]any)["actions"] = []string{action}
			before = p21JournalEventCount(t, a.workspace)
			response = a.p22Run(t, "envelope", "create", "--actor", a.actor, "--input", p21Input(t, a.workspace, "reject-grant-"+action, draft), "--idempotency-key", "reject-grant-"+action)
			p22AssertActionRejection(t, response)
			if p21JournalEventCount(t, a.workspace) != before {
				t.Fatal("rejected grant appended Event")
			}
		})
	}
}

func p22AssertActionRejection(t *testing.T, response p22ConformanceResponse) {
	t.Helper()
	if response.ExitCode != 2 || response.Stdout != "" || !strings.Contains(response.Stderr, `"code":"invalid-argument"`) {
		t.Fatalf("closed action rejection = %#v", response)
	}
}
