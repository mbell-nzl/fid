# Phase 22 - Task 6: Recovery CLI and early process journey

## Objective

Expose inspect, checkpoint, observation, expiry, and recover through canonical JSON. Prove an interrupted one-Issue journey with the built binary before later features extend it.

## Validation Contract Coverage

- VAL-022: A built-binary recovery journey works by task 6.
- VAL-023: Recovery leaves enforce their closed wire schemas.
- VAL-024: Recovery failures have deterministic typed process exits.

Full sources, pass/fail criteria, and validator ownership are in
`project-docs/docs/phases/phase-22/qa/validation-contract.md`.

## Binding Constraints

- DEC-184: retain --project; cover required observation IDs, duplicate executor
  details, and rejection-without-Stop behavior in the shared corpus. (VAL-023/024)
- APPROVED under DEC-188; execute after the Runner's required QA readiness check.
- public-interface.md run/claim recovery leaves; contracts.md sections 2-4.
- Run a real subprocess, never replace this journey with an in-process command test.
- Introduce the shared semantic case runner and test-only fake for implemented
  recovery cases here; run the same expected outcomes against fake and fid.
- The early journey proves only delivered recovery with null/unavailable budget counters until task 7; no stubs or fake success for future behavior.
- Cover every new leaf's positive result and authority, invalid input, stale version, and liveness error exits.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.
- DEC-194 amends this scope with the four `internal/workflow` domain files
  below. They are task 4's accepted but undelivered scope; `run.checkpoint` and
  `run.inspect` have no domain producer or read API, so no CLI leaf here can
  conform without them. Implement them to the already accepted
  `contracts.md` section 3/3.2, `public-interface.md` section 3, and
  `schema.md` section 5 rows. Change no wire shape, error code, exit code, DDL,
  Effect, or Event semantic. Record task-4 provenance for each file; do not
  relabel the accepted owner. Task 4 stays COMPLETE and its report is not edited.

## Scope

### Files To Create Or Modify

- `internal/workflow/run_checkpoint.go` - ADDED under DEC-194; the
  `run.checkpoint` mutation and `run.checkpointed` Event producer, structured
  exactly like `ObserveRun` in `internal/workflow/run_observation.go` and
  emitting only the `schema.md` section 5 Effects (insert one
  `run_observations` row with `kind='checkpoint'`; increment
  `run_runtime.version`). Re-observe under the write lock and reject a
  mismatched `reconciliation_digest` with `recovery-state-conflict`.
- `internal/workflow/run_checkpoint_test.go` - ADDED under DEC-194; positive,
  stale-version, inactive-Run, authority-denied, conflicting-observation, and
  idempotent-replay coverage.
- `internal/workflow/run_inspect.go` - ADDED under DEC-194; the read-only
  `{run,reconciliation,liveness}` composition over existing `loadRun`,
  `project.ObserveRepository`, `ReconcileRepository`, the latest accepted
  checkpoint as comparison base, and the same-package `latestLiveness` reused
  directly, never duplicated. Appends no Event and writes nothing.
- `internal/workflow/run_inspect_test.go` - ADDED under DEC-194; repeated
  unchanged inspection yields the same reconciliation digest, a checkpoint moves
  the comparison base, and inspection leaves refs/index/stash/worktree byte-inert.
- Only the minimum additional `internal/workflow` edit needed to register
  these operations. No other domain change.
- `cmd/fid/run.go`
- `cmd/fid/run_recovery.go`
- `cmd/fid/claim.go`
- `cmd/fid/response_errors.go`
- `internal/cli/`
- `cmd/fid/integration/p22_recovery_journey_test.go`
- `cmd/fid/integration/p22_fake_test.go` - shared semantic fake for implemented cases only.
- `cmd/fid/integration/testdata/p22_conformance/`
- Matching in-package tests and `cmd/fid/integration/` contract fixtures for
  the behavior above, including retained rows invalidated by this change.
- `project-docs/docs/phases/phase-22/tasks/task-6-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-22/contracts.md` - only the sections named above.
- `project-docs/docs/phases/phase-22/public-interface.md` - owned input/result/error types.
- `project-docs/docs/phases/phase-22/schema.md` - owned DDL and Event rows.
- `project-docs/docs/phases/phase-22/execution-scope.md` and `v1-convergence-inventory.md`.
- Canonical source sections cited by this task's validation rows.
- `project-docs/docs/phases/phase-22/tasks/task-5-done.md` - actual predecessor hand-back.

### Out Of Scope

- Budget/Attention CLI leaves, import, final capability enablement, and future-feature stubs.
- Any `internal/workflow` behavior beyond the four DEC-194 files; reopening,
  re-reporting, or editing task 4; and any change to VAL-013..016 ownership.
- Product policy, protected gate definitions, workflow controls, dependencies,
  and unrelated implementation changes beyond this accepted task.

## Acceptance Criteria

- [ ] A built-binary recovery journey works by task 6. (VAL-022)
- [ ] Recovery leaves enforce their closed wire schemas. (VAL-023)
- [ ] Recovery failures have deterministic typed process exits. (VAL-024)
- [ ] Full required commands terminate; preserve red/green and failed-attempt
  evidence with output tokens and observed exit status.
- [ ] Every converged retained row identifies introducing commit/phase/task,
  invalidating change, and the exact command that executes it.
- [ ] `run.checkpoint` has a real domain producer whose accepted Effects and
  preconditions are proven in `internal/workflow`, and `fid run checkpoint`
  exposes it. (VAL-022/023, DEC-194)
- [ ] `run.inspect` has a real read-only domain composition returning
  `{run,reconciliation,liveness}` with a stable repeated digest, and
  `fid run inspect <run-id>` exposes it. (VAL-022/023, DEC-194)
- [ ] Each of the four DEC-194 files records accepted owner phase-22/task-4,
  that no commit ever created it, the DEC-194 transfer, and the exact command
  that executes it.

## First Failing Test

TestP22EarlyRecoveryJourney starts a Run, records termination, expires its lease, and requires a linked successor from the real CLI.
Capture the failure before product edits; then implement the smallest behavior
that satisfies it. Build failures alone do not prove the intended rejection.

## Documentation Updates

- `project-docs/docs/phases/phase-22/tasks/task-6-done.md`: implementation,
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

- Start after task 5 closes.
- One task per Developer session; recovery may continue the same task.
- QA Review window 2 follows this task; wait for its required disposition.
- Commit using Conventional Commits with `Refs: phase-22/task-6`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
