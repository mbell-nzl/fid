# Phase 22 - Task 3: Lease lifecycle and executor observations

## Objective

Complete explicit expiry, current-version lease checks, and trusted Run executor observations. Prevent a new Claim from bypassing an unresolved prior executor.

## Validation Contract Coverage

- VAL-008: Lease validity changes exactly at expiry and cannot be revived by renewal.
- VAL-009: Concurrent same-Issue mutation acquisition has at most one winner.
- VAL-010: Only a trusted bound observer can establish executor liveness.
- VAL-011: Claim renewal preserves Run lineage while fencing obsolete lease versions.
- VAL-012: A Run binds exact executor, optional Obligation, and originating authority context.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- DEC-184/public-interface section 2.1: add exactly claim.expire and run.observe.
  Executor IDs are project-global immutable identities; conflicting new start
  or imported identity reuse returns the declared duplicate-entity shape with
  no effects. Same-key replay returns the original result first. (VAL-012)
- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 2; schema.md section 5.
- Unknown liveness is not termination; expiry equality is expired.
- Run start adds executor/optional Obligation binding immediately; update current command parsing and every affected ordinary/tagged fixture in this task.
- Renewal preserves immutable Run binding and authorizes only the current Claim version.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/claims/`
- `internal/workflow/run.go`
- `internal/workflow/run_read.go`
- `internal/workflow/run_observation.go`
- `internal/workflow/run_observation_test.go`
- `internal/workflow/activation.go`
- `cmd/fid/run.go`
- `cmd/fid/claim.go`
- `internal/journal/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-3-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-2-done.md` - actual predecessor hand-back.

### Out Of Scope

- Repository reconciliation, recovery takeover, and new public observe/expire leaves (task 6).
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Lease validity changes exactly at expiry and cannot be revived by renewal. (VAL-008)
- [ ] Concurrent same-Issue mutation acquisition has at most one winner. (VAL-009)
- [ ] Only a trusted bound observer can establish executor liveness. (VAL-010)
- [ ] Claim renewal preserves Run lineage while fencing obsolete lease versions. (VAL-011)
- [ ] A Run binds exact executor, optional Obligation, and originating authority context. (VAL-012)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestExpiredClaimWithUnresolvedExecutorCannotBeReacquired fails against the current automatic-expiry acquisition path.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-3-done.md`: implementation,
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

- Start after task 2 closes.
- One task per Developer session; recovery may continue the same task.
- QA Review window 1 follows this task; wait for its required disposition.
- Commit using Conventional Commits with `Refs: phase-22/task-3`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
