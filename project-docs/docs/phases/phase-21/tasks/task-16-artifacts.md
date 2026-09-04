# Phase 21 - Task 16: Safe Referenced Artifacts

## Objective

Retain declared safe regular files atomically by sanitized content digest, reject every unsafe candidate, and observe availability changes without identity loss.

## Validation Contract Coverage

- VAL-055: owning observable in the Phase 21 validation contract.
- VAL-056: owning observable in the Phase 21 validation contract.
- VAL-057: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Resolve and revalidate containment/symlinks before and during copy.
- Rejected candidates never satisfy required material.
- Availability changes are observation Events, not identity mutation.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/config/`
- `internal/journal/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-16-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §9`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/specs/04_TRD_CONFIG.md`

### Out Of Scope

- Evidence qualification and pruning.

## Acceptance Criteria

- [ ] Safe Artifact bytes/digest/metadata publish atomically. (VAL-055)
- [ ] Every unsafe candidate rejects and cannot satisfy. (VAL-056)
- [ ] Loss/restoration toggles availability only. (VAL-057)

## First Failing Test

Add failing safe-content, adversarial candidate, and loss/restoration fixtures.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-16-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/config ./internal/journal ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after task 15 closes.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-16-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
