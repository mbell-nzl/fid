# Phase 21 - Task 13: Portable Executor Policy

## Objective

Implement and validate the standard-library executor policy boundary without spawning gates yet. Resolve immutable Registration policy into structured argv, exact cwd, a closed environment, configured limits, credential names, and safe paths while rejecting every unsupported restriction.

## Validation Contract Coverage

- VAL-028: the independent Group 1 checkpoint passes before Group 2 dispatch.
- VAL-029: supported structured-argv/no-shell policy resolves exactly.
- VAL-030: unsupported containment/network/write/limit/path values reject.
- VAL-031: the environment is closed and credential values are not persisted.
- VAL-032: cwd and Artifact paths remain within approved roots.

## Binding Constraints

- Task 13 may start only after the recorded Group 1 checkpoint PASS.
- Accept only `sandbox:none`, `network:allowed`, and `write:allowed`; make no containment claim.
- A dependency or platform-specific isolation need is a Human decision, not an implementation adaptation.

## Scope

### Files To Create Or Modify

- `internal/evidence/` - Registration policy and executor-port implementation.
- `internal/project/`
- `internal/cli/`
- `project-docs/docs/phases/phase-21/tasks/task-13-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§5,9`
- `project-docs/docs/phases/phase-21/qa/group-1-checkpoint-review.md`
- `project-docs/specs/04_TRD_CONFIG.md`

### Out Of Scope

- Process spawn, live output, Receipts, and OS containment.

## Acceptance Criteria

- [ ] The independent Group 1 review and fresh feasibility checkpoint is
  recorded PASS before this task starts. (VAL-028)
- [ ] Exact accepted policy resolves to the intended argv/cwd/limits. (VAL-029)
- [ ] Every stronger/unsafe value rejects at Registration with the exact reason. (VAL-030, VAL-032)
- [ ] Closed environment assembly stores names but never credential values. (VAL-031)

## First Failing Test

Add a failing table-driven policy resolver test with one accepted case and one row per unsupported value/path class.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-13-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/project
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-13-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after the recorded VAL-028 Group 1 checkpoint PASS closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-13-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
