# Phase 22 - Task 11: Activation CLI and config boundary

## Objective

Expose activation resource observation and recovery with exact errors and crash-point behavior. Pin whole-config rejection to operations that consume configuration.

## Validation Contract Coverage

- VAL-045: Activation recovery commands expose closed results and typed direction errors.
- VAL-046: Activation recovers at every durable boundary through the built binary.
- VAL-047: Whole-config validation occurs before a consuming operation uses values.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- public-interface.md activation leaves; contracts.md section 12.
- CF-035 uses consuming-path validation; help and archive validation must remain usable with malformed project config.
- Remove activation-recovery-unavailable only when direction/agreement errors are wired.
- Extend the shared case runner/fake with activation recovery and run the same
  semantic expectations against fake and fid before task close.
- Converge existing activation request/resource fixtures in the same task that changes their adapter.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `cmd/fid/activation.go`
- `cmd/fid/activation_test.go`
- `internal/config/project_config.go`
- `internal/config/project_config_test.go`
- `internal/cli/`
- `cmd/fid/integration/p22_activation_recovery_test.go`
- `cmd/fid/integration/testdata/p22_conformance/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-11-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-10-done.md` - actual predecessor hand-back.

### Out Of Scope

- Global config parsing on every command, repair of project files, and promotion of bootstrap attribution to protected Actor authority.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Activation recovery commands expose closed results and typed direction errors. (VAL-045)
- [ ] Activation recovers at every durable boundary through the built binary. (VAL-046)
- [ ] Whole-config validation occurs before a consuming operation uses values. (VAL-047)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestActivationRecoveryAtEveryDurableBoundary uses the real binary to reject rollback after commit.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-11-done.md`: implementation,
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

- Start after task 10 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 4.
- Commit using Conventional Commits with `Refs: phase-22/task-11`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
