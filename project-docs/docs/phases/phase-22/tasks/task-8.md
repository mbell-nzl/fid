# Phase 22 - Task 8: Interrupted Attempts and derived Run closure

## Objective

Persist only observed terminal gate facts, preserve abandoned Attempts, and derive Run closure from trusted evidence. Add intentional abandonment without creating success proof.

## Validation Contract Coverage

- VAL-031: Observed unsuccessful termination cannot produce a passing Receipt.
- VAL-032: An abandoned Attempt cannot be replayed into an automatic second spawn.
- VAL-033: Run success is derived from qualifying trusted facts.
- VAL-034: Intentional abandonment is distinct from interruption and proof.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 7; public-interface.md run.close/run.abandon.
- Cancellation, timeout, signal, spawn failure, and harness failure are distinct non-success observations.
- A lost executor response cannot manufacture a Receipt or imply pass.
- Run success requires shared evaluator success and qualifying executor/Claim authority; abandonment requires terminal proof or exact force authority.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/evidence/execution.go`
- `internal/evidence/executor.go`
- `internal/evidence/live.go`
- `internal/evidence/execution_test.go`
- `internal/evidence/executor_test.go`
- `internal/workflow/run.go`
- `internal/workflow/run_close.go`
- `internal/workflow/run_close_test.go`
- `internal/validation/`
- `internal/journal/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-7-done.md` - actual predecessor hand-back.

### Out Of Scope

- Caller-reported terminal outcomes, external attested Receipts, process supervision, and removal of the CLI finish leaf before task 9.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Observed unsuccessful termination cannot produce a passing Receipt. (VAL-031)
- [ ] An abandoned Attempt cannot be replayed into an automatic second spawn. (VAL-032)
- [ ] Run success is derived from qualifying trusted facts. (VAL-033)
- [ ] Intentional abandonment is distinct from interruption and proof. (VAL-034)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestKilledFidLeavesAttemptWithoutReceiptAndCannotReplaySpawn fails when an in-flight gate key can be rerun.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-8-done.md`: implementation,
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

- Start after task 7 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 3.
- Commit using Conventional Commits with `Refs: phase-22/task-8`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

