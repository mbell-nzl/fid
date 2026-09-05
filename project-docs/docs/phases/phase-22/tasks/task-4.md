# Phase 22 - Task 4: Read-only repository reconciliation

## Objective

Capture and compare repository state against Run start or a durable checkpoint. Produce exact adoption candidates without changing repository content.

## Validation Contract Coverage

- VAL-013: Clean unchanged repository state qualifies for ordinary reconciliation.
- VAL-014: Dirty in-scope work requires exact recorded adoption.
- VAL-015: Unsafe or ambiguous repository state blocks recovery.
- VAL-016: Recovery never changes source files or Git state.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- APPROVED under DEC-188; execute after the Runner's required QA readiness check.
- contracts.md section 3; public-interface.md Reconciliation.
- DEC-185/contracts sections 3.1/3.2 are the exact dirty/reconciliation digest
  algorithms. Add independently pinned canonical-byte/hash vectors for staged,
  unstaged, mode-only, deleted, renamed and untracked paths, symlink payloads,
  exclusions, stable repeated capture, and rejected unsupported state.
- Add exactly run.checkpoint to the activated action set in this task.
- Parse structured NUL-delimited Git output; bound path count and input size.
- Checkpoint is an observed fact, never a Git commit; bind version and digest.
- Snapshot refs/index/stash/worktree before and after every positive/negative scenario.
- Source-file comparison excludes fid's own canonical journal/projection writes;
  Git metadata and protected .fid internals cannot be adopted by force.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/project/repository.go`
- `internal/project/repository_test.go`
- `internal/workflow/reconciliation.go`
- `internal/workflow/reconciliation_test.go`
- `internal/workflow/run_checkpoint.go`
- `internal/workflow/run_checkpoint_test.go`
- `internal/workflow/activation.go` - only the run.checkpoint action addition.
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-4-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-3-done.md` - actual predecessor hand-back.

### Out Of Scope

- Git mutation, force takeover, caller-supplied repository truth, and public CLI wiring.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Clean unchanged repository state qualifies for ordinary reconciliation. (VAL-013)
- [ ] Dirty in-scope work requires exact recorded adoption, with canonical
  digest vectors independent of the observer and invariant under branch rename
  and operational .fid writes. (VAL-014)
- [ ] Unsafe or ambiguous repository state blocks recovery. (VAL-015)
- [ ] Recovery never changes source files or Git state. (VAL-016)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestRecoveryObservationDetectsChangedUntrackedFileWithNULDelimitedPaths fails before observation is implemented.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-4-done.md`: implementation,
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

- Start after task 3 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 2.
- Commit using Conventional Commits with `Refs: phase-22/task-4`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
