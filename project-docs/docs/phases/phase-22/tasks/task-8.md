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

- APPROVED under DEC-188; execute after the Runner's required QA readiness check.
- contracts.md section 7; public-interface.md run.close/run.abandon.
- Cancellation, timeout, signal, spawn failure, and harness failure are distinct non-success observations.
- DEC-184: add exactly run.close/run.abandon actions; abandonment requires a
  non-null existing observation just as recovery does. Keep run.finish only
  until task 9 removes its leaf, action, and current Event acceptance.
- DEC-186/R-14: repair the existing signal/timeout cases in
  TestGateRunRecordsExactNonSuccessTerminalOutcomes while adding cancellation.
  After self-signalling, the helper must block until observed termination;
  it must not return into os.Exit(0). The timeout fixture uses a blocking child,
  verifies it spawned, and leaves adequate startup time before its actual gate
  deadline. Synchronize on child readiness, not arbitrary sleep or signal order.
  Include bounded CPU-load coverage with cleanup of helpers. Preserve exact
  outcome/exit/signal assertions and production timeout semantics; never accept
  spawn-failed as a timeout, disable parallel tests, or retry until green.
- A lost executor response cannot manufacture a Receipt or imply pass.
- Run success requires shared evaluator success and qualifying executor/Claim authority; abandonment requires terminal proof or exact force authority.
- `execution-scope.md` is incorporated into this scope: it assigns only
  subject-preserving convergence of assertions displaced by this task's accepted
  change. All owned ordinary and tagged code executes under the commands below.
- DEC-192 progressive fence: `internal/journal/` is already in scope. Delete
  exactly the `"run.closed"` and `"run.abandoned"` keys from
  `unavailableRecoveryEvents` in `internal/journal/project_effects.go` as part
  of implementing those Events, and prove live-apply/replay parity for each.
  Leave every other key fenced.

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
- `internal/workflow/activation.go` - only close/abandon action admission.
- `cmd/fid/gate_run_test.go` - cancelled outcome and bounded R-14 fixture repair.
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
- [ ] The signal and timeout regressions retain their exact subjects and pass
  with synchronized subprocesses under load. Record per-row provenance using
  DEC-186's corrected b694991/76ad781 history, red/green, and every attempt.
  The focused repeated command and two consecutive full benchmark-gate runs
  terminate successfully; no discarded first failure. (VAL-031)
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
go test -race -run '^TestGateRunRecordsExactNonSuccessTerminalOutcomes$' -count=15 ./cmd/fid
go test -bench . -benchmem ./...
go test -bench . -benchmem ./...
```

No exception to the full task gate. Commands whose expected pass is nonzero:
none. Focused red/green selections do not replace these commands. Report each
attempt using the command's own output token and observed termination, as
required by AGENTS.md. No interrupted/hung command is PASS; preserve all failed
attempts and rerun the byte-identical command after authorized causal repair.
Report lint is a Runner control after the report commit, not a Developer gate.

## Notes For Developer

- Start after task 7 closes.
- If an earlier task integrated the bounded R-14 repair under DEC-186, inspect
  its actual commit and retain its evidence; verify it here and converge the
  cancellation extension. Do not reintroduce the defect to manufacture red.
- One task per Developer session; recovery may continue the same task.
- This task belongs to QA Review window 3.
- Commit using Conventional Commits with `Refs: phase-22/task-8`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
