# Phase 21 - Task 13: Portable Executor Policy

## Objective

Implement the standard-library process adapter for the exact enforceable policy and reject every unsupported restriction at Registration. Establish no OS containment claim.

## Validation Contract Coverage

- VAL-044: owning observable in the Phase 21 validation contract.
- VAL-045: owning observable in the Phase 21 validation contract.
- VAL-046: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Structured argv/no shell, exact cwd/env, limits, credentials, redaction inputs, and Artifact declarations only.
- Accept only `sandbox:none`, `network:allowed`, and `write:allowed`.
- No runtime dependency addition.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/config/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-13-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §5`
- `project-docs/docs/phases/phase-21/public-interface.md §§2,7`
- `project-docs/specs/04_TRD_CONFIG.md`

### Out Of Scope

- Live framing, Receipt persistence, OS isolation, and platform-specific sandboxing.

## Acceptance Criteria

- [ ] Accepted policy is observably enforced. (VAL-044)
- [ ] Every stronger/invalid value rejects without persistence. (VAL-045)
- [ ] No containment claim exists. (VAL-046)

## First Failing Test

Add failing positive policy observation and one-field unsupported policy tables.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-13-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/config ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-13-done.md
```

Every command is required and must terminate with session-captured evidence. The report-lint command runs before the next dispatch. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start only after Group 1 checkpoint PASS (VAL-042/043 and feasibility probes).
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-13-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
