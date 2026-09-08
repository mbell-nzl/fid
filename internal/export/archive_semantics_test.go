package export

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mbell-nzl/fid/internal/bootstrap"
	"github.com/mbell-nzl/fid/internal/journal"
)

func TestArchiveWithValidHashesButInvalidRecoveryReferenceRejects(t *testing.T) {
	workspace := t.TempDir()
	if _, err := bootstrap.InitializeWorkspace(bootstrap.InitOptions{
		WorkspaceRoot: workspace,
		MigrationsDir: testMigrationsDir(t),
		Prefix:        "FID",
	}); err != nil {
		t.Fatalf("initialize workspace: %v", err)
	}
	addRecoveredExportFixtureFacts(t, workspace)
	archivePath := filepath.Join(t.TempDir(), "invalid-recovery-reference.zip")
	if _, err := Create(CreateOptions{WorkspaceRoot: workspace, ArchivePath: archivePath}); err != nil {
		t.Fatalf("create archive: %v", err)
	}
	if _, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)}); err != nil {
		t.Fatalf("valid recovered-history control: %#v", err)
	}

	entries := archiveContents(t, archivePath)
	journalPath := archivePathWithPrefix(t, entries, "journal/")
	var terminalHash string
	entries[journalPath], terminalHash = replaceRecoveredObservationReference(t, entries[journalPath])

	manifest := decodeManifestMap(t, entries["manifest.json"])
	manifest["terminal_event_hash"] = terminalHash
	updateManifestEntry(t, manifest, journalPath, entries[journalPath])
	entries["manifest.json"] = encodeManifestMap(t, manifest)
	writeTestArchive(t, archivePath, entries)

	_, err := Validate(ValidateOptions{ArchivePath: archivePath, MigrationsDir: testMigrationsDir(t)})
	if err == nil {
		t.Fatal("Validate() accepted valid-hash recovery Event with a missing observation reference")
	}
	assertArchiveCheck(t, err, "reference-invalid")
}

func replaceRecoveredObservationReference(t *testing.T, data []byte) ([]byte, string) {
	t.Helper()
	lines := bytes.Split(bytes.TrimSuffix(data, []byte{'\n'}), []byte{'\n'})
	changed := false
	for index, line := range lines {
		var event journal.Event
		if err := event.UnmarshalJSON(line); err != nil {
			t.Fatalf("decode journal Event %d: %v", index+1, err)
		}
		if event.EventType != "run.recovered" {
			continue
		}
		for effectIndex := range event.Effects {
			effect := &event.Effects[effectIndex]
			if effect.Table == "run_recoveries" && effect.Op == "insert" {
				effect.Values["observation_id"] = "OBS-missing"
				changed = true
			}
		}
		encoded, err := journal.MarshalCanonicalJSON(event)
		if err != nil {
			t.Fatalf("encode tampered Recovery Event: %v", err)
		}
		lines[index] = encoded
	}
	if !changed {
		t.Fatal("recovered journal Event was absent")
	}
	var terminal journal.Event
	if err := terminal.UnmarshalJSON(lines[len(lines)-1]); err != nil {
		t.Fatalf("decode terminal Event: %v", err)
	}
	terminalHash, err := journal.HashEvent(terminal)
	if err != nil {
		t.Fatalf("hash terminal Event: %v", err)
	}
	return append(bytes.Join(lines, []byte{'\n'}), '\n'), terminalHash
}

func recoveryEventWithMissingReference() journal.Event {
	event := journal.Event{
		EventID:     "EVT-invalid-recovery-reference",
		EventSchema: journal.EventSchemaV2,
		EventType:   "run.recovered",
		OccurredAt:  time.Date(2026, time.September, 8, 0, 0, 0, 0, time.UTC),
		ActorV2: journal.EventActor{
			ActorID:       "agent:recovery",
			ActorType:     "agent",
			Producer:      "fid",
			ProvenanceRef: "test:archive-semantic-reference",
		},
		Aggregate:   journal.Aggregate{ID: "RUN-prior-missing", Sequence: 1, Type: "run"},
		CommandName: "run.recover",
		Effects: []journal.Effect{
			{Op: "insert", Table: "claims", PK: map[string]any{"id": "CLM-successor-missing"}, Values: map[string]any{
				"issue_id": "ISS-missing", "envelope_id": "ENV-missing", "envelope_revision": 1, "actor_id": "agent:recovery", "worktree_id": "WT-missing",
				"scope_json": "{}", "mode": "mutation", "version": 1, "status": "active", "acquired_at": "2026-09-08T00:00:00Z", "expires_at": "2026-09-08T01:00:00Z", "released_at": nil, "release_reason": nil,
			}},
			{Op: "insert", Table: "runs", PK: map[string]any{"id": "RUN-successor-missing"}, Values: map[string]any{
				"issue_id": "ISS-missing", "envelope_id": "ENV-missing", "envelope_revision": 1, "claim_id": "CLM-successor-missing", "claim_version": 1, "actor_id": "agent:recovery",
				"start_repository_json": "{}", "terminal_repository_json": nil, "status": "active", "started_at": "2026-09-08T00:00:00Z", "ended_at": nil,
			}},
			{Op: "insert", Table: "run_runtime", PK: map[string]any{"run_id": "RUN-successor-missing"}, Values: map[string]any{"version": 1, "executor_id": "EXEC-successor", "obligation_id": nil}},
			{Op: "insert", Table: "run_recoveries", PK: map[string]any{"id": "REC-invalid-reference"}, Values: map[string]any{
				"prior_run_id": "RUN-prior-missing", "successor_run_id": "RUN-successor-missing", "prior_claim_id": "CLM-prior-missing", "successor_claim_id": "CLM-successor-missing", "observation_id": "OBS-missing",
				"reconciliation_json": "{}", "adopted_paths_json": "[]", "authority_json": "{}", "reason": "semantic reference test", "created_event_id": "EVT-invalid-recovery-reference",
			}},
		},
	}
	event.Idempotency = &journal.EventIdempotency{Key: "missing-reference", Scope: "run.recover:missing-project", RequestDigest: strings.Repeat("1", 64)}
	event.Effects = append(event.Effects, journal.Effect{Op: "insert", Table: "idempotency_results", PK: map[string]any{"operation_scope": event.Idempotency.Scope, "idempotency_key": event.Idempotency.Key}, Values: map[string]any{"request_digest": event.Idempotency.RequestDigest, "result_entity_id": "REC-invalid-reference", "result_event_id": event.EventID, "aggregate_sequence": 1, "result_json": `{"replayed":false}`, "created_at": event.OccurredAt.Format(time.RFC3339Nano)}})
	return event
}

func TestSemanticValidatorRejectsNestedDuplicateJSONKey(t *testing.T) {
	validator, err := newSemanticValidator(testMigrationsDir(t), 16)
	if err != nil {
		t.Fatalf("new semantic validator: %v", err)
	}
	defer func() { _ = validator.close() }()

	raw := []byte(`{"actor":{"actor_id":"agent:one","actor_id":"agent:two","actor_type":"agent","producer":"fid","provenance_ref":"test:duplicate","session_ref":null},"aggregate":{"id":"PROJ-1","sequence":1,"type":"project"},"causation_id":null,"command":"init","correlation_id":null,"effects":[],"event_id":"EVT-duplicate","event_schema":2,"event_type":"project.initialized","idempotency":null,"occurred_at":"2026-09-08T00:00:00Z","repository":null}`)
	if err := validator.accept(raw); err == nil || !strings.Contains(err.Error(), `duplicate key "actor_id"`) {
		t.Fatalf("accept nested duplicate key error = %v, want duplicate actor key rejection", err)
	}
}

func TestSemanticValidatorRejectsMissingRecoveryReference(t *testing.T) {
	validator, err := newSemanticValidator(testMigrationsDir(t), 16)
	if err != nil {
		t.Fatalf("new semantic validator: %v", err)
	}
	defer func() { _ = validator.close() }()

	raw, err := journal.MarshalCanonicalJSON(recoveryEventWithMissingReference())
	if err != nil {
		t.Fatalf("encode recovery Event: %v", err)
	}
	err = validator.accept(raw)
	if err == nil || !strings.Contains(err.Error(), "missing reference") {
		t.Fatalf("accept missing recovery reference error = %v, want foreign-key rejection", err)
	}
	if got := semanticCheckCode(err); got != "reference-invalid" {
		t.Fatalf("semantic check code = %q, want reference-invalid", got)
	}
}

func TestSemanticValidatorRejectsRecoveryEffectFamilyTampering(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		mutate func(*journal.Event)
	}{
		{
			name: "wrong operation",
			mutate: func(event *journal.Event) {
				event.Effects[0].Op = "update"
			},
		},
		{
			name: "reordered successor run",
			mutate: func(event *journal.Event) {
				event.Effects[0], event.Effects[1] = event.Effects[1], event.Effects[0]
			},
		},
		{
			name: "wrong aggregate",
			mutate: func(event *journal.Event) {
				event.Aggregate.Type = "issue"
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			validator, err := newSemanticValidator(testMigrationsDir(t), 16)
			if err != nil {
				t.Fatalf("new semantic validator: %v", err)
			}
			defer func() { _ = validator.close() }()

			event := recoveryEventWithMissingReference()
			testCase.mutate(&event)
			raw, err := journal.MarshalCanonicalJSON(event)
			if err != nil {
				t.Fatalf("encode tampered recovery Event: %v", err)
			}
			err = validator.accept(raw)
			if err == nil {
				t.Fatal("accept tampered recovery Event error = nil")
			}
			if got := semanticCheckCode(err); got != "effect-invalid" {
				t.Fatalf("semantic check code = %q, want effect-invalid; error = %v", got, err)
			}
		})
	}
}
