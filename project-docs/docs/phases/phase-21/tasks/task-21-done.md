# Phase 21 - Task 21: Final-Schema Maintenance, Init, Rebuild, And Export Parity - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: None.
- Recovery attempts: 0
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added final-binary integration coverage for role-free fresh initialization and delete/rebuild capability parity.
- Added read-only operational integrity checks for journal sequences/bindings, attempt receipts, artifacts, frozen state, and residues of removed tables.
- Migrated the 100k replay fixture to Event-v2 role-free facts and added named-read/export byte-parity coverage after rebuild.
- Added a bounded SQLite busy timeout for concurrent idempotent commits and stabilized the retained signaled-process test under the full race/coverage load.

## Files created or modified

- `cmd/fid/integration/task21_maintenance_rebuild_test.go`
- `cmd/fid/task21_maintenance_rebuild_test.go`
- `cmd/fid/gate_run_test.go`
- `internal/db/db.go`
- `internal/journal/rebuild_scale_test.go`
- `internal/maintenance/check.go`
- `internal/maintenance/operational.go`
- `internal/maintenance/check_test.go`
- `project-docs/docs/phases/phase-21/tasks/task-21-done.md`

## Documentation updated

- `project-docs/docs/phases/phase-21/tasks/task-21-done.md`

## Validation Contract Coverage

- VAL-076: PASS - `TestFinalBinaryInitAndRebuildRemainRoleFree` verifies a final binary `init`, role-free bootstrap event, and absent `roles` table before rebuild at `cmd/fid/integration/task21_maintenance_rebuild_test.go:15`.
- VAL-077: PASS - the same real-process test deletes `.fid/fid.sqlite`, runs `rebuild`, then compares JSON capabilities and reasserts role-free bootstrap/table absence at `cmd/fid/integration/task21_maintenance_rebuild_test.go:31`.
- VAL-078: PASS - `TestCheckReportsFinalSchemaIntegrityClassesWithoutRepair` asserts exact `event-sequence`, `aggregate-sequence`, and `event-binding` findings at `internal/maintenance/check_test.go:150`.
- VAL-079: PASS - the same test asserts exact `attempt-receipt-integrity` and `artifact-integrity` findings at `internal/maintenance/check_test.go:165`.
- VAL-080: PASS - the same test asserts `frozen-state-present` and `removed-state-residue`, plus unchanged journal bytes and SQLite `total_changes`, at `internal/maintenance/check_test.go:169`.
- VAL-081: PASS - `TestFinalSchemaNamedReadsAndExportMatchAfterRebuild` compares capabilities, envelope get, evidence evaluation, and resume JSON before and after a deleted projection rebuild at `cmd/fid/task21_maintenance_rebuild_test.go:10`.
- VAL-082: PASS - the same named-read test compares export validation and exact export bytes before and after rebuild at `cmd/fid/task21_maintenance_rebuild_test.go:10`.
- VAL-083: PASS - `TestScaleReplayFixtureUsesRoleFreeEventV2Facts` requires Event-v2 role-free scale facts at `internal/journal/rebuild_scale_test.go:82`; three fresh-process 100k samples completed with parity assertions in 27,804 ms, 27,866 ms, and 27,906 ms per rebuild, each below 30 seconds.

## Tests run

- Red: `go test ./internal/maintenance -run '^TestCheckReportsFinalSchemaIntegrityClassesWithoutRepair$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the new test exposed no operational findings before implementation.
- Green: `go test ./internal/maintenance -run '^TestCheckReportsFinalSchemaIntegrityClassesWithoutRepair$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/maintenance 0.041s` (terminated `exit_code: 0`).
- Red: `go test ./internal/journal -run '^TestScaleReplayFixtureUsesRoleFreeEventV2Facts$' -count=1` - FAIL - `FAIL` (terminated `exit_code: 1`); the retained scale event reported schema `1`, not Event-v2.
- Green: `go test ./internal/journal -run '^(TestScaleReplayFixtureUsesRoleFreeEventV2Facts|TestRebuildWorkspace100kEvents)$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/internal/journal 1.039s` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/... -run '^TestFinalBinaryInitAndRebuildRemainRoleFree$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid/integration 1.935s` (terminated `exit_code: 0`); this builds and invokes the final `fid` binary for the changed runtime entrypoint.
- `go test ./cmd/fid -run '^TestFinalSchemaNamedReadsAndExportMatchAfterRebuild$' -count=1` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 0.684s` (terminated `exit_code: 0`).
- `go test -race ./internal/journal -run '^TestCommitIdempotentSerializesConcurrentSameKeyCalls$' -count=100` - PASS - `ok github.com/mbell-nzl/fid/internal/journal 52.909s` (terminated `exit_code: 0`).
- `go test -race ./cmd/fid -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=10` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 99.840s` (terminated `exit_code: 0`).

## Exact quality-gate results

- `go test -race ./internal/bootstrap ./internal/maintenance ./internal/journal ./internal/db ./internal/export ./internal/cli ./cmd/fid/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid 339.436s` (terminated `exit_code: 0`).
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok github.com/mbell-nzl/fid/cmd/fid/integration 7.125s` (terminated `exit_code: 0`).
- `go version` - PASS - `go version go1.26.2 linux/amd64` (terminated `exit_code: 0`).
- `go env GOOS GOARCH` - PASS - `amd64` (terminated `exit_code: 0`).
- `uname -srm` - PASS - `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64` (terminated `exit_code: 0`).
- `nproc` - PASS - `6` (terminated `exit_code: 0`).
- `grep -m1 '^model name' /proc/cpuinfo` - PASS - `model name\t: AMD Ryzen 7 5700G with Radeon Graphics` (terminated `exit_code: 0`).
- `grep -m1 '^MemTotal:' /proc/meminfo` - PASS - `MemTotal:        8133236 kB` (terminated `exit_code: 0`).
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` (fresh process sample 1) - PASS - `ok github.com/mbell-nzl/fid/internal/journal 30.506s` (terminated `exit_code: 0`); `BenchmarkRebuildWorkspace100kEvents-6 1 27804406580 ns/op 27804 ms/rebuild`.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` (fresh process sample 2) - PASS - `ok github.com/mbell-nzl/fid/internal/journal 30.543s` (terminated `exit_code: 0`); `BenchmarkRebuildWorkspace100kEvents-6 1 27865944299 ns/op 27866 ms/rebuild`.
- `go test -run '^$' -bench '^BenchmarkRebuildWorkspace100kEvents$' -benchmem -benchtime 1x -timeout 30m ./internal/journal` (fresh process sample 3) - PASS - `ok github.com/mbell-nzl/fid/internal/journal 30.650s` (terminated `exit_code: 0`); `BenchmarkRebuildWorkspace100kEvents-6 1 27906022596 ns/op 27906 ms/rebuild`.
- `go test -bench . -benchmem ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow 2.260s` (terminated `exit_code: 0`).
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`).
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`).
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 27.5% of statements` (terminated `exit_code: 0`).
- `go tool cover -func=coverage.out` - PASS - `total:\t\t\t\t\t\t\t\t\t(statements)\t\t\t\t61.3%` (terminated `exit_code: 0`).

Reference identity for the DEC-158 samples: `go1.26.2 linux/amd64`; `Linux 6.6.114.1-microsoft-standard-WSL2 x86_64`; `6` logical processors; `AMD Ryzen 7 5700G with Radeon Graphics`; `MemTotal 8133236 kB`.

## Deviations from spec

- The first required focused race command failed with `FAIL` (terminated `exit_code: 1`). It exposed a retained V0.1 fixture that still expected no report for now-observed frozen `task_completions` state and two concurrent idempotent writers returning `SQLITE_BUSY`. The accepted final-state contract requires frozen state to be reported, so the retained assertion was converged; `busy_timeout(5000)` now serializes those SQLite writers. The byte-identical rerun passed with `ok github.com/mbell-nzl/fid/cmd/fid 339.436s` (terminated `exit_code: 0`).
- The first required full race/coverage command failed with `FAIL` (terminated `exit_code: 1`). Under coverage load, the retained signaled-process fixture's default 1,000 ms deadline expired before its child could be signaled, producing the distinct timed-out outcome. The fixture now supplies a 10,000 ms deadline only for the signaled case; its exact signal assertion remains unchanged. The byte-identical rerun passed with `ok github.com/mbell-nzl/fid/internal/workflow (cached) coverage: 27.5% of statements` (terminated `exit_code: 0`).
- No product, schema, or public-interface deviation was made. DEC-173 remains intact: no role lock/registry/fixture was restored and no dropped-table Effect tolerance was added.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `76ad781727b9e8545e069c4ac0e9f926012e6af9`

## Commit subject

- `test(gate): stabilize signaled process outcome`
