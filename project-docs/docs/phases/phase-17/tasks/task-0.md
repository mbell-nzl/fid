# Phase 17 - Task 0: Restore the Journal Coverage Gate

## Objective

Add focused tests for observable `internal/journal` behavior until the package
meets the existing COV-001 threshold. This is an entry remediation for the
78.9% baseline found by QA-BLK-01; it changes no product behavior or policy.

## Validation Contract Coverage

- VAL-024: restore `internal/journal` to at least 80% so the unchanged full
  phase gate has an authorized path to pass.

## Binding Constraints

- `08_QUALITY_PLAN.md` COV-001 remains binding at 80%; no waiver or threshold
  change is authorized.
- Add behavioral tests only. Do not modify non-test Go files, exported APIs,
  journal schemas, event bytes, replay semantics, or persistence behavior.
- Tests must assert observable success, failure, or boundary behavior and must
  be capable of failing. Do not add calls whose only purpose is line coverage.
- Use isolated temporary workspaces and the existing journal test helpers.
- If a test exposes a product defect or requires a production seam, stop and
  report it for re-planning instead of expanding this task.

## Scope

### Files To Create Or Modify

- `internal/journal/*_test.go` - add the smallest focused table or regression
  tests needed to cover existing, observable journal behavior.

### Files To Read For Context

- `coverage.out` - QA-BLK-01 baseline branch evidence; do not commit it.
- `internal/journal/*.go` - select uncovered behavior worth asserting.
- `internal/journal/*_test.go` - reuse established fixtures and test style.
- `project-docs/specs/08_QUALITY_PLAN.md` sections 3.1-3.4 and 4.3.
- `project-docs/docs/phases/phase-17/qa/qa-plan.md` QA-BLK-01.

### Out Of Scope

- Any production-code change or refactor.
- A coverage-policy waiver, exclusion, generated profile, or threshold change.
- Coverage work in packages other than `internal/journal`.
- Phase 17 maintenance, formatting, governance, or closure behavior.

## Acceptance Criteria

- [ ] New tests assert meaningful existing journal behavior and fail if that
  behavior regresses. (VAL-024)
- [ ] `internal/journal` reports at least 80.0% in the focused package run and
  in the unchanged full-suite coverage gate. (VAL-024)
- [ ] All production files and compatibility surfaces are byte-identical to
  the pre-task state. (VAL-024)
- [ ] The completion report cites the exact focused and full coverage outputs,
  including the `internal/journal` percentage. (VAL-024)

## First Failing Test

The red gate is QA-BLK-01's reproducible 78.9% package result against COV-001.
Select uncovered observable behavior, add assertions for it, and keep the
production implementation unchanged while the coverage gate moves to green.

## Documentation Updates

- None expected beyond `task-0-done.md`; the quality policy is unchanged.

## Exact Commands To Run

Run exactly:

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=/tmp/fid-journal-coverage.out ./internal/journal/...
go tool cover -func=/tmp/fid-journal-coverage.out
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Required-command rules live in `AGENTS.md`. Every command is reported `PASS`,
or `BLOCKED` with reason; the full gate is required before task close.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- This task runs before Task 1, but only after the QA plan authorizes the entry
  remediation and ESC-009 is resolved.
- A QA Reviewer entry review must accept the test quality and exact coverage
  evidence before Task 1 starts. QA Lead retains terminal VAL-024 ownership.
- Completion report: `project-docs/docs/phases/phase-17/tasks/task-0-done.md`.
