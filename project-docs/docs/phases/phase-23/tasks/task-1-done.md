# Phase 23 - Task 1: Review and judgment schema baseline - Completion Report

## Workflow outcome

- Task outcome: RECOVERY_REQUIRED
- Category: session-continuation
- Condition: The tagged integration and rebuild benchmark commands remain to be run after the completed repository-wide race gate.
- Recovery attempts: 0
- Next action: Run the remaining exact commands byte-for-byte, record their observed termination evidence, exercise the changed integration entrypoint, then complete the ordinary two-commit report sequence.
- Human attention required: no

## Summary of what was implemented

- Added schema-17 migration DDL, the seven accepted schema-baseline hunks, the nine-event progressive fence and closed journal/archive registrations, schema-17 archive validation, and integration-binary reuse for re-executed integration children.

## Files created or modified

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0015_p22_recovery.up.sql`
- `migrations/0017_p23_review_judgment.up.sql`
- `migrations/0017_p23_review_judgment.down.sql`
- `internal/db/`, `internal/journal/`, `internal/export/`, and `cmd/fid/integration/main_test.go`
- `project-docs/specs/03_TRD_SCHEMA.md`

## Documentation updated

- `project-docs/specs/03_TRD_SCHEMA.md` adds Phase 23 section 8.4, incorporating the accepted phase schema authority.

## Validation Contract Coverage

- VAL-001: IN PROGRESS - `TestMigration0017CreatesExactlyTheAcceptedObjects` passed after the migration was added; exact inventory expansion remains for recovery.
- VAL-002: IN PROGRESS - accepted discriminator DDL is present; behavior matrix remains for recovery.
- VAL-003: IN PROGRESS - all nine fence keys are registered; live/replay rejection proof remains for recovery.
- VAL-004: IN PROGRESS - migration and package suite passed; explicit pragma proof remains for recovery.
- VAL-005: IN PROGRESS - accepted link staging is present; seeded in-place migration proof remains for recovery.
- VAL-075: IN PROGRESS - seven accepted source hunks are applied; explicit baseline-diff proof remains for recovery.

## Tests run

- `go test ./internal/db -run '^TestMigration0017CreatesExactlyTheAcceptedObjects$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); expected red before migration 0017 existed.
- `go test ./internal/db -run '^TestMigration0017CreatesExactlyTheAcceptedObjects$' -count=1` - PASS - `ok  github.com/mbell-nzl/fid/internal/db 0.063s` (terminated `exit_code: 0`)
- `go test ./internal/db ./internal/journal ./internal/export` - PASS - `ok  github.com/mbell-nzl/fid/internal/export 15.296s` (terminated `exit_code: 0`)

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -timeout 30m -coverprofile=coverage.out ./...` - PASS - `ok  github.com/mbell-nzl/fid/internal/workflow 595.167s coverage: 62.6% of statements` (terminated `exit_code: 0`)
- Wall times: `cmd/fid` 581.395s; `internal/workflow` 595.167s.
- `go tool cover -func=coverage.out` - PASS - `total: (statements) 65.5%` (terminated `exit_code: 0`)
- `go test -tags integration -race -timeout 30m ./cmd/fid/integration/...` - BLOCKED - not started before session checkpoint.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchtime 1x ./internal/journal/` - BLOCKED - not started before session checkpoint.

## Deviations from spec

- First full-package attempt failed before retained schema-16 expectations were converged. The failure showed `active migration = 17, want16` and the old 0013 inventory; the accepted schema-17 replacement caused it. The in-scope retained expectations were converged and `go test ./internal/db ./internal/journal ./internal/export` later passed.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `4aa2320`

## Commit subject

- `feat(phase-23): add review judgment schema baseline`
