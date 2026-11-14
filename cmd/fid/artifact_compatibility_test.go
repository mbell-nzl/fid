package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestArtifactVerifyObservesLossAndRestorationWithoutIdentityMutation(t *testing.T) {
	t.Parallel()
	f := leanJourneyFixture(t, leanHelperDefinition(t, "pass", "1m", []string{"coverage.out"}))
	checked := f.run("--json", "check", "lean", "--hold", "HLD-aaaaaa", "--actor", "agent:claude-builder", "--idempotency-key", "artifact-source")
	if checked.err != nil {
		t.Fatalf("check: %#v", checked)
	}
	before := latestLeanRecord(t, f)
	artifact := before.Result.Artifacts[0]
	path := filepath.Join(config.LedgerDir(f.workspace), "artifacts", artifact.Digest[len("sha256:"):])
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(t.TempDir(), "verify.json")
	writeActivationInput(t, input, map[string]any{"artifact_id": artifact.ID})
	var identity map[string]any
	for _, row := range []struct {
		key, actor, availability string
		changed                  bool
	}{
		{"missing", "agent:builder", "missing", true}, {"restored", "human:mark", "available", true}, {"unchanged", "agent:builder", "available", false},
	} {
		if row.key == "missing" {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		}
		if row.key == "restored" {
			if err := os.WriteFile(path, body, 0o600); err != nil {
				t.Fatal(err)
			}
		}
		args := []string{"--json", "artifact", "verify", "--actor", row.actor, "--input", input, "--idempotency-key", row.key}
		count := startEventCount(t, f.workspace)
		got := f.run(args...)
		if got.err != nil || got.stderr != "" {
			t.Fatalf("%s: %#v", row.key, got)
		}
		result := decodeCanonicalResult(t, got.stdout)
		metadata := result["artifact"].(map[string]any)
		assertExactKeys(t, metadata, []string{"id", "digest", "media_type", "size_bytes", "retention_class", "availability", "location_policy", "created_at", "expires_at"})
		wantMetadata := map[string]any{"id": artifact.ID, "digest": artifact.Digest, "size_bytes": json.Number(strconv.FormatInt(artifact.SizeBytes, 10)), "media_type": "application/octet-stream", "retention_class": "evidence", "availability": row.availability, "location_policy": "content-addressed", "created_at": before.Result.FinishedAt, "expires_at": nil}
		if !reflect.DeepEqual(metadata, wantMetadata) {
			t.Fatalf("Artifact metadata=%#v, want %#v", metadata, wantMetadata)
		}
		if metadata["availability"] != row.availability || result["changed"] != row.changed || result["replayed"] != false {
			t.Fatalf("%s: %s", row.key, got.stdout)
		}
		delete(metadata, "availability")
		if identity == nil {
			identity = metadata
		} else if !reflect.DeepEqual(identity, metadata) {
			t.Fatalf("identity changed: %v / %v", identity, metadata)
		}
		if after := latestLeanRecord(t, f); !reflect.DeepEqual(before, after) {
			t.Fatalf("check result changed: %#v / %#v", before, after)
		}
		if startEventCount(t, f.workspace) != count+1 {
			t.Fatal("observation must write exactly one Event")
		}
		var envelope struct {
			Event struct {
				ID string `json:"id"`
			} `json:"event"`
		}
		if err := json.Unmarshal([]byte(got.stdout), &envelope); err != nil {
			t.Fatal(err)
		}
		event, err := journal.FindEvent(filepath.Join(config.LedgerDir(f.workspace), "journal"), envelope.Event.ID)
		if err != nil {
			t.Fatal(err)
		}
		if event.Event.EventType != "artifact.availability-observed" {
			t.Fatalf("event=%#v", event.Event)
		}
		for _, effect := range event.Event.Effects {
			if effect.Table != "referenced_artifacts" && effect.Table != "idempotency_results" {
				t.Fatalf("kernel side effect: %#v", effect)
			}
		}
		if !row.changed && (len(event.Event.Effects) != 1 || event.Event.Effects[0].Table != "idempotency_results") {
			t.Fatalf("unchanged event=%#v", event.Event)
		}
		replay := f.run(args...)
		if replay.err != nil {
			t.Fatalf("replay=%#v", replay)
		}
		if !decodeCanonicalResult(t, replay.stdout)["replayed"].(bool) || startEventCount(t, f.workspace) != count+1 {
			t.Fatal("retry did not replay without a write")
		}
	}
	human := f.run("artifact", "verify", "--input", input, "--actor", "human:mark", "--idempotency-key", "text")
	if human.err != nil || human.stdout != artifact.ID+"\tavailable\n" || human.stderr != "" {
		t.Fatalf("text=%#v", human)
	}
}
