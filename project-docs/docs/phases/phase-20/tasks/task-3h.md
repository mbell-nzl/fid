# Phase 20 - Task 3h: Reusable pristine workspace fixture

## Objective

Create a real initialized pristine workspace once and clone it for `cmd/fid`
command tests whose subject is not initialization. Reduce repeated migrations
and rebuild verification without mocking journal or SQLite state.

## Validation Contract Coverage

- No new VAL-NNN applies. Existing command behavior retains its assertions;
  fixture parity is delivery evidence.

## Binding Constraints

- The template is produced by real `InitializeWorkspace`, immutable after
  creation, and copied into a fresh per-test root.
- Clone config, journal, projection, permissions, and layout.
- Rebuild parity and cross-test mutation isolation are behavioral tests.
- Init/bootstrap/migration/rebuild/durability/compatibility tests and tests of
  init output/failure keep real initialization.
- No global cache outside package-owned test lifecycle.

## Scope

### Files To Create Or Modify

- `cmd/fid/workspace_fixture_test.go`
- Selected non-init `cmd/fid/*_test.go` call sites
- `project-docs/docs/phases/phase-20/infra-improvement-pack.md`
- `project-docs/docs/phases/phase-20/tasks/task-3h-done.md`

### Files To Read For Context

- `cmd/fid/init_test.go`
- `internal/bootstrap/init.go`
- `internal/journal/rebuild.go`
- Tasks 3e and 3g reports

### Out Of Scope

- Changing production initialization or migrating init-focused tests.
- Parallel tests; Task 3j owns concurrency.

## Acceptance Criteria

- [ ] Fixture creation uses the complete production initialization path.
- [ ] Two clones start byte-equivalent, mutate independently, and leave the
      template and sibling unchanged. (no VAL)
- [ ] A cloned workspace rebuilds to byte-equivalent projection state.
- [ ] Non-init command tests adopt the helper; a protected inventory proves
      init/durability owners still use real initialization.
- [ ] Ordinary package timing is recorded before and after migration.

## First Failing Test

Add `TestPristineWorkspaceClonePreservesRebuildParityAndIsolation`; confirm it
fails because the clone fixture does not yet exist.

## Documentation Updates

- `infra-improvement-pack.md` - fixture coverage, protected inventory, timing.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go build -o /tmp/fid ./cmd/fid
```

Required-command rules live in `AGENTS.md` (Commands).

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Use stdlib file APIs and preserve modes; copy SQLite/journal only while the
  pristine workspace is closed.
- Completion report: `tasks/task-3h-done.md`.
