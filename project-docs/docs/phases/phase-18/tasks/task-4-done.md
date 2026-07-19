# Phase 18 - Task 4: Migrate issues.workflow_id / status to real FKs (table rebuild) - Completion Report

## Summary of what was implemented

- Added migration 0008, which rebuilds `issues` with a real `workflow_id` → `workflows(id)` FK and a composite `(workflow_id, status)` → `workflow_states(workflow_id, state)` FK.
- Removed the v0.1 workflow/status enum CHECK constraints while preserving the kind, priority, epic-token, sequence, and sibling-UNIQUE checks plus every row value.
- Recreated the `issues_enriched` view and the five triggers defined `ON issues`; retained all other-table triggers and supporting indexes unchanged. The migration and its reversible down migration use deferred FK enforcement plus legacy rename behavior so SQLite can rebuild a referenced table without dropping retained dependents.
- Added regression coverage for FK insert/update enforcement, clean `foreign_key_check`, data and view snapshots, schema/PRAGMA fingerprints through down/up, retained-trigger probes, autoindexes, and rendered issue get/list parity.
- Updated test-only maintenance corruption fixtures and expected integrity findings so they deliberately establish invalid workflow pairs under the new real FKs; no production maintenance or transition consumer changed.

## Files created or modified

- `migrations/0008_v02a_issues_workflow_fk.up.sql`
- `migrations/0008_v02a_issues_workflow_fk.down.sql`
- `internal/db/issues_workflow_fk_test.go`
- `internal/db/issues_workflow_fk_cli_test.go`
- `internal/db/db_test.go`
- `cmd/fid/maintenance_test.go`
- `internal/maintenance/check_test.go`
- `project-docs/docs/phases/phase-18/tasks/task-4-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-18/tasks/task-4-done.md`

## Validation Contract Coverage

- VAL-006: PASS — `TestMigration0008EnforcesWorkflowAndStateForeignKeys` proves unknown workflow and unknown state-pair INSERTs are rejected, valid custom workflow/state pairs are accepted after enum removal, invalid workflow/state UPDATEs are rejected, and `PRAGMA foreign_key_list(issues)` plus `PRAGMA foreign_key_check` verify both requested FKs and clean integrity.
- VAL-007: PASS — `TestMigration0008PreservesIssuesAndDependentObjectsThroughDownUp` seeds v7-compatible epic/task/subtask rows covering all `issues` columns, compares canonical base-table and `issues_enriched` snapshots before/after 0008, and repeats the base-table comparison after down/up.
- VAL-018: PASS — the same test compares a complete `sqlite_schema` fingerprint excluding only the intentionally rebuilt `issues` table before/after and after down/up; it includes the recreated view and five `ON issues` triggers, unchanged supporting indexes, every retained 0004/0005 trigger, and structural `PRAGMA index_list`/`index_info` checks for the PK and `(parent_id, sequence)` UNIQUE autoindexes. It also functionally probes all five retained kind triggers and all five retained ref-integrity triggers.
- VAL-019: PASS — `TestMigration0008PreservesIssueGetAndListRenderedBytes` freezes JSON bytes from the shared `issues.Get`/`issues.List` plus formatter path against a migration-0007-only fixture, migrates it through 0008, and requires byte-identical get and list output. The real CLI exercise below independently confirmed `fid issue get`/`list` output before and after `fid rebuild`.

## Tests run

- `go test ./internal/db -run '^TestMigration0008EnforcesWorkflowAndStateForeignKeys$' -count=1` - expected RED: valid custom workflow/state pair was rejected by the v0.1 kind→workflow enum CHECK before migration 0008 existed.
- `go test ./internal/db -run '^TestMigration0008EnforcesWorkflowAndStateForeignKeys$' -count=1` - PASS after migration 0008.
- `go test ./internal/db -run '^TestMigration0008' -count=1` - PASS; FK, preservation, dependent-object, down/up, and rendered-read parity coverage.
- `go test ./internal/db -count=1` - PASS.
- `go test ./cmd/fid -run '^TestMaintenanceCheckCommandRenders(HealthyAndFailingResults|AllFormats)$' -count=1` - PASS; corrupt-fixture bypass and resulting FK findings covered.
- `go test ./internal/maintenance -run '^TestCheckErrorsReportsAllSeededStructuralViolationFamilies$' -count=1` - PASS; seeded workflow violations contribute the expected FK-integrity findings.
- `go test ./... -count=1` - PASS.
- Runtime entrypoint exercise: built `/tmp/fid`, then in `/tmp/fid-task4-DzfKml` ran `fid init`; `fid issue create --kind task --title 'Task 4 CLI FK migration exercise' --by developer`; `fid issue get FID-TASK4-DZFKML-94b516 --format json`; `fid issue list --format json`; `fid rebuild --by tl`; and the same `fid issue get` again — PASS. Rebuild replayed 2 events and the post-rebuild JSON was byte-identical to the pre-rebuild JSON. (`fid rebuild --by developer` was correctly denied by the unchanged policy; retrying as authorized `tl` succeeded.)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS (exit 0; no files listed).
- `golangci-lint run ./...` - PASS (`0 issues.`).
- `go test -race -coverprofile=coverage.out ./...` - PASS (exit 0; coverage profile produced).
- `go tool cover -func=coverage.out` - PASS (exit 0; total statement coverage: 74.5%).
- `go test -tags integration -race ./...` - PASS (exit 0).

## Deviations from spec

- Test-only maintenance fixtures were updated because real FK enforcement rejects their formerly CHECK-only corrupt workflow rows. The production maintenance and transition consumers remain unchanged, as required; the fixtures now deliberately disable FKs only while creating their corrupt test state and assert the resulting integrity findings.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
