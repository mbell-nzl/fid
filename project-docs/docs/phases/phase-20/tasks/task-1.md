# Phase 20 - Task 1: Response contract foundation

## Objective

Create the domain-neutral `fid-cli.v1` Response Envelope and closed reason-code
registry without migrating a command. This isolates transport mechanics from
project discovery and domain error mapping.

## Validation Contract Coverage

- VAL-001 and VAL-002 (foundation): exact success/failure key sets, channels,
  and serialization are proven at the `internal/cli` boundary. Tasks 2-3 prove
  them through real commands.

## Binding Constraints

- Implement `../contracts.md` sections 2.2-2.5 literally.
- `internal/cli` has no domain or Cobra imports.
- `internal-error` is a backstop, never an expected specified failure.
- The reason-code registry is closed and validates code/category/retryability,
  subject type, exact details object, and redaction policy.

## Scope

### Files To Create Or Modify

- `internal/cli/envelope.go`, `internal/cli/envelope_test.go`.
- `internal/cli/codes.go`, `internal/cli/codes_test.go`.

### Files To Read For Context

- `project-docs/docs/phases/phase-20/contracts.md` sections 2.2-2.5.
- `project-docs/specs/06_DESIGN_PLAN_CLI.md` sections 2 and 5.

### Out Of Scope

- Cobra flags, project discovery, command migration, and domain error mapping.

## Acceptance Criteria

- [ ] Success serializes exactly `contract_version`, `ok`, `operation`, `result`, `event`. (VAL-001)
- [ ] Failure serializes exactly `contract_version`, `ok`, `operation`, `error`. (VAL-002)
- [ ] Required null/object fields are not omitted and one trailing newline is emitted. (VAL-001, VAL-002)
- [ ] Every registered code has one category, retryability value, subject rule,
      details schema, and redaction rule; unknown codes reject. (VAL-002)

## First Failing Test

`internal/cli/envelope_test.go` - `TestWriteSuccessUsesExactContractKeys`.

## Documentation Updates

- None expected.

## Exact Commands To Run

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go build -o /tmp/fid ./cmd/fid
```

## Notes For Developer

- Use `encoding/json`; do not introduce a formatter abstraction or dependency.
- Completion report: `task-1-done.md`.
