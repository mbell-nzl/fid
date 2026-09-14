# Phase 23 - Task 6: Review evidence, review cycles, and rework exhaustion

## Objective

Make a passing Review usable as Evidence for a `review` Obligation, charge
exactly one `review-cycles` sample per `changes-required` Review, and refuse the
next rework atomically when that budget is exhausted.

## Validation Contract Coverage

- VAL-024: `evidence.link` accepts a `pass` Review and rejects
  `changes-required` and `blocked`.
- VAL-025: An evidence type outside the Obligation's accepted list rejects.
- VAL-026: An unresolved blocking Finding keeps the Obligation unsatisfied even
  after a later `pass` Review.
- VAL-027: Exactly one `review-cycles` charge per `changes-required` Review,
  keyed on the Review ID, never double-charged.
- VAL-028: `run.observe-budget` rejects the fid-owned `review-cycles` unit.
- VAL-029: Exhaustion commits the Review, charge, Stop, and Attention
  atomically and refuses the next rework.
- VAL-030: `review-cycles` and `review-rework-exhausted` require each other and
  a `review-cycles` budget requires a `review` Obligation.

## Binding Constraints

- `contracts.md` sections 6 and 13 and kit dispositions C-03 and C-04 are the
  exact authority. Only `review-cycles` and `review-rework-exhausted` are added;
  generic `retry-exhausted` stays Phase 24.
- The charge source sample identity is the Review ID, so replay, linked
  recovery, and import cannot charge the same Review twice. Use the existing
  `budget_samples` uniqueness; add no new table or column.
- One `changes-required` Review is one cycle. Never derive a charge from a
  Finding count, a rework Claim, or a Run.
- Admission of the last allowed cycle is atomic with charging it and raising its
  Stop and Attention. The Review still commits; the Stop blocks the next rework.
- The exhaustion Stop uses code `review-rework-exhausted`, subject Envelope,
  resolver capability `stop.resolve`, blocked actions exactly `claim.acquire`
  and `run.start`, and resolution condition `budget-observed` on that budget.
- A `review` Obligation with no `review-cycles` budget is valid: unlimited
  rework is a recorded policy choice, not an error.
- `evidence.link` appends one immutable record per (Obligation, subject) pair
  and never overwrites prior proof. `evidence.evaluate` gains its optional
  Obligation ID here.
- Do not implement staling or Human Outcome evidence; those are tasks 7 and 10.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/evidence/evidence.go`, `policy.go`, `execution.go`.
- `internal/workflow/budget.go`, `run_budget.go`, `recovery_budget.go`,
  `evaluation.go`, `resolve.go`.
- `internal/review/` - the verdict effects that charge and raise.
- `internal/journal/` - the `review.recorded` effect-rule optional rows.
- `internal/cli/codes.go` - `review-changes-required`,
  `review-rework-exhausted`, `finding-unresolved`, `evidence-type-not-accepted`.
- `cmd/fid/evidence.go`, `run.go`, `claim.go` - denial responses only.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-6-done.md`.

### Files To Read For Context

- `contracts.md` sections 2, 4, 5, 6, 13; `public-interface.md` sections 2, 4, 5;
  `schema.md` sections 1, 6.2, 6.3.
- `project-docs/docs/phases/phase-22/contracts.md` section 6 for the retained
  budget arithmetic and denial rules.
- `task-5-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- Freshness and staling, Human Outcomes, protected Decisions, and Human text.

## Acceptance Criteria

- [ ] `pass` links and can satisfy the Obligation; `changes-required` and
  `blocked` reject `review-changes-required` with the blocking Finding IDs.
  (VAL-024)
- [ ] Only the three matching evidence-type and Obligation-type pairs link; the
  six mismatches reject with the accepted list in details. (VAL-025)
- [ ] A blocking unresolved Finding keeps the Obligation unsatisfied with reason
  `finding-unresolved` even after a later `pass` Review. (VAL-026)
- [ ] Consumption rises by exactly one per distinct `changes-required` Review
  and by zero on replay, recovery, and import. (VAL-027)
- [ ] A harness `run.observe-budget` against `review-cycles` rejects with no
  sample. (VAL-028)
- [ ] The exhausting Review, its charge, its Stop, and its Attention commit in
  one Event, and the next `claim.acquire` and `run.start` return
  `review-rework-exhausted` with no Claim or Run. (VAL-029)
- [ ] The three invalid budget and Obligation pairings reject; the matched pair
  and the unbudgeted `review` Obligation are accepted. (VAL-030)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestChangesRequiredReviewChargesExactlyOneReviewCycle` fails against the
absent unit and charging path.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-6-done.md` in the full format.
- Flag spec impact; only Tech Lead accepts a contract or schema change.
## Exact Commands To Run

Working directory: repository root `/home/mark/projects/fid`, or the isolated
same-task worktree root explicitly recorded by the Runner. No command is a
placeholder. Apply gofmt and goimports to edited Go files before the final gate.

After implementation, every command below is required before task close.
The tagged integration gate is mandatory and is never substituted.

```text
gofmt -l cmd internal
golangci-lint run ./...
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race ./cmd/fid/integration/...
go test -race -cover ./internal/evidence/
go test -race -cover ./internal/workflow/
```

Required-command rules live in `AGENTS.md` (Commands): every required command
is reported `PASS`, or `BLOCKED` with reason; missing, silently skipped, or
narrower substituted commands are failed gates. No exception to the full task
gate. Focused red/green selections never replace these commands. Report each
attempt with the command's own output token and observed termination. An
interrupted or hung command is never `PASS`; preserve all failed attempts and
rerun the byte-identical command after an authorized causal repair. Report lint
is a Runner control after the report commit, not a Developer gate.

List any conditional command required by this task and why:

- `go test -race -cover ./internal/evidence/` and `go test -race -cover ./internal/workflow/` before and after the edits, because this task adds branch-heavy evaluation and budget paths to two existing packages.

## Notes For Developer

- This task belongs to QA Review window 2.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-6`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
