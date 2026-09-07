//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestP22DualAdapterActivationConformance(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join(p22CorpusDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		ActivationCases []string `json:"activation_cases"`
	}
	if err := json.Unmarshal(contents, &manifest); err != nil || len(manifest.ActivationCases) != 1 {
		t.Fatalf("activation corpus manifest: %v %#v", err, manifest)
	}
	contents, err = os.ReadFile(filepath.Join(p22CorpusDir, manifest.ActivationCases[0]))
	if err != nil {
		t.Fatal(err)
	}
	var cases map[string][]p22ConformanceCase
	if err := json.Unmarshal(contents, &cases); err != nil {
		t.Fatal(err)
	}
	for _, status := range []string{"prepared", "committed"} {
		for _, kind := range []string{"fake", "fid"} {
			t.Run(status+"/"+kind, func(t *testing.T) {
				var adapter p22RecoveryAdapter
				if kind == "fake" {
					adapter = p22NewFakeActivationAdapter(status)
				} else {
					adapter = p22NewFIDActivationAdapter(t, status)
				}
				var selected p22ConformanceCase
				var selectedResponse map[string]any
				for _, fixture := range cases[status] {
					t.Run(fixture.Name, func(t *testing.T) {
						before, events := p22ActivationSnapshot(t, adapter)
						response := adapter.Execute(t, fixture)
						assertP22ConformanceResponse(t, fixture, response)
						after, afterEvents := p22ActivationSnapshot(t, adapter)
						if response.ExitCode != 0 {
							if !reflect.DeepEqual(before, after) || events != afterEvents {
								t.Fatalf("rejected %s mutated state/Events", fixture.Name)
							}
							assertP22ActivationFailure(t, fixture, before, response)
							return
						}
						var envelope map[string]any
						if err := json.Unmarshal([]byte(response.Stdout), &envelope); err != nil {
							t.Fatal(err)
						}
						assertP22ActivationSuccess(t, fixture, before, envelope)
						if fixture.Event {
							if afterEvents != events+1 {
								t.Fatalf("%s Event count: %d -> %d", fixture.Name, events, afterEvents)
							}
							assertP22Replay(t, fixture, envelope, adapter.Execute(t, fixture))
							replayedState, replayedEvents := p22ActivationSnapshot(t, adapter)
							if !reflect.DeepEqual(after, replayedState) || afterEvents != replayedEvents {
								t.Fatalf("%s replay mutated state/Events", fixture.Name)
							}
							if fixture.Operation == "activation.recover" {
								selected, selectedResponse = fixture, envelope
							}
						} else if !reflect.DeepEqual(before, after) || events != afterEvents {
							t.Fatalf("%s read mutated state/Events", fixture.Name)
						}
					})
				}
				if selectedResponse == nil {
					t.Fatal("scenario has no successful recovery selection")
				}
				before, events := p22ActivationSnapshot(t, adapter)
				assertP22Replay(t, selected, selectedResponse, adapter.Execute(t, selected))
				after, afterEvents := p22ActivationSnapshot(t, adapter)
				if !reflect.DeepEqual(before, after) || events != afterEvents {
					t.Fatal("selection replay after completion mutated current state")
				}
			})
		}
	}
}

func p22ActivationSnapshot(t *testing.T, adapter p22RecoveryAdapter) (map[string]any, int) {
	t.Helper()
	if fake, ok := adapter.(*p22FakeRecoveryAdapter); ok {
		return p22CloneMap(fake.model.activation), fake.model.sequence
	}
	fid := adapter.(*p22FIDRecoveryAdapter)
	status := p21Run(t, fid.binary, fid.workspace, "--json", "activation", "status", fid.activation["id"].(string))
	entries, err := os.ReadDir(filepath.Join(fid.workspace, ".fid", "journal"))
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		contents, err := os.ReadFile(filepath.Join(fid.workspace, ".fid", "journal", entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		count += bytes.Count(contents, []byte{'\n'})
	}
	return p21Path(t, status, "result", "activation").(map[string]any), count
}

func assertP22ActivationSuccess(t *testing.T, fixture p22ConformanceCase, before, envelope map[string]any) {
	t.Helper()
	p22AssertKeys(t, envelope, "contract_version event ok operation result")
	result := envelope["result"].(map[string]any)
	fields := "activation"
	if fixture.Event {
		fields += " idempotency_key replayed"
	}
	if fixture.Operation == "activation.resource" {
		fields += " resource"
	}
	p22AssertKeys(t, result, fields)
	a := result["activation"].(map[string]any)
	p22AssertKeys(t, a, "id status kit_identity backend_identity phase direction point_of_no_return version prior_activation_id component_digests resources recovery created_at updated_at")
	if a["id"] != before["id"] || !reflect.DeepEqual(a["prior_activation_id"], before["prior_activation_id"]) || !reflect.DeepEqual(a["component_digests"], before["component_digests"]) {
		t.Fatalf("mutation changed frozen activation identity/content: %#v", a)
	}
	if fixture.Event && p22Number(a["version"]) != p22Number(before["version"])+1 {
		t.Fatalf("activation version did not increment exactly once: %#v", a)
	}
	resources := a["resources"].([]any)
	if len(resources) != 5 {
		t.Fatalf("resource count = %d", len(resources))
	}
	for _, item := range resources {
		r := item.(map[string]any)
		p22AssertKeys(t, r, "component version old_identity new_identity observed_identity state evidence_ref error_code")
		old := p22ActivationResource(before, r["component"].(string))
		if old == nil || !reflect.DeepEqual(old["old_identity"], r["old_identity"]) || old["new_identity"] != r["new_identity"] {
			t.Fatalf("mutation changed resource binding: %#v", r)
		}
	}
	if resource, ok := result["resource"].(map[string]any); ok {
		stored := p22ActivationResource(a, resource["component"].(string))
		if !reflect.DeepEqual(resource, stored) {
			t.Fatalf("resource result and activation disagree: %#v", result)
		}
	}
	if recovery, ok := a["recovery"].(map[string]any); ok {
		p22AssertKeys(t, recovery, "direction actor_attribution reason selected_event_id completed_event_id")
		if fixture.Operation == "activation.recover" && recovery["selected_event_id"] != p22ResultPath(envelope, "event.id") {
			t.Fatal("selection does not bind its Event")
		}
		if fixture.Operation == "activation.recover-complete" {
			if recovery["completed_event_id"] != p22ResultPath(envelope, "event.id") || recovery["selected_event_id"] != p22ResultPath(before, "recovery.selected_event_id") {
				t.Fatal("completion changed selection or omitted completion Event")
			}
		}
	}
}

func assertP22ActivationFailure(t *testing.T, fixture p22ConformanceCase, before map[string]any, response p22ConformanceResponse) {
	t.Helper()
	var envelope map[string]any
	if err := json.Unmarshal([]byte(response.Stderr), &envelope); err != nil {
		t.Fatal(err)
	}
	p22AssertKeys(t, envelope, "contract_version error ok operation")
	failure := envelope["error"].(map[string]any)
	p22AssertKeys(t, failure, "category code details message retryable subject")
	if failure["retryable"] != false {
		t.Fatal("activation failure unexpectedly retryable")
	}
	details := failure["details"].(map[string]any)
	switch fixture.FailureCode {
	case "activation-version-conflict":
		p22AssertKeys(t, details, "activation_id expected_version actual_version")
		expected, actual := p22ExpectedActivationVersionConflict(fixture, before)
		if details["activation_id"] != before["id"] || p22Number(details["expected_version"]) != expected || p22Number(details["actual_version"]) != actual || failure["category"] != "conflict" {
			t.Fatalf("version conflict = %#v", failure)
		}
	case "activation-direction-conflict":
		p22AssertKeys(t, details, "activation_id point_of_no_return required_direction")
		if details["activation_id"] != before["id"] || details["point_of_no_return"] != before["point_of_no_return"] || details["required_direction"] != p22ActivationDirection(before) || failure["category"] != "precondition" {
			t.Fatalf("direction conflict = %#v", failure)
		}
	case "activation-incomplete":
		p22AssertKeys(t, details, "activation_id components reasons")
		if details["activation_id"] != before["id"] || !reflect.DeepEqual(details["components"], []any{"backend", "kit", "overlay", "project", "provider"}) || failure["category"] != "precondition" {
			t.Fatalf("incomplete agreement = %#v", failure)
		}
		reasons, ok := details["reasons"].([]any)
		if !ok || len(reasons) == 0 {
			t.Fatal("incomplete agreement omitted typed Reasons")
		}
		for _, item := range reasons {
			reason := item.(map[string]any)
			p22AssertKeys(t, reason, "code detail subject_id subject_type")
			if reason["code"] != "activation-incomplete" || reason["subject_id"] != before["id"] || reason["subject_type"] != "activation" || reason["detail"] == "" {
				t.Fatalf("incomplete agreement Reason = %#v", reason)
			}
		}
	}
	if fixture.FailureCode != "invalid-argument" && (p22ResultPath(failure, "subject.id") != before["id"] || p22ResultPath(failure, "subject.type") != "activation") {
		t.Fatalf("activation error subject = %#v", failure["subject"])
	}
}

func p22ExpectedActivationVersionConflict(fixture p22ConformanceCase, before map[string]any) (int, int) {
	if expected, ok := fixture.Input["expected_activation_version"]; ok {
		return p22Number(expected), p22Number(before["version"])
	}
	if fixture.Operation == "activation.resource" {
		component := "backend"
		if value, ok := fixture.Input["component"].(string); ok {
			component = value
		}
		return p22Number(fixture.Input["expected_version"]), p22Number(p22ActivationResource(before, component)["version"])
	}
	if values, ok := fixture.Input["resource_versions"].([]any); ok {
		for _, item := range values {
			v := item.(map[string]any)
			r := p22ActivationResource(before, v["component"].(string))
			if p22Number(v["version"]) != p22Number(r["version"]) {
				return p22Number(v["version"]), p22Number(r["version"])
			}
		}
	}
	return p22Number(fixture.Input["expected_version"]), p22Number(before["version"])
}

func p22AssertKeys(t *testing.T, object map[string]any, fields string) {
	t.Helper()
	want := strings.Fields(fields)
	got := make([]string, 0, len(object))
	for field := range object {
		got = append(got, field)
	}
	sort.Strings(want)
	sort.Strings(got)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("closed fields = %v, want %v", got, want)
	}
}
