//go:build integration

package integration

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase23LeavesRejectUnknownRequestFields(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
		[]string{"claim.acquire", "envelope.activate", "envelope.create", "gate.register", "review.record", "run.start"}, "mutation")

	before := p21JournalEventCount(t, a.workspace)
	input := p23ReviewInput(a, a.runID, a.runVersion, "pass", []any{}, []any{})
	input["unaccepted_field"] = true
	rejected := p23ReviewRecord(t, a, a.actor, "unknown-review-field", input)
	assertP22ProcessFailure(t, "review.record", "invalid-argument", 2, rejected)
	if got := p21JournalEventCount(t, a.workspace); got != before {
		t.Fatalf("unknown field appended Event(s): %d -> %d", before, got)
	}

	listed := a.p22Run(t, "review", "list")
	p23AssertReviewSuccess(t, "review list", listed)
	result := p23Envelope(t, listed.Stdout)
	if result["operation"] != "review.list" {
		t.Fatalf("review list operation = %#v, want review.list", result["operation"])
	}
}

func TestPhase23MutationsRejectMalformedObjectsWithoutEvents(t *testing.T) {
	t.Parallel()
	a := p22NewFIDRecoveryAdapterWithClaimModeAt(t, integrationFIDBinary(t), t.TempDir(), "agent", 3600, []string{},
		[]string{"claim.acquire", "concern.dispose", "concern.record", "decision.record-reference", "deviation.record", "envelope.activate", "envelope.complete", "envelope.create", "evidence.link", "finding.address", "gate.register", "human-outcome.record", "review.record", "run.start"}, "mutation")
	operations := [][2]string{
		{"review", "record"}, {"finding", "address"}, {"concern", "record"}, {"concern", "dispose"},
		{"deviation", "record"}, {"decision", "record-reference"}, {"human-outcome", "record"},
		{"requirement", "record-reference"}, {"requirement", "link"},
		{"envelope", "create"}, {"envelope", "activate"}, {"envelope", "complete"},
		{"claim", "acquire"}, {"evidence", "link"},
	}
	cases := []struct{ name, payload string }{
		{"unknown", `{"unaccepted_field":true}`},
		{"wrong-nullability", `null`},
		{"trailing-json", `{} {}`},
		{"duplicate-key", `{"id":"one","id":"two"}`},
	}
	for _, operation := range operations {
		for _, testCase := range cases {
			name := operation[0] + "." + operation[1] + "/" + testCase.name
			t.Run(name, func(t *testing.T) {
				path := filepath.Join(t.TempDir(), "input.json")
				if err := os.WriteFile(path, []byte(testCase.payload), 0o600); err != nil {
					t.Fatal(err)
				}
				before := p21JournalEventCount(t, a.workspace)
				stdout, stderr, err := closureRunCommand(a.workspace, a.binary, "--project", a.workspace, "--json", operation[0], operation[1], "--input", path, "--actor", a.actor, "--idempotency-key", "invalid-"+testCase.name)
				if err == nil || stdout != "" || !p23HasJSON(stderr, `"code":"invalid-argument"`) {
					t.Fatalf("%s malformed input = stdout %q stderr %q error %v", name, stdout, stderr, err)
				}
				if after := p21JournalEventCount(t, a.workspace); after != before {
					t.Fatalf("%s appended Event(s): %d -> %d", name, before, after)
				}
			})
		}
	}
}
