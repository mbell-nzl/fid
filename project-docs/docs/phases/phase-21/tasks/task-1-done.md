# Phase 21 - Task 1: Operational Schema And Projection - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: qa-rework
- Condition: Review-1 MF-1 required exhaustive, canonical migration-0013 inventory assertions for VAL-001.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the accepted additive Phase 21 operational-kernel migration 0013 and its reverse-order development rollback.
- Added integration coverage for the schema inventory, representative accepted and rejected rows, and the 0013 up/down/up round trip.
- Reworked the migration-0013 inventory into a fixed, exhaustive 23-table/four-index schema-object delta that compares each canonical DDL definition, including every CHECK and foreign key.

## Files created or modified

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0013_p21_operational_kernel.down.sql`
- `internal/db/operational_kernel_migration_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-1-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-1-done.md`

## Validation Contract Coverage

- VAL-001: PASS - `TestMigration0013OperationalKernelInventoryAndRoundTrip` in `internal/db/operational_kernel_migration_test.go:14` compares the complete 0013-only schema-object delta against fixed canonical table/index DDL at `:88` and `:145`; this covers all tables, named indexes, CHECKs, and foreign keys, while retaining representative valid/rejected rows at `:429` and the up/down/up comparisons at `:35` and `:42`. `TestOperationalKernelInventoryComparisonRejectsSchemaDrift` at `:47` proves extra table, index, CHECK, and foreign-key drift is rejected. The focused green run ended `ok  	github.com/mbell-nzl/fid/internal/db	1.679s` (terminated `exit_code: 0`).

## Tests run

- `go test -race ./internal/db/... -run '^TestMigration0013OperationalKernelInventoryAndRoundTrip$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); before migration 0013 existed, the test reported `0013 schema missing table "activation_actor_profiles"`.
- `go test -race ./internal/db/... -run '^TestMigration0013OperationalKernelInventoryAndRoundTrip$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/internal/db	1.649s` (terminated `exit_code: 0`).
- `go test -race ./internal/db/... -run '^TestOperationalKernelInventoryComparisonRejectsUnexpectedObjects$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); `undefined: compareOperationalKernelInventory`.
- `go test -race ./internal/db/... -run '^(TestOperationalKernelInventoryComparisonRejectsUnexpectedObjects|TestMigration0013OperationalKernelInventoryAndRoundTrip)$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/internal/db	2.371s` (terminated `exit_code: 0`).
- `go test -race ./internal/db/... -run '^(TestOperationalKernelInventoryComparisonRejectsSchemaDrift|TestMigration0013OperationalKernelInventoryAndRoundTrip)$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/internal/db	1.679s` (terminated `exit_code: 0`); the comparator rejects an extra table, index, CHECK, and foreign key.
- `go run ./cmd/fid init --project /tmp/fid-task-1-init.wdCmO4 --prefix FID` - PASS - `initialized fid project in .fid/` (terminated `exit_code: 0`)

## Exact quality-gate results

- `go test -race ./internal/db/...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/db	9.880s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 75.5% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%` (terminated `exit_code: 0`)

## Deviations from spec

- Review-1 MF-1 recovery added the exhaustive inventory assertion without changing VAL-001 or migration 0013. The first focused exact-inventory attempt failed with `FAIL` (terminated `exit_code: 1`) because `compareOperationalKernelInventory` was not yet implemented; the root cause was the previous presence-only test. The repair added fixed canonical DDL comparison and an unexpected-object rejection test; the focused rerun passed with `ok  	github.com/mbell-nzl/fid/internal/db	2.371s` (terminated `exit_code: 0`).
- The prior Phase Runner report-lint failure remains preserved in the run ledger: it terminated `exit_code: 1` because one PASS termination-evidence line ended with punctuation after the required closing parenthesis. The report-only correction retained the command, token, status, and observed exit-code evidence; its later two-argument report-lint rerun returned `OUTCOME COMPLETE` (terminated `exit_code: 0`). Report lint is not a Developer task-gate command for this recovery.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `a094e690815c9940ca926bcad68aedacea663c43`

## Commit subject

- `feat(db): add phase 21 operational kernel schema`
