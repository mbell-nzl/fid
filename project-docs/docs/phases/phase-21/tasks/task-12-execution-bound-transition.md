# Phase 21 - Task 12: Shared Evaluator And Nonterminal Transition

## Objective

Implement the execution-bound nonterminal transition predicate and one shared Obligation evaluator used fail-closed by both terminal-success transition and completion. Positive trusted Evidence remains Group 2 work.

## Validation Contract Coverage

- VAL-038: owning observable in the Phase 21 validation contract.
- VAL-039: owning observable in the Phase 21 validation contract.
- VAL-040: owning observable in the Phase 21 validation contract.
- VAL-041: owning observable in the Phase 21 validation contract.

## Binding Constraints

- Terminal-success cannot succeed in this task because qualifying Evidence does not exist yet.
- Return all reasons in fixed kebab-case order with atomic failure.
- Apply CF-032 invalid revision-effect tests first.

## Scope

### Files To Create Or Modify

- `internal/workflow/`
- `internal/issues/`
- `internal/policy/`
- `internal/validation/`
- `internal/cli/`
- `internal/dev/` (retained-test convergence only; see DEC-167)
- `cmd/fid/` (including the `cmd/fid/integration/` package)
- `project-docs/docs/phases/phase-21/tasks/task-12-done.md`

### Files To Read For Context

- `project-docs/docs/phases/phase-21/contracts.md §§6,7,10`
- `project-docs/docs/phases/phase-21/public-interface.md`
- `project-docs/docs/phases/phase-21/v1-convergence-inventory.md`
- `project-docs/docs/phases/phase-21/qa/validation-contract.md` VAL-030

### Out Of Scope

- Evidence linking, positive terminal transition, and completion mutation.
- Generic Claim/Run implementation and semantics; final V1 removal and
  migration 0014; removal of the `internal/dev` source package.
- This does **not** exclude converging a retained assertion that an accepted
  Phase 21 change invalidated, wherever that assertion lives. See DEC-167.

### Convergence Ownership (DEC-162, DEC-165, DEC-166, DEC-167)

Convergence ownership follows the invalidating accepted change, not the fixture
file's primary subject or package. Two obligations apply to this task.

**1. This task's own collateral.** Task 12 replaced the retained non-bound Issue
transition's V1 Claim precondition with attribution-only structural transitions.
Converge every assertion that change invalidated, including
`internal/dev/start_test.go` `TestStartRejectsInvalidWorkflowEdgeBeforeAppend`.
Drop only the stale *"want missing claim rejection"* expectation on
`issues.Transition`. Keep the `Start()` subject, the
`*issues.InvalidTransitionError` assertion, the `CurrentState` check, and the
journal-unchanged check byte-identical. Leave every other `internal/dev`
assertion untouched for task 20.

**2. Designated owner for the tagged integration package.** Task 12 is the first
task whose required gate runs `cmd/fid/integration/`, so DEC-167 assigns it the
accumulated retained-V1 convergence in that package. Converge every retained row
there that an accepted Phase 21 change invalidated, until
`go test -tags integration -race ./cmd/fid/integration/...` is green. Known rows
at the time of this amendment - a navigation aid, not a limit; the gate is the
authority:

- `e2e_v01_closure_test.go:38` - `issue ac verify --status pass --by qa_lead`
  runs as a success call. VAL-030 (task 9) makes this leaf unconditionally
  `capability-unavailable` (`cmd/fid/issue_ac.go:109-111`). Convert it to an
  explicit fail-closed rejection with no journal append. Changing the flags alone
  does not fix it.
- `e2e_v01_closure_helpers_test.go:263` and `:266` - `concern raise` and
  `concern address` are the same VAL-030 fail-closed family
  (`cmd/fid/concern.go:101` and `:83`). Same conversion. Note that
  `closureRaiseAndAddressConcern` returns IDs the journey later consumes;
  restructure the journey rather than fabricating a value.
- Preserved task-8 Event-actor and `issue create --by` assertions in
  `TestP20StandaloneIssueJourneyJSONAndRebuildParity`,
  `TestV01ClosureJourneyAndRebuildParity`, and
  `TestV02aWorkflowAuthoringExitScenarioAndRebuildParity`, named in
  `task-9-done.md:122`.

Rules for both obligations:

- Record per-row provenance in the completion report: file and line, the
  invalidating accepted change, and the task that made it. A row whose
  provenance you cannot establish as an accepted Phase 21 change is not
  converged - record it as a concern instead.
- Preserve every subject and assertion strength. Deleting an assertion requires
  a named surviving test that already owns the accepted equivalent; name both in
  the report.
- Never restore `--by`, `--status`, V1 Claim behaviour, `fid policy`, role
  authority, or any removed writer, and never change production behaviour to
  make a retained fixture green.
- `t.Skip`, build tags, `-run` narrowing, threshold changes, and silent
  assertion loss are never convergence.
- Historical task numbering in a path, test name, or prose label is not
  current-phase ownership. See DEC-166.
- The tagged journey reveals one blocking step at a time. Reaching a later step
  is real progress; keep repairing rather than routing the next step elsewhere.
  If a required gate fails on a file outside this scope list, do not route it to
  another task in your report - establish provenance and, if an accepted Phase 21
  change invalidated it, name the scope gap in a `RECOVERY_REQUIRED` report.

## Acceptance Criteria

- [ ] Only all-qualified nonterminal transition commits. (VAL-038)
- [ ] Failures return all reasons with no effect. (VAL-039)
- [ ] Both guarded actions share the same unsatisfied evaluator result. (VAL-040)
- [ ] Invalid revision effects reject. (VAL-041)
- [ ] Every retained assertion invalidated by an accepted Phase 21 change is
      converged with its subject and strength preserved, and
      `go test -tags integration -race ./cmd/fid/integration/...` plus
      `go test -race -coverprofile=coverage.out ./...` are green. (No VAL: this
      is retained-fixture convergence under DEC-167, not a new observable. The
      converged rows preserve VAL-030's fail-closed contract rather than
      restating it.)

## First Failing Test

Add CF-032 tests first, then the nonterminal predicate and dual-call fail-closed evaluator goldens.

## Documentation Updates

- `project-docs/docs/phases/phase-21/tasks/task-12-done.md` - record red/green evidence, deviations, required commands, and spec impact.
- The completion report also records a per-row convergence table: file and line,
  the invalidating accepted change, and the task that made it. Name any deleted
  assertion's surviving accepted equivalent. Report any row whose provenance you
  cannot establish as a concern instead of converging it.

## Exact Commands To Run

Run exactly:

```text
go test -race ./internal/workflow ./internal/issues ./internal/policy ./internal/validation ./internal/cli ./cmd/fid/...
go test -tags integration -race ./cmd/fid/integration/...
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Every command is required and must terminate with session-captured evidence. Report lint is not a Developer task-gate command; after the report commit, the Phase Runner runs the phase brief's two-argument pre-dispatch control. A Ready QA plan may add a task-specific command; the Tech Lead must amend this section before authorization if it does.

## Notes For Developer

- Start after tasks 7 and 10-11 close. QA Review window 4 and the independent Group 1 checkpoint follow; task 13 cannot start without PASS.
- Completion report: `project-docs/docs/phases/phase-21/tasks/task-12-done.md`.
- Stop for a protected contract change; ordinary in-scope red/green repair remains this task.
- Do not read held-out scenario material.
- VAL-038..041 are already satisfied and their evidence is recorded in the
  existing report; preserve it. The remaining work is the convergence above.
- The tagged closure journey blocks one step at a time. Expect to fix a step,
  rerun, and meet the next one. That is progress, not a new blocker - keep going
  until the gate is green or you hit a genuinely protected condition.
- Required commands are unchanged and rerun byte-identically. They already
  exercise every file scope this task owns, so no command amendment is needed.
