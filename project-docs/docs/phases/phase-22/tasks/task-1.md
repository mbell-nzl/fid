# Phase 22 - Task 1: Recovery schema and replay foundation

## Objective

Install the proposed recovery baseline and closed Event validation before domain writers depend on it. Keep initial capability availability false.

## Validation Contract Coverage

- VAL-001: Fresh init installs the recovery schema and enforces its identity constraints.
- VAL-002: Disposable full down/up restores the same accepted empty baseline.
- VAL-003: The same closed Event semantics govern live apply and replay.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- schema.md sections 1-3,5-6.
- Apply exactly the DDL in schema.md; no migration of old fid journals.
- Converge active schema reporting in rebuild/export and their fixtures now; do not hard-code final schema 16.
- Do not enable domain operations or mark recovery capabilities available.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `migrations/0013_p21_operational_kernel.up.sql`
- `migrations/0015_p22_recovery.up.sql`
- `migrations/0015_p22_recovery.down.sql`
- `internal/db/`
- `internal/journal/`
- `internal/export/export.go`
- `cmd/fid/rebuild.go`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-1-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.

### Out Of Scope

- Business operations, capability enablement, activation resource behavior, and escalation deletion.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Fresh init installs the recovery schema and enforces its identity constraints. (VAL-001)
- [ ] Disposable full down/up restores the same accepted empty baseline. (VAL-002)
- [ ] The same closed Event semantics govern live apply and replay. (VAL-003)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestRecoverySchemaFreshInitAndReplay rejects an unknown recovery Effect and proves the new tables are absent before implementation.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-1-done.md`: implementation,
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

- Start only after the Ready QA plan and Human approval.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 1.
- Commit using Conventional Commits with `Refs: phase-22/task-1`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

