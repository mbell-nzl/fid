# Phase 20 - Task 17: Delete promised-command stubs

## Objective

Delete the unimplemented decision, escalation, and triage command stubs so
unsupported product promises are absent from the live command tree.

## Validation Contract Coverage

- VAL-034.

## Binding Constraints

- Delete only the command adapters, registrations, stub tests, and version-gate helper.
- Retain `decisions`, `escalations`, and `triage_events` tables until Phase 25
  because maintenance still audits them.
- Unknown-command behavior comes from Cobra; do not add replacement tombstones.

## Scope

### Files To Create Or Modify

- Delete `cmd/fid/governance_stubs.go` and its test.
- `cmd/fid/root.go` and root command-tree tests.

### Files To Read For Context

- `../contracts.md` section 5.4; roadmap Phase 20/25 removal ownership.

### Out Of Scope

- Persistence tables, real future Decision References, or any new command family.

## Acceptance Criteria

- [ ] `fid decision`, `fid escalation`, and `fid triage` are unknown commands. (VAL-034)
- [ ] `versionGateError` and `ships in v0.2` are absent from `cmd/`. (VAL-034)
- [ ] Maintenance tests still pass with retained tables. (VAL-034)

## First Failing Test

`cmd/fid/root_test.go` - `TestPromisedGovernanceCommandsAreAbsent`.

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

- Completion report: `task-17-done.md`.
