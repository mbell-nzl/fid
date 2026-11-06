package journal

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestCurrentActivationCompletionRejectsPendingResource(t *testing.T) {
	t.Parallel()
	conn := openReplayDB(t, filepath.Join(t.TempDir(), "current.sqlite"))
	defer closeProjectDB(t, conn)
	seedActivationRecoveryProjection(t, conn, "activation.recovery-completed")
	if _, err := conn.Exec(`UPDATE activation_resources SET state='pending',observed_identity=NULL,evidence_ref=NULL,error_code=NULL WHERE activation_id='ACT-projection' AND component='kit'`); err != nil {
		t.Fatal(err)
	}
	event := activationRecoveryProjectionEvent("activation.recovery-completed")
	for i := range event.Effects {
		if event.Effects[i].Table != "idempotency_results" {
			continue
		}
		result := projectionFixtureDecode(event.Effects[i].Values["result_json"].(string)).(map[string]any)
		for _, value := range result["activation"].(map[string]any)["resources"].([]any) {
			r := value.(map[string]any)
			if r["component"] == "kit" {
				r["state"], r["observed_identity"], r["evidence_ref"], r["error_code"] = "pending", nil, nil, nil
			}
		}
		event.Effects[i].Values["result_json"] = projectionFixtureJSON(result)
	}
	if err := DryRunProjectEvent(conn, event); err == nil {
		t.Fatal("completion admitted without recorded resource agreement")
	}
}

func TestCurrentValueDomains(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, kind, raw string }{
		{"byte-bound", "Text", `"` + strings.Repeat("界", 1400) + `"`},
		{"utc-time", "Timestamp", `"2026-09-09T00:00:00+02:00"`},
		{"invalid-utf8", "Text", "\"\xff\""},
		{"ordinary-reference", "Authority", `{"mode":"ordinary","actor_id":"actor","reason":"reason","reference":"force-only","grant_ref":null}`},
		{"unknown-activated-capability", "ActorProfile", `{"actor_id":"actor","actor_type":"agent","producer":"fid","provenance_ref":"test","session_ref":null,"capabilities":["invented.action"]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v, err := decodeCurrentJSON(tc.raw)
			if err == nil && validateCurrentValue(tc.kind, v) == nil {
				t.Fatal("invalid current scalar domain accepted")
			}
		})
	}
}

func TestCurrentNonJSONDomainsRejectMalformedIdentity(t *testing.T) {
	t.Parallel()
	for column, value := range map[string]any{"actor_id": "", "content_digest": "not-a-digest", "semantic_key": "plain-text", "declared_path": "../escape", "executor_id": "bad\nidentity"} {
		t.Run(column, func(t *testing.T) {
			if validateCurrentColumnDomain(column, value) == nil {
				t.Fatal("invalid non-JSON domain accepted")
			}
		})
	}
}

func TestCurrentReconciliationDigestBindsCapturedContent(t *testing.T) {
	t.Parallel()
	value := projectionReconciliation()
	raw := projectionFixtureJSON(value)
	decoded, err := decodeCurrentJSON(raw)
	if err != nil || validateCurrentValue("Reconciliation", decoded) != nil {
		t.Fatal("valid control rejected")
	}
	value["digest"] = "sha256:" + strings.Repeat("0", 64)
	decoded, err = decodeCurrentJSON(projectionFixtureJSON(value))
	if err == nil && validateCurrentValue("Reconciliation", decoded) == nil {
		t.Fatal("unbound reconciliation digest accepted")
	}
}

func TestCurrentClosedValuesRejectIncompleteAndInvalidFacts(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name, family, table, column, raw string
		values                           map[string]any
	}{
		{"checkpoint-required-fields", "run.checkpointed", "run_observations", "observation_json", `{"id":"OBS-test","run_id":"RUN-test","kind":"checkpoint"}`, map[string]any{"kind": "checkpoint", "run_id": "RUN-test"}},
		{"authority-required-nullables", "run.recovered", "run_recoveries", "authority_json", `{"mode":"ordinary","actor_id":"actor-test","reason":"recovery"}`, nil},
		{"authority-unknown-replaces-required", "run.recovered", "run_recoveries", "authority_json", `{"mode":"ordinary","actor_id":"actor-test","reason":"recovery","extra":null,"reference":null}`, nil},
		{"reconciliation-required-fields", "run.recovered", "run_recoveries", "reconciliation_json", `{"repository":{}}`, nil},
		{"condition-required-nullables", "stop.raised", "stops", "resolution_condition_json", `{"type":"work-ended","subject":{"type":"run","id":"RUN-test"}}`, nil},
		{"condition-subject-binding", "stop.raised", "stops", "resolution_condition_json", `{"type":"work-ended","subject":{"type":"run","id":"RUN-other"},"budget_name":null,"capability_id":null}`, map[string]any{"subject_type": "run", "subject_id": "RUN-test"}},
		{"claim-scope-closed", "run.recovered", "claims", "scope_json", `{"paths":[],"extra":true}`, nil},
		{"repository-closed", "run.started", "runs", "start_repository_json", `{"extra":true}`, nil},
		{"gate-definition-closed", "activation.component-staged", "activation_gate_definitions", "definition_json", `{"extra":true}`, nil},
		{"activation-evidence-closed", "activation.component-staged", "activation_components", "evidence_json", `[{"extra":true}]`, map[string]any{"component_type": "authority-profile"}},
		{"result-other-family-closed", "run.recovered", "idempotency_results", "result_json", `{"replayed":false,"extra":true}`, nil},
		{"stop-details-types", "activation.recovery-selected", "stops", "details_json", `{"activation_id":"ACT-test","point_of_no_return":"false","components":42}`, map[string]any{"code": "activation-incomplete"}},
		{"budget-wrong-scalar-must-not-panic", "run.started", "run_budget_bindings", "budget_json", `{"name":"attempts","scope":"envelope","unit":"attempts","maximum":"5","warning_threshold":null,"exhaustion_action":"budget-exhausted","consumed":null,"remaining":null}`, nil},
		{"budget-positive-maximum", "run.started", "run_budget_bindings", "budget_json", `{"name":"attempts","scope":"envelope","unit":"attempts","maximum":-1,"warning_threshold":null,"exhaustion_action":"budget-exhausted","consumed":null,"remaining":null}`, nil},
		{"adoption-parent-path", "run.recovered", "run_recoveries", "adopted_paths_json", `[".."]`, nil},
		{"adoption-nul-path", "run.recovered", "run_recoveries", "adopted_paths_json", `["file\u0000.txt"]`, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if value := recover(); value != nil {
					t.Errorf("invalid captured fact panicked instead of rejecting: %v", value)
				}
			}()
			effect := Effect{Table: tc.table, Op: "insert", PK: map[string]any{"id": "OBS-test"}, Values: tc.values}
			if err := validateCurrentJSON(Event{EventType: tc.family}, effect, tc.column, tc.raw); err == nil {
				t.Fatal("incomplete or invalid captured fact accepted")
			}
		})
	}
}

// A .fid/ in a worktree is an ordinary path (D4 / F1, decided 6 November
// 2026), so recovery may adopt it. Paths in .git stay protected.
func TestCurrentAdoptedPathsAcceptALeftoverFidDirectory(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, raw string
		accepted  bool
	}{
		{"fid-directory", `[".fid"]`, true},
		{"fid-files", `[".fid/config.toml",".fid/fid.sqlite"]`, true},
		{"fid-and-other-files", `[".fid/journal/2026-10-01.jsonl","base.txt"]`, true},
		{"git-directory", `[".git"]`, false},
		{"git-ledger-file", `[".git/fid/fid.sqlite"]`, false},
		{"fid-and-git-files", `[".fid/config.toml",".git/fid/config.toml"]`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateCurrentJSON(Event{EventType: "run.recovered"}, Effect{Table: "run_recoveries", Op: "insert", PK: map[string]any{"id": "REC-test"}}, "adopted_paths_json", tc.raw)
			if tc.accepted && err != nil {
				t.Fatalf("adopted paths %s rejected: %v", tc.raw, err)
			}
			if !tc.accepted && (err == nil || !strings.Contains(err.Error(), "protected adopted path")) {
				t.Fatalf("adopted paths %s error = %v, want protected adopted path", tc.raw, err)
			}
		})
	}
}

const task9ConditionReference = "74c83b7033729785eb99403ecec9dcde30ac4280"

func TestTask9ConditionDispatchRetainsPinnedPhase22Bytes(t *testing.T) {
	t.Parallel()
	canonical := map[string]string{
		"executor-terminated":  `{"type":"executor-terminated","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null}`,
		"run-recovered":        `{"type":"run-recovered","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null}`,
		"activation-recovered": `{"type":"activation-recovered","subject":{"type":"activation","id":"ACT-pinned"},"budget_name":null,"capability_id":null}`,
		"capability-available": `{"type":"capability-available","subject":{"type":"activation","id":"ACT-pinned"},"budget_name":null,"capability_id":"attention.v1"}`,
		"budget-observed":      `{"type":"budget-observed","subject":{"type":"run","id":"RUN-pinned"},"budget_name":"attempts","capability_id":null}`,
		"work-ended":           `{"type":"work-ended","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null}`,
	}
	if task9ConditionReference == "" {
		t.Fatal("pre-change reference is required")
	}
	for conditionType, raw := range canonical {
		t.Run(conditionType, func(t *testing.T) {
			if err := validateCurrentJSON(Event{EventType: "stop.raised"}, Effect{Table: "stops", Values: map[string]any{"subject_type": subjectTypeForPinnedCondition(conditionType), "subject_id": subjectIDForPinnedCondition(conditionType)}}, "resolution_condition_json", raw); err != nil {
				t.Fatalf("pinned %s Condition rejected: %v", conditionType, err)
			}
			mutated := strings.TrimSuffix(raw, "}") + `,"decision_class":"release"}`
			if err := validateCurrentJSON(Event{EventType: "stop.raised"}, Effect{Table: "stops", Values: map[string]any{"subject_type": subjectTypeForPinnedCondition(conditionType), "subject_id": subjectIDForPinnedCondition(conditionType)}}, "resolution_condition_json", mutated); err == nil {
				t.Fatalf("mutated %s Condition accepted", conditionType)
			}
		})
	}
}

func TestTask9ConditionDispatchPreservesNullableBindings(t *testing.T) {
	t.Parallel()
	cases := []string{
		`{"type":"work-ended","subject":{"type":"run","id":"RUN-pinned"},"budget_name":"attempts","capability_id":null}`,
		`{"type":"budget-observed","subject":{"type":"run","id":"RUN-pinned"},"budget_name":null,"capability_id":null}`,
		`{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":"attempts","capability_id":null,"decision_class":"release"}`,
	}
	for _, raw := range cases {
		value, err := decodeCurrentJSON(raw)
		if err != nil {
			t.Fatal(err)
		}
		if err := validateCurrentValue("Condition", value); err == nil || err.Error() != "condition nullable binding mismatch" {
			t.Fatalf("binding error = %v", err)
		}
	}
}

func TestTask9DecisionConditionShapeAndNonStringDiscriminant(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		raw  string
		ok   bool
	}{
		{"canonical", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"decision_class":"release"}`, true},
		{"wrong-field", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"outcome_class":"release"}`, false},
		{"null-own-field", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"decision_class":null}`, false},
		{"two-new-fields", `{"type":"decision-satisfied","subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"decision_class":"release","outcome_class":"release"}`, false},
		{"non-string-type", `{"type":42,"subject":{"type":"envelope","id":"ENV-pinned"},"budget_name":null,"capability_id":null,"decision_class":"release"}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Fatalf("invalid Condition panicked: %v", recovered)
				}
			}()
			value, err := decodeCurrentJSON(tc.raw)
			if err != nil {
				t.Fatal(err)
			}
			err = validateCurrentValue("Condition", value)
			if (err == nil) != tc.ok {
				t.Fatalf("validation error = %v, want accepted=%v", err, tc.ok)
			}
		})
	}
}

func subjectTypeForPinnedCondition(conditionType string) string {
	if conditionType == "activation-recovered" || conditionType == "capability-available" {
		return "activation"
	}
	return "run"
}

func subjectIDForPinnedCondition(conditionType string) string {
	if subjectTypeForPinnedCondition(conditionType) == "activation" {
		return "ACT-pinned"
	}
	return "RUN-pinned"
}
