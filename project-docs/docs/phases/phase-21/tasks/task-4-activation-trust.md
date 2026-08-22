# Phase 21 - Task 4: Activation Lifecycle And Trust Profiles

## Objective

Implement the bounded Activation lifecycle, component staging, finalized trust source, and Actor profile resolution. Establish that caller-declared roles/capabilities/provenance have no authority and that recovery/respawn remain absent.

## Validation Contract Coverage

- VAL-009: Activation follows the accepted lifecycle and blocks normal work while incomplete.
- VAL-010: kit Actor authority resolves only from the active Activation.
- VAL-011: recovery, continuation, and respawn behavior are absent.

## Binding Constraints

- Commit is the point of no return; abort after commit returns `activation_recovery_unavailable`.
- Only authority-profile and gate-definitions components exist.
- Do not advertise `kit.activation.v1` as available.

## Scope

### Files To Create Or Modify

- `internal/workflow/` - Activation aggregate files in the work-domain owner.
- `internal/cli/`
- `cmd/fid/`
- `internal/journal/`
- `internal/db/`
- `project-docs/docs/phases/phase-21/tasks/task-4-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §4`
- `project-docs/specs/01_PRD.md PRD-OPS-004..006 and PRD-AUTH-001..007`

### Out Of Scope

- Post-commit recovery, rollback, automatic respawn, and protected Human authority.

## Acceptance Criteria

- [ ] Lifecycle and one-active/nonterminal invariants pass with atomic Events. (VAL-009)
- [ ] Resolved Actor data is activation-authored; forged caller fields fail closed. (VAL-010)
- [ ] No recovery/respawn leaf or side effect exists. (VAL-011)

## First Failing Test

Add a failing lifecycle table test covering prepare through finalize and every illegal transition.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-4-done.md` - record red/green evidence, deviations, required command results, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/cli ./cmd/fid/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
project-docs/tools/report-lint.sh project-docs/docs/phases/phase-21/tasks/task-4-done.md
```

Every command above is required and must terminate with captured evidence. The report-lint command runs in this task before the next task is dispatched.

## Notes For Developer

- Start only after tasks 1-3 and QA Review 1 closes with required evidence.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-4-done.md`.
- Stop for protected contract changes; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
