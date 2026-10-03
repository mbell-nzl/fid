# Phase 24 - Task 7: Run leaves

## Objective

Converge the Run family under DEC-273:
- `run start`, `run close` and `run abandon` replace `repository` with a
  required `expected_repository`, and record fid's observation;
- `run start` compares the Envelope binding as *same worktree*;
- `run observe-budget` takes `measurement_ref` as a `SourceInput`;
- `run observe`, `run checkpoint` and `run recover` get descriptors (seven
  leaves leave the fence).

## Validation Contract Coverage

- VAL-009, VAL-011..016, VAL-017, VAL-018, VAL-019 (`measurement_ref`),
  VAL-022, VAL-023 (these leaves).
- VAL-024 (`run.start`, `run.close`, `run.abandon`): expectation check.
- VAL-025 (run leaves): expectation mismatch returns `repository-mismatch`;
  unsafe observation still returns `recovery-state-conflict`.
- VAL-029 (`run.start`): imported-destination equivalence.
- VAL-030: `run.start` is admitted after edits and commits since Envelope
  creation.
- VAL-031: close and abandon record the observed terminal repository.
- VAL-035, VAL-052.

Assertions live in `project-docs/docs/phases/phase-24/qa/validation-contract.md`.

## Binding Constraints

- Contract sections 4, 5 and 7; DEC-273.
- **`run.start`** checks `expected_repository` against the observation
  (task 5 helper). It checks the Envelope scope with the task 6 *same
  worktree* predicate, never the stored HEAD, tree or dirty state.
  `start_repository` is the observation. The Claim worktree check uses the
  same predicate.
- **`run.close` / `run.abandon`** replace `validateRunTerminalRepository`'s
  caller assertion with the expectation check. A mismatch is
  `repository-mismatch`. An unobservable or unsafe repository keeps the
  Phase 22 `recovery-state-conflict`. `terminal_repository` is the
  observation with `base_oid` = HEAD (the caller's `base_oid` and its
  existence check are gone). Close evaluates Obligations against that
  observation.
- **`run.observe-budget`**: `measurement_ref` is `?SourceInput`.
- **`run.recover`** already observes, so its input is unchanged.

## Scope

### Files To Create Or Modify

- `internal/workflow/{run,run_close,run_budget,resume}.go` (the
  `resumeObservedRepository` callers move to the shared observation
  function)
- `cmd/fid/{run,run_recovery}.go` and descriptor files for the seven leaves
- Retained tests (`generic_run_test.go` `runRepositoryForTest`,
  `run_close_*`, `p22_close_proof_test.go`, recovery journeys, fake
  adapters): convergence with provenance

### Files To Read For Context

- `project-docs/docs/phases/phase-24/public-interface.md` sections 4, 5
- `project-docs/docs/phases/phase-24/input-inventory.md` sections 1-4
- `project-docs/docs/phases/phase-24/decisions.md` DEC-273
- `project-docs/docs/phases/phase-22/public-interface.md` section 3 (Run
  leaves)

### Out Of Scope

- `gate.run` and Evidence (task 8).
- Recovery semantics, liveness and reconciliation (unchanged).

## Acceptance Criteria

- [ ] After the Envelope is created, commit and edit files; `run.start` with
  a matching expectation is accepted and records the observation, not the
  Envelope snapshot. (VAL-030, VAL-022)
- [ ] For each of start, close and abandon: match accepted; HEAD-differs and
  dirty-differs return `repository-mismatch` with no Event. (VAL-024)
- [ ] Close/abandon with an unsafe repository still return
  `recovery-state-conflict`. (VAL-025)
- [ ] The close and abandon `terminal_repository` equals the observation
  with `base_oid` = HEAD, and close evaluates against it. (VAL-031)
- [ ] `run.start` in an imported destination is accepted, and refused in a
  plain copy. (VAL-029)
- [ ] `measurement_ref` passes VAL-018/019; retries replay. (VAL-018,
  VAL-019, VAL-009)
- [ ] `run.start` and `run.observe-budget` with an unobservable repository
  return `repository-unavailable` with no Event. Close and abandon are
  covered by VAL-025. (VAL-023)
- [ ] The seven leaves leave the fence; fake/fid agree. (VAL-011..016,
  VAL-017, VAL-035)
- [ ] Wall times inside budget. (VAL-052)

## First Failing Test

`run start` after a tracked-file edit since Envelope creation, with a matching `expected_repository`, is accepted. Today the input requires the Envelope's snapshot block.

## Documentation Updates

- `project-docs/specs/06_DESIGN_PLAN_CLI.md` and
  `07_DESIGN_PLAN_FLOWS.md`: Run inputs and the rework-Run flow under
  DEC-273.

## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root the Runner records. No command is a placeholder.
Apply gofmt and goimports to edited Go files before the final gate. Every
command below is required before task close. The tagged integration gate is
mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -gcflags=modernc.org/...=-d=checkptr=0 -timeout 30m ./cmd/fid/integration/...
```

The `-timeout 30m` bounds follow DEC-227/259. Record the wall time `go test`
reports for `cmd/fid/integration`, `cmd/fid` and `internal/workflow`. The
budgets are tagged integration under 900 s and no ordinary package over 600 s
(VAL-052). A breach caused by this task's own new journeys is a bounded speed
repair owned by this task; otherwise stop and escalate. Never record a
timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands). Every required command
is reported `PASS` or `BLOCKED` with reason, with its own output token and
observed termination. Focused red/green selections never replace these
commands. Report lint is a Runner control, not a Developer gate.

List any conditional command required by this task and why:

- None.

## Notes For Developer

- Report path: `project-docs/docs/phases/phase-24/tasks/task-7-done.md`.
  Commit with `Refs: phase-24/task-7`.
- Red first: record the first failing test before changing product code.
- Retained-test convergence follows DEC-162/165/166. Record each displaced
  assertion with its introducing commit, phase and task. Never infer
  ownership from a filename or task label. Never delete or weaken an
  assertion to pass a gate.
- Any missing registry member, contract ambiguity, or need to change an
  accepted contract shape or assertion is a `BLOCKED` escalation candidate,
  not a local decision (DEC-195).
- Remove this task's leaves from the pending-descriptor fence (task 3). Each
  adopted leaf gets a descriptor, help, `--example`, typed template,
  exclusivity and structural-validation corpus rows (VAL-011..016).
