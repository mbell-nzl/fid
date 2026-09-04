# Phase 22 - Task 2: Typed Stops and Attention domain

## Objective

Implement immutable Stop facts, explicit resolution, and deduplicated Attention with atomic Event effects. Give workflow recovery branches direct package-level evidence.

## Validation Contract Coverage

- VAL-004: A Stop's original fact remains immutable after resolution.
- VAL-005: Equivalent Attention requests share one open identity without losing distinct Stops.
- VAL-006: Attention resolution requires qualifying facts for every linked Stop.
- VAL-007: Routine status and warnings do not create Attention.
- VAL-079: Reopened workflow recovery branches have direct package tests and measured coverage.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md section 5; public-interface.md sections 2,4.
- CF-034: capture workflow package coverage before edits and after completion; add direct in-package positive and rejection branch tests. No numeric threshold.
- Resolving Attention never deletes Stops or accepts an unqualified reference.
- Extend the activated capability allowlist only for accepted Stop/Attention actions.
- DEC-184/public-interface section 2.1 enumerates the exact actions and the
  activation-subject profile-only case; do not invent an Envelope for activation.
- Attention input cites disposition-containing Events for every linked Stop,
  not underlying facts. Missing/extra/duplicate/wrong-subject sets reject;
  resolving Attention never inserts a Stop disposition.
- Rejected recovery errors do not produce Stops. Exercise explicit Stop
  production separately; automatic budget/activation producers have later owners.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `internal/workflow/stop.go`
- `internal/workflow/attention.go`
- `internal/workflow/stop_test.go`
- `internal/workflow/attention_test.go`
- `internal/workflow/authority.go`
- `internal/workflow/activation.go`
- `internal/journal/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-2-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-1-done.md` - actual predecessor hand-back.

### Out Of Scope

- Public CLI wiring, budget producers, Human Outcomes, and Review semantics.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] A Stop's original fact remains immutable after resolution. (VAL-004)
- [ ] Equivalent Attention requests share one open identity without losing distinct Stops. (VAL-005)
- [ ] Attention resolution requires the exact disposition Event set for every
  linked Stop and revalidates their qualifying facts; invalid sets have no
  effects. (VAL-006)
- [ ] Routine status and warnings do not create Attention. (VAL-007)
- [ ] Reopened workflow recovery branches have direct package tests and measured coverage. (VAL-079)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestEquivalentStopsShareAttentionButKeepDistinctFacts fails before the new domain operation exists.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-2-done.md`: implementation,
  red/green, every exact command, observed termination, assertion evidence,
  convergence provenance, dependencies, and deviations.
- Flag spec impact against the named contracts; only Tech Lead accepts a
  contract/schema change. Record no spec change when implementation conforms.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt/goimports to edited Go files before the final gate.

Before any test or production edit, record this baseline command once:

```text
go test -race -cover ./internal/workflow
```

After implementation, every command below is required before task close.
The integration gate is mandatory from task 1, including for retained fixtures.

```text
go test -race -cover ./internal/workflow
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

- Start after task 1 closes.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 1.
- Commit using Conventional Commits with `Refs: phase-22/task-2`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
