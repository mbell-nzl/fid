# Phase 22 - Task 10: Activation recovery domain

## Objective

Persist old/new resource agreement, choose recovery direction at the commit boundary, and recover under the same activation identity. Prevent mixed trust during installation and recovery.

## Validation Contract Coverage

- VAL-040: Activation prepare records the exact old and new resource sets.
- VAL-041: Activation recovery direction respects the point of no return.
- VAL-042: Activation cannot finalize until every required resource agrees.
- VAL-043: Incomplete activation blocks normal work and preserves recovery inspection.
- VAL-044: Activation recovery retries preserve one durable direction and identity.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 8; schema.md activation resources/recovery.
- DEC-184: version mismatch uses activation-version-conflict; competing prepare
  retains its separate activation-conflict shape. Recovery leaves retain
  bootstrap attribution and add no activated action or Envelope grant.
- DEC-187/kit-contract-mapping.md: fid persists component agreement; the kit
  coordinator owns durable recovery-export binding. No fid archive field or
  claim that portable availability proves that external binding is authorized.
- Five resource identities plus authority/Gate content are immutable prepare inputs; observations use expected versions.
- Kit materializes files; fid records and validates agreement without editing provider/project files.
- Converge the existing activation.prepare adapter input and all resource fixtures now; task 11 adds the new recovery leaves, not delayed parser compatibility.
- Prepare with a predecessor, ordinary finalize, abort, and recovery all preserve exactly one active trust selection.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/workflow/activation.go`
- `internal/workflow/activation_recovery.go`
- `internal/workflow/activation_recovery_test.go`
- `internal/workflow/authority.go`
- `internal/evidence/registration.go`
- `internal/journal/`
- `cmd/fid/activation.go`
- `cmd/fid/response_errors.go` - exact version-conflict mapping for amended leaves.
- `internal/cli/` - register activation-version-conflict with its closed shape.
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-10-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-9-done.md` - actual predecessor hand-back.

### Out Of Scope

- Provider generation, rollback after commit, remote authentication, and dependency changes.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Activation prepare records the exact old and new resource sets. (VAL-040)
- [ ] Activation recovery direction respects the point of no return. (VAL-041)
- [ ] Activation cannot finalize until every required resource agrees. (VAL-042)
- [ ] Incomplete activation blocks normal work and preserves recovery inspection. (VAL-043)
- [ ] Activation recovery retries preserve one durable direction and identity. (VAL-044)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestCommittedActivationCanOnlyRecoverForward fails against the walking-skeleton recovery-unavailable path.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-10-done.md`: implementation,
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

- Start after task 9 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 4.
- Commit using Conventional Commits with `Refs: phase-22/task-10`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
