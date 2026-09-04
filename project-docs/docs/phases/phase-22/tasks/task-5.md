# Phase 22 - Task 5: Atomic linked Run recovery

## Objective

Implement ordinary and authorized force recovery as one logical mutation. Preserve predecessor identity and create exactly one successor Claim/Run after reconciliation.

## Validation Contract Coverage

- VAL-017: Recovery creates a distinct linked successor and keeps the predecessor terminal.
- VAL-018: Recovery is atomic and idempotent across response loss and projection interruption.
- VAL-019: Unknown liveness always blocks ordinary and force recovery.
- VAL-020: Known-live takeover requires exact current Human force authority.
- VAL-021: Recovery preserves the Issue and immutable Envelope execution contract.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 4; schema.md run.recovered.
- Closed liveness is terminated/live/unknown: unknown always rejects, live requires exact Human force authority.
- Re-observe repository under the write lock and compare the supplied digest before any Effect.
- No nested write lock: the recovery orchestrator combines Claim and Run effects under one lock.
- Source, revision, scopes, Obligations, and budget policy cannot change during recovery.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/workflow/recovery.go`
- `internal/workflow/recovery_test.go`
- `internal/workflow/run.go`
- `internal/workflow/run_read.go`
- `internal/workflow/authority.go`
- `internal/workflow/activation.go`
- `internal/claims/`
- `internal/project/`
- `internal/journal/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-5-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-4-done.md` - actual predecessor hand-back.

### Out Of Scope

- Automatic process spawn, multi-Issue recovery, budget counter implementation, and activation recovery.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Recovery creates a distinct linked successor and keeps the predecessor terminal. (VAL-017)
- [ ] Recovery is atomic and idempotent across response loss and projection interruption. (VAL-018)
- [ ] Unknown liveness always blocks ordinary and force recovery. (VAL-019)
- [ ] Known-live takeover requires exact current Human force authority. (VAL-020)
- [ ] Recovery preserves the Issue and immutable Envelope execution contract. (VAL-021)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestRecoverRunCommitsExactlyOneSuccessorAfterLostResponse fails against the absent recovery operation.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-5-done.md`: implementation,
  red/green, every exact command, observed termination, assertion evidence,
  convergence provenance, dependencies, and deviations.
- Flag spec impact against the named contracts; only Tech Lead accepts a
  contract/schema change. Record no spec change when implementation conforms.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt/goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The integration gate is mandatory from task 1, including for retained fixtures.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
```

No exception to the full task gate. Commands whose expected pass is nonzero:
none. Focused red/green selections do not replace these commands. Report each
attempt using the command's own output token and observed termination, as
required by AGENTS.md. No interrupted/hung command is PASS; preserve all failed
attempts and rerun the byte-identical command after authorized causal repair.
Report lint is a Runner control after the report commit, not a Developer gate.

## Notes For Developer

- Start after task 4 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 2.
- Commit using Conventional Commits with `Refs: phase-22/task-5`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

