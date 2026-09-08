//go:build integration

package integration

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"
)

func TestP22QualifiedCapabilityMatrixMatchesFakeDiscoveryAndExport(t *testing.T) {
	want := p22FakeQualifiedCapabilities(t)
	adapter := p22NewFIDRecoveryAdapter(t)
	discovery := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "capabilities")
	if got := p22CapabilityItems(t, discovery["result"].(map[string]any)["items"]); !reflect.DeepEqual(got, want) {
		t.Fatalf("fid capability discovery = %#v, want shared fake matrix %#v", got, want)
	}

	archive := filepath.Join(t.TempDir(), "recovery.zip")
	p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "export", "create", archive)
	manifest := p22CapabilityManifest(t, p22ArchiveEntries(t, archive)["manifest.json"])
	if got := p22CapabilityItems(t, manifest["capabilities"]); !reflect.DeepEqual(got, want) {
		t.Fatalf("export capability snapshot = %#v, want shared fake matrix %#v", got, want)
	}

	resume := p21Run(t, adapter.binary, adapter.workspace, "--project", adapter.workspace, "--json", "resume", adapter.issueID)
	actions := resume["result"].(map[string]any)["resume"].(map[string]any)["next_actions"].([]any)
	wantResumeCapabilities := p22FakeResumeRequiredCapabilities()
	for _, action := range actions {
		item := action.(map[string]any)
		if item["operation"] == "run.observe" && reflect.DeepEqual(item["required_capabilities"], wantResumeCapabilities) {
			return
		}
	}
	t.Fatalf("resume next_actions = %#v, want run.observe with shared fake required capabilities %#v", actions, wantResumeCapabilities)
}

func p22FakeQualifiedCapabilities(t *testing.T) []any {
	t.Helper()
	var items []any
	if err := json.Unmarshal([]byte(`[
		{"available":true,"capability_id":"attention.v1","limits":{"human_outcome_resolution":false,"state":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},
		{"available":true,"capability_id":"evidence.freshness.v1","limits":{"branch_labels":false,"target":"repository-tree"},"schema_max":1,"schema_min":1,"trust_source":"fid-git-observation"},
		{"available":true,"capability_id":"evidence.gate-definitions.v1","limits":{"executor_policy":"standard-library","network":"allowed","sandbox":"none","write":"allowed"},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},
		{"available":true,"capability_id":"evidence.receipts.v1","limits":{"attestation":false,"producer":"fid","terminal_receipts":true},"schema_max":1,"schema_min":1,"trust_source":"fid-executed"},
		{"available":true,"capability_id":"export.portable.v1","limits":{"create":true,"import":true,"import_mode":"empty-target","validate":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},
		{"available":false,"capability_id":"human.outcomes.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},
		{"available":true,"capability_id":"kit.activation.v1","limits":{"abort_before_commit":true,"automatic_respawn":false,"finalize":true,"recover":true},"schema_max":1,"schema_min":1,"trust_source":"active-kit-activation"},
		{"available":false,"capability_id":"reviews.coverage.v1","limits":{"record":false},"schema_max":1,"schema_min":1,"trust_source":"phase-23-unavailable"},
		{"available":false,"capability_id":"telemetry.trajectory.v1","limits":{"ingest":false},"schema_max":1,"schema_min":1,"trust_source":"phase-25-unavailable"},
		{"available":false,"capability_id":"work.claims.v1","limits":{"max_active_per_issue":1,"parallelism":false,"recover":true},"schema_max":1,"schema_min":1,"trust_source":"phase-24-unavailable"},
		{"available":true,"capability_id":"workflow.state.v1","limits":{"max_envelope_revisions":1,"max_primary_issues":1,"related_issues":false},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"},
		{"available":true,"capability_id":"workflow.transitions.v1","limits":{"execution_bound":true,"terminal_evidence_guard":true},"schema_max":1,"schema_min":1,"trust_source":"fid-journal"}
	]`), &items); err != nil {
		t.Fatalf("decode fake capability matrix: %v", err)
	}
	return items
}

func p22FakeResumeRequiredCapabilities() []any { return []any{"run.observe"} }

func p22CapabilityManifest(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var manifest map[string]any
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("decode export manifest: %v", err)
	}
	return manifest
}

func p22CapabilityItems(t *testing.T, value any) []any {
	t.Helper()
	items, ok := value.([]any)
	if !ok {
		t.Fatalf("capability items = %#v, want array", value)
	}
	return items
}
