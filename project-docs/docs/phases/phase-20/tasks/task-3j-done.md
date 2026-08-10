# Phase 20 - Task 3j: Bounded command-test parallelism - Completion Report

## Summary of what was implemented

- Added a goroutine-safe command execution result helper and a race-enabled
  concurrent regression proving two cloned workspaces cannot cross journal
  state or mutate the pristine template.
- Added `t.Parallel()` to 104 top-level tests in the 25 cloned-workspace files
  plus the fixture owner.
- Retained serial execution for protected real-init, journal, rebuild,
  maintenance durability, Issue/AC compatibility, and claim-list rebuild tests.
- Used only Go's default scheduler bound: six available CPUs on this host with
  `GOMAXPROCS` unset; no custom worker pool or goroutine fan-out was added.
- Reduced uncached package race time from 88.748s to 72.130s (18.7%) and the
  exact five-count gate from 431.955s to 353.465s.

## Files created or modified

- `cmd/fid/workspace_fixture_test.go`
- 25 cloned-workspace command-test files receiving top-level `t.Parallel()`
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3j-done.md`

## Documentation updated

- `infra-improvement-pack.md` records parallel/serial ownership, excluded
  global-state categories, scheduler bound, concurrency proof, and timings.

## Validation Contract Coverage

- No VAL-NNN applies. Existing assertions are unchanged.
  `TestConcurrentClonedWorkspaceCommandsRemainIsolated` proves concurrent
  command, journal, sibling-clone, and template isolation under the race
  detector.

## Tests run

- `go test -run '^TestConcurrentClonedWorkspaceCommandsRemainIsolated$' ./cmd/fid` - EXPECTED RED - `FAIL`.
- `go test -race -run '^TestConcurrentClonedWorkspaceCommandsRemainIsolated$' ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	1.493s`.
- Post-parallel `go test -race -count=1 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	72.130s`.
- `nproc && go env GOMAXPROCS` - PASS - `6` (the following empty value means
  `GOMAXPROCS` is unset).

## Exact quality-gate results

- `go test -race -count=5 ./cmd/fid` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid	353.465s`.
- `gofmt -l cmd internal` - PASS - `exit 0 (no output)`.
- `golangci-lint run ./...` - PASS - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					77.7%`.
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.406s`.
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)`.

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- Parallelism cannot reduce the protected serial initialization owners. The
  current boundary favors evidence integrity over a broader but riskier speedup.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `338c5dc9289aa85d1592fe0817f5cbb10d4cf609`

## Commit subject

- `test(cli): parallelize isolated command tests`
