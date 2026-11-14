package export

import (
	"database/sql"
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
	retained := retainedArtifactMetadata(t, conn, file.ID)
	options := artifactverify.Options{WorkspaceRoot: fixture.workspace, MigrationsDir: testMigrationsDir(t), ActorID: "agent:builder", IdempotencyKey: "available", Request: artifactverify.Request{ArtifactID: file.ID}}
	first, err := artifactverify.Verify(options)
	if err != nil {
		t.Fatal(err)
	}
	assertArtifactMetadata(t, first.Artifact, retained, "available")
	if first.Artifact.Digest != file.Digest || first.Artifact.SizeBytes != file.SizeBytes || first.Changed || first.Replayed || first.EventID == "" || first.EventOccurredAt == "" {
		t.Fatalf("available=%#v", first)
	}
	snapshot := workspaceTreeDigest(t, fixture.workspace)
	initialReplay, err := artifactverify.Verify(options)
	if err != nil || !initialReplay.Replayed || initialReplay.EventID != first.EventID || initialReplay.EventOccurredAt != first.EventOccurredAt || initialReplay.Changed || workspaceTreeDigest(t, fixture.workspace) != snapshot {
		t.Fatalf("initial replay=%#v/%v", initialReplay, err)
	}
	assertArtifactMetadata(t, initialReplay.Artifact, retained, "available")
	path := filepath.Join(config.LedgerDir(fixture.workspace), "artifacts", strings.TrimPrefix(file.Digest, "sha256:"))
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var missing artifactverify.Result
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
			assertArtifactMetadata(t, result.Artifact, retained, step.availability)
			if step.name == "missing" {
				missing = result
			}
			snapshot := workspaceTreeDigest(t, fixture.workspace)
			replay, err := artifactverify.Verify(options)
			if err != nil || !replay.Replayed || replay.EventID != result.EventID || replay.EventOccurredAt != result.EventOccurredAt || replay.Changed != step.changed || workspaceTreeDigest(t, fixture.workspace) != snapshot {
				t.Fatalf("replay=%#v/%v", replay, err)
			}
			assertArtifactMetadata(t, replay.Artifact, retained, step.availability)
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
	snapshot = workspaceTreeDigest(t, fixture.workspace)
	stale, err := artifactverify.Verify(options)
	if err != nil || !stale.Replayed || !stale.Changed || stale.EventID != missing.EventID || stale.EventOccurredAt != missing.EventOccurredAt || workspaceTreeDigest(t, fixture.workspace) != snapshot {
		t.Fatalf("historical replay=%#v/%v", stale, err)
	}
	assertArtifactMetadata(t, stale.Artifact, retained, "missing")
	options.Request.ArtifactID = "ART-other"
	snapshot = workspaceTreeDigest(t, fixture.workspace)
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

// Read the retained row before verification, without using the verifier's
// loader or Artifact type to define the expected JSON keys or values.
func retainedArtifactMetadata(t *testing.T, conn *sql.DB, id string) map[string]any {
	t.Helper()
	var raw string
	if err := conn.QueryRow(`SELECT json_object(
		'id', id, 'digest', digest, 'media_type', media_type,
		'size_bytes', size_bytes, 'retention_class', retention_class,
		'availability', availability, 'location_policy', location_policy,
		'created_at', created_at, 'expires_at', expires_at
	) FROM referenced_artifacts WHERE id = ?`, id).Scan(&raw); err != nil {
		t.Fatalf("read retained Artifact metadata: %v", err)
	}
	var fields map[string]any
	if err := json.Unmarshal([]byte(raw), &fields); err != nil {
		t.Fatal(err)
	}
	return fields
}

func assertArtifactMetadata(t *testing.T, artifact artifactverify.Artifact, retained map[string]any, availability string) {
	t.Helper()
	want := make(map[string]any, len(retained))
	for key, value := range retained {
		want[key] = value
	}
	want["availability"] = availability
	raw, err := json.Marshal(artifact)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(raw, &fields); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(fields, want) {
		t.Fatalf("Artifact JSON = %s, want all nine retained fields: %#v", raw, want)
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
