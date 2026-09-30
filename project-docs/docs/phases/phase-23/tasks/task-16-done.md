# Phase 23 - Task 16: Portable recovery for the review and judgment domain - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- The archive semantic validator accepts the existing boolean representation of `review_findings.blocking`. This lets `export validate` and `export import` replay real Phase 23 Review Events against projection schema 18.
- A built-binary integration journey exports, validates, rebuilds, and imports four populated Phase 23 workspaces. It compares every application table row and every canonical archive entry between live, rebuilt, and imported state, and verifies the schema-18 projection integrity. The four workspaces collectively populate every table named below and emit all nine Phase 23 Event types.
- The journey verifies non-mutating validation, exact schema-18 import and schema-17 rejection, one `review-cycles` sample keyed to its Review ID, and rejection with no projection change for three corrupt payloads per Phase 23 Event type plus an unknown type.

## Files created or modified

- `internal/export/archive_semantics_values.go`: accept `review_findings.blocking` as a JSON boolean in the archive semantic replay.
- `cmd/fid/integration/p23_portable_recovery_test.go`: portable parity, corruption, schema, and charge journey.
- `project-docs/docs/phases/phase-23/tasks/task-16-done.md`: this report.

## Documentation updated

- This report records the Phase 23 table and Event inventory, validation evidence, exact command results, and provenance. No governing contract or Runner-owned ledger was edited.

## Validation Contract Coverage

- VAL-004: PASS - `TestMigration0018RoundTripIntegrity` checks `foreign_key_check` and `integrity_check` on the final schema during the ordinary gate; `TestPortableRoundTripCoversEveryPhase23Table` calls `p23AssertRebuiltProjectionIntegrity` after each real CLI rebuild and import (`p23_portable_recovery_test.go:233` and `:246`). The helper checks no foreign-key rows and `integrity_check = ok`.
- VAL-027 (import leg): PASS - `TestPortableRoundTripCoversEveryPhase23Table/review-exhaustion` asserts one `review-cycles` row, consumed value 1, and a `source_sample_id` matching `review_records.id` (`p23_portable_recovery_test.go:88`). Its exact `budget_samples` row is compared after rebuild and import (`:219`); the imported state has no duplicate charge. Task 6's recording, replay, and linked-recovery legs remain covered by `TestP23ReviewCycleChargeSurvivesLinkedRecoveryAndRebuild` in the full tagged gate.
- VAL-064: PASS - `TestPortableRoundTripCoversEveryPhase23Table` compares every application table's complete sorted rows (`p23_portable_recovery_test.go:254`, `:288`) and canonical archive entries including ordered Event bytes (`:219`) across live, rebuilt, and imported workspaces. It checks `.fid` file bytes before and after `export validate` (`:300`), and asserts every named table has populated rows (`:22`, `:34`). Archive entry equality also preserves project, entity, and Event identities and history. The direct CLI create, validate, and import below additionally returned the same project ID and terminal Event hash.
- VAL-065: PASS - the current archive imports and reports `projection_schema:18` (`p23_portable_recovery_test.go:219`); `p23PortableAssertPriorSchemaRejected` changes only the manifest projection schema to 17, then checks `archive-invalid` with `schema-unsupported` and absence of a published `.fid` (`:173`). The current manifest remains `fid.export-manifest.v2`.
- VAL-076: PASS - `p23PortableAssertCorruptEventsRejected` (`p23_portable_recovery_test.go:104`) uses real exported Events for all nine named types. Each is tried with an unknown Effect table, a missing required Effect, and a malformed scalar body; the unknown type is also tried. Each rejection comes from the current Event validator, and a complete projection snapshot is unchanged after each attempt.
- VAL-079: PASS - ordinary `cmd/fid` took `246.504s` and `internal/workflow` took `100.101s` (each under 600s); tagged `cmd/fid/integration` took `269.692s` (under 900s). All three are below 80% of budget; neither gate timed out.

The table provenance is the accepted Phase 23 task that introduced or changed the projected row, as recorded in `contract-coverage.md`. Each row below is executed by `TestPortableRoundTripCoversEveryPhase23Table` under `go test -tags integration -run '^TestPortableRoundTripCoversEveryPhase23Table$' -count=1 -v ./cmd/fid/integration/` (PASS, `exit_code: 0`) and the mandatory tagged race gate. The named subcase supplies non-empty rows; `p23PortableRoundTrip` compares the entire row set across all three projections, and compares the complete archive.

| Table | Row provenance | Executing subcase |
|---|---|---|
| `review_records` | Phase 23 task 4, Review record | `mixed-findings` |
| `review_findings` | Phase 23 task 4, Review Findings | `mixed-findings` |
| `finding_dispositions` | Phase 23 task 5, Finding disposition | `mixed-findings` |
| `concerns` | Phase 23 task 8 generic Concern; task 15 rename from `judgment_concerns` | `references-and-judgment` |
| `deviations` | Phase 23 task 8, Deviation | `references-and-judgment` |
| `decision_refs` | Phase 23 task 9, Decision Reference | `protected-decision` |
| `human_outcomes` | Phase 23 task 10, Human Outcome | `protected-decision` |
| `envelope_protected_decisions` | Phase 23 tasks 2 and 9, protected Decision binding | `protected-decision` |
| `requirement_refs` | Phase 23 task 12, Requirement Reference | `references-and-judgment` |
| `requirement_issue_links` | Phase 23 task 12, traceability link | `references-and-judgment` |
| `envelope_obligations` | Phase 23 task 2, Review and judgment Obligation generalization | `mixed-findings` |
| `evidence_records` | Phase 23 tasks 1, 6, and 10, three Evidence subject types | `protected-decision` |
| `budget_samples` | Phase 23 task 6, Review-cycle charge | `review-exhaustion` |
| `activity_log` | Phase 23 tasks 8, 13, and 15, Review and judgment Activity Effects | `references-and-judgment` |

The same executing command checks these nine exported Event types in `p23PortableJournalEvents` (`p23_portable_recovery_test.go:324`) and the three corruption subcases in `p23PortableAssertCorruptEventsRejected` (`:104`): `review.recorded` (task 4), `finding.addressed` (task 5), `concern.recorded` and `concern.disposed` (task 8, rename converged by task 15), `deviation.recorded` (task 8), `decision.referenced` (task 9), `human-outcome.recorded` (task 10), `requirement.referenced` and `requirement.linked` (task 12). An `unknown.phase23` type is rejected by that command too. No retained test row was changed or converged; this task adds new rows and preserves the existing assertions.

## Tests run

- Focused red: `go test -tags integration -run '^TestPortableRoundTripCoversEveryPhase23Table$' -count=1 -v ./cmd/fid/integration/` - FAIL - `FAIL` (terminated `exit_code: 1`). The diagnostic was `review_findings.blocking requires an integer, not a boolean` in `export validate`; this is the first expected TDD failure.
- Focused test-fixture attempt: the same command - FAIL - `FAIL` (terminated `exit_code: 1`). The new `references-and-judgment` fixture had not created its Git directory. The fixture now creates it.
- Focused corruption-test attempt: the same command - FAIL - `FAIL` (terminated `exit_code: 1`). A zero-Effect test Event could not be hashed; the test now removes only its required Effect, leaving a serializable malformed Event.
- Focused final: `go test -tags integration -run '^TestPortableRoundTripCoversEveryPhase23Table$' -count=1 -v ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	9.125s` (terminated `exit_code: 0`)
- Pre-gate `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 ./internal/export/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/export	74.008s` (terminated `exit_code: 0`)
- Pre-gate `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -count=1 -run '^TestPortableRoundTripCoversEveryPhase23Table$' ./cmd/fid/integration/` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	13.085s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid-task16 ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- Direct CLI `/tmp/fid-task16 --project /tmp/fid-task16-cli-KjxBpO/source --json export create /tmp/fid-task16-cli-KjxBpO/archive.zip` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.create","result":{"export":{"archive_path":"/tmp/fid-task16-cli-KjxBpO/archive.zip","archive_digest":"sha256:bfc37ca738c91111eaae92d7b200ac2622e50584a2e3d9f238a14721304bb0b0","manifest_schema":"fid.export-manifest.v2","project_id":"project-0b7bf012d95afc763c9b5a706b4c8a47","event_count":1,"terminal_event_hash":"sha256:826a110a565da6901e0f173e7d85164f044c60d2c6f0ba0a63d944116627a8bd","projection_schema":18,"entry_count":2}}}` (terminated `exit_code: 0`)
- Direct CLI `/tmp/fid-task16 --project /tmp/fid-task16-cli-KjxBpO/source --json export validate /tmp/fid-task16-cli-KjxBpO/archive.zip` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.validate","result":{"validation":{"archive_path":"/tmp/fid-task16-cli-KjxBpO/archive.zip","archive_digest":"sha256:bfc37ca738c91111eaae92d7b200ac2622e50584a2e3d9f238a14721304bb0b0","status":"valid","manifest_schema":"fid.export-manifest.v2","entry_count":2}}}` (terminated `exit_code: 0`)
- Direct CLI `/tmp/fid-task16 --json export import /tmp/fid-task16-cli-KjxBpO/archive.zip --repository /tmp/fid-task16-cli-KjxBpO/target --idempotency-key task16-direct-cli` - PASS - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"export.import","result":{"import":{"project_id":"project-0b7bf012d95afc763c9b5a706b4c8a47","archive_digest":"sha256:bfc37ca738c91111eaae92d7b200ac2622e50584a2e3d9f238a14721304bb0b0","event_count":1,"terminal_event_hash":"sha256:826a110a565da6901e0f173e7d85164f044c60d2c6f0ba0a63d944116627a8bd","projection_schema":18,"target_root":"/tmp/fid-task16-cli-KjxBpO/target","replayed":false}}}` (terminated `exit_code: 0`)

## Exact quality-gate results

Working directory for all six commands: `/home/mark/projects/fid`. These are the task-spec commands in order and were run on the final implementation tree.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	100.101s	coverage: 62.1% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:									(statements)			60.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	269.692s` (terminated `exit_code: 0`)
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	25.376s` (terminated `exit_code: 0`)

Benchmark row: `22689 ms/rebuild` (reported by the same command).

## Deviations from spec

- None. The initial focused failure was the missing archive boolean allowance; adding `review_findings.blocking` to the accepted semantic boolean columns produced the passing focused and full gates. The two intermediate focused failures were new test-fixture defects and were repaired in the same test file. All failed attempts above remain recorded; the exact same focused command later passed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `d02614e6c26757076a91d8a12e08dfa328fe699d`

## Commit subject

- `fix(export): validate portable Phase 23 review recovery`
