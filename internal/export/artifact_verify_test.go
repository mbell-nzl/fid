package export

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mbell-nzl/fid/internal/artifactverify"
	"github.com/mbell-nzl/fid/internal/checkread"
	"github.com/mbell-nzl/fid/internal/config"
	fiddb "github.com/mbell-nzl/fid/internal/db"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestLeanArtifactVerificationPreservesIdentityAndCheckResult(t *testing.T) {
	t.Parallel()
	var task string
	fixture := newExportSourceFixture(t, func(t *testing.T, root string) string { task = addLeanExportFacts(t, root); return task })
	conn, err := fiddb.Open(filepath.Join(config.LedgerDir(fixture.workspace), config.DatabaseFileName))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
	records, err := checkread.Read(conn, task)
	if err != nil || len(records) != 1 || records[0].Result == nil {
		t.Fatalf("checks=%#v/%v", records, err)
	}
	before, err := json.Marshal(records)
	if err != nil {
		t.Fatal(err)
	}
	file := records[0].Result.Artifacts[0]
	options := artifactverify.Options{WorkspaceRoot: fixture.workspace, MigrationsDir: testMigrationsDir(t), ActorID: "agent:builder", IdempotencyKey: "available", Request: artifactverify.Request{ArtifactID: file.ID}}
	first, err := artifactverify.Verify(options)
	if err != nil {
		t.Fatal(err)
	}
	original := first.Artifact
	raw, err := json.Marshal(original)
	var fields map[string]any
	if err != nil || json.Unmarshal(raw, &fields) != nil || len(fields) != 9 || original.Digest != file.Digest || original.SizeBytes != file.SizeBytes || original.Availability != "available" || first.Changed || first.EventID == "" || first.EventOccurredAt == "" {
		t.Fatalf("available=%#v fields=%s/%v", first, raw, err)
	}
	path := filepath.Join(config.LedgerDir(fixture.workspace), "artifacts", strings.TrimPrefix(file.Digest, "sha256:"))
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, step := range []struct {
		name, actor, availability string
		body                      []byte
		changed                   bool
	}{
		{"missing", "human:mark", "missing", nil, true},
		{"corrupt", "agent:builder", "missing", []byte("wrong bytes"), false},
		{"restored", "agent:builder", "available", body, true},
	} {
		t.Run(step.name, func(t *testing.T) {
			if step.body == nil {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
			} else if err := os.WriteFile(path, step.body, 0o600); err != nil {
				t.Fatal(err)
			}
			options.ActorID, options.IdempotencyKey = step.actor, step.name
			result, err := artifactverify.Verify(options)
			if err != nil || result.Changed != step.changed || result.Artifact.Availability != step.availability || result.Replayed {
				t.Fatalf("verify=%#v/%v", result, err)
			}
			identity := result.Artifact
			identity.Availability = original.Availability
			if !reflect.DeepEqual(identity, original) {
				t.Fatalf("identity changed: %#v -> %#v", original, result.Artifact)
			}
			snapshot := workspaceTreeDigest(t, fixture.workspace)
			replay, err := artifactverify.Verify(options)
			if err != nil || !replay.Replayed || replay.EventID != result.EventID || replay.EventOccurredAt != result.EventOccurredAt || !reflect.DeepEqual(replay.Artifact, result.Artifact) || replay.Changed != result.Changed || workspaceTreeDigest(t, fixture.workspace) != snapshot {
				t.Fatalf("replay=%#v/%v", replay, err)
			}
			event, err := journal.FindEvent(filepath.Join(config.LedgerDir(fixture.workspace), "journal"), result.EventID)
			if err != nil {
				t.Fatal(err)
			}
			wantEffects := 1
			if step.changed {
				wantEffects++
			}
			if event.Event.EventType != "artifact.availability-observed" || len(event.Event.Effects) != wantEffects {
				t.Fatalf("availability event=%#v", event.Event)
			}
		})
	}
	// A retry replays the historical observation even when bytes have returned.
	options.ActorID, options.IdempotencyKey = "human:mark", "missing"
	stale, err := artifactverify.Verify(options)
	if err != nil || stale.Artifact.Availability != "missing" || !stale.Replayed {
		t.Fatalf("historical replay=%#v/%v", stale, err)
	}
	options.Request.ArtifactID = "ART-other"
	snapshot := workspaceTreeDigest(t, fixture.workspace)
	_, err = artifactverify.Verify(options)
	var conflict *journal.IdempotencyConflictError
	if !errors.As(err, &conflict) || workspaceTreeDigest(t, fixture.workspace) != snapshot {
		t.Fatalf("conflict=%v", err)
	}
	after, err := checkread.Read(conn, task)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(after)
	if err != nil || string(encoded) != string(before) {
		t.Fatalf("immutable check changed: %s -> %s / %v", before, encoded, err)
	}
	count := 0
	for _, line := range leanJournalLines(t, fixture.workspace) {
		if strings.Contains(string(line), `"event_type":"artifact.availability-observed"`) {
			count++
		}
	}
	if count != 4 {
		t.Fatalf("availability Events=%d, want 4", count)
	}
}

func TestLeanArtifactVerificationRefusesWithoutMutation(t *testing.T) {
	t.Parallel()
	fixture := newExportSourceFixture(t, addLeanExportFacts)
	for _, row := range []struct {
		name, actor, id, key, message string
		denied                        bool
	}{
		{"untyped actor", "builder", "ART-absent", "bad-actor", "invalid actor", false},
		{"service actor", "service:bot", "ART-absent", "service", "can't verify Artifacts", true},
		{"missing Artifact", "human:mark", "ART-absent", "missing", `Artifact "ART-absent" was not found`, false},
		{"missing key", "agent:builder", "ART-absent", "", "idempotency key must contain", false},
	} {
		t.Run(row.name, func(t *testing.T) {
			before := workspaceTreeDigest(t, fixture.workspace)
			_, err := artifactverify.Verify(artifactverify.Options{WorkspaceRoot: fixture.workspace, MigrationsDir: testMigrationsDir(t), ActorID: row.actor, Request: artifactverify.Request{ArtifactID: row.id}, IdempotencyKey: row.key})
			var denied *artifactverify.DeniedActorError
			if err == nil || !strings.Contains(err.Error(), row.message) || errors.As(err, &denied) != row.denied || workspaceTreeDigest(t, fixture.workspace) != before {
				t.Fatalf("verify refusal=%v", err)
			}
		})
	}
}
