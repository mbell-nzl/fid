# Phase 21 - Task 17: Trusted Evidence Evaluation

## Objective

Link only exact fid-executed Receipts and evaluate validity, repository-tree freshness, required Artifact availability, deterministic all-reasons order, and no Attention.

## Validation Contract Coverage

- VAL-058: owning observable in the Phase 21 validation contract.
- VAL-059: owning observable in the Phase 21 validation contract.
- VAL-060: owning observable in the Phase 21 validation contract.
- VAL-061: owning observable in the Phase 21 validation contract.
- VAL-062: owning observable in the Phase 21 validation contract.
- VAL-063: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Caller-authored pass, command text, import, attestation, and role-looking records cannot qualify.
- Validity, freshness, and availability remain independent.
- Return every reason; never short-circuit or raise Attention.

## Scope

### Files To Create Or Modify

- `internal/evidence/`
- `internal/validation/`
- `internal/workflow/`
- `internal/cli/`
- `cmd/fid/`
- `project-docs/docs/phases/phase-21/tasks/task-17-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§9,10`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/specs/01_PRD.md PRD-EVID-001..006`

### Out Of Scope

- Terminal Issue transition and Envelope completion.

## Acceptance Criteria

- [ ] Only exact Receipt binding links. (VAL-058)
- [ ] Validity is independent. (VAL-059)
- [ ] Tree freshness is exact. (VAL-060)
- [ ] Artifact dimension is independent. (VAL-061)
- [ ] All reasons are stable/complete. (VAL-062)
- [ ] No Attention state/effect exists. (VAL-063)

## First Failing Test

Add the binding matrix, three paired dimension tests, multi-reason golden, and Attention absence proof.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-17-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/evidence ./internal/validation ./internal/workflow ./internal/cli ./cmd/fid/...
rg -n -i --glob '*.go' --glob '!**/*_test.go' '\battention\b|attention_items' cmd/fid internal
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Expected Non-Zero Commands

- `rg -n -i --glob '*.go' --glob '!**/*_test.go' '\battention\b|attention_items' cmd/fid internal` - expected `exit_code: 1` - no Phase 21 production Attention state or effect

## Notes For Developer

- Start after tasks 15-16 close.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-17-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
