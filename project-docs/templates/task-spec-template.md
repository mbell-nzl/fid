# Phase <N> - Task <M>: <Title>

## Objective

<What this task changes and why.>

## Requirements

- PRD refs: <IDs or sections>
- TRD refs: <IDs or sections>
- Design refs: <sections or "N/A">
- Quality refs: <sections or "N/A">

## Validation Contract Coverage

List the `VAL-NNN` assertions this task advances. Assertions live in
`project-docs/docs/phases/phase-<N>/qa/validation-contract.md`.

- VAL-NNN: <one-line restatement of behavior or invariant>

If this task does not advance a validation assertion, state why explicitly
(for example, `documentation-only task; no behavior assertion`).

## Scope

### Files To Create Or Modify

- `<path>`

### Files To Read For Context

- `<path>`

### Out Of Scope

- <explicit non-goal>

## Acceptance Criteria

- [ ] <observable/testable criterion> (VAL-NNN or no VAL)
- [ ] <observable/testable criterion> (VAL-NNN or no VAL)

## TDD Plan

- First failing test:
- Additional tests:
- Regression risks:

## Documentation Updates

- `<doc path>` - <why>

## Scenario / Parity / Audit Notes

- Held-out scenario impact: <None / scenario theme only; do not paste hidden scenario content>
- Cross-surface parity impact: <None / Phase Brief Parity Matrix axes to check>
- Audit or release-readiness impact: <None / template or review to update>

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

The full task gate above is required before task close unless this task
explicitly states otherwise here. A focused test command may be used during the
red/green loop, but it does not substitute for the full task gate at close.
Every required command must be recorded in the completion report as `PASS`, or
as `BLOCKED` with reason in a blocker report. Missing, silently skipped, or
narrower substituted commands are failed gates.

Conditional commands:

```text
go test -tags integration -race ./...
go test -bench . -benchmem ./...
go build -o /tmp/fid ./cmd/fid
```

List any conditional command required by this task and why:

- None.

## Notes For Developer

- <constraints, ordering, known pitfalls>
