package journal

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/config"
)

func artifactAvailabilityEvent(t *testing.T) Event {
	t.Helper()
	event := removedFamilyEvent("artifact.availability-observed", EventSchemaV2)
	event.Aggregate.Type, event.Aggregate.ID, event.CommandName = "artifact", "ART-1", "artifact.verify"
	request := testIdempotencyRequest()
	request.OperationScope = "artifact.verify:project-1"
	event.Idempotency = &EventIdempotency{Key: request.IdempotencyKey, Scope: request.OperationScope, RequestDigest: request.RequestDigest}
	event.Effects = []Effect{{Op: "update", Table: "referenced_artifacts", PK: map[string]any{"id": "ART-1"}, Values: map[string]any{"availability": "missing"}, Preconditions: map[string]any{"availability": "available"}, Prior: map[string]any{"availability": "available"}}}
	result := map[string]any{"artifact": map[string]any{"id": "ART-1", "digest": "sha256:" + strings.Repeat("a", 64), "media_type": "text/plain", "size_bytes": 3, "retention_class": "evidence", "availability": "missing", "location_policy": "content-addressed", "created_at": "2026-11-15T00:00:00Z", "expires_at": nil}, "changed": true, "idempotency_key": request.IdempotencyKey, "replayed": false, "EventID": "", "EventOccurredAt": ""}
	raw, err := MarshalCanonicalJSON(result)
	if err != nil {
		t.Fatal(err)
	}
	event.Effects = append(event.Effects, idempotencyResultEffect(request, event, "ART-1", string(raw)))
	return event
}

func TestArtifactAvailabilityStrictLiveAndReplay(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		mutate func(*Event)
		valid  bool
	}{
		{"valid", func(*Event) {}, true},
		{"wrong command", func(e *Event) { e.CommandName = "issue.update" }, false},
		{"service actor", func(e *Event) { e.ActorV2.ActorType = "service" }, false},
		{"wrong aggregate", func(e *Event) { e.Aggregate.ID = "ART-foreign" }, false},
		{"missing result", func(e *Event) { e.Effects = e.Effects[:1] }, false},
		{"duplicate update", func(e *Event) { e.Effects = append([]Effect{e.Effects[0]}, e.Effects...) }, false},
		{"duplicate result", func(e *Event) { e.Effects = append(e.Effects, e.Effects[1]) }, false},
		{"reordered effects", func(e *Event) { e.Effects[0], e.Effects[1] = e.Effects[1], e.Effects[0] }, false},
		{"missing prior", func(e *Event) { e.Effects[0].Prior = nil }, false},
		{"stale prior", func(e *Event) { e.Effects[0].Prior["availability"] = "missing" }, false},
		{"unknown availability", func(e *Event) { e.Effects[0].Values["availability"] = "unknown" }, false},
		{"immutable metadata", func(e *Event) { e.Effects[0].Values["digest"] = "sha256:" + strings.Repeat("b", 64) }, false},
		{"no change update", func(e *Event) { e.Effects[0].Values["availability"] = "available" }, false},
		{"forged result metadata", func(e *Event) {
			e.Effects[1].Values["result_json"] = strings.ReplaceAll(e.Effects[1].Values["result_json"].(string), "text/plain", "text/forged")
		}, false},
		{"forged changed", func(e *Event) {
			e.Effects[1].Values["result_json"] = strings.ReplaceAll(e.Effects[1].Values["result_json"].(string), `"changed":true`, `"changed":false`)
		}, false},
		{"duplicate nested key", func(e *Event) {
			e.Effects[1].Values["result_json"] = strings.ReplaceAll(e.Effects[1].Values["result_json"].(string), `"changed":true`, `"changed":false,"changed":true`)
		}, false},
		{"wrong field digest", func(e *Event) { e.Effects[1].Values["request_field_digests_json"] = `{"--actor":"plain"}` }, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			for _, mode := range []string{"project", "replay"} {
				event := artifactAvailabilityEvent(t)
				tc.mutate(&event)
				workspace := newReplayWorkspace(t)
				conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
				defer closeProjectDB(t, conn)
				_, err := conn.Exec(`INSERT INTO referenced_artifacts (id,digest,media_type,size_bytes,retention_class,availability,location_policy,created_at,expires_at) VALUES ('ART-1',?,'text/plain',3,'evidence','available','content-addressed','2026-11-15T00:00:00Z',NULL)`, "sha256:"+strings.Repeat("a", 64))
				if err != nil {
					t.Fatal(err)
				}
				before := currentProjectionSnapshot(t, conn)
				if mode == "project" {
					err = ProjectEvent(conn, ProjectionInput{Event: event, EventHash: strings.Repeat("a", 64)})
				} else {
					dir := filepath.Join(config.LedgerDir(workspace), "journal")
					writeReplayEvents(t, filepath.Join(dir, "2026-11-15.jsonl"), event)
					_, err = ReplayJournal(conn, dir, nil)
				}
				if tc.valid {
					if err != nil {
						t.Fatalf("%s control: %v", mode, err)
					}
					continue
				}
				if err == nil {
					t.Fatalf("%s accepted malformed observation at %s", mode, event.OccurredAt.Format(time.RFC3339))
				}
				if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
					t.Fatalf("%s rejection changed projection", mode)
				}
			}
		})
	}
}

func TestArtifactAvailabilityRejectsLegacyWire(t *testing.T) {
	t.Parallel()
	workspace := newReplayWorkspace(t)
	conn := openReplayDB(t, filepath.Join(config.LedgerDir(workspace), "fid.sqlite"))
	defer closeProjectDB(t, conn)
	_, err := conn.Exec(`INSERT INTO referenced_artifacts (id,digest,media_type,size_bytes,retention_class,availability,location_policy,created_at,expires_at) VALUES ('ART-1',?,'text/plain',3,'evidence','available','content-addressed','2026-11-15T00:00:00Z',NULL)`, "sha256:"+strings.Repeat("a", 64))
	if err != nil {
		t.Fatal(err)
	}
	event := artifactAvailabilityEvent(t)
	event.EventSchema = EventSchemaV1
	before := currentProjectionSnapshot(t, conn)
	if err := ProjectEvent(conn, ProjectionInput{Event: event, EventHash: strings.Repeat("a", 64)}); err == nil {
		t.Fatal("version-1 wire bypassed strict Artifact validation")
	}
	dir := filepath.Join(config.LedgerDir(workspace), "journal")
	writeReplayEvents(t, filepath.Join(dir, "2026-11-15.jsonl"), event)
	if _, err := ReplayJournal(conn, dir, nil); err == nil {
		t.Fatal("version-1 wire replay bypassed strict Artifact validation")
	}
	if !reflect.DeepEqual(before, currentProjectionSnapshot(t, conn)) {
		t.Fatal("legacy-wire refusal changed projection")
	}
}
