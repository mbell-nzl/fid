# Phase 23 - Task 2: Typed multi-Obligation Envelope contract

## Objective

Replace the single hard-coded `gate-evidence` Obligation with 1..16 typed
Obligations of three kinds, add protected-Decision authoring storage, and make
Obligation evaluation return one independent result per Obligation.

## Validation Contract Coverage

- VAL-006: Envelope accepts 1..16 Obligations and 0..16 protected Decisions and
  enforces each type's exact field combination and the protected-Decision
  authoring rules (0, 16, and 17 entries; duplicate class; duplicate ordinal;
  an affected action outside the guarded set).
- VAL-007: Duplicate IDs, a second `gate-evidence` Obligation, duplicate
  `review_profile` or `outcome_class`, and a non-primary subject each reject.
- VAL-008: Evaluation returns one result per Obligation in ID order with
  independent dimensions and every unsatisfied reason. Evidence here is the
  shared domain evaluator and `envelope.activate`'s `Evaluation[]`; the
  `evidence.evaluate` leaf re-check belongs to task 6, which owns that leaf's
  optional Obligation ID.
- VAL-009: A guarded action requires every applicable Obligation to be satisfied.
  This task owns the blocking half only - a satisfied `gate-evidence` Obligation
  with an unsatisfied `review` Obligation - because review and Human Outcome
  evidence cannot be linked before tasks 6 and 10. Task 10 completes the
  all-satisfied half.
- VAL-010: `artifact_availability:"not-applicable"` is returned exactly for
  `metadata-sufficient` Obligations.

Full sources and criteria are in
`project-docs/docs/phases/phase-23/qa/validation-contract.md`.

## Binding Constraints

- `contracts.md` section 2 is the exact authority for the per-type field
  combinations; `public-interface.md` section 2 is the exact wire shape.
- Do not widen `guarded_actions` beyond `envelope.complete` and
  `issue.transition-terminal-success`. At most one Obligation of each type may
  guard `issue.transition-terminal-success`.
- `envelope_obligations.review_profile` and `outcome_class` are typed columns
  written for their own Obligation type only; the full Obligation value still
  lives in `obligation_json`.
- Evaluation must keep absence, record validity and freshness, and Artifact
  availability serialized as three distinct dimensions (PRD-EVID-002). Never
  collapse them into one reason code, and always return every unsatisfied
  reason from every Obligation (PRD-EVID-006).
- `envelope_protected_decisions` rows are written here as accepted authoring
  data. Their blocking, Stop, and satisfaction semantics belong to task 9; do
  not implement or assert them in this task.
- Do not assert the absence of review, Human Outcome, or review-cycle behavior;
  those arrive in tasks 4, 6, and 10. Assert only what this task delivers.
- Envelope revision stays 1 and the primary Issue stays single; multi-revision
  and multi-Issue remain Phase 24.
- `execution-scope.md` is incorporated: subject-preserving convergence only for
  assertions this accepted change displaces.

## Scope

### Files To Create Or Modify

- `internal/workflow/envelope.go`, `evaluation.go`, `resolve.go`, `errors.go`.
- `internal/evidence/evidence.go`, `policy.go`, `execution.go`.
- `internal/journal/` - the `envelope.created` effect rule row only.
- `cmd/fid/envelope.go`, `evidence.go` - request and result shapes only.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-2-done.md`.

### Files To Read For Context

- `project-docs/docs/phases/phase-23/contracts.md` sections 1, 2, 13.
- `project-docs/docs/phases/phase-23/public-interface.md` sections 2, 3, 4.
- `project-docs/docs/phases/phase-23/schema.md` sections 1, 2, 6.2.
- `project-docs/docs/phases/phase-23/tasks/task-1-done.md` - the actual
  predecessor hand-back, including which fence keys remain closed.
- `execution-scope.md`, `contract-coverage.md`.

### Out Of Scope

- Review, Finding, Concern, Deviation, Decision, and Human Outcome behavior.
- Protected-Decision blocking, Stops, and satisfaction.
- Claim modes, budgets, freshness changes, and Human text.

## Acceptance Criteria

- [ ] Every accepted per-type field combination succeeds and every prohibited
  one rejects with no Event. (VAL-006)
- [ ] All five negative authoring cases reject with an unchanged projection.
  (VAL-007)
- [ ] Evaluation returns one ordered result per Obligation with independent
  dimensions and the complete reason set. (VAL-008)
- [ ] A guarded action is permitted only when every applicable Obligation is
  satisfied. (VAL-009)
- [ ] `not-applicable` appears exactly for `metadata-sufficient` Obligations
  and never substitutes for `available` or `missing`. (VAL-010)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records introducing commit, invalidating
  change, and its executing command.

## First Failing Test

`TestEnvelopeCreateAcceptsTypedObligationSet` fails against the current
single-Obligation validator. Capture the failure before product edits.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-2-done.md` in the full format.
- Flag spec impact against the named contracts; only Tech Lead accepts a
  contract or schema change. Record no spec change when implementation conforms.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
go test -race -timeout 30m -cover ./internal/workflow/
```

Two gates run for minutes and are explicitly bounded (DEC-227): the `-timeout
30m` above replaces Go's implicit 10-minute per-binary default, which the tagged
suite already approaches at 560 s. Record in the completion report the wall time
`go test` reports for `cmd/fid/integration`, `cmd/fid`, and `internal/workflow`,
so a regression is attributed to the task that caused it. Budgets: tagged
integration under 900 s and no ordinary package over 600 s (VAL-079). A breach is
a bounded speed repair owned by this task when this task's own new journeys
caused it; otherwise stop and escalate. Never record a timed-out gate as `PASS`.

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- `go test -race -timeout 30m -cover ./internal/workflow/` before and after the edits, because this task adds branch-heavy validation to an existing package.

## Notes For Developer

- This task belongs to QA Review window 1.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-2`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.

## Repair Amendment - DEC-249 (2026-09-26, ESC-23-004 BC-1)

This is a bounded conformance repair after task close. It runs as step 2 of
DEC-255's order.

- **Change:** Make `issue.transition-terminal-success` (`evaluation.go`) and
  `envelope.complete` (`envelope.go`) evaluate only the Obligations whose
  `guarded_actions` contain that action. Reads keep reporting every
  Obligation. Do not change readiness raising, `run.close`, reason codes,
  precedence or denial shapes.
- **Evidence (VAL-009 task-2 half):** Use one Envelope with a satisfied
  `gate-evidence` Obligation and an unsatisfied `review` Obligation guarding
  both actions.
  - Both actions are denied, naming the review Obligation, with no Event and
    no projection change.
  - Contrast: an unsatisfied Obligation guarding only `envelope.complete`
    leaves terminal success unblocked and still blocks completion.
  - Provide both an in-package test and a built-binary observation.
- **Convergence:** This change displaces retained assertions, including the
  task-10 journey's stage 1-3 `qa_pass` rows. Converge each one here, with
  introducing-commit provenance and its executing command.
- **Scope and gates:** Use the existing Scope. The exact commands above are
  unchanged and execute every touched file.
- **Report:** Add a dated section to `task-2-done.md`. Replace the
  unsupported `:34` probe claim with the new citations.
