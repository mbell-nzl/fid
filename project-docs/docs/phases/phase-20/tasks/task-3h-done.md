# Phase 20 - Task 3h: Reusable pristine workspace fixture - Completion Report

## Workflow outcome

- Task outcome: COMPLETE
- Category: none
- Condition: MF-01 report conformance repaired with fresh current-session Task 3h gate evidence.
- Recovery attempts: 1
- Next action: None.
- Human attention required: no

## Summary of what was implemented

- Added a package-lifecycle pristine workspace created once through the full
  production `bootstrap.InitializeWorkspace` path and closed before tests run.
- Added mode-preserving recursive clone support for directories, regular files,
  and symlinks, with a fresh `t.TempDir()` destination for every test.
- Added a behavioral regression proving byte-and-mode clone parity, sibling and
  template mutation isolation, and byte-identical SQLite projection after a
  cloned workspace rebuild.
- Migrated 116 exact default-init setups in 25 non-protected command-test files;
  retained 50 real-init calls in protected init, rebuild, journal, maintenance
  durability, golden/rebuild scenario, and compatibility-heavy Issue/AC owners.
- Reduced an uncached `cmd/fid` race run from 150.540s to 88.748s (41.0%).

## Files created or modified

- `cmd/fid/workspace_fixture_test.go`
- 25 non-protected command-test files using
  `clonePristineWorkspaceInto`, from `activity_test.go` through
  `workflow_list_test.go`, enumerated by the committed diff and anchor
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3h-done.md`

## Documentation updated

- `infra-improvement-pack.md` records fixture semantics, migrated/protected
  setup counts, parity/isolation evidence, and directly comparable timing.

## Validation Contract Coverage

- No VAL-NNN applies. `TestPristineWorkspaceClonePreservesRebuildParityAndIsolation`
  proves the test-delivery invariants. Product initialization, journal,
  projection, migration, CLI, and durability code is unchanged.

## Tests run

- `go test -run '^TestPristineWorkspaceClonePreservesRebuildParityAndIsolation$' ./cmd/fid` - EXPECTED RED - `FAIL`.
- `go test -run '^TestPristineWorkspaceClonePreservesRebuildParityAndIsolation$' ./cmd/fid` - UNVERIFIED - `ok  	github.com/mbell-nzl/fid/cmd/fid	0.062s`.
- Pre-migration `go test -race -count=1 ./cmd/fid` - UNVERIFIED - `ok  	github.com/mbell-nzl/fid/cmd/fid	150.540s`.
- Post-migration `go test -count=1 ./cmd/fid` - UNVERIFIED - `ok  	github.com/mbell-nzl/fid/cmd/fid	9.809s`.
- Post-migration `go test -race -count=1 ./cmd/fid` - UNVERIFIED - `ok  	github.com/mbell-nzl/fid/cmd/fid	88.748s`.

The historical PASS rows above did not retain observed process termination and
are therefore preserved as UNVERIFIED evidence. The current-session exact gate
results below are the passing evidence for this report continuation.

## Exact quality-gate results

- `gofmt -l cmd internal` - UNVERIFIED - `exit 0 (no output)`.
- `golangci-lint run ./...` - UNVERIFIED - `0 issues.`.
- `go test -race -coverprofile=coverage.out ./...` - UNVERIFIED - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.5% of statements`.
- `go tool cover -func=coverage.out` - UNVERIFIED - `total:										(statements)					77.7%`.
- `go test -tags integration -race ./cmd/fid/integration/...` - UNVERIFIED - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	8.325s`.
- `go build -o /tmp/fid ./cmd/fid` - UNVERIFIED - `exit 0 (no output)`.

- `gofmt -l cmd internal` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)
- `golangci-lint run ./...` - PASS - `0 issues.` (terminated `exit_code: 0`)
- `go test -race -coverprofile=coverage.out ./...` - PASS - `ok  	github.com/mbell-nzl/fid/internal/workflow	(cached)	coverage: 80.6% of statements` (terminated `exit_code: 0`)
- `go tool cover -func=coverage.out` - PASS - `total:										(statements)					78.7%` (terminated `exit_code: 0`)
- `go test -tags integration -race ./cmd/fid/integration/...` - PASS - `ok  	github.com/mbell-nzl/fid/cmd/fid/integration	4.515s` (terminated `exit_code: 0`)
- `go build -o /tmp/fid ./cmd/fid` - PASS - `exit 0 (no output)` (terminated `exit_code: 0`)

## Deviations from spec

- None.

## Concerns, technical debt, or notes for the Tech Lead

- The protected 50-call inventory is deliberately conservative. Later tasks
  may optimize an owner only with a narrower behavioral proof; this task does
  not trade init/rebuild/durability evidence for additional speed.

## Dependency-budget note

- No dependency changes.

## Commit hash

- `bca799f52f55286e8a45cc68c341cb539ba50481`

## Commit subject

- `test(cli): reuse pristine workspace fixtures`
