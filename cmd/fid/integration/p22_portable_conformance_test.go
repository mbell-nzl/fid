//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type p22PortableCase struct {
	State    string         `json:"state"`
	Expected map[string]any `json:"expected"`
}

func TestP22PortableRecoveryConformance(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join(p22CorpusDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Cases []string `json:"portable_cases"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if len(manifest.Cases) != 8 {
		t.Fatalf("portable corpus has %d cases, want 8", len(manifest.Cases))
	}
	for _, path := range manifest.Cases {
		data, err := os.ReadFile(filepath.Join(p22CorpusDir, path))
		if err != nil {
			t.Fatal(err)
		}
		var fixture p22PortableCase
		if err := json.Unmarshal(data, &fixture); err != nil {
			t.Fatal(err)
		}
		t.Run(fixture.State, func(t *testing.T) {
			adapter, activations := p22PortableSource(t, fixture.State)
			entries, live := p22PortableThreeWay(t, adapter, activations)
			fid, fake := p22PortableNamedSummary(live, entries), p22PortableFakeSummary(t, entries, adapter.runID)
			fields := "activation_statuses artifact_availability artifact_body_count attention_count attention_link_count run_statuses stop_count tokens_consumed tokens_remaining tokens_state"
			p22AssertKeys(t, fid, fields)
			p22AssertKeys(t, fake, fields)
			p22AssertKeys(t, fixture.Expected, fields)
			if !reflect.DeepEqual(fid, fake) {
				t.Fatalf("complete independent portable summaries differ: fid=%#v fake=%#v", fid, fake)
			}
			for name, actual := range map[string]map[string]any{"fid": fid, "fake": fake} {
				t.Run(name, func(t *testing.T) {
					for field, expected := range fixture.Expected {
						if !reflect.DeepEqual(actual[field], expected) {
							t.Fatalf("%s %s = %#v, want %#v", name, field, actual[field], expected)
						}
					}
				})
			}
		})
	}
}

func p22PortableSource(t *testing.T, state string) (*p22FIDRecoveryAdapter, []string) {
	t.Helper()
	paths := []string{}
	if state == "missing-artifact" {
		paths = []string{"base.txt"}
	}
	a := p22NewFIDRecoveryAdapterWithArtifacts(t, "agent", 3600, paths)
	activations := []string{a.activation["id"].(string)}
	accepted := func(name, operation string, input map[string]any) {
		r := a.Execute(t, p22ConformanceCase{Name: name, Operation: operation, Kind: "accepted", Input: input})
		if r.ExitCode != 0 {
			t.Fatalf("prepare portable %s/%s: %#v", state, operation, r)
		}
	}
	switch state {
	case "active":
		accepted("portable-baseline", "run.observe-budget", nil)
	case "unknown-budget":
	case "interrupted":
		a.p22Observe(t, "terminated")
		a.p22Expire(t)
		r := a.p22Recover(t)
		if r.ExitCode != 0 {
			t.Fatalf("portable recovery: %#v", r)
		}
		a.p22Advance(t, "run.recover", r.Stdout)
	case "abandoned":
		a.p22Observe(t, "terminated")
		accepted("portable-abandon", "run.abandon", nil)
	case "exhausted-budget":
		accepted("portable-exhausted", "run.observe-budget", map[string]any{"value": 10})
	case "deduplicated-attention":
		accepted("portable-stop-first", "stop.raise", nil)
		accepted("portable-stop-second", "stop.raise", nil)
	case "partial-activation":
		prior := p21Run(t, a.binary, a.workspace, "--json", "activation", "status", activations[0])
		resources := []any{}
		for _, value := range p21Path(t, prior, "result", "activation", "resources").([]any) {
			r := value.(map[string]any)
			resources = append(resources, map[string]any{"component": r["component"], "old_identity": r["new_identity"], "new_identity": r["new_identity"].(string) + ":next"})
		}
		input := map[string]any{"backend_identity": "fid", "direction": "activate", "kit_identity": "portable-upgrade", "phase": 21, "expected_prior_activation_id": activations[0], "expected_prior_version": p21Path(t, prior, "result", "activation", "version"), "resources": resources}
		prepared := p21Run(t, a.binary, a.workspace, "--json", "activation", "prepare", "--input", p21Input(t, a.workspace, "portable-upgrade", input), "--idempotency-key", "portable-upgrade")
		activations = append(activations, p21PathString(t, prepared, "result", "activation", "id"))
	case "missing-artifact":
		accepted("portable-artifact-budget", "run.observe-budget", nil)
		gate := p22GateRun(t, a.binary, a.workspace, a.actor, p21Input(t, a.workspace, "portable-artifact-gate", map[string]any{"run_id": a.runID, "registration_id": a.registrationID, "registration_version": a.registrationVersion, "obligation_id": "OBL-p22-corpus", "expected_repository": map[string]any{"head_oid": a.repository["head_oid"], "dirty": a.repository["dirty"]}}), "portable-artifact-gate")
		artifacts := p21Path(t, gate, "result", "receipt", "artifacts").([]any)
		if len(artifacts) != 1 {
			t.Fatalf("portable artifact fixture = %#v", artifacts)
		}
		artifact := artifacts[0].(map[string]any)
		if err := os.Remove(filepath.Join(a.workspace, ".fid", "artifacts", strings.TrimPrefix(artifact["digest"].(string), "sha256:"))); err != nil {
			t.Fatal(err)
		}
		verified := p21Run(t, a.binary, a.workspace, "--json", "artifact", "verify", "--actor", a.actor, "--input", p21Input(t, a.workspace, "portable-missing", map[string]any{"artifact_id": artifact["id"]}), "--idempotency-key", "portable-missing")
		if p21Path(t, verified, "result", "artifact", "availability") != "missing" {
			t.Fatalf("missing fixture = %#v", verified)
		}
	default:
		t.Fatalf("unknown portable state %s", state)
	}
	return a, activations
}

func TestP22PortableCorpusCoversRecoveryStates(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile(filepath.Join(p22CorpusDir, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Portable []string `json:"portable_cases"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	want := []string{"abandoned", "active", "deduplicated-attention", "exhausted-budget", "interrupted", "missing-artifact", "partial-activation", "unknown-budget"}
	var states []string
	for _, path := range manifest.Portable {
		data, err := os.ReadFile(filepath.Join(p22CorpusDir, path))
		if err != nil {
			t.Fatal(err)
		}
		var fixture struct {
			State string `json:"state"`
		}
		if err := json.Unmarshal(data, &fixture); err != nil {
			t.Fatal(err)
		}
		states = append(states, fixture.State)
	}
	sort.Strings(states)
	if !reflect.DeepEqual(states, want) {
		t.Fatalf("portable recovery states = %v, want exact %v", states, want)
	}
}
