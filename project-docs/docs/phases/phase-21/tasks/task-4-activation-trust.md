# Phase 21 - Task 4: Activation Trust Lifecycle

## Objective

Implement prepare/component/commit/abort/finalize/status for one Activation and resolve Actor/Gate components only from finalized values. Reject caller-self-declared identity, capabilities, or provenance.

## Validation Contract Coverage

- VAL-011: owning observable in the Phase 21 validation contract.
- VAL-012: owning observable in the Phase 21 validation contract.
- VAL-013: owning observable in the Phase 21 validation contract.
- VAL-014: owning observable in the Phase 21 validation contract.
- VAL-015: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Only the accepted lifecycle and pre-commit abort exist.
- Commit is the point of no return; recovery and respawn remain unavailable.
- Bootstrap attribution is recorded but never grants protected authority.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/policy/`
- `internal/cli/`
- `cmd/fid/` (including the `cmd/fid/integration/` package, DEC-169)
- `project-docs/docs/phases/phase-21/tasks/task-4-done.md`

### Checkpoint Repair Obligation (DEC-169, MF-1)

No lifecycle precondition may preempt an exact same-key, same-digest idempotent
replay anywhere in the Activation family: `prepare`, `component`, and the
`commit`/`abort`/`finalize` transitions. A replay returns the original Event,
entity, and domain fields with only `replayed:true` and appends no Event.
New-key lifecycle conflicts and same-key/different-digest conflicts still fail
with byte-identical envelopes. Every error path emits its typed failure envelope
on the accepted channel; a silent nonzero exit is a defect. The red test records
both streams verbatim for the pre-repair second call. No Event schema,
persistence, public-interface, result-shape, or command-grammar change.

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §4`
- `project-docs/docs/phases/phase-21/public-interface.md §§2-4`
- `project-docs/specs/02_TRD.md §3.1`

### Out Of Scope

- Gate Registration and recovery.

## Acceptance Criteria

- [ ] Happy lifecycle and abort states are exact. (VAL-011, VAL-012)
- [ ] Recovery/respawn remain unavailable. (VAL-013)
- [ ] Only finalized profiles resolve trusted Actor values; self-declaration rejects. (VAL-014, VAL-015)

## First Failing Test

Add failing lifecycle and caller-augmentation tables before state mutations.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-4-done.md` - record red/green evidence, deviations, required commands, and spec impact.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/policy ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"recover|newActivationRecover|activation\.recover' cmd/fid internal
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Expected Non-Zero Commands

- `rg -n --glob '*.go' --glob '!**/*_test.go' 'Use:\s*"recover|newActivationRecover|activation\.recover' cmd/fid internal` - expected `exit_code: 1` - no activation recovery or automatic-respawn command/path

## Notes For Developer

- Start after tasks 2-3 close.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-4-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
