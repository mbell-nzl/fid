# Phase 24 - Task 1: Idempotency field-digest storage - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Migration 0019 adds the accepted non-null JSON column and down restores every schema-18 `sqlite_master` object byte for byte. The active migration and export manifest use projection schema 19; projection meta format remains 5.
- Every central idempotency insert and the direct `gate.execution-started` insert carry a canonical field-digest map. Current call sites write `{}`. The `gate.receipt-recorded` Effect guards the stored map beside `request_digest`; its update values do not include the map. Replay validates an object of digest values and fails closed on malformed or missing maps. Conflict errors retain stored and supplied maps for task 2, without changing public failure details.
- Retained schema, journal, portable recovery, Stop, and archive fixtures were converged to the accepted new storage shape. The old Phase 23 journal fixture retains its original bytes and digest and now asserts the required fail-closed refusal.

## Files created or modified

- `migrations/0019_p24_idempotency_fields.up.sql`; `migrations/0019_p24_idempotency_fields.down.sql`
- `internal/db/projection_meta.go`; `internal/db/migration0019_test.go`; `internal/db/escalation_removal_test.go`
- `internal/journal/idempotency.go`; `internal/journal/idempotency_test.go`; `internal/journal/current_event_idempotency.go`; `internal/journal/current_event_columns.go`; `internal/journal/current_event_json.go`; `internal/journal/current_event_updates.go`; `internal/journal/current_event_validation_test.go`; `internal/journal/checkpoint_projection_test.go`; `internal/journal/task9_prechange_replay_test.go`
- `internal/evidence/execution.go`
- `internal/export/archive_semantics.go`; `internal/export/archive_semantics_test.go`; `internal/export/export_test.go`
- `internal/maintenance/check_test.go`; `internal/workflow/stop_test.go`
- `cmd/fid/integration/p22_portable_schema_test.go`; `cmd/fid/integration/p23_portable_recovery_test.go`
- `project-docs/specs/03_TRD_SCHEMA.md`; `project-docs/docs/phases/phase-24/tasks/task-1-done.md`

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md` section 8.5 records schema 19, exact DDL, immutable receipt guard, down migration, and unsupported older journals.

## Validation Contract Coverage

- VAL-001: PASS - `TestMigration0019AddsOnlyFieldDigestColumn` in `internal/db/migration0019_test.go:10` compares every `sqlite_master` object and the full schema-19 table SQL against the accepted single-column addition.
- VAL-058: PASS - the same test migrates down and compares every schema-18 `sqlite_master` object byte for byte.
- VAL-002 (storage half): PASS - `TestPortableRoundTripCoversEveryPhase23Table` and `p23PortableCollectFieldDigestFamilies` in `cmd/fid/integration/p23_portable_recovery_test.go:31` include keyed issue, claim, run, and review rows with `{}` and compare all live projection rows to rebuild and import; `TestCommitIdempotentReplaysStoredResultAfterRebuild` in `internal/journal/idempotency_test.go:14` compares the stored column byte for byte before and after rebuild.
- VAL-003: PASS - `TestCurrentIdempotencyFieldDigestsReplayFailClosed` in `internal/journal/current_event_validation_test.go:14` replays the valid object and rejects missing, array, and non-digest maps at line 1 without projecting a row. `TestTask9PrechangeBudgetObservedJournalFailsClosedWithoutMutation` preserves the old fixture bytes and rejects its missing map at line 2. The built `fid --json rebuild` also returned `journal-replay-stopped`, line 2, exit 6 on that fixture. `internal/journal/current_event_updates.go:46` requires the receipt guard; `internal/journal/project_effects.go:158` rejects any guard value that differs from the stored row before update.
- VAL-004: PASS - `TestSchema19ArchiveImportsAndSchema18ManifestFailsClosed` in `internal/export/export_test.go:75` creates, validates, and imports schema 19; a copy declaring 18 fails `schema-unsupported` with unchanged target tree. `internal/export/archive_validation.go:130` binds this check to `manifest.json`.
- VAL-005: PASS - the required 100k benchmark completed its parity self-assert and reported `22373 ms/rebuild`, below 30000.
- VAL-052: PASS - the final ordinary gate reported `cmd/fid` 215.296 s and `internal/workflow` 72.492 s. The first full attempt reported the slowest other ordinary package, `internal/issues`, at 194.345 s; that package was cached on the final attempt after only a workflow fixture changed. Tagged `cmd/fid/integration` took 253.353 s. No package reached the 480 s headroom warning or its budget.

## Tests run

- First red: `go test -run '^TestMigration0019AddsOnlyFieldDigestColumn$' -count=1 ./internal/db` - FAIL - `FAIL	github.com/mbell-nzl/fid/internal/db	0.078s` (terminated `exit_code: 1`); migration 19 was absent.
- Green: `go test -run '^TestMigration0019AddsOnlyFieldDigestColumn$' -count=1 ./internal/db` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	0.088s` (terminated `exit_code: 0`)
- `go test -run '^TestCurrentIdempotencyFieldDigestsReplayFailClosed$' -count=1 ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	0.303s` (terminated `exit_code: 0`)
- `go test -run '^TestSchema19ArchiveImportsAndSchema18ManifestFailsClosed$' -count=1 ./internal/export` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	0.828s` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestPortableRoundTripCoversEveryPhase23Table$' ./cmd/fid/integration` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	14.165s` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/workflow` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	58.494s` (terminated `exit_code: 0`)
- Built entrypoint: `go build -o /tmp/fid-p24-task1 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`). A keyed `fid --json issue create` then returned `"replayed":false` (terminated `exit_code: 0`); `sqlite3` read `request_field_digests_json` as `{}`. `/tmp/fid-p24-task1 --project /tmp/fid-p24-task1-smoke-nKX1AY --json rebuild` reported `"events_replayed":2` and `"schema_version":19` (terminated `exit_code: 0`); the same `sqlite3` query still returned `{}`.
- Old-journal entrypoint: `/tmp/fid-p24-task1 --project /tmp/fid-p24-task1-old-oZxm2g --json rebuild` observed the intended fail-closed refusal: `"code":"journal-replay-stopped"` with `"line":2` (terminated `exit_code: 6`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- First attempt: `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - FAIL - `FAIL	github.com/mbell-nzl/fid/internal/workflow	103.620s` (terminated `exit_code: 1`).
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	72.492s	coverage: 62.7% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:											(statements)					60.4%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	253.353s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	25.020s` (terminated `exit_code: 0`)
  - Benchmark line: `22373 ms/rebuild`.
- `rg -n "idempotency_results" migrations internal -g '!*_test.go'` - PASS - `internal/workflow/activation.go:845:	if err := conn.QueryRow(`SELECT COUNT(*) FROM idempotency_results WHERE result_entity_id = ?`, activationID).Scan(&count); err != nil {` (terminated `exit_code: 0`)

## Deviations from spec

- No accepted-contract deviation. The first required full gate failed with `FAIL	github.com/mbell-nzl/fid/internal/workflow	103.620s` (terminated `exit_code: 1`): the shared retained Stop fixture omitted the new required map, so five workflow tests stopped at the insert shape. Adding `{}` in `internal/workflow/stop_test.go:256` preserved those tests' original assertions. The byte-identical required command then passed with `ok  	github.com/mbell-nzl/fid/internal/workflow	72.492s	coverage: 62.7% of statements` (terminated `exit_code: 0`).

## Concerns, technical debt, or notes for the Tech Lead

- `internal/export/archive_updates.go` remained unchanged: the field map is never a receipt update value. The schema-19 semantic validator admission in `internal/export/archive_semantics.go` was required for archive create/validate/import acceptance.
- Retained-test provenance was checked from introducing commits, not filenames: `checkpoint_projection_test.go` 59c62b9 (phase-22/task-4); `task9_prechange_replay_test.go` 828eb18 (phase-23/task-9); `archive_semantics_test.go` acb569b (phase-22/task-13); `check_test.go` d2a087d (phase-17/task-3); `p22_portable_schema_test.go` 5ec69e6 (phase-22/task-15); `p23_portable_recovery_test.go` d02614e (phase-23/task-16); `escalation_removal_test.go` 644f30e (phase-22/task-16); `stop_test.go` ac82795 (phase-22/task-2). Each retains its original subject; the accepted schema-19 change displaced only its fixture or version expectation. The historical migration-0013 test in `internal/db/operational_kernel_migration_test.go` was left unchanged because it runs at schema 13.
- An existing Phase 23 workspace with keyed rows cannot catch up in place: SQLite refuses the accepted non-null `ADD COLUMN` on its populated table. Rebuilding its old keyed journal also stops closed at the first missing map; the retained old-journal fixture and CLI run returned `journal-replay-stopped`, line 2, exit 6. No compatibility path was added under DEC-270.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `58352c73edf81ac8e89ffb7d7d824a67def3e3c4`

## Commit subject

- `feat(idempotency): store request field digests in schema 19`
