# Phase 21 - Task 1: Operational Schema And Projection - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added the accepted additive Phase 21 operational-kernel migration 0013 and its reverse-order development rollback.
- Added integration coverage for the schema inventory, representative accepted and rejected rows, and the 0013 up/down/up round trip.

## Files created or modified

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0013_p21_operational_kernel.down.sql`
- `internal/db/operational_kernel_migration_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-1-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-1-done.md`

## Validation Contract Coverage

- VAL-001: PASS - `TestMigration0013OperationalKernelInventoryAndRoundTrip` in `internal/db/operational_kernel_migration_test.go:13` validates the 0013 inventory and its checks/FKs/indexes, exercises representative valid and rejected rows at `:136`, compares the 0013-down state to the pre-0013 snapshot at `:31`, and compares the second up state to the first at `:38`. The focused green run ended `ok  \tgithub.com/mbell-nzl/fid/internal/db\t1.649s` (terminated `exit_code: 0`).

## Tests run

- `go test -race ./internal/db/... -run '^TestMigration0013OperationalKernelInventoryAndRoundTrip$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); before migration 0013 existed, the test reported `0013 schema missing table "activation_actor_profiles"`.
- `go test -race ./internal/db/... -run '^TestMigration0013OperationalKernelInventoryAndRoundTrip$' -count=1` - GREEN - `ok  \tgithub.com/mbell-nzl/fid/internal/db\t1.649s` (terminated `exit_code: 0`).
- `go run ./cmd/fid init --project /tmp/fid-task-1-init.wdCmO4 --prefix FID` - PASS - `initialized fid project in .fid/` (terminated `exit_code: 0`).

## Exact quality-gate results

- `go test -race ./internal/db/...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/db\t9.963s` (terminated `exit_code: 0`)
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  \tgithub.com/mbell-nzl/fid/internal/workflow\t24.759s\tcoverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t\t78.7%` (terminated `exit_code: 0`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- Pending commit

## Commit subject

- Pending commit
