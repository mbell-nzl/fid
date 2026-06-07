## Summary of what was implemented

Implemented Phase 2 Task 1 by adding migration `0002` for the remaining v0.1 base and stub tables, bumping `internal/db.ProjectionSchemaVersion` from `1` to `2`, and extending package-level `internal/db` and `internal/journal` tests to prove schema presence, 0002 step-down/up behavior, representative journal projection reachability, and stale-version rebuild behavior.

## Files created or modified

- `migrations/0002_v01_tables.up.sql`
- `migrations/0002_v01_tables.down.sql`
- `internal/db/projection_meta.go`
- `internal/db/db_test.go`
- `internal/db/migrate_v01_tables_test.go`
- `internal/journal/v01_tables_test.go`
- `project-docs/docs/phases/phase-2/tasks/task-1-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-2/tasks/task-1-completion.md`

## Validation Contract Coverage

- `VAL-001` PASS: `internal/db` migration introspection now asserts all Phase 2 Task 1 tables exist, expected representative columns are present, `gate_runs` carries the static `CHECK (attempt_number > 0)`, and representative foreign keys resolve correctly for `epic_metadata`, `acceptance_criteria`, `ac_verifications`, `gate_runs`, `task_completions`, `concerns`, `agent_runs`, `issue_labels`, `decisions`, `escalations`, and `triage_events`.
- `VAL-002` PASS: `TestMigration0002DownLeavesPhaseOneSchemaAndUpRestoresPhaseTwoTables` steps the migrator back one version, confirms `projection_meta`/`roles`/`issues` remain while Task 1 tables disappear, then migrates back up and confirms the Task 1 tables return.
- `VAL-003` PASS: `TestProjectEventPopulatesPhaseTwoRepresentativeTablesViaJournalEffects` in `internal/journal` projects one event with ordered insert effects into `requirements`, `labels`, `epic_metadata`, `acceptance_criteria`, `task_completions`, and `concerns`, then reads the projected rows back and matches the inserted values.
- `VAL-004` PASS: `ProjectionSchemaVersion` is now `2`; `TestDecideProjectionFreshnessRequiresPhaseTwoSchemaVersion` proves stored version `1` rebuilds and stored version `2` stays fresh, and `TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion` proves a persisted prior-version checkpoint triggers a real rebuild and returns the checkpoint to version `2`.

## Tests run

- `go test ./internal/db -run 'TestMigrateUpCreatesBaseSchema|TestMigrateUpCreatesPhaseTwoV01TablesAndConstraints|TestMigration0002DownLeavesPhaseOneSchemaAndUpRestoresPhaseTwoTables|TestDecideProjectionFreshnessRequiresPhaseTwoSchemaVersion'`
  Red: failed before implementation because Phase 2 tables were absent and `ProjectionSchemaVersion` was still `1`.
  Green: passed after adding migration `0002` and bumping the schema version.
- `go test ./internal/journal -run 'TestProjectEventPopulatesPhaseTwoRepresentativeTablesViaJournalEffects|TestEnsureProjectionFreshRebuildsStoredPriorProjectionSchemaVersion'`
  Red: first failed in test compilation due an incorrect helper signature, then failed on SQLite date representation until the assertion was normalized with `date(...)`.
  Green: passed after fixing the test helper and canonicalizing the `DATE` read-back assertion.
- `go test ./internal/db ./internal/journal`

## Exact quality-gate results

- `gofmt -l cmd internal` — PASS
- `golangci-lint run ./...` — PASS (`0 issues.`)
- `go test -race -coverprofile=coverage.out ./...` — PASS
- `go tool cover -func=coverage.out` — PASS (`total: 71.9%`; `internal/journal: 80.6%`; `internal/issues: 85.2%`; `cmd/fid: 70.5%`)
- `go test -tags integration -race ./...` — NOT NEEDED; no `//go:build integration` tests exist in the repo for Task 1.
- `go build -o /tmp/fid ./cmd/fid` — PASS (conditional support check run because the task spec recommends verifying the binary after the version bump)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- PENDING

## Commit subject

- PENDING
