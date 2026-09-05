# Phase 22 - Task 1: Recovery schema and replay foundation - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Replaced the accepted development-baseline CHECK values in migration 0013 and added the exact schema-15 recovery up/down migrations.
- Added fresh-init, full down/up, constraint, live-apply, and replay coverage for recovery tables and rejected Effects without enabling recovery domain operations.
- Centralized active migration reporting at schema 15 for rebuild and export while preserving the separate projection-freshness schema value.
- Converged retained Phase 21 schema and Attention-boundary assertions displaced by the accepted Phase 22 schema surface.

## Files created or modified

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0015_p22_recovery.up.sql`
- `migrations/0015_p22_recovery.down.sql`
- `internal/db/projection_meta.go`
- `internal/db/operational_kernel_migration_test.go`
- `internal/db/recovery_migration_test.go`
- `internal/journal/project_effects.go`
- `internal/journal/recovery_schema_test.go`
- `internal/export/export.go`
- `internal/export/export_test.go`
- `cmd/fid/rebuild.go`
- `cmd/fid/rebuild_test.go`
- `cmd/fid/contract_surface_test.go`
- `cmd/fid/integration/task22_recovery_schema_test.go`
- `project-docs/docs/phases/phase-22/tasks/task-1-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-22/tasks/task-1-done.md`

## Validation Contract Coverage

- VAL-001: PASS - `TestMigration0015RecoverySchemaConstraintsAndRoundTrip` in `internal/db/recovery_migration_test.go:9` verifies all recovery tables, the non-null `run_recoveries.observation_id` column at `:69`, every accepted Stop/Attention code and both invalid direct inserts at `:94`; `TestRecoverySchemaFreshBinaryInitAndRebuild` in `cmd/fid/integration/task22_recovery_schema_test.go:13` proves the built binary installs the recovery schema.
- VAL-002: PASS - `TestMigration0015RecoverySchemaConstraintsAndRoundTrip` runs full `MigrateDown` then `MigrateUp` at `internal/db/recovery_migration_test.go:32`, verifies the same tables at `:43`, and proves every recovery table is empty at `:121`.
- VAL-003: PASS - `TestRecoverySchemaFreshInitAndReplay` in `internal/journal/recovery_schema_test.go:14` rejects the unavailable `run.recovered` Effect through live apply at `:34` and replay at `:49`; it also rejects an invalid Stop code in live apply at `:60` and replay at `:76`.

## Tests run

- `go test -race ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); the first attempt was a test-harness compile failure: `internal/journal/recovery_schema_test.go:40:16: undefined: openProjectDB`.
- `go test -race ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); after correcting only the test harness, the named test produced `run_recoveries table absent after fresh init` before product edits.
- `go test -race ./internal/journal -run '^TestRecoverySchemaFreshInitAndReplay$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/internal/journal	2.670s` (terminated `exit_code: 0`).
- `go test -race ./internal/db -run '^TestMigration0015RecoverySchemaConstraintsAndRoundTrip$' -count=1` - RED - `FAIL` (terminated `exit_code: 1`); before the migration, `table "run_runtime" does not exist`.
- `go test -race ./internal/db -run '^TestMigration0015RecoverySchemaConstraintsAndRoundTrip$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/internal/db	3.751s` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestRecoverySchemaFreshBinaryInitAndRebuild$' -count=1` - GREEN - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	3.352s` (terminated `exit_code: 0`).
- `/tmp/fid-phase22-task1 --project /tmp/fid-phase22-task1.Sw07ej --json init --prefix P22` - GREEN - `{"contract_version":"fid-cli.v1","event":{"id":"EVT-20260905-000001","occurred_at":"2026-09-05T02:05:09.402846Z"},"ok":true,"operation":"project.init","result":{"journal_path":".fid/journal/2026-09-05.jsonl","project":{"id":"project-252b9f9aeebb071de1f086771f08c54f","issue_prefix":"P22","root":"/tmp/fid-phase22-task1.Sw07ej"}}}` (terminated `exit_code: 0`).
- `/tmp/fid-phase22-task1 --project /tmp/fid-phase22-task1.Sw07ej --json rebuild` - GREEN - `{"contract_version":"fid-cli.v1","event":null,"ok":true,"operation":"project.rebuild","result":{"events_replayed":1,"files_replayed":1,"schema_version":15,"upgraded_events":0}}` (terminated `exit_code: 0`).

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Attempt 1.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Attempt 2.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
  - Attempt 3.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 1.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 2.
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
  - Attempt 3.
- `go test -race -coverprofile=coverage.out ./...` - FAIL - `FAIL` (terminated `exit_code: 1`).
  - Attempt 1.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	35.791s	coverage: 27.4% of statements` (terminated `exit_code: 0`)
  - Attempt 2.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 27.4% of statements` (terminated `exit_code: 0`)
  - Attempt 3.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.2%` (terminated `exit_code: 0`)
  - Attempt 1.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)				61.2%` (terminated `exit_code: 0`)
  - Attempt 2.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.145s` (terminated `exit_code: 0`)
  - Attempt 1.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	7.933s` (terminated `exit_code: 0`)
  - Attempt 2.
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`)
  - Attempt 1.
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`)
  - Attempt 2.
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`)
  - Attempt 1.
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`)
  - Attempt 2.
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`)
  - Attempt 1.
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`)
  - Attempt 2.
- `nproc` - PASS - `6` (terminated `exit_code: 0`)
  - Attempt 1.
- `nproc` - PASS - `6` (terminated `exit_code: 0`)
  - Attempt 2.
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name	: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`)
  - Attempt 1.
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name	: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`)
  - Attempt 2.
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`)
  - Attempt 1.
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`)
  - Attempt 2.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	31.083s` (terminated `exit_code: 0`)
  - Attempt 1; sample `28341 ms/rebuild`, below the unchanged `<30000 ms` assertion.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` - PASS - `ok  	github.com/mbell-nzl/fid/internal/journal	30.587s` (terminated `exit_code: 0`)
  - Attempt 2; sample `27805 ms/rebuild`, below the unchanged `<30000 ms` assertion.

## Deviations from spec

- No contract or command definition changed. The first required race/coverage attempt failed with `FAIL` (terminated `exit_code: 1`) because the accepted schema-15 changes invalidated two retained Phase 21 assertions: `TestDecideProjectionFreshnessRequiresPhase19SchemaVersion` had frozen the separate projection-freshness value, and `TestPhase21AttentionBoundary` prohibited any Attention schema/applier surface. The causal repair preserved projection freshness at 5, introduced a separate active migration value of 15 for rebuild/export, and constrained the retained Attention inventory to its continuing CLI/workflow no-domain boundary. The byte-identical later attempts passed with `ok  	github.com/mbell-nzl/fid/internal/workflow	35.791s	coverage: 27.4% of statements` and then the final cached workflow line, both terminated `exit_code: 0`.
- Convergence provenance: `internal/db/migrate_v01_tables_test.go` was introduced by `b63a3cd3f237d48990404d6b9c7099db69ce998e` (`test(convergence): finalize task-20 V1 fixture removal`, Refs: phase-21/task-20); `cmd/fid/contract_surface_test.go` lines 70-98 were introduced by `7277139ee5d5dde4cd3627335d39b14f616f46a1` (`feat(evidence): evaluate trusted receipt evidence`, Refs: phase-21/task-17). The accepted task-1 schema and unavailable-Event foundation displaced their older assertions. Both are executed by `go test -race -coverprofile=coverage.out ./...`.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `72c5974f6698fa08a4fe1a24099b3a93b3882f1d`

## Commit subject

- `feat(recovery): add schema and replay foundation`
