# Phase 3 - Task 4: Rebuild Performance Verification At Single-Project Scale - Completion Report

## Summary of what was implemented

- Added a scale-focused `internal/journal` determinism test that seeds a valid mixed journal, replays it once, rebuilds after deleting the DB, and rebuilds a second time to prove stable projection parity across all three states.
- Added reusable scale-journal seed helpers that generate valid `issue.created` events plus representative acceptance-criteria, label-assignment, link, and status-transition events without changing production replay behavior.
- Added `go test -bench` coverage for ~1k and 10k event rebuilds against the real `RebuildWorkspace` path and reported explicit elapsed `ms/rebuild` metrics for VAL-006 evidence.

## Files created or modified

- `internal/journal/rebuild_scale_test.go`
- `internal/journal/rebuild_bench_test.go`
- `project-docs/docs/phases/phase-3/tasks/task-4-completion.md`

## Documentation updated

- `project-docs/docs/phases/phase-3/tasks/task-4-completion.md`

## Validation Contract Coverage

- `VAL-006` - PASS. `internal/journal/rebuild_bench_test.go` adds `BenchmarkRebuildWorkspace1kEvents` and `BenchmarkRebuildWorkspace10kEvents` over the real `RebuildWorkspace` path. The required `go test -bench . -benchmem ./...` run reported `355.4 ms/rebuild` for ~1k events and `3483 ms/rebuild` for 10k events, satisfying the task thresholds of `< 1s` at ~1k and `< 5s` at 10k.
- `VAL-007` - PASS. `internal/journal/rebuild_scale_test.go` seeds a valid 1k-event mixed journal and proves three invariants on the real replay path: the replayed event count matches the seeded count, the projection state after delete-DB + rebuild matches the initial replay state, and a second rebuild produces the same projection state again.

## Tests run

- `go test ./internal/journal -run TestRebuildWorkspaceIsDeterministicAtSingleProjectScale` - FAIL (expected red: missing `seedScaleReplayWorkspace` / `snapshotScaleProjectionState` helpers)
- `go test ./internal/journal -run TestRebuildWorkspaceIsDeterministicAtSingleProjectScale` - PASS
- `go test ./internal/journal -run '^$' -bench 'BenchmarkRebuildWorkspace(1kEvents|10kEvents)$' -benchmem -benchtime=1x` - PASS

## Exact quality-gate results

- `gofmt -l cmd internal` - PASS
- `golangci-lint run ./...` - PASS
- `go test -race -coverprofile=coverage.out ./...` - PASS
- `go tool cover -func=coverage.out` - PASS (`internal/journal` 80.6%, `cmd/fid` 72.5%, total 74.9%)
- `go test -bench . -benchmem ./...` - PASS (`BenchmarkRebuildWorkspace1kEvents` `355374632 ns/op`, `355.4 ms/rebuild`; `BenchmarkRebuildWorkspace10kEvents` `3483337397 ns/op`, `3483 ms/rebuild`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- None.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `6a68b1e`

## Commit subject

- `test(journal): add scale rebuild verification`
