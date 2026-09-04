# Phase 22 - Task 18: Final recovery conformance and Human validation fixture

## Objective

Complete the shared recovery case corpus and independently exercisable built-binary journey. Provide reproducible Human recovery, activation, and import checkpoints for final QA.

## Validation Contract Coverage

- VAL-053: Human and JSON recovery views use the same semantic projection.
- VAL-073: One shared semantic case set passes against fake and fid recovery adapters.
- VAL-074: The full one-Issue failure/recovery journey completes through public operations.
- VAL-075: Every amended wire variant has a non-vacuous deterministic conformance fixture.
- VAL-076: A Human can choose the next recovery action from a fresh resume view.
- VAL-077: A Human can distinguish restore-prior from complete-forward recovery.
- VAL-078: A Human can recognize restored work and its remaining intervention need.
- VAL-080: Independent review finds no unresolved cross-product recovery contract mismatch.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- PROPOSED: execution requires Human acceptance of DEC-180 and a Ready QA plan.
- contracts.md sections 1-12; public-interface.md sections 1-6.
- The same cases execute against a test-only semantic fake and the fid JSON adapter; the fake never opens fid SQLite/journal/.fid files.
- Pin kit backend-port.v1 at c5229617220f97ff9228df261e24718b1d1d6d89 for semantic comparison; any mismatch requires explicit TL disposition.
- No new domain behavior is assigned to this final assembly task. Return causal findings for minimum-scope amendment per DEC-175.
- Developer produces runnable Human fixtures; only Human may judge VAL-076..078.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.

## Scope

### Files To Create Or Modify

- `cmd/fid/integration/p22_conformance_test.go`
- `cmd/fid/integration/p22_recovery_journey_test.go`
- `cmd/fid/integration/p22_activation_recovery_test.go`
- `cmd/fid/integration/p22_import_recovery_test.go`
- `cmd/fid/integration/p22_fake_test.go`
- `cmd/fid/integration/testdata/p22_conformance/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-18-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-17-done.md` - actual predecessor hand-back.

### Out Of Scope

- Held-out evaluator material, production kit adapter code, generic fake fallback, and unapproved contract changes.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] Human and JSON recovery views use the same semantic projection. (VAL-053)
- [ ] One shared semantic case set passes against fake and fid recovery adapters. (VAL-073)
- [ ] The full one-Issue failure/recovery journey completes through public operations. (VAL-074)
- [ ] Every amended wire variant has a non-vacuous deterministic conformance fixture. (VAL-075)
- [ ] A Human can choose the next recovery action from a fresh resume view. (VAL-076)
- [ ] A Human can distinguish restore-prior from complete-forward recovery. (VAL-077)
- [ ] A Human can recognize restored work and its remaining intervention need. (VAL-078)
- [ ] Independent review finds no unresolved cross-product recovery contract mismatch. (VAL-080)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.

## First Failing Test

TestP22ConformanceCorpusHasEveryAcceptedOperationAndFailure detects the first missing amended leaf without vacuous success.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-18-done.md`: implementation,
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
go build -o /tmp/fid ./cmd/fid
```

No exception to the full task gate. Commands whose expected pass is nonzero:
none. Focused red/green selections do not replace these commands. Report each
attempt using the command's own output token and observed termination, as
required by AGENTS.md. No interrupted/hung command is PASS; preserve all failed
attempts and rerun the byte-identical command after authorized causal repair.
Report lint is a Runner control after the report commit, not a Developer gate.

## Notes For Developer

- Start after task 17 closes.
- One task per Developer session; recovery may continue the same task.
- QA Review window 6 follows this task; wait for its required disposition.
- Commit using Conventional Commits with `Refs: phase-22/task-18`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

