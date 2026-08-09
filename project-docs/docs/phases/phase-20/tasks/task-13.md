# Phase 20 - Task 13: Maintenance inspection contract

## Objective

Migrate `maintenance check` to its exact Response Envelope schema and typed
integrity behavior. `journal inspect` deliberately remains legacy and moves to
Task 14's direct-JSON conversion.

## Validation Contract Coverage

- VAL-002, VAL-008, VAL-017, VAL-027, the maintenance row of VAL-044, and the
  maintenance ordering/limit cases of VAL-045.

## Binding Constraints

- Maintenance findings are identical with and without `--warn-as-error`; only
  status/envelope branch/channel/exit behavior changes. A fail result is a
  `projection-divergence` failure whose `error.details` carries the exact
  maintenance result schema.
- `maintenance vacuum` remains a guarded legacy command.
- Remove `--format` only from `maintenance check`; do not touch legacy journal inspection.
- Read the resolved project root from Task 2 command context.
- This task contracts the currently useful journal/projection/reference checks;
  it does not claim the later Attempt/Claim/binding/artifact/activation checks.
- Implement exactly the closed error/warning check-name registry in section
  3.5. `started_without_claim` uses a temporary fixed 300-second grace and must
  not read or reinterpret target `abandoned_attempt_grace_seconds`.

## Scope

### Files To Create Or Modify

- `cmd/fid/maintenance.go`, direct result structs, and focused tests.
- Typed error mapping needed for projection/journal failures.

### Files To Read For Context

- `../contracts.md` sections 2.4-2.5, 3.2, 3.5, 3.7, and 6; playbook integrity guidance.

### Out Of Scope

- New maintenance checks, vacuum redesign, revert execution, and recovery work.

## Acceptance Criteria

- [ ] Maintenance accepts only its closed request and result keys/types are exact; warn-only succeeds on
      stdout and warn-as-error fails on stderr with identical findings in
      `error.details`. (VAL-002, VAL-027, VAL-044)
- [ ] Findings use only the closed check registry; the started-without-claim
      boundary is fixed at 300 seconds and independent of Attempt config. (VAL-027)
- [ ] Findings obey exact deterministic order and bound without truncation. (VAL-045)
- [ ] The read returns `event:null`; Human success uses stdout. (VAL-008, VAL-017)
- [ ] `maintenance vacuum --json` rejects as legacy rather than accepting an ignored flag. (VAL-044)

## First Failing Test

`cmd/fid/maintenance_test.go` - `TestMaintenanceCheckJSONMatchesContract`.

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

- Completion report: `task-13-done.md`.
