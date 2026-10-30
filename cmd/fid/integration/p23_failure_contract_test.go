//go:build integration

package integration

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

// These expectations are copied from public-interface.md section 5 and 5.1.
// Each existing real-process producer probe calls this through
// assertP22ProcessFailure, so the code, envelope, and closed details travel
// through the same binary boundary.
var p23FailureContracts = map[string]struct {
	category, subject, details, producers string
}{
	"review-not-independent":      {"precondition", "", "review_actor_id conflicting_subject conflicting_actor_id", "review.record evidence.link"},
	"review-verdict-invalid":      {"invalid-input", "", "verdict reason", "review.record"},
	"review-changes-required":     {"precondition", "", "review_id obligation_id blocking_finding_ids", "evidence.link"},
	"review-rework-exhausted":     {"precondition", "envelope", "budget stop_ids attention_ids event replayed", "claim.acquire run.start"},
	"finding-already-disposed":    {"conflict", "finding", "finding_id state", "finding.address review.record"},
	"finding-not-found":           {"not-found", "finding", "id", "finding.address review.record concern.dispose"},
	"review-not-found":            {"not-found", "review", "id", "review.get evidence.link"},
	"concern-not-found":           {"not-found", "concern", "id", "concern.dispose"},
	"deviation-not-found":         {"not-found", "deviation", "id", "deviation.record"},
	"outcome-not-found":           {"not-found", "human-outcome", "id", "evidence.link"},
	"requirement-not-found":       {"not-found", "requirement", "id", "requirement.get requirement.link"},
	"concern-already-disposed":    {"conflict", "concern", "concern_id disposed_at", "concern.dispose"},
	"claim-mode-invalid":          {"precondition", "claim", "claim_id mode required_mode", "review.record finding.address run.recover run.abandon"},
	"protected-decision-required": {"authority", "envelope", "envelope_id decision_class required_capabilities affected_actions stop_ids attention_ids event replayed", "envelope.complete issue.transition"},
	"human-outcome-required":      {"precondition", "envelope", "envelope_id obligation_id outcome_class stop_ids attention_ids event replayed", "envelope.complete issue.transition"},
	"outcome-not-permitted":       {"precondition", "", "obligation_id obligation_type outcome_id result", "evidence.link"},
	"evidence-type-not-accepted":  {"precondition", "", "obligation_id evidence_type accepted_evidence_types", "evidence.link"},
	"source-digest-mismatch":      {"precondition", "", "display_label path expected_digest observed_digest", "decision.record-reference requirement.record-reference human-outcome.record concern.record deviation.record review.record run.observe-budget"},
	"repository-mismatch":         {"conflict", "", "expected observed", "review.record run.start run.close run.abandon gate.run"},
}

func assertP23AddedFailureContract(t *testing.T, operation, code string, envelope map[string]any) {
	t.Helper()
	want, found := p23FailureContracts[code]
	if !found {
		return
	}
	if !strings.Contains(" "+want.producers+" ", " "+operation+" ") {
		t.Fatalf("%s is not an accepted producer of %s", operation, code)
	}
	if !p23SameKeys(envelope, "contract_version error ok operation") {
		t.Fatalf("%s envelope has extra or missing keys: %#v", code, envelope)
	}
	value := envelope["error"].(map[string]any)
	if !p23SameKeys(value, "category code details message retryable subject") || value["category"] != want.category || value["retryable"] != false {
		t.Fatalf("%s failure envelope differs from accepted category, retryability, or keys: %#v", code, value)
	}
	details, ok := value["details"].(map[string]any)
	if !ok || !p23SameKeys(details, want.details) {
		t.Fatalf("%s failure details are not closed: %#v", code, value["details"])
	}
	if want.subject == "" {
		if value["subject"] != nil {
			t.Fatalf("%s has unexpected subject: %#v", code, value["subject"])
		}
		return
	}
	subject, ok := value["subject"].(map[string]any)
	if !ok || !p23SameKeys(subject, "type id") || subject["type"] != want.subject || subject["id"] == nil {
		t.Fatalf("%s subject = %#v, want %s identity", code, value["subject"], want.subject)
	}
}

func p23SameKeys(value map[string]any, names string) bool {
	want := strings.Fields(names)
	got := make([]string, 0, len(value))
	for name := range value {
		got = append(got, name)
	}
	sort.Strings(got)
	sort.Strings(want)
	return reflect.DeepEqual(got, want)
}
