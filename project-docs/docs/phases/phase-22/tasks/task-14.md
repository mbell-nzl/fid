# Phase 22 - Task 14: Atomic import into an empty target

## Objective

Validate and stage a recovery archive, reconstruct retained state, and atomically publish a target project. Preserve identity/history and support a retry after a lost publication response.

## Validation Contract Coverage

- VAL-058: Import preserves stable entity and Event identity/history.
- VAL-059: Import publishes all reconstructed state or no target project.
- VAL-060: Import response-loss retry does not duplicate facts.
- VAL-061: Import cannot overwrite or merge a populated target.
- VAL-062: Imported active work does not automatically regain mutation authority.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- APPROVED under DEC-188; execute after the Runner's required QA readiness check.
- contracts.md section 10; public-interface.md export.import.
- Use target-specific exclusive locking and same-filesystem publication; do not replace any existing .fid project.
- Preserve original Event bytes/order and idempotency results; publication adds no domain Event.
- Active imported work remains historical until explicit liveness/reconciliation; no automatic authorization or spawn.
- DEC-202 allocates the import-specific destination reconciliation route. The
  target may contain Git metadata and project source files; empty means no
  existing `.fid` project. Preserve those files and all imported history.
- Tests inject faults before/after durability and publication boundaries with fresh scratch targets.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/export/import.go`
- `internal/export/import_test.go`
- `internal/export/`
- `internal/project/`
- `internal/config/`
- `internal/journal/`
- `internal/db/`
- `internal/workflow/authority.go`
- `internal/workflow/reconciliation.go`, `internal/workflow/run_inspect.go`,
  `internal/workflow/recovery.go`, `internal/workflow/run_observation.go`,
  `internal/workflow/run.go`, `internal/workflow/evaluation.go` - DEC-202:
  qualify new destination facts through the installed import binding only;
  preserve ordinary foreign-project/worktree rejection and historical values.
- `internal/claims/generic.go`
- `internal/evidence/execution.go`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-14-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-13-done.md` - actual predecessor hand-back.

### Out Of Scope

- Populated-project merge, distributed synchronization, identity remapping, old-data migration, and source-code restoration.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Import preserves stable entity and Event identity/history. (VAL-058)
- [ ] Import publishes all reconstructed state or no target project. (VAL-059)
- [ ] Import response-loss retry does not duplicate facts. (VAL-060)
- [ ] Import cannot overwrite or merge a populated target. (VAL-061)
- [ ] Imported active work does not automatically regain mutation authority. (VAL-062)
- [ ] A target containing Git/source files but no `.fid` imports without
  changing those files. Fresh local liveness and exact destination
  reconciliation can produce a successor Claim/Run bound to that destination;
  new successor execution qualifies while imported proof/history stays
  immutable. Unqualified imported work and ordinary foreign-worktree requests
  remain rejected. (VAL-058/061/062; DEC-202)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestImportCrashBeforePublicationLeavesNoPartialProject fails before staged publication exists.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-14-done.md`: implementation,
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

- Start after task 13 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 5.
- Commit using Conventional Commits with `Refs: phase-22/task-14`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
