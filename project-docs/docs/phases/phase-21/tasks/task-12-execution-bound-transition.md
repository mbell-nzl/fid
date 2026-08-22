# Phase 21 - Task 12: Execution-Bound Transition Foundation

## Objective

Enforce the complete active Envelope, Actor, Claim, Run, repository, expected Issue revision, workflow edge, and applicable-Obligation predicate for execution-bound transitions. Return every independent reason without partial effects and expose the durable facts needed by later resume.

## Validation Contract Coverage

- VAL-026: execution-bound transition succeeds only when every accepted condition qualifies.
- VAL-027: rejected transitions are atomic and return every reason in fixed order.

## Binding Constraints

- Use the shared authority and evaluation vocabulary; do not invent a second readiness model.
- Keep deterministic reason precedence from contract §10.
- CF-032 applies because `BuildIssueRevisionEffects` is reopened: add invalid-input tests before implementation changes.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/issues/`
- `internal/policy/`
- `internal/validation/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§7,10`
- `project-docs/specs/01_PRD.md PRD-WF-003..005`
- `project-docs/docs/current/decisions_active.md CF-032`

### Out Of Scope

- Evidence linking/freshness, final completion, and resume rendering.

## Acceptance Criteria

- [ ] Exactly the all-qualifying condition row advances Issue once. (VAL-026)
- [ ] Every failing row returns all independent reasons in exact order and no Event/revision change. (VAL-027)
- [ ] CF-032 invalid revision-effect inputs are covered before code change. (VAL-027)

## First Failing Test

First add CF-032 failing invalid-input tests, then the full execution-bound transition condition matrix.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-12-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-12-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 7 and 10-11 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-12-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
