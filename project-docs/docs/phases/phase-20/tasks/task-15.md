# Phase 20 - Task 15: Remaining format caller and test convergence

## Objective

Remove remaining `--format` invocations and assumptions from broad command,
integration, and parity tests after Tasks 8-14 have converted every live
handler. Keep behavioral assertions intact so Task 16 can delete the formatter
framework against a green, caller-free tree.

## Validation Contract Coverage

- The whole-tree test/caller portion of VAL-031; VAL-029 regression coverage
  must remain green while callers are rewritten.

## Binding Constraints

- Update command-family tests and remaining V1/e2e invocations; do not delete
  behavioral assertions merely to remove a flag.
- Leave `internal/format` in place until Task 16 so source guards and the
  compiler identify missed callers.

## Scope

### Files To Create Or Modify

- Remaining command/e2e tests that pass `--format` after Tasks 8-14.
- `internal/db/issues_workflow_fk_cli_test.go` and retained V1/V2a parity suites.
- Root command-tree guard test for zero runtime format flags.

### Files To Read For Context

- `../contracts.md` sections 5.2, 5.3, and 6.

### Out Of Scope

- Deleting `internal/format`, YAML dependency removal, legacy result redesign, and command stubs.

## Acceptance Criteria

- [ ] Live Cobra tree contains no universal format flag. (VAL-031)
- [ ] All remaining tests invoke the new surface without weakening semantic assertions. (VAL-029, VAL-031)
- [ ] Source/test search finds no live `--format` invocation outside historical docs. (VAL-031)
- [ ] Explicit root flags remain rejected on legacy leaves in retained regression coverage. (VAL-029)

## First Failing Test

`cmd/fid/contract_surface_test.go` -
`TestLiveTestsContainNoUniversalFormatInvocation` (repository source guard).

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./...
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- A discovered live product caller belongs to its owning Task 8-14 and is a
  blocker here, not permission to widen this mechanical cleanup task.
- Completion report: `task-15-done.md`.
