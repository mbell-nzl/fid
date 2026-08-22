# Phase 21 - Task 16: Safe Content-Addressed Artifacts

## Objective

Collect only declared safe regular-file Artifacts after gate termination, sanitize and atomically publish them by SHA-256 identity, and implement explicit availability verification without erasing provenance. Fail every unsafe/missing/racy candidate closed.

## Validation Contract Coverage

- VAL-040: safe Artifacts publish atomically by sanitized content identity.
- VAL-041: unsafe candidates reject and cannot satisfy material policy.
- VAL-042: availability changes preserve Artifact/Receipt/Evidence identity.

## Binding Constraints

- Revalidate containment and file type after termination and while copying.
- Availability is only `available` or `missing`; pruning/expiry automation is deferred.
- `artifact verify` may append an availability Event only when observed state changes.

## Scope

### Files To Create Or Modify

- `internal/evidence/` - Artifact storage and executor collection files.
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-16-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §9`
- `project-docs/specs/04_TRD_CONFIG.md §4`

### Out Of Scope

- Retention automation, external Artifact stores, and Evidence evaluation.

## Acceptance Criteria

- [ ] Declared safe files publish atomically with exact digest/metadata. (VAL-040)
- [ ] Every unsafe path/type/race/size case rejects and stays unsatisfied. (VAL-041)
- [ ] Verify loss/restoration changes availability only. (VAL-042)

## First Failing Test

Add a failing Artifact table test for safe content followed by traversal, symlink, special-file, size, and changed-during-copy cases.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-16-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-16-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after task 15 and QA Review 5 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-16-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
