# Phase 23 - Task 11: Attention resolution through Human Outcomes

## Objective

Complete the closed four-type judgment resolution Condition set and make a
qualifying Human Outcome resolve exactly the Stops whose recorded condition
matches, without changing the `attention.resolve` input shape or erasing the
original request.

## Validation Contract Coverage

- VAL-047: `attention.v1` reports `human_outcome_resolution:true` only once the
  resolution path exists.
- VAL-048: A `human-outcome-recorded` condition resolves only with a matching
  class and result.
- VAL-078: All four Condition types - `review-recorded`, `finding-resolved`,
  `decision-satisfied`, `human-outcome-recorded` - resolve only on their own
  qualifying durable fact, and exactly one Condition field is non-null per type.
- VAL-049: `attention.resolve` keeps its exact input shape and retains history.
- VAL-050: A qualifying Human Outcome resolves only matching Stops.

## Binding Constraints

- `contracts.md` section 11 is the exact authority. The closed Condition type
  set is exactly `review-recorded`, `finding-resolved`, `human-outcome-recorded`,
  and `decision-satisfied`, each with only its matching non-null field. Tasks 9
  and 10 register `decision-satisfied` and `human-outcome-recorded` when they
  raise the Stops that carry them; this task registers the remaining two and
  owns the resolution behavior of all four (R-14).
- Stop resolution still requires a durable qualifying fact Event; a caller
  assertion never resolves a Stop.
- `attention.resolve` keeps its exact existing input: the sorted unique set of
  containing disposition Event IDs for all linked Stops. Do not add a field, a
  bulk mode, or an Outcome-keyed variant.
- Resolution references the qualifying Outcome. Never erase, edit, or back-date
  the original raise facts, Stop links, or deduplication key (PRD-HUM-003).
- Recording a qualifying Outcome must not close unrelated Attention. Prove the
  negative with a second open Stop carrying a different condition.
- `review-recorded` and `finding-resolved` are registered here but no Phase 23
  Event raises a Stop carrying them. `validateStopRaise` checks code and
  condition type independently, so VAL-078's per-type probes are constructible
  through the existing `fid stop raise` leaf with an accepted code and the
  condition under test.
- This task implements the path. The capability limit flips to true only in
  task 17, after every replacement proof passes; do not change the matrix here.
- `execution-scope.md` is incorporated on the usual convergence terms.

## Scope

### Files To Create Or Modify

- `internal/workflow/resolve.go`, `evaluation.go`, `attention` readers.
- `internal/review/` - the resolution effects emitted with a Human Outcome.
- `internal/journal/` - the `human-outcome.recorded` optional effect rows.
- `internal/workflow/stop.go` - the closed Stop registries (R-14). Add the two
  remaining Condition types, `review-recorded` and `finding-resolved`, to
  `conditionTypes`, and the `finding_id` field to the `Condition` struct with
  its exactly-one-non-null-field rule in `validateCondition`. Tasks 9 and 10
  already added `decision-satisfied` and `human-outcome-recorded`, because each
  must exist before its own raising Event can commit; this task completes the
  closed four-type set and owns their resolution behavior (VAL-078). Add no
  Stop code.
- `internal/journal/current_event_value_catalog.go` - the closed `Condition`
  value grammar gaining `review-recorded`, `finding-resolved`, and
  `finding_id:?ID` (R-14).
- `cmd/fid/attention.go`, `stop.go` - result shapes only.
- Matching in-package tests and `cmd/fid/integration/` fixtures invalidated by
  this change.
- `project-docs/docs/phases/phase-23/tasks/task-11-done.md`.

### Files To Read For Context

- `contracts.md` sections 9, 10, 11; `public-interface.md` sections 2 and 4;
  `schema.md` sections 1 and 6.2.
- `project-docs/docs/phases/phase-22/contracts.md` section 5 for the retained
  Stop and Attention semantics.
- `task-10-done.md`; `execution-scope.md`; `contract-coverage.md`.

### Out Of Scope

- The final capability matrix and every other availability change.
- New Stop codes beyond those already added by tasks 6, 9, and 10.

## Acceptance Criteria

- [ ] A wrong-class and a wrong-result Outcome both leave the Stop open and
  reject; the matching Outcome resolves it. (VAL-048)
- [ ] The `attention.resolve` input shape is unchanged and all original raise
  facts and Stop links remain readable after resolution. (VAL-049)
- [ ] With two open Stops carrying different conditions, a qualifying Outcome
  resolves exactly one. (VAL-050)
- [ ] Capability discovery still reports `human_outcome_resolution:false` at the
  end of this task, with the working path proven separately. (VAL-047, first
  half)
- [ ] Full required commands terminate with preserved red/green evidence.
- [ ] Every converged retained row records its provenance and executing command.

## First Failing Test

`TestHumanOutcomeResolvesOnlyMatchingStopCondition` fails against the absent
judgment resolution conditions.

## Documentation Updates

- `project-docs/docs/phases/phase-23/tasks/task-11-done.md` in the full format.
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
go test -race -timeout 30m -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go test -tags integration -race -timeout 30m ./cmd/fid/integration/...
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

- None.

## Notes For Developer

- The independent QA Lead checkpoint runs after this task and must PASS
  before task 12 is dispatched. Leave the workspace and evidence in a state
  an independent reviewer can reproduce without your session.
- This task belongs to QA Review window 4.
- One task per Developer session; recovery may continue the same task.
- Commit using Conventional Commits with `Refs: phase-23/task-11`.
- Never infer ownership from a historical filename or task label. No-progress
  budgets exclude interruptions and attempts forbidden from repairing the cause.
- Do not read held-out scenario material or activate another agent.
