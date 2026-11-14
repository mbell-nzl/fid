package export

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/gittest"
	"github.com/mbell-nzl/fid/internal/issues"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestCreateReplayAndArchiveValidateRetryLane(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "epic"} {
		for _, row := range []struct {
			name  string
			value any
			valid bool
		}{
			{"unset", nil, true}, {"missing", nil, true},
			{"routine", "routine", true}, {"standard", "standard", true}, {"assured", "assured", true},
			{"empty", "", false}, {"unknown", "urgent", false}, {"number", 42, false},
			{"boolean", true, false}, {"array", []any{}, false}, {"object", map[string]any{}, false},
		} {
			for _, mode := range []string{"replay", "archive"} {
				t.Run(kind+"/"+row.name+"/"+mode, func(t *testing.T) {
					t.Parallel()
					validator := startArchiveValidator(t)
					valid := row.valid && (kind == "task" || row.value == nil)
					event := archiveCreateResultEvent(t, kind, row.value, valid, row.name == "missing")
					err := acceptRetryEvent(t, validator, event, mode)
					if (err == nil) != valid || !valid && !strings.Contains(err.Error(), "Issue create result") {
						t.Fatalf("%s creation retry = %v, valid=%v", mode, err, valid)
					}
					var issues, metadata, activities, bindings int
					if err := validator.conn.QueryRow(`SELECT (SELECT count(*) FROM issues WHERE id='FID-new'), (SELECT count(*) FROM epic_metadata WHERE issue_id='FID-new'), (SELECT count(*) FROM activity_log), (SELECT count(*) FROM idempotency_results)`).Scan(&issues, &metadata, &activities, &bindings); err != nil {
						t.Fatal(err)
					}
					if valid {
						wantMetadata := 0
						if kind == "epic" {
							wantMetadata = 1
						}
						if issues != 1 || metadata != wantMetadata || activities != 1 || bindings != 1 {
							t.Fatalf("valid creation = %d/%d/%d/%d", issues, metadata, activities, bindings)
						}
					} else if issues != 0 || metadata != 0 || activities != 0 || bindings != 0 {
						t.Fatalf("rejected creation left partial state = %d/%d/%d/%d", issues, metadata, activities, bindings)
					}
				})
			}
		}
	}
}

func TestArchiveValidateRejectsMalformedCreationLaneWithValidIntegrity(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "epic"} {
		t.Run(kind, func(t *testing.T) {
			t.Parallel()
			workspace := gittest.NewWorkspace(t)
			if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Prefix: "FID"}); err != nil {
				t.Fatal(err)
			}
			scheduling := ""
			if kind == "epic" {
				scheduling = "prioritised"
			}
			if _, err := issues.Create(issues.CreateOptions{WorkspaceRoot: workspace, MigrationsDir: testMigrationsDir(t), Kind: kind, Title: "Planned work", Scheduling: scheduling, Actor: "agent:builder", IdempotencyKey: "create"}); err != nil {
				t.Fatal(err)
			}
			archive := filepath.Join(t.TempDir(), "create.zip")
			if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archive}); err != nil {
				t.Fatal(err)
			}
			if _, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)}); err != nil {
				t.Fatalf("valid control: %v", err)
			}
			contents := archiveContents(t, archive)
			path := archivePathWithPrefix(t, contents, "journal/")
			data, hash := mutateArchiveEvent(t, contents[path], "issue.created", func(event *journal.Event) {
				for i := range event.Effects {
					if event.Effects[i].Table == "idempotency_results" {
						event.Effects[i].Values["result_json"] = strings.Replace(event.Effects[i].Values["result_json"].(string), `"lane":null`, `"lane":""`, 1)
					}
				}
			})
			if bytes.Equal(data, contents[path]) {
				t.Fatal("creation lane was not mutated")
			}
			contents[path] = data
			manifest := decodeManifestMap(t, contents["manifest.json"])
			manifest["terminal_event_hash"] = hash
			updateManifestEntry(t, manifest, path, data)
			contents["manifest.json"] = encodeManifestMap(t, manifest)
			writeTestArchive(t, archive, contents)
			before := workspaceTreeDigest(t, workspace)
			_, err := Validate(ValidateOptions{ArchivePath: archive, MigrationsDir: testMigrationsDir(t)})
			assertArchiveCheck(t, err, "reference-invalid")
			invalid := err.(*InvalidArchiveError)
			if len(invalid.Checks) != 1 || !strings.Contains(invalid.Checks[0].Detail, "Issue create result") {
				t.Fatalf("integrity-valid malformed creation lane checks = %#v", invalid.Checks)
			}
			if after := workspaceTreeDigest(t, workspace); after != before {
				t.Fatal("archive validation changed workspace")
			}
		})
	}
}

func archiveCreateResultEvent(t *testing.T, kind string, lane any, valid, missing bool) journal.Event {
	t.Helper()
	now := time.Date(2026, 11, 10, 9, 5, 0, 0, time.UTC)
	values := map[string]any{"kind": kind, "title": "Planned work", "workflow_id": kind + "-default", "status": "todo", "revision": 1, "owner_role": "developer", "created_at": now.Format(time.RFC3339), "updated_at": now.Format(time.RFC3339)}
	if valid && lane != nil {
		values["lane"] = lane
	}
	event := journal.Event{EventSchema: journal.EventSchemaV2, EventID: "EVT-20261110-000003", EventType: "issue.created", CommandName: "issue.create", OccurredAt: now, ActorV2: journal.EventActor{ActorID: "builder", ActorType: "agent", Producer: "fid-standalone", ProvenanceRef: "project-local-attribution"}, Aggregate: journal.Aggregate{ID: "FID-new", Type: "issue", Sequence: 1}, Effects: []journal.Effect{{Op: "insert", Table: "issues", PK: map[string]any{"id": "FID-new"}, Values: values}}}
	if kind == "epic" {
		event.Effects = append(event.Effects, journal.Effect{Op: "insert", Table: "epic_metadata", PK: map[string]any{"issue_id": "FID-new"}, Values: map[string]any{"scheduling": "prioritised"}})
	}
	event.Effects = append(event.Effects, journal.Effect{Op: "activity", Table: "activity_log", PK: map[string]any{"id": "creation"}, Values: map[string]any{"journal_event_id": event.EventID, "ref_type": "issue", "ref_id": "FID-new", "event": "issue-created"}})
	takeoverArchiveRetry(&event)
	event.Idempotency.Scope = "issue.create:project-test"
	last := &event.Effects[len(event.Effects)-1]
	last.PK["operation_scope"] = event.Idempotency.Scope
	issue := map[string]any{"id": "FID-new", "kind": kind, "title": "Planned work", "status": "todo", "revision": 1, "lane": lane}
	if missing {
		delete(issue, "lane")
	}
	raw, err := journal.MarshalCanonicalJSON(map[string]any{"issue": issue, "event_occurred_at": now.Format(time.RFC3339), "idempotency_key": event.Idempotency.Key, "replayed": false})
	if err != nil {
		t.Fatal(err)
	}
	last.Values["result_json"] = string(raw)
	return event
}

func acceptRetryEvent(t *testing.T, validator *semanticValidator, event journal.Event, mode string) error {
	t.Helper()
	if mode == "replay" {
		return journal.ProjectEvent(validator.conn, journal.ProjectionInput{Event: event, EventHash: "sha256:retry-result", File: "journal", Line: 2, Offset: 2})
	}
	raw, err := journal.MarshalCanonicalJSON(event)
	if err != nil {
		t.Fatal(err)
	}
	return validator.accept(raw)
}
