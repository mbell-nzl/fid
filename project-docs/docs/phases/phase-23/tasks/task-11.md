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
  qualifying durable fact, and each carries exactly its own discriminator field
  and no other: `finding_id` for `finding-resolved`, `decision_class` for
  `decision-satisfied`, `outcome_class` for `human-outcome-recorded`, and
  **none** for `review-recorded` (DEC-236).
- VAL-049: `attention.resolve` keeps its exact input shape and retains history.
- VAL-050: A qualifying Human Outcome resolves only matching Stops.

## Binding Constraints

- `contracts.md` section 11 is the exact authority. The closed Condition type
  set is exactly `review-recorded`, `finding-resolved`, `human-outcome-recorded`,
  and `decision-satisfied`. Three of the four name one discriminator field each
  and only that one; **`review-recorded` names none** - section 11 enumerates
  the pairing and assigns `finding_id` to `finding-resolved` alone (DEC-236,
  resolving QA R-26). Four types and three fields do not divide evenly, and
  `review-recorded` is the type without one: its qualifying fact is a
  `review.recorded` Event scoped by `subject`, which needs no further
  discriminant. Do not add a field to it to make the arithmetic even. Tasks 9
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
- **No `current_event_catalog.go` edit in this task** (I-28, DEC-229): task 1
  registered `human-outcome.recorded`'s complete `schema.md` 6.2 rule including
  every optional row this task's resolution effects use, and task 10 owns that
  family's typed-result registries (DEC-231). If a needed optional row is
  missing, that is a task 1 defect - stop and escalate rather than editing a map
  task 1 owns. This task unfences nothing and registers no typed result.
- `internal/workflow/stop.go` - the closed Stop registries (R-14). Add the two
  remaining Condition types, `review-recorded` and `finding-resolved`, to
  `conditionTypes`, and the `finding_id` field to the `Condition` struct with
  its discriminator rule in `validateCondition`. **Tag the new
  field `omitempty`, and add `omitempty` to nothing else** (DEC-235); tasks 9
  and 10 set this pattern. `budget_name` and `capability_id` keep
  present-and-null treatment and gain no `omitempty`: tagging either changes
  every stored `semantic_key` and `request_digest` in existence, and task 9's
  six-row byte-identity bar is the check that catches it (I-39). `finding_id`
  belongs to **`finding-resolved` only**; `review-recorded` carries no
  discriminator field at all (DEC-236). Completing the closed set therefore
  makes `validateCondition`'s rule exact rather than uniform: **at most one**
  discriminator field is non-null, it is exactly the one matching the type for
  `finding-resolved`, `decision-satisfied`, and `human-outcome-recorded`, and
  **none** of the three is set for `review-recorded`, whose Condition carrying
  any of them is rejected. Tasks 9 and 10
  already added `decision-satisfied` and `human-outcome-recorded`, because each
  must exist before its own raising Event can commit; this task completes the
  closed four-type set and owns their resolution behavior (VAL-078). Add no
  Stop code.
- `internal/journal/current_event_value_relations.go` - **do not touch
  `case "Condition"`** (DEC-235, restating inventory row 8). ESC-23-001 is
  resolved: exclusivity is structural, enforced by the per-type shapes task 9
  created, so there is no relation work for these fields. The branch stays as it
  is because it still carries the Phase 22 `budget_name` / `capability_id`
  binding, and the accepted rows keep both fields `?ID` on every type, so this
  relation is their only enforcer (I-40). **Confirm in the completion report
  that the relation still fires** from task 9's dispatch: a `finding-resolved`
  Condition carrying a stray `budget_name` is rejected `condition nullable
  binding mismatch`. No test in the repository asserted that message before task
  9, so a dispatch that returned after shape validation would have deleted the
  binding silently (QA R-24); if it does not fire, that is a task 9 defect -
  stop and escalate rather than adding the call here. This task closes the
  four-type set, so also confirm that all four new types and all three new
  fields are covered **by shape rows**, and that the six Phase 22 rows are
  unchanged.
- `internal/journal/current_event_value_catalog.go` - the per-type `Condition`
  shape map task 9 created, gaining its **last two** rows (R-14, DEC-235):
  `type:=review-recorded subject:Subject budget_name:?ID capability_id:?ID` -
  four fields, with **no** `finding_id`, because `contracts.md` section 11 gives
  that type no discriminator (DEC-236) - and
  `type:=finding-resolved subject:Subject budget_name:?ID capability_id:?ID
  finding_id:ID`. The shapes are closed and field-count-exact, so this is what
  makes a `review-recorded` Condition carrying `finding_id` reject with `closed
  object has wrong shape`. Ten rows exist after
  this task - the six Phase 22 types unchanged and the four judgment types - and
  the map is closed. Add rows; do not touch the six Phase 22 rows.
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
- [ ] All four Condition types resolve only on their own qualifying durable
  fact: for each of `review-recorded`, `finding-resolved`, `decision-satisfied`,
  and `human-outcome-recorded`, a resolution attempt with a non-qualifying fact
  Event leaves the Stop open and rejects, and the qualifying Event resolves it.
  **And** the closed-schema half, per type: for `finding-resolved`,
  `decision-satisfied`, and `human-outcome-recorded`, a Condition carrying a
  second discriminator field, the wrong one, or none of them rejects with no
  Event; and for `review-recorded`, which carries none, a Condition carrying
  `finding_id` - or either other discriminator - rejects with no Event, while
  the bare four-field form is accepted (DEC-236). Under DEC-235 every one of
  those is a **shape** rejection, not a nullable-binding relation rejection -
  `closed object has wrong shape` for a field the row does not declare,
  `missing required field` for the wrong field, and `invalid ID` for an own
  field that is null. Evidence citing the relation message would mean the
  encoding was implemented wrong. (VAL-078)
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
