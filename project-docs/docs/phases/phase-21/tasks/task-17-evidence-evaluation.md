# Phase 21 - Task 17: Receipt Evidence And Deterministic Evaluation

## Objective

Link only exact fid-executed Receipts to gate-evidence Obligations and evaluate absence, validity, repository-tree freshness, and required Artifact availability independently. Return all reasons in the fixed order and never raise Attention.

## Validation Contract Coverage

- VAL-043: only an exact trusted Receipt can link as Evidence.
- VAL-044: evidence dimensions serialize and evaluate independently.
- VAL-045: repository-tree target changes stale bound Evidence.
- VAL-046: evaluation returns every reason deterministically.
- VAL-047: unsatisfied routine Obligations never raise Attention.

## Binding Constraints

- No caller pass flag, attestation, Review, Human Outcome, or external evidence type.
- Use repository object identity, not branch names or timestamps, for freshness.
- The same evaluator must serve guarded transitions and completion.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/validation/`
- `internal/gitctx/`
- `internal/workflow/`
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-17-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§9,10`
- `project-docs/specs/01_PRD.md PRD-EVID-001..006`

### Out Of Scope

- Completion, resume, Reviews/Human Outcomes, Attention, and external freshness.

## Acceptance Criteria

- [ ] Only exact fid-executed binding links. (VAL-043)
- [ ] Four dimensions remain independent and tree changes stale correctly. (VAL-044, VAL-045)
- [ ] All reasons are complete/stable and no Attention state appears. (VAL-046, VAL-047)

## First Failing Test

Add a failing evidence binding matrix, then a four-dimension multi-reason golden evaluation test.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-17-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/validation ./internal/gitctx ./internal/workflow ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-17-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 15-16 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-17-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
